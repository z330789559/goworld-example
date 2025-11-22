# 事件与协议草稿

- `EnterLobby`：请求进入大厅，服务器返回绑定 lobby_id 及初始订阅版本号。
- `SubscribeRoomList`：订阅房间列表，网关转发到当前 sticky lobby。
- `RoomEvent`：房间快照/增量推送，包含房间座位、人数、状态版本号。
- `CreateRoom` / `JoinRoom` / `SetReady`：房间服务 RPC，经网关路由。
