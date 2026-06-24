---
type: Reference
title: 测试
description: Go 单元测试栈、约定与覆盖范围。
tags: [testing, go, testify]
timestamp: 2026-06-24T00:00:00Z
resource: file://tests/
---

# Schema

## 技术栈

| 工具 | 用途 |
|------|------|
| `testing` | 测试入口、`t.Run` 表驱动 |
| `testify/require` | 断言（失败即停，替代 `t.Fatal`） |
| `go-cmp/cmp` | struct / slice 深度比较，输出 `-want +got` diff |

## 约定

- 测试包使用 `package xxx_test`（外部视角）
- 表驱动测试 + `t.Run` 子用例
- 复杂结构用 `cmp.Diff`；简单值用 `require.Equal`

## 覆盖

| 路径 | 覆盖 |
|------|------|
| `tests/config/config_test.go` | ProxyConfig / APIConfig 默认值 |
| `tests/recordstore/store_test.go` | SQLite 读写 |
| `tests/notify/client_test.go` | API 通知客户端 |
| `tests/apiserver/server_test.go` | 独立 API 路由与 notify |
| `tests/proxy/proxy_test.go` | Proxy 初始化 |
| `tests/api/handlers_test.go` | REST/WS handler |
| `tests/httpstream/*` | 流解析 |

# Examples

```bash
pnpm exec nx test tap
# 或
go test ./tests/... 
```

## 相关

- [开发工作流](/development/workflow.md)

# Citations

[1] [testify](https://github.com/stretchr/testify)
[2] [go-cmp](https://github.com/google/go-cmp)
