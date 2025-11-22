# 部署与运行指南

本指南面向本仓库的 Gateway / Lobby / Room / Mahjong Table 栈，默认使用本地 Docker 提供的 PostgreSQL 与 Redis。Rust 组件基于 Tokio 异步运行时。

## 依赖
- Rust 1.76+（建议使用 `rustup`）
- Docker & Docker Compose（用于 Postgres + Redis）
- 可选：Cocos Creator（用于客户端接入）

## 快速启动

1. 启动基础依赖

```bash
./scripts/devstack.sh up
```

该脚本会启动：
- Redis: `6379`
- PostgreSQL: `5432`（用户名 `mahjong`，密码 `mahjong`，数据库 `mahjong`）

2. 编译全部服务

```bash
cargo build
```

3. 分别运行示例流程

- Gateway 心跳与注册演示：

```bash
cargo run -p gateway
```

- Lobby 注册 + 心跳循环 + Redis 房间装载：

```bash
cargo run -p lobby
```

- Room 创建/占座并写入 Redis 快照：

```bash
cargo run -p room
```

- Mahjong Table 对局模拟：

```bash
cargo run -p table
```

4. 关闭依赖

```bash
./scripts/devstack.sh down
```

## 目录结构

- `crates/common`：跨服务协议、负载均衡算法。
- `crates/gateway`：注册表、心跳、健康检查与 sticky session 演示。
- `crates/lobby`：注册/心跳上报与 Redis 房间列表加载。
- `crates/room`：房间创建、占座、快照写入 Redis。
- `crates/table`：麻将桌逻辑占位符与一局模拟。
- `client/cocos`：Cocos 客户端事件草稿与 mock handler。

## 生产落地提示

- 将 `MockGatewayClient` 替换为真实 RPC 客户端（gRPC/HTTP 均可）。
- Gateway 需实现 WebSocket 连接管理与实际推送。
- Room Service 建议将房间状态快照落盘至 PostgreSQL，Redis 作为高频索引层。
- 增加 OpenTelemetry / Prometheus 监控，绑定设计中的指标项。
- 在 Kubernetes 中部署时，可将 Gateway/Lobby/Room/Table 分别作为 Deployment，利用 Service 暴露 RPC/WS 端口。
