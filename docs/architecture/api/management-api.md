---
type: API Endpoint
title: 管理 API
description: tap-api 独立进程 9090 端口的 REST 与 WebSocket 管理接口。
tags: [api, rest, websocket]
timestamp: 2026-06-24T00:00:00Z
resource: http://localhost:9090
---

路由定义：`internal/apiserver/server.go` + `internal/api/handlers.go`

# Schema

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/status` | `{"status":"running"}` |
| GET | `/api/stats` | 会话统计 + `ws_clients` + `record_count` |
| GET | `/api/ca/cert` | 下载 CA 证书 |
| GET | `/api/records?limit=N` | 从 SQLite 读取最近 N 条（1–1000，默认 100） |
| WS | `/ws/records` | 实时 JSON record 流 |
| POST | `/internal/notify` | Proxy 推送 `latest_id`（localhost only） |
| GET | `/api/sessions` | **未实现** — CLI `sessions` 命令会 404 |

## 其他行为

- CORS：`Access-Control-Allow-Origin: *`（开发模式）
- 由 **tap-api** 进程写入 `~/.cursor-tap/api.addr`

## 相关

- [API 服务](/modules/api-server.md) — 独立进程
- [tap 代理](/modules/tap.md) — SQLite 写入与 notify
- [Record 类型契约](/record-types.md)
- [已知问题](/known-issues.md)
