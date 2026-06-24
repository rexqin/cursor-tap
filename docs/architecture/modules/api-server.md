---
type: Module
title: API 服务
description: 独立 tap-api 进程，提供 REST/WebSocket 并从 SQLite 读取 record。
tags: [go, api, websocket, sqlite]
timestamp: 2026-06-24T00:00:00Z
resource: file://apps/api/main.go
---

独立管理 API 进程，入口 `apps/api/main.go`，实现位于 `internal/apiserver/`。

# Schema

## 子命令

| 命令 | 功能 |
|------|------|
| `start` | 启动 API HTTP 服务 |

## 默认端口

| 端口 | 协议 | 用途 |
|------|------|------|
| 9090 | HTTP | REST + WebSocket |

## 关键 flags

| Flag | 说明 |
|------|------|
| `--port` | API 监听端口（默认 9090） |
| `--cert-dir` | CA 证书目录（用于 `/api/ca/cert`） |
| `--record-db` | SQLite 数据库路径（默认 `{cert-dir}/data/records.db`） |

## 内部端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/internal/notify` | 仅 localhost；Proxy 推送 `{latest_id}` 触发增量广播 |

启动后写入 `~/.cursor-tap/api.addr`。

## 相关

- [管理 API](/api/management-api.md) — 公开 REST/WS 端点
- [tap 代理](/modules/tap.md) — SQLite 写入与 notify
- [数据流](/data-flow.md)
