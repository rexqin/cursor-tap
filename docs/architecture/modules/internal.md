---
type: Module
title: internal 包
description: Go 核心业务库，包含 CA、MITM、HTTP 流解析、SQLite 存储与 API 服务。
tags: [go, internal, mitm, httpstream]
timestamp: 2026-06-24T00:00:00Z
resource: file://internal/
---

不可独立运行的 Go 核心业务库，由 [tap 代理](/modules/tap.md) 与 [API 服务](/modules/api-server.md) 调用。

# Schema

```
internal/
├── apiserver/    # 独立 API HTTP 服务
├── recordstore/  # SQLite record 持久化
├── notify/       # Proxy → API 更新通知客户端
├── ca/           # 自签 CA、动态域名证书
├── config/       # ProxyConfig / APIConfig
├── proxy/        # HTTP / SOCKS5 代理编排
├── mitm/         # TLS MITM、HTTP/2 桥接
├── httpstream/   # 流解析、gRPC/SSE 解码、SQLite Recorder
└── api/          # REST + WebSocket 路由与 Hub
```

## 包职责

| 包 | 职责 |
|----|------|
| `apiserver/` | 独立 API 进程 HTTP 服务 |
| `recordstore/` | SQLite 读写 |
| `notify/` | 异步 POST notify |
| `ca/` | 自签 CA 与动态域名证书 |
| `config/` | ProxyConfig / APIConfig |
| `proxy/` | HTTP/SOCKS5 代理 |
| `mitm/` | TLS 拦截、SNI 检测、HTTP/2 桥接 |
| `httpstream/` | 流解析、gRPC/SSE 解码、SQLite 录制 |
| `api/` | REST 路由与 WebSocket Hub |

详见 [数据流](/data-flow.md)、[管理 API](/api/management-api.md)。
