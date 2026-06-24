---
type: Reference
title: 测试
description: Go 与 Web 单元测试栈、约定与覆盖范围。
tags: [testing, go, vitest, playwright]
timestamp: 2026-06-24T00:00:00Z
resource: file://tests/
---

# Schema

## Go 技术栈

| 工具 | 用途 |
|------|------|
| `testing` | 测试入口、`t.Run` 表驱动 |
| `testify/require` | 断言（失败即停，替代 `t.Fatal`） |
| `go-cmp/cmp` | struct / slice 深度比较，输出 `-want +got` diff |

### Go 约定

- 测试包使用 `package xxx_test`（外部视角）
- 表驱动测试 + `t.Run` 子用例
- 复杂结构用 `cmp.Diff`；简单值用 `require.Equal`

### Go 覆盖

| 路径 | 覆盖 |
|------|------|
| `tests/config/config_test.go` | ProxyConfig / APIConfig 默认值 |
| `tests/recordstore/store_test.go` | SQLite 读写 |
| `tests/notify/client_test.go` | API 通知客户端 |
| `tests/apiserver/server_test.go` | 独立 API 路由与 notify |
| `tests/proxy/proxy_test.go` | Proxy 初始化 |
| `tests/api/handlers_test.go` | REST/WS handler |
| `tests/httpstream/*` | 流解析 |

## Web 技术栈

| 工具 | 用途 |
|------|------|
| **Vitest** | 单元 / 组件 / Hook 测试 |
| **React Testing Library** | 组件与 `renderHook` |
| **MSW** | Mock REST API（`/api/records`） |
| **Playwright** | E2E（页面加载、API mock 交互） |

### Web 约定

- 纯逻辑放在 `apps/web/src/lib/record-utils.ts`，优先单元测试
- 测试文件与源码 colocate：`*.test.ts` / `*.test.tsx`
- 共享 fixture：`apps/web/src/test/fixtures/`
- MSW setup：`apps/web/src/test/setup.ts`
- 通过 public API / DOM 断言，不测内部 state

### Web 覆盖

| 路径 | 覆盖 |
|------|------|
| `src/lib/record-utils.test.ts` | gRPC 解析、session 聚合、筛选、去重 |
| `src/lib/ws-client.test.ts` | WebSocket 连接、消息、重连 |
| `src/hooks/use-records.test.ts` | Hook + MSW 集成 |
| `e2e/home.spec.ts` | 页面加载、mock API 展示 |

# Examples

```bash
# Go
pnpm exec nx test tap
go test ./tests/...

# Web 单元测试
pnpm exec nx test web

# Web E2E（需 Chromium，首次运行 playwright install chromium）
pnpm exec nx e2e web
```

## 相关

- [开发工作流](/development/workflow.md)

# Citations

[1] [testify](https://github.com/stretchr/testify)
[2] [go-cmp](https://github.com/google/go-cmp)
[3] [Vitest](https://vitest.dev/)
[4] [Playwright](https://playwright.dev/)
