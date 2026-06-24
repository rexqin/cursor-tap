---
type: Reference
title: Record 类型契约
description: 前后端通过 httpstream.Record JSON 对齐的消息类型定义。
tags: [record, json, grpc, contract]
timestamp: 2026-06-24T00:00:00Z
---

前后端通过 `httpstream.Record` JSON 对齐，前端定义见 `apps/web/src/lib/types.ts`。

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

- [internal/httpstream](/modules/internal.md) — 后端 Record 生成
- [Web Inspector](/modules/web.md) — 前端消费与展示
- [管理 API](/api/management-api.md) — 传输通道
