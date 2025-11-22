use anyhow::Result;
use common::{RoomSnapshot, Seat};
use tokio::time::Duration;
use tracing::info;
use tracing_subscriber::EnvFilter;

/// Very small simulation of a mahjong table loop.
async fn simulate_round(room: RoomSnapshot) -> Result<()> {
    info!("starting round for room {}", room.summary.room_id);
    tokio::time::sleep(Duration::from_millis(200)).await;
    let players: Vec<u64> = room
        .seats
        .iter()
        .filter_map(|seat: &Seat| seat.role_id)
        .collect();
    info!("round finished, players={:?}", players);
    Ok(())
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let snapshot = RoomSnapshot {
        summary: common::RoomSummary {
            room_id: 500,
            owner_role_id: 1001,
            seat_taken: 4,
            max_rounds: 8,
            pay_type: common::PayType::Owner,
            allow_spectator: true,
            version: 1,
        },
        seats: [
            Seat {
                index: 0,
                role_id: Some(1001),
                status: common::SeatStatus::Occupied,
                ready: true,
            },
            Seat {
                index: 1,
                role_id: Some(1002),
                status: common::SeatStatus::Occupied,
                ready: true,
            },
            Seat {
                index: 2,
                role_id: Some(1003),
                status: common::SeatStatus::Occupied,
                ready: true,
            },
            Seat {
                index: 3,
                role_id: Some(1004),
                status: common::SeatStatus::Occupied,
                ready: true,
            },
        ],
        state: common::RoomState::Playing,
        created_at: chrono::Utc::now(),
        updated_at: chrono::Utc::now(),
    };

    simulate_round(snapshot).await?;
    Ok(())
}
