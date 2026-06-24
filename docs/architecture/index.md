---
okf_version: "0.1"
---

# Cursor-Tap 项目架构

MITM 代理 + gRPC 解析 + Web Inspector 工具链的 OKF 知识 Bundle。

## 概览

* [系统概览](/overview.md) - 整体架构、组件关系与端口解耦
* [Monorepo 结构](/monorepo.md) - 目录布局与 Nx 项目
* [数据流](/data-flow.md) - 从代理接入到 WebSocket 推送的完整链路

## 模块

* [tap 代理](/modules/tap.md) - Go MITM 代理、SQLite 录制与 notify
* [API 服务](/modules/api-server.md) - 独立 tap-api 进程（9090）
* [internal 包](/modules/internal.md) - CA、MITM、httpstream、recordstore
* [Proto 包](/modules/proto.md) - 逆向 proto 定义与 buf 生成
* [Web Inspector](/modules/web.md) - Next.js 四栏 UI 与状态管理

## API 与契约

* [管理 API](/api/management-api.md) - REST 与 WebSocket 端点
* [Record 类型契约](/record-types.md) - 前后端 JSON 对齐

## 开发与运维

* [开发工作流](/development/workflow.md) - 安装、启动、构建与测试命令
* [Cursor 配置](/development/cursor-config.md) - 代理与 CA 证书环境变量
* [开发工具](/tools.md) - proto 提取、日志还原等辅助工具
* [测试](/testing.md) - 测试栈、约定与覆盖范围
* [技术依赖](/dependencies.md) - Go、前端与工具链依赖
* [已知问题](/known-issues.md) - 当前限制与待修复项

## 外部参考

* [Cursor 逆向笔记](../cursor-reverse-notes-1.md) - MITM 原理、proto 提取、Connect Protocol 解析
* [README](../../README.md) - 快速开始与使用说明
