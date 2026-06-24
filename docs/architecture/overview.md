---
type: Architecture
title: 系统概览
description: Cursor-Tap MITM 代理 + 独立 API 服务 + Web Inspector 的整体架构。
tags: [architecture, overview, mitm, grpc]
timestamp: 2026-06-24T00:00:00Z
---

Cursor-Tap 是一个 **MITM 代理 + 独立 API 服务 + Web Inspector** 工具链，用于拦截、解密并可视化 Cursor IDE 与后端之间的 Connect/gRPC 通信。

```
┌─────────────┐     HTTP_PROXY      ┌──────────────────────────┐
│  Cursor IDE │ ──────────────────► │  tap (Go) :8080 / :1080  │
│             │                     │  MITM + gRPC 解析        │
└─────────────┘                     │  └─ SQLite 写入 + notify │
                                    └────────────┬─────────────┘
                                                 │ POST /internal/notify
                                                 │ 共享 records.db
                                                 ▼
                                    ┌──────────────────────────┐
                                    │  tap-api (Go) :9090      │
                                    │  REST + WebSocket        │
                                    └────────────┬─────────────┘
                                                 │ WS + REST
                                                 ▼
                                    ┌──────────────────────────┐
                                    │  web (Next.js) :3000     │
                                    │  四栏 Inspector UI       │
                                    └──────────────────────────┘
```

Proxy 与 API 为**独立进程**：Web UI 连接 API（9090），Cursor 仅感知代理（8080/1080）。

## 组件关系

| 组件 | 路径 | 说明 |
|------|------|------|
| tap | [apps/tap](../../apps/tap/) | MITM 代理 CLI |
| tap-api | [apps/api](../../apps/api/) | 管理 API 独立进程 |
| internal | [internal/](../../internal/) | 核心业务库 |
| web | [apps/web](../../apps/web/) | Next.js Inspector UI |
| proto | [packages/proto](../../packages/proto/) | Protobuf 定义 |

详见 [数据流](/data-flow.md)、[Monorepo 结构](/monorepo.md)。

# Citations

[1] [项目 README](../../README.md)
[2] [Cursor 逆向笔记](../cursor-reverse-notes-1.md)
