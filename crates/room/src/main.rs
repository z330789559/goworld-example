use anyhow::{Context, Result};
use common::{PayType, RoomSnapshot, RoomState, RoomSummary, Seat, SeatStatus};
use redis::AsyncCommands;
use tokio::time::Instant;
use tracing::info;
use tracing_subscriber::EnvFilter;

#[derive(Clone)]
struct RoomService {
    redis: redis::Client,
}

impl RoomService {
    async fn create_room(&self, owner_role_id: u64, pay_type: PayType) -> Result<RoomSnapshot> {
        let mut conn = self.redis.get_async_connection().await?;
        let room_id: u64 = conn.incr("room:id_seq", 1).await?;
        let summary = RoomSummary {
            room_id,
            owner_role_id,
            seat_taken: 1,
            max_rounds: 8,
            pay_type,
            allow_spectator: true,
            version: 1,
        };
        let snapshot = RoomSnapshot {
            summary,
            seats: [
                Seat {
                    index: 0,
                    role_id: Some(owner_role_id),
                    status: SeatStatus::Occupied,
                    ready: false,
                },
                Seat {
                    index: 1,
                    role_id: None,
                    status: SeatStatus::Empty,
                    ready: false,
                },
                Seat {
                    index: 2,
                    role_id: None,
                    status: SeatStatus::Empty,
                    ready: false,
                },
                Seat {
                    index: 3,
                    role_id: None,
                    status: SeatStatus::Empty,
                    ready: false,
                },
            ],
            state: RoomState::Waiting,
            created_at: chrono::Utc::now(),
            updated_at: chrono::Utc::now(),
        };
        let payload = serde_json::to_string(&snapshot)?;
        conn.set(format!("room:snapshot:{room_id}"), payload).await?;
        info!("room {room_id} created");
        Ok(snapshot)
    }

    async fn join_room(&self, room_id: u64, role_id: u64) -> Result<RoomSnapshot> {
        let mut conn = self.redis.get_async_connection().await?;
        let key = format!("room:snapshot:{room_id}");
        let raw: String = conn.get(&key).await.context("room not found")?;
        let mut snapshot: RoomSnapshot = serde_json::from_str(&raw)?;
        if let Some(seat) = snapshot.seats.iter_mut().find(|seat| seat.status == SeatStatus::Empty)
        {
            seat.status = SeatStatus::Occupied;
            seat.role_id = Some(role_id);
            snapshot.summary.seat_taken += 1;
            snapshot.summary.version += 1;
            snapshot.state = RoomState::ReadyCheck;
            snapshot.updated_at = chrono::Utc::now();
            let payload = serde_json::to_string(&snapshot)?;
            conn.set(&key, payload).await?;
            Ok(snapshot)
        } else {
            anyhow::bail!("no empty seat")
        }
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let redis = redis::Client::open("redis://127.0.0.1/")?;
    let service = RoomService { redis };

    let room = service.create_room(1001, PayType::Owner).await?;
    info!("created room {}", room.summary.room_id);

    let updated = service.join_room(room.summary.room_id, 1002).await?;
    info!("player joined room {} -> {} seats", updated.summary.room_id, updated.summary.seat_taken);

    Ok(())
}
