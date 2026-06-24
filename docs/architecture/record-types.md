---
type: Reference
title: Record 类型契约
description: 前后端通过 httpstream.Record JSON 对齐的消息类型定义。
tags: [record, json, grpc, contract, sqlite]
timestamp: 2026-06-24T00:00:00Z
---

前后端通过 `httpstream.Record` JSON 对齐，前端定义见 `apps/web/src/lib/types.ts`。

## 持久化与传输

| 阶段 | 位置 | 说明 |
|------|------|------|
| 写入 | Proxy → SQLite | `{dataDir}/records.db`，由 `recordstore` 存储 JSON payload |
| 读取 | API ← SQLite | `GET /api/records` 从 DB 读取最近 N 条 |
| 实时 | API → Web | Proxy `POST /internal/notify` 触发 API 增量广播至 `/ws/records` |

# Schema

| `type` | 含义 |
|--------|------|
| `request` | HTTP 请求 |
| `response` | HTTP 响应 |
| `grpc` | Connect/gRPC 帧（含 `grpc_service`, `grpc_method`, `grpc_data`） |
| `sse` | Server-Sent Events |
| `body` | 原始 body 片段 |
| `error` | 解析错误 |
| `debug` | 调试信息 |

## 相关

- [internal/httpstream](/modules/internal.md) — Proxy 侧 Record 生成
- [API 服务](/modules/api-server.md) — SQLite 读取与 WS 广播
- [Web Inspector](/modules/web.md) — 前端消费与展示
- [管理 API](/api/management-api.md) — REST/WS 端点
- [数据流](/data-flow.md)
