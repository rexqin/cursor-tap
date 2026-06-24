---
type: Module
title: Proto 包
description: 从 Cursor 客户端逆向提取的 Protobuf 定义，由 buf 生成 Go 代码。
tags: [protobuf, buf, reverse-engineering]
timestamp: 2026-06-24T00:00:00Z
resource: file://packages/proto/
---

从 Cursor 客户端 JS（`protobuf-es` 编译产物）逆向提取的 `.proto` 定义。

# Schema

| 文件 | package | 说明 |
|------|---------|------|
| `aiserver_v1.proto` | `aiserver.v1` | 主 AI 服务 RPC |
| `agent_v1.proto` | `agent.v1` | Agent 消息/工具协议 |
| `anyrun_v1.proto` | `anyrun.v1` | Anyrun 相关 |
| `internapi_v1.proto` | `internapi.v1` | 内部 API |

生成输出：`packages/proto/gen/{aiserver,agent,anyrun,internapi}/v1/*.pb.go`

构建链：`tap:build` → dependsOn `proto:build` → `buf generate`

## 相关

- [开发工具](/tools.md) — `tools/ext` proto 提取
- [internal/httpstream](/modules/internal.md) — gRPC 帧反序列化

# Citations

[1] [Cursor 逆向笔记](../cursor-reverse-notes-1.md)
