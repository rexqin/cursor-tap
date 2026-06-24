# Cursor-Tap 项目架构

> 本文档已按 [OKF v0.1](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md) 重组为知识 Bundle。

请从 **[architecture/index.md](./architecture/index.md)** 进入完整架构文档。

## 快速链接

* [系统概览](./architecture/overview.md) — 三进程架构（tap / tap-api / web）
* [Monorepo 结构](./architecture/monorepo.md)
* [数据流](./architecture/data-flow.md) — SQLite + notify
* [tap 代理](./architecture/modules/tap.md) — MITM :8080 / :1080
* [API 服务](./architecture/modules/api-server.md) — 独立 :9090
* [开发工作流](./architecture/development/workflow.md)

## 启动顺序

1. `tap-api start` — 管理 API（REST + WebSocket）
2. `tap start --http-parse` — MITM 代理，写入 SQLite 并通知 API
3. Web UI — 连接 `localhost:9090`

或使用 `pnpm dev` 并行启动三者。
