# 架构文档更新日志

## 2026-06-24

* **Update**: 文档同步 Proxy/API 拆分：更新 cursor-config、record-types、known-issues、web 等概念文档。
* **Update**: Proxy 与 API 拆分为独立进程；record 持久化改为 SQLite；新增 `apps/api`（tap-api）。
* **Update**: `DefaultConfig().APIPort` 统一为 9090；移除 `nx.json` 中过时的 `go.work` 引用；已知问题列表同步更新。
* **Creation**: 将 `docs/architecture.md` 单体文档迁移为 OKF v0.1 知识 Bundle。
* **Update**: 拆分为概览、模块、API、开发与依赖等独立概念文档。
