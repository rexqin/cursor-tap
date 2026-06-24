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

# 并行启动代理 + WebUI
pnpm dev

# 分别启动
pnpm dev:tap    # Go 代理（:8080 / :1080 / :9090）
pnpm dev:web    # Next.js（:3000）

# 构建
pnpm build:tap  # → dist/tap
pnpm build:web

# Proto 生成
pnpm exec nx run proto:generate

# 测试
pnpm exec nx test tap
```

## 相关

- [Cursor 配置](/development/cursor-config.md) — 配置 IDE 走代理
- [tap 代理](/modules/tap.md) — 端口与 flags
- [测试](/testing.md)
