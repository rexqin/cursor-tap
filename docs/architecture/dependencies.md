---
type: Reference
title: 技术依赖
description: cursor-tap Go、前端与工具链主要依赖。
tags: [dependencies, go, nextjs, nx]
timestamp: 2026-06-24T00:00:00Z
---

# Schema

## Go（go 1.25）

| 依赖 | 用途 |
|------|------|
| `spf13/cobra` | tap / tap-api CLI |
| `modernc.org/sqlite` | record 持久化（纯 Go SQLite，WAL 模式） |
| `gorilla/websocket` | WebSocket Hub（tap-api） |
| `google.golang.org/protobuf` + `grpc` | protobuf 反序列化 |
| `bufbuild/protocompile` | tools/inline proto 编译 |
| `andybalholm/brotli` | Brotli 解压 |
| `golang.org/x/net/http2` | HTTP/2 MITM 桥接 |

## 前端

| 依赖 | 版本 | 用途 |
|------|------|------|
| Next.js | 16.1.6 | App Router + Turbopack |
| React | 19.2.3 | UI |
| Tailwind CSS | v4 | 样式 |
| shadcn/ui + Radix | — | 组件库 |

## 工具链

| 工具 | 用途 |
|------|------|
| pnpm 10.28 | 包管理 |
| Nx 23 | Monorepo 编排 |
| buf | proto lint / generate |

## 相关

- [Monorepo 结构](/monorepo.md)
- [Web Inspector](/modules/web.md)
