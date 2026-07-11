# Conductor 浏览器桌面端开发说明

本文说明如何在本地开发、调试和构建 Conductor 浏览器桌面端。

普通用户使用方式见 [README.md](../README.md#conductor-浏览器桌面端)。

## 环境要求

- Go（与 `go.mod` 一致）
- [Wails CLI](https://wails.io/)（当前仓库按 v2 使用）
- Node.js（仅用于从 logo 生成 exe 图标）
- Windows 10/11 + WebView2（运行桌面窗口需要）

## 本地运行

### 桌面窗口（默认）

```bash
wails dev
# 或
go run -tags desktop,dev .
```

桌面模式默认浏览 **可执行文件同目录下的 `output`**。用 `go run` / `wails dev` 时，可执行文件往往在临时目录，建议显式指定仓库内产出目录：

```bash
go run -tags desktop,dev . -root output
```

### 纯 HTTP 模式

适合用系统浏览器调试页面与 API：

```bash
go run . -web
```

默认地址：`http://127.0.0.1:8080`

指定端口或目录：

```bash
go run . -web -addr 127.0.0.1:8090
go run . -web -root output
```

## 构建便携 exe

```bash
npm install
npm run gen-icon   # 用 assets/favicon.svg 生成 exe 图标
wails build
Copy-Item build/bin/Conductor.exe .\Conductor.exe
```

产物：

- 仓库根目录 `Conductor.exe`（便于分发）
- 同时也会生成在 `build/bin/Conductor.exe`

## 自动构建

推送到 `main`，且变更命中程序相关路径（如 `main.go`、`internal/`、`assets/`、`wails.json` 等）时，GitHub Actions 会构建。纯 Markdown / Agent 文档变更不会触发。流水线会：

1. 根据 `assets/favicon.svg` 生成图标
2. 执行 `wails build`
3. 将 `Conductor.exe` 放到仓库根目录并提交
4. 同时上传 Actions Artifact

工作流文件：[`.github/workflows/build-desktop.yml`](../.github/workflows/build-desktop.yml)

## 相关目录

```text
main.go                 桌面 / -web 入口
internal/browser/       成果浏览 HTTP UI 与 API
wails.json              Wails 应用配置
frontend/               Wails 前端占位（页面由 Go Handler 提供）
assets/                 logo / favicon
scripts/gen-icon.mjs    SVG 转 appicon / ico
build/windows/          Windows 图标与清单
```
