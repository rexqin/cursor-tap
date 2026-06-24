---
type: Module
title: internal 包
description: Go 核心业务库，包含 CA、MITM、HTTP 流解析与管理 API。
tags: [go, internal, mitm, httpstream]
timestamp: 2026-06-24T00:00:00Z
resource: file://internal/
---

不可独立运行的 Go 核心业务库，由 [tap 代理](/modules/tap.md) 调用。

# Schema

```
internal/
├── ca/           # 自签 CA、按域名动态签发证书、磁盘缓存
├── config/       # Config 结构体（端口、目录、HTTP 解析选项）
├── proxy/        # 三服务编排：HTTP / SOCKS5 / API
├── mitm/
│   ├── interceptor.go   # TLS MITM 核心
│   ├── detect.go        # TLS/SNI 检测（Peek ClientHello）
│   ├── dialer.go        # 上游代理连接
│   ├── h2bridge.go      # HTTP/2 桥接 + 流式解析
│   └── keylog.go        # TLS KeyLog → sslkeys.log
├── httpstream/
│   ├── parser.go        # 双向 HTTP 流解析（零阻塞透传 + 异步镜像）
│   ├── decoder.go       # Content-Encoding 解压（gzip/deflate/brotli）
│   ├── sse.go           # SSE 事件解析
│   ├── grpc.go          # Connect/gRPC 帧解析 + protobuf→JSON
│   ├── grpc_registry.go # protoregistry 消息类型发现
│   ├── recorder.go      # JSONL 录制 + 内存 ring buffer + WS 回调
│   ├── logger.go        # 控制台彩色日志
│   └── types.go         # HTTPMessage, Record, LogLevel 等
└── api/
    ├── handlers.go      # REST + WebSocket 路由
    └── hub.go           # WebSocket Hub 广播
```

## 包职责

| 包 | 职责 |
|----|------|
| `ca/` | 自签 CA 与动态域名证书 |
| `config/` | 配置结构与默认值 |
| `proxy/` | HTTP/SOCKS5/API 三服务编排 |
| `mitm/` | TLS 拦截、SNI 检测、HTTP/2 桥接 |
| `httpstream/` | 流解析、gRPC/SSE 解码、录制与推送 |
| `api/` | REST 路由与 WebSocket Hub |

详见 [数据流](/data-flow.md)、[管理 API](/api/management-api.md)。
