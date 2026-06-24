---
type: Reference
title: Monorepo 结构
description: cursor-tap 仓库目录布局与 Nx 项目配置。
tags: [monorepo, nx, structure]
timestamp: 2026-06-24T00:00:00Z
---

# Schema

## 目录树

```
cursor-tap/
├── apps/
│   ├── tap/              # Go MITM 代理 CLI
│   ├── api/              # 独立管理 API 服务
│   └── web/              # Next.js Web Inspector
├── internal/             # Go 核心业务库
├── packages/
│   └── proto/            # Protobuf 定义 + buf 生成代码
├── tests/                # Go 单元测试
├── tools/                # 开发/逆向辅助工具
├── docs/                 # 文档与逆向笔记
├── assets/               # 从 Cursor 提取的配置快照
├── nx.json               # Nx 工作区配置
├── go.mod                # Go 模块（github.com/burpheart/cursor-tap）
├── package.json          # 根 workspace scripts
└── pnpm-workspace.yaml   # pnpm monorepo（apps/*）
```

## Nx 项目

| 项目 | 路径 | 类型 | 说明 |
|------|------|------|------|
| `tap` | `apps/tap` | application (Go) | MITM 代理 CLI，`build` 依赖 `proto:build` |
| `api` | `apps/api` | application (Go) | 管理 API 服务 |
| `web` | `apps/web` | application (TS) | Next.js 前端，由 `@nx/next/plugin` 推断 targets |
| `proto` | `packages/proto` | library | buf lint / generate |

## 相关模块

- [tap 代理](/modules/tap.md)
- [API 服务](/modules/api-server.md)
- [internal 包](/modules/internal.md)
- [Proto 包](/modules/proto.md)
- [Web Inspector](/modules/web.md)
