---
type: Module
title: tap 代理
description: Go CLI 入口，提供 HTTP/SOCKS5 MITM 代理与 SQLite 录制。
tags: [go, cli, proxy, cobra]
timestamp: 2026-06-24T00:00:00Z
resource: file://apps/tap/main.go
---

基于 Cobra 的 MITM 代理 CLI，入口文件 `apps/tap/main.go`。不再内嵌管理 API。

# Schema

## 子命令

| 命令 | 功能 |
|------|------|
| `start` | 启动 HTTP/SOCKS5 代理 |
| `ca info/export/regenerate/clean-certs` | 自签 CA 证书管理 |
| `sessions` | 查询活跃会话（调用 API `GET /api/sessions`，**未实现**） |
| `stats` | 查询统计（调用 API `GET /api/stats`） |

## 默认端口（`start`）

| 端口 | 协议 | 用途 |
|------|------|------|
| 8080 | HTTP CONNECT | HTTP/HTTPS 代理（Cursor 配置 `HTTP_PROXY`） |
| 1080 | SOCKS5 | SOCKS5 代理（可选） |

## 关键 flags

| Flag | 说明 |
|------|------|
| `--http-parse` | 启用 HTTP 流解析与 SQLite 录制 |
| `--http-log N` | 控制台日志级别（0–4） |
| `--record-db PATH` | SQLite 数据库（默认 `{data-dir}/records.db`） |
| `--api-notify-url URL` | API 更新通知地址（默认 `http://127.0.0.1:9090/internal/notify`） |
| `--upstream URL` | 上游代理（如 `socks5://127.0.0.1:7890`） |

Nx 开发模式（`pnpm dev:tap`）自动附加：`--http-parse --http-log 4`

## 相关概念

- [API 服务](/modules/api-server.md) — 独立 9090 进程
- [internal 包](/modules/internal.md) — 核心业务实现
- [数据流](/data-flow.md) — 运行时链路
- [Cursor 配置](/development/cursor-config.md) — 代理环境变量
