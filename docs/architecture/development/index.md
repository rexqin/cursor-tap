# 开发

* [开发工作流](/development/workflow.md) - 三进程启动（tap-api + tap + web）、构建与测试
* [Cursor 配置](/development/cursor-config.md) - 代理与 CA 证书环境变量

## 启动顺序

1. `pnpm dev:api` 或 `tap-api start`
2. `pnpm dev:tap` 或 `tap start --http-parse`
3. `pnpm dev:web`

或一条命令：`pnpm dev`
