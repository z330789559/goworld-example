export type GatewayEvent =
  | { type: "RoomEvent"; roomId: number; version: number }
  | { type: "Kick"; reason: string }
  | { type: "Heartbeat" };

/**
 * 非生产用：用于演示如何在 Cocos 中处理来自 Gateway 的事件。
 */
export function handleGatewayEvent(evt: GatewayEvent): void {
  switch (evt.type) {
    case "RoomEvent":
      console.log("room update", evt.roomId, evt.version);
      break;
    case "Kick":
      console.warn("kicked", evt.reason);
      break;
    case "Heartbeat":
      break;
  }
}
