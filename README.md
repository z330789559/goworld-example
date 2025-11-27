# GoWorld 商业项目骨架

面向 GoWorld 风格的游戏后端，提供一套可直接运行的模块化脚手架，方便与 [gforgame](https://github.com/free-city/gforgame/tree/main) 对比和快速迭代。

## 主要特性
- **模块齐全**：account、player、bag、item、shop、mail、notice、chat、room、match、dao、redis、log 等核心能力。
- **清晰依赖注入**：集中构造服务，方便替换 DAO/缓存实现或接入真实中间件。
- **即开即用**：内置内存存储与示例数据，`go run ./cmd/server` 即可启动。
- **可扩展性**：模块之间通过明确定义的 service 层交互，便于迁移到数据库、消息队列等生产设施。

## 目录结构
```
.
├── cmd/server          # 程序入口
├── internal
│   ├── config          # 配置默认值
│   ├── dao             # 内存数据层（可替换为数据库）
│   ├── log             # 日志封装
│   ├── redis           # 内存缓存（模拟 Redis）
│   ├── server          # 路由聚合
│   └── modules         # 业务模块
│       ├── account
│       ├── bag
│       ├── chat
│       ├── item
│       ├── mail
│       ├── match
│       ├── notice
│       ├── player
│       ├── room
│       └── shop
└── go.mod
```

## 快速开始
```bash
go run ./cmd/server
```

可选接口示例：
- `POST /api/account/register` 注册账号
- `POST /api/account/login` 登录并获取 token
- `GET  /api/player/:id` 查询角色
- `GET  /api/bag/:playerID` 查询背包
- `GET  /api/items/` 道具表
- `GET  /api/shop/items` 商城列表
- `GET  /api/mail/:playerID` 邮件与附件
- `GET  /api/notice/` 公告
- `POST /api/chat/` 发送聊天
- `GET  /api/chat/` 获取聊天记录
- `POST /api/room/create` 创建房间（麻将/斗地主等）
- `GET  /api/room/` 房间列表
- `POST /api/match/enqueue` 匹配示例

## 对比 gforgame 时的参考点
- **模块覆盖度**：本模板聚焦核心玩法前置功能，使用内存实现方便快速试用。
- **可替换性**：DAO、Redis、日志的接口化设计方便在对比时替换为 gforgame 的实现，验证性能或工程风格。
- **学习路径**：通过逐个模块替换，逐步迁移到生产级依赖，便于团队评估两套框架的上手成本。
