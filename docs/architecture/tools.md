---
type: Reference
title: 开发工具
description: tools/ 目录下的 proto 提取、日志还原与调试辅助工具。
tags: [tools, reverse-engineering, proto]
timestamp: 2026-06-24T00:00:00Z
resource: file://tools/
---

# Schema

| 工具 | 路径 | 用途 |
|------|------|------|
| `ext` | `tools/ext/` | 从 Cursor 打包 JS 提取 protobuf-es 结构，生成 `.proto` |
| `restore` | `tools/restore/` | 从 record 日志还原/重组 Agent 消息（待迁移 SQLite） |
| `inline` | `tools/inline/` | protocompile 工具：查看 proto 消息结构 |
| `debug_bidi` | `tools/debug_bidi/` | 调试 BidiAppend 日志条目 |

## 相关

- [Proto 包](/modules/proto.md) — `ext` 生成的 proto 定义
- [Cursor 逆向笔记](../cursor-reverse-notes-1.md)
