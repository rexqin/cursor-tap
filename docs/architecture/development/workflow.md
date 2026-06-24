---
type: Playbook
title: 开发工作流
description: cursor-tap 本地开发、构建与测试命令。
tags: [development, workflow, nx, pnpm]
timestamp: 2026-06-24T00:00:00Z
---

# Examples

```bash
pnpm install

# 并行启动代理 + API + WebUI
pnpm dev

# 分别启动
pnpm dev:api    # API 服务（:9090）
pnpm dev:tap    # Go 代理（:8080 / :1080）
pnpm dev:web    # Next.js（:3000）

# 构建
pnpm build:tap  # → dist/tap
pnpm build:api  # → dist/tap-api
pnpm build:web

# Proto 生成
pnpm exec nx run proto:generate

# 测试
pnpm exec nx test tap
pnpm exec nx test api
pnpm exec nx test web
pnpm exec nx e2e web   # 首次需 playwright install chromium
```

## 相关

- [Cursor 配置](/development/cursor-config.md) — 配置 IDE 走代理
- [tap 代理](/modules/tap.md) — 代理 flags
- [API 服务](/modules/api-server.md) — API flags
- [测试](/testing.md)
