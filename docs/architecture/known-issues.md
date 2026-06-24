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
| API 需独立启动 | `tap start` 不再监听 9090；未启动 `tap-api` 时 Web UI 与 CLI `stats` 不可用 |
| `tools/restore` 未适配 SQLite | 仍面向旧 JSONL 格式；record 已迁移至 `{dataDir}/records.db` |
| Proxy notify 可丢失 | API 未运行时 notify 被丢弃；API 启动后 Web 可通过 `GET /api/records` 从 DB 补偿 |

## 相关

- [管理 API](/api/management-api.md)
- [tap 代理](/modules/tap.md)
- [API 服务](/modules/api-server.md)
- [开发工具](/tools.md)
