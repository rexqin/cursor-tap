---
type: Reference
title: 数据流
description: 从代理接入、TLS MITM 到 JSONL 录制与 WebSocket 推送的完整链路。
tags: [data-flow, mitm, grpc, websocket]
timestamp: 2026-06-24T00:00:00Z
---

# Schema

| 步骤 | 组件 | 说明 |
|------|------|------|
| 1 | `proxy.Server.Start()` | 初始化 CA、KeyLog、WebSocket Hub；若启用 `--http-parse`，创建 `httpstream.Recorder`，`OnRecord` 回调广播到 Hub |
| 2 | HTTP/SOCKS5 接入 | HTTP `CONNECT` 或 SOCKS5 `CONNECT` → `mitm.Interceptor.InterceptAuto()` |
| 3 | TLS 检测 | `detect.go` Peek 首字节判断 TLS，提取 SNI；`ca.CA` 签发目标域名证书，完成双 TLS 握手 |
| 4 | 协议桥接 | ALPN 协商 `h2` 或 `http/1.1`；HTTP/2 走 `h2bridge.go` |
| 5 | 流解析 | `parser.go` 镜像双向 body，异步解析 request/response/SSE/gRPC；`grpc.go` 解析 Connect framing，经 `MessageRegistry` 反序列化为 JSON |
| 6 | 录制与推送 | 每条 `Record` 写入 JSONL（`access.log`），存入内存 ring buffer（默认 10000 条），实时 JSON 广播到 WebSocket 客户端 |

## 涉及包

- [internal 包](/modules/internal.md) — `proxy/`、`mitm/`、`httpstream/`
- [管理 API](/api/management-api.md) — `/ws/records` 推送
- [Record 类型契约](/record-types.md) — JSON 结构
