use std::{collections::HashMap, sync::Arc};

use anyhow::Result;
use common::{select_lobby, GatewayError, LobbyHeartbeat, LobbyInfo, LobbyStatus, RegisterLobbyRequest, RegisterLobbyResponse, Session, SessionStatus, UnregisterLobbyRequest};
use tokio::{sync::RwLock, time::{Duration, Instant}};
use tracing::{error, info, warn};
use tracing_subscriber::EnvFilter;

#[derive(Default)]
struct GatewayState {
    lobby_registry: HashMap<u64, LobbyInfo>,
    sessions: HashMap<u64, Session>,
    next_lobby_id: u64,
}

#[derive(Clone)]
struct Gateway {
    state: Arc<RwLock<GatewayState>>,
}

impl Gateway {
    fn new() -> Self {
        Self {
            state: Arc::new(RwLock::new(GatewayState::default())),
        }
    }

    async fn handle_register_lobby(&self, req: RegisterLobbyRequest) -> Result<RegisterLobbyResponse> {
        let mut state = self.state.write().await;
        let lobby_id = state.next_lobby_id;
        state.next_lobby_id += 1;

        let info = LobbyInfo {
            lobby_id,
            address: req.rpc_address,
            max_capacity: req.max_capacity,
            current_players: 0,
            current_rooms: 0,
            cpu_usage: 0.0,
            memory_mb: 0,
            version: req.version,
            weight: req.weight,
            last_heartbeat: Instant::now(),
            status: LobbyStatus::Healthy,
        };
        state.lobby_registry.insert(lobby_id, info);
        info!("registered lobby {}", lobby_id);
        Ok(RegisterLobbyResponse { lobby_id })
    }

    async fn handle_lobby_heartbeat(&self, heartbeat: LobbyHeartbeat) -> Result<()> {
        let mut state = self.state.write().await;
        if let Some(info) = state.lobby_registry.get_mut(&heartbeat.lobby_id) {
            info.current_players = heartbeat.player_count;
            info.current_rooms = heartbeat.room_count;
            info.cpu_usage = heartbeat.cpu_usage;
            info.memory_mb = heartbeat.memory_mb;
            info.last_heartbeat = Instant::now();
            info.status = LobbyStatus::Healthy;
        } else {
            warn!("heartbeat for unknown lobby {}", heartbeat.lobby_id);
        }
        Ok(())
    }

    async fn handle_unregister_lobby(&self, req: UnregisterLobbyRequest) -> Result<()> {
        let mut state = self.state.write().await;
        state.lobby_registry.remove(&req.lobby_id);
        info!("unregistered lobby {}", req.lobby_id);
        Ok(())
    }

    async fn select_or_reuse_lobby(&self, role_id: u64) -> Result<u64, GatewayError> {
        let mut state = self.state.write().await;
        if let Some(session) = state.sessions.get(&role_id) {
            if let Some(lobby_id) = session.lobby_id {
                if state
                    .lobby_registry
                    .get(&lobby_id)
                    .map(|info| info.is_healthy())
                    .unwrap_or(false)
                {
                    return Ok(lobby_id);
                }
            }
        }

        let registry: Vec<_> = state.lobby_registry.values().cloned().collect();
        let lobby_id = select_lobby(&registry)?;

        let entry = state.sessions.entry(role_id).or_insert(Session {
            role_id,
            lobby_id: None,
            last_heartbeat: Instant::now(),
            status: SessionStatus::Connected,
        });
        entry.lobby_id = Some(lobby_id);
        Ok(lobby_id)
    }

    async fn health_check_task(self) {
        let mut ticker = tokio::time::interval(Duration::from_secs(5));
        loop {
            ticker.tick().await;
            let mut state = self.state.write().await;
            let now = Instant::now();
            for info in state.lobby_registry.values_mut() {
                let elapsed = now.duration_since(info.last_heartbeat);
                if elapsed > Duration::from_secs(15) {
                    if info.status == LobbyStatus::Healthy {
                        warn!("lobby {} heartbeat timeout", info.lobby_id);
                    }
                    info.status = LobbyStatus::Unhealthy;
                }
            }
        }
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let gateway = Gateway::new();
    tokio::spawn(gateway.clone().health_check_task());

    // Demo run: register two lobbies and select one for a player.
    let lobby_a = gateway
        .handle_register_lobby(RegisterLobbyRequest {
            server_name: "lobby-a".into(),
            rpc_address: "127.0.0.1:7001".into(),
            max_capacity: 10000,
            version: "0.1.0".into(),
            weight: 50,
        })
        .await?;
    let lobby_b = gateway
        .handle_register_lobby(RegisterLobbyRequest {
            server_name: "lobby-b".into(),
            rpc_address: "127.0.0.1:7002".into(),
            max_capacity: 10000,
            version: "0.1.0".into(),
            weight: 50,
        })
        .await?;

    let lobby_id = gateway.select_or_reuse_lobby(42).await?;
    info!("player 42 assigned to lobby {}", lobby_id);

    // Heartbeat one lobby to keep it healthy; leave another to be marked unhealthy.
    gateway
        .handle_lobby_heartbeat(LobbyHeartbeat {
            lobby_id: lobby_a.lobby_id,
            player_count: 100,
            room_count: 25,
            cpu_usage: 10.5,
            memory_mb: 512,
            timestamp: chrono::Utc::now(),
        })
        .await?;

    tokio::time::sleep(Duration::from_secs(1)).await;
    gateway
        .handle_unregister_lobby(UnregisterLobbyRequest {
            lobby_id: lobby_b.lobby_id,
        })
        .await?;

    info!("gateway bootstrap complete");
    Ok(())
}
