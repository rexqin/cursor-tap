---
type: Playbook
title: Cursor 配置
description: 配置 Cursor IDE 使用 tap MITM 代理与自签 CA 证书。
tags: [cursor, proxy, ca, configuration]
timestamp: 2026-06-24T00:00:00Z
---

## 前置条件

开发/抓包前需先启动：

1. [API 服务](/modules/api-server.md) — `tap-api start`（:9090）
2. [tap 代理](/modules/tap.md) — `tap start --http-parse`（:8080 / :1080）

CA 证书在首次 `tap start` 时自动生成，存放于 `~/.cursor-tap/ca/`。

# Examples

```bash
# Windows
set HTTP_PROXY=http://localhost:8080
set HTTPS_PROXY=http://localhost:8080
set NODE_EXTRA_CA_CERTS=C:\Users\<user>\.cursor-tap\ca\ca.crt

# macOS/Linux
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
export NODE_EXTRA_CA_CERTS=~/.cursor-tap/ca/ca.crt
```

Cursor 仅需配置 **代理端口 8080**；Web Inspector 连接 **API 9090**，与 Cursor 无关。

## 相关

- [tap 代理](/modules/tap.md) — 8080/1080 端口
- [API 服务](/modules/api-server.md) — 9090 端口
- [开发工作流](/development/workflow.md) — 启动命令
- [管理 API](/api/management-api.md) — CA 证书下载 `/api/ca/cert`
