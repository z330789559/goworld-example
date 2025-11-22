use anyhow::Result;
use async_trait::async_trait;
use common::{LobbyHeartbeat, RegisterLobbyRequest, RegisterLobbyResponse, RoomSnapshot};
use redis::AsyncCommands;
use tokio::time::{self, Duration};
use tracing::{info, warn};
use tracing_subscriber::EnvFilter;

#[async_trait]
pub trait GatewayApi: Send + Sync {
    async fn register_lobby(&self, req: RegisterLobbyRequest) -> Result<RegisterLobbyResponse>;
    async fn lobby_heartbeat(&self, heartbeat: LobbyHeartbeat) -> Result<()>;
}

/// Minimal in-memory gateway client used for demos.
pub struct MockGatewayClient;

#[async_trait]
impl GatewayApi for MockGatewayClient {
    async fn register_lobby(&self, req: RegisterLobbyRequest) -> Result<RegisterLobbyResponse> {
        info!("registering lobby {}", req.server_name);
        Ok(RegisterLobbyResponse { lobby_id: 1 })
    }

    async fn lobby_heartbeat(&self, heartbeat: LobbyHeartbeat) -> Result<()> {
        info!(
            "heartbeat: lobby={}, players={}, rooms={}",
            heartbeat.lobby_id, heartbeat.player_count, heartbeat.room_count
        );
        Ok(())
    }
}

struct LobbyService<G: GatewayApi> {
    gateway: G,
    lobby_id: u64,
    redis_client: redis::Client,
}

impl<G: GatewayApi> LobbyService<G> {
    async fn register(&mut self) -> Result<()> {
        let response = self
            .gateway
            .register_lobby(RegisterLobbyRequest {
                server_name: "lobby-1".into(),
                rpc_address: "127.0.0.1:7001".into(),
                max_capacity: 10_000,
                version: env!("CARGO_PKG_VERSION").into(),
                weight: 50,
            })
            .await?;
        self.lobby_id = response.lobby_id;
        Ok(())
    }

    async fn heartbeat_loop(&self) -> Result<()> {
        let mut interval = time::interval(Duration::from_secs(5));
        loop {
            interval.tick().await;
            let heartbeat = LobbyHeartbeat {
                lobby_id: self.lobby_id,
                player_count: 0,
                room_count: 0,
                cpu_usage: 0.0,
                memory_mb: 0,
                timestamp: chrono::Utc::now(),
            };
            if let Err(e) = self.gateway.lobby_heartbeat(heartbeat).await {
                warn!("heartbeat failed: {e}");
            }
        }
    }

    async fn load_rooms(&self) -> Result<Vec<RoomSnapshot>> {
        let mut conn = self.redis_client.get_async_connection().await?;
        let keys: Vec<String> = conn.keys("room:snapshot:*").await?;
        let mut rooms = Vec::new();
        for key in keys {
            let raw: String = conn.get(key).await?;
            if let Ok(snapshot) = serde_json::from_str::<RoomSnapshot>(&raw) {
                rooms.push(snapshot);
            }
        }
        Ok(rooms)
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let redis_client = redis::Client::open("redis://127.0.0.1/")?;
    let mut lobby = LobbyService {
        gateway: MockGatewayClient,
        lobby_id: 0,
        redis_client,
    };

    lobby.register().await?;
    tokio::spawn(lobby.heartbeat_loop());

    let rooms = lobby.load_rooms().await.unwrap_or_default();
    info!("bootstrapped {} rooms from redis", rooms.len());

    Ok(())
}
