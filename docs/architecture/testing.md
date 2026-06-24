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
| `tests/config/config_test.go` | 配置默认值 |
| `tests/httpstream/decoder_test.go` | Content-Encoding 解码 |
| `tests/httpstream/grpc_test.go` | gRPC 帧解析 |
| `tests/httpstream/sse_test.go` | SSE 解析 |
| `tests/httpstream/types_test.go` | 类型 |

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
