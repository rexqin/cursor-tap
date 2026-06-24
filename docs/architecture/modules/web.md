---
type: Module
title: Web Inspector
description: Next.js 四栏 Inspector UI，通过 tap-api WebSocket 实时展示 gRPC 流量。
tags: [nextjs, react, frontend, inspector]
timestamp: 2026-06-24T00:00:00Z
resource: file://apps/web/
---

Web UI 连接 **tap-api**（:9090），不直接与 tap 代理通信。

# Schema

## 技术栈

| 依赖 | 版本 | 用途 |
|------|------|------|
| Next.js | 16.1.6 | App Router + Turbopack |
| React | 19.2.3 | UI |
| Tailwind CSS | v4 | 样式 |
| shadcn/ui + Radix | — | 组件库（new-york 风格） |
| react-resizable-panels | — | 四栏拖拽布局 |

## 目录结构

```
apps/web/src/
├── app/
│   ├── layout.tsx          # 根布局
│   ├── page.tsx            # 主页面（四栏 Inspector）
│   └── globals.css         # Tailwind + CSS 变量
├── components/
│   ├── filter-sidebar.tsx  # 左栏：服务/方法树过滤
│   ├── session-list.tsx    # 第二栏：RPC 调用列表
│   ├── record-list.tsx     # 第三栏：帧列表
│   ├── detail-panel.tsx    # 右栏：详情（JSON/headers/body）
│   ├── filter-bar.tsx      # 顶栏：暂停/清空
│   ├── resizable-panels.tsx
│   ├── json-viewer.tsx
│   ├── header-table.tsx
│   └── ui/                 # shadcn 组件
├── hooks/
│   ├── use-records.ts      # records/sessions/filters 状态
│   └── use-websocket.ts    # WebSocket 生命周期
└── lib/
    ├── types.ts            # Record / SessionInfo（与后端 JSON 对齐）
    ├── ws-client.ts        # 自动重连 + 指数退避
    └── utils.ts
```

## UI 布局

`page.tsx` 使用 `ResizablePanels`，默认比例 `[15, 22, 22, 41]`：

1. **FilterSidebar** — gRPC 服务/方法多选过滤
2. **SessionList** — 按 session 聚合的 RPC 调用
3. **RecordList** — 选中 session 内的帧列表
4. **DetailPanel** — 单条 record 详情

## 状态与通信

**`use-records.ts`：**

- 浏览器端最多保留 2000 条 records
- 从 records 客户端聚合 `SessionInfo`
- 初始加载：`GET http://localhost:9090/api/records?limit=100`（tap-api 从 SQLite 读取）
- 重连后 `fetchAndMergeRecords()` 去重合并

**`ws-client.ts`：**

- 默认 `ws://localhost:9090/ws/records`（tap-api 广播）
- 断线自动重连（1s → 30s 指数退避）

## 环境变量

| 变量 | 默认值 | 用途 |
|------|--------|------|
| `NEXT_PUBLIC_WS_URL` | `ws://localhost:9090/ws/records` | WebSocket（tap-api） |
| `NEXT_PUBLIC_API_URL` | `http://localhost:9090` | REST API（tap-api） |

## 相关

- [API 服务](/modules/api-server.md) — 数据源进程
- [Record 类型契约](/record-types.md)
- [管理 API](/api/management-api.md)
- [开发工作流](/development/workflow.md)
