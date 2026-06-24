---
type: Reference
title: 已知问题
description: cursor-tap 当前限制与待修复项。
tags: [known-issues, limitations]
timestamp: 2026-06-24T00:00:00Z
---

# Schema

| 问题 | 说明 |
|------|------|
| `/api/sessions` 未实现 | CLI `sessions` 命令会 404；Web UI 在客户端从 records 聚合 session |
| `DefaultConfig().APIPort` 为 8888 | CLI 默认 `--api-port` 为 9090，以 CLI 为准 |
| pnpm monorepo + Turbopack | `apps/web/next.config.ts` 需设置 `turbopack.root` 指向 workspace 根目录 |
| `nx.json` 引用 `go.work` | 仓库中暂无 `go.work` 文件，不影响当前构建 |

## 相关

- [管理 API](/api/management-api.md)
- [tap 代理](/modules/tap.md)
