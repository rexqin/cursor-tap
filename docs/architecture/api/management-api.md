---
type: API Endpoint
title: 管理 API
description: tap 代理 9090 端口的 REST 与 WebSocket 管理接口。
tags: [api, rest, websocket]
timestamp: 2026-06-24T00:00:00Z
resource: http://localhost:9090
---

路由定义：`internal/proxy/server.go`（status/stats/ca）+ `internal/api/handlers.go`（records/ws）

# Schema

| 方法 | 路径 | 条件 | 说明 |
|------|------|------|------|
| GET | `/api/status` | 始终 | `{"status":"running"}` |
| GET | `/api/stats` | 始终 | 会话统计 + `ws_clients` 数 |
| GET | `/api/ca/cert` | 始终 | 下载 CA 证书 |
| GET | `/api/records?limit=N` | 需 `--http-record` | 最近 N 条 record（1–1000，默认 100） |
| WS | `/ws/records` | 需 `--http-record` | 实时 JSON record 流 |
| GET | `/api/sessions` | **未实现** | CLI `sessions` 命令会 404 |

## 其他行为

- CORS：`Access-Control-Allow-Origin: *`（开发模式）
- 启动后写入 `~/.cursor-tap/api.addr`，供 CLI 子命令读取 API 地址

## 相关

- [tap 代理](/modules/tap.md) — 9090 端口服务
- [Record 类型契约](/record-types.md) — `/api/records` 与 `/ws/records`  payload
- [已知问题](/known-issues.md) — `/api/sessions` 未实现
