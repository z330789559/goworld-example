//! Shared data structures and utility helpers for the Mahjong hallway/room stack.

use std::time::Duration;

use chrono::{DateTime, Utc};
use rand::Rng;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use tokio::time::Instant;

/// Session lifecycle managed by the Gateway.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Session {
    pub role_id: u64,
    pub lobby_id: Option<u64>,
    pub last_heartbeat: Instant,
    pub status: SessionStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum SessionStatus {
    Connected,
    Disconnected,
}

/// Lobby metadata stored inside the Gateway registry.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LobbyInfo {
    pub lobby_id: u64,
    pub address: String,
    pub max_capacity: u32,
    pub current_players: u32,
    pub current_rooms: u32,
    pub cpu_usage: f64,
    pub memory_mb: u64,
    pub version: String,
    pub weight: u32,
    #[serde(skip)]
    pub last_heartbeat: Instant,
    pub status: LobbyStatus,
}

impl LobbyInfo {
    /// Composite load score: lower values are preferred during balancing.
    pub fn load_score(&self) -> u64 {
        let cpu_score = (self.cpu_usage * 60.0) as u64;
        let player_ratio =
            (self.current_players as f64 / self.max_capacity as f64 * 40.0).round() as u64;
        cpu_score + player_ratio
    }

    pub fn is_healthy(&self) -> bool {
        self.status == LobbyStatus::Healthy
            && self.current_players < self.max_capacity
            && Instant::now().duration_since(self.last_heartbeat) < Duration::from_secs(15)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum LobbyStatus {
    Healthy,
    Unhealthy,
}

/// Registry RPC contract shared between services.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterLobbyRequest {
    pub server_name: String,
    pub rpc_address: String,
    pub max_capacity: u32,
    pub version: String,
    pub weight: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterLobbyResponse {
    pub lobby_id: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LobbyHeartbeat {
    pub lobby_id: u64,
    pub player_count: u32,
    pub room_count: u32,
    pub cpu_usage: f64,
    pub memory_mb: u64,
    pub timestamp: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UnregisterLobbyRequest {
    pub lobby_id: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RoomSummary {
    pub room_id: u64,
    pub owner_role_id: u64,
    pub seat_taken: u8,
    pub max_rounds: u16,
    pub pay_type: PayType,
    pub allow_spectator: bool,
    pub version: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PayType {
    Owner,
    Aa,
    Diamond,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RoomSnapshot {
    pub summary: RoomSummary,
    pub seats: [Seat; 4],
    pub state: RoomState,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Seat {
    pub index: u8,
    pub role_id: Option<u64>,
    pub status: SeatStatus,
    pub ready: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum SeatStatus {
    Empty,
    Occupied,
    Locked,
    Offline,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum RoomState {
    Waiting,
    ReadyCheck,
    Playing,
    Settling,
    Closed,
}

/// Simplified load balancer error domain.
#[derive(Debug, Error)]
pub enum GatewayError {
    #[error("no healthy lobby available")]
    NoHealthyLobby,
    #[error("session missing for role {0}")]
    SessionMissing(u64),
}

/// Gray release helper.
pub fn weighted_select_lobby(registry: &[LobbyInfo]) -> Option<u64> {
    let healthy: Vec<_> = registry.iter().filter(|info| info.is_healthy()).collect();
    if healthy.is_empty() {
        return None;
    }
    let total: u32 = healthy.iter().map(|info| info.weight).sum();
    if total == 0 {
        return None;
    }
    let mut rng = rand::thread_rng();
    let mut acc = 0u32;
    let target = rng.gen_range(0..total);
    for info in healthy {
        acc += info.weight;
        if target < acc {
            return Some(info.lobby_id);
        }
    }
    None
}

/// Gateway selection strategy for sticky sessions.
pub fn select_lobby(registry: &[LobbyInfo]) -> Result<u64, GatewayError> {
    registry
        .iter()
        .filter(|info| info.is_healthy())
        .min_by_key(|info| info.load_score())
        .map(|info| info.lobby_id)
        .ok_or(GatewayError::NoHealthyLobby)
}
