---
type: Architecture
title: 系统概览
description: Cursor-Tap MITM 代理 + gRPC 解析 + Web Inspector 的整体架构。
tags: [architecture, overview, mitm, grpc]
timestamp: 2026-06-24T00:00:00Z
---

Cursor-Tap 是一个 **MITM 代理 + gRPC 解析 + Web Inspector** 工具链，用于拦截、解密并可视化 Cursor IDE 与后端之间的 Connect/gRPC 通信。

```
┌─────────────┐     HTTP_PROXY      ┌──────────────────────────────────┐
│  Cursor IDE │ ──────────────────► │  tap (Go)                        │
│             │     :8080 / :1080   │  ├─ HTTP/SOCKS5 代理             │
└─────────────┘                     │  ├─ TLS MITM + 动态证书          │
                                    │  ├─ HTTP/2 桥接                  │
                                    │  ├─ Connect/gRPC 帧解析          │
                                    │  └─ 管理 API + WebSocket :9090   │
                                    └──────────────┬───────────────────┘
                                                   │ WS + REST
                                                   ▼
                                    ┌──────────────────────────────────┐
                                    │  web (Next.js) :3000             │
                                    │  四栏 Inspector UI               │
                                    └──────────────────────────────────┘
```

前后端通过 **9090 端口** 解耦：Web UI 不嵌入代理进程，Cursor 不感知 Web UI 的存在。

## 组件关系

| 组件 | 路径 | 说明 |
|------|------|------|
| tap | [apps/tap](../../apps/tap/) | Go MITM 代理 CLI |
| internal | [internal/](../../internal/) | 核心业务库 |
| web | [apps/web](../../apps/web/) | Next.js Inspector UI |
| proto | [packages/proto](../../packages/proto/) | Protobuf 定义 |

详见 [数据流](/data-flow.md)、[Monorepo 结构](/monorepo.md)。

# Citations

[1] [项目 README](../../README.md)
[2] [Cursor 逆向笔记](../cursor-reverse-notes-1.md)
