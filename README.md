# 麻将大厅与房间制示例栈

本仓库实现了一个面向 Gateway / Lobby / Room / Mahjong Table 分层的参考实现。代码以 Rust 为主，依赖 Redis 提供房间索引与快照存储，并预留 PostgreSQL 作为后续账户/支付/战绩落库的存储层。仓库同时包含一个 cocos 客户端占位符目录，方便后续对接。

## 组件列表

- `crates/common`：跨服务共享的数据结构、错误定义与负载均衡算法。
- `crates/gateway`：网关层模拟实现，维护 Lobby 注册表、Session 表与健康检查。
- `crates/lobby`：大厅服务原型，负责注册/心跳与 Redis 房间快照装载。
- `crates/room`：房间服务原型，演示创建房间与占座写入 Redis 快照。
- `crates/table`：麻将对局模拟器，占位一局对局的启动与结束。
- `client/cocos`：cocos 客户端占位目录，提供事件/协议草稿。

更多部署细节请参考 `docs/DEPLOY.md`。
