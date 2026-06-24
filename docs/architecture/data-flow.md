---
type: Reference
title: 数据流
description: 从代理接入、TLS MITM 到 SQLite 持久化与 API WebSocket 推送的完整链路。
tags: [data-flow, mitm, grpc, websocket, sqlite]
timestamp: 2026-06-24T00:00:00Z
---

# Schema

| 步骤 | 组件 | 说明 |
|------|------|------|
| 1 | `proxy.Server.Start()` | 初始化 CA、KeyLog；若启用 `--http-parse`，创建 SQLite `Recorder` |
| 2 | HTTP/SOCKS5 接入 | HTTP `CONNECT` 或 SOCKS5 `CONNECT` → `mitm.Interceptor.InterceptAuto()` |
| 3 | TLS 检测 | `detect.go` Peek 首字节判断 TLS，提取 SNI；`ca.CA` 签发目标域名证书，完成双 TLS 握手 |
| 4 | 协议桥接 | ALPN 协商 `h2` 或 `http/1.1`；HTTP/2 走 `h2bridge.go` |
| 5 | 流解析 | `parser.go` 镜像双向 body，异步解析 request/response/SSE/gRPC；`grpc.go` 解析 Connect framing，经 `MessageRegistry` 反序列化为 JSON |
| 6 | 持久化 | 每条 `Record` 写入 SQLite（`{dataDir}/records.db`） |
| 7 | 通知 API | Proxy `POST /internal/notify` 发送 `latest_id`；API 从 SQLite 读取增量并 `Hub.Broadcast` |
| 8 | Web 消费 | Web UI 通过 `GET /api/records` 初始加载，`/ws/records` 实时接收 |

## 涉及包

- [tap 代理](/modules/tap.md) — 写入与 notify
- [API 服务](/modules/api-server.md) — 读取与广播
- [internal 包](/modules/internal.md) — `recordstore/`、`notify/`
- [管理 API](/api/management-api.md) — REST/WS 端点
- [Record 类型契约](/record-types.md) — JSON 结构
