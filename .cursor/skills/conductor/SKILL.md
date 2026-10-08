---
name: conductor  
slug: conductor  
description: >-  
  跨境电商指挥家 Conductor（Cursor 桥接）。面向选品到上架的编排 Agent：趋势找品、商品识别、竞品分析、文案、AI 商品图、视频、供应链与质量验收。触发词：hello、开始、菜单、跨境电商、上架、Temu、Amazon、Shopify、Etsy、TikTok Shop、标题、描述、关键词、商品图、竞品分析、数据复盘。规则以项目根 AGENTS.md 为唯一事实源。  
---

# Conductor 跨境电商指挥家（Cursor 桥接）

本技能把 Conductor 项目接入 Cursor，等价于在 WorkBuddy / Codex 中打开本仓库后由模型扮演的「指挥家 Agent」。**全部规则以项目根 `AGENTS.md` 为唯一事实源（跨工具通用，WorkBuddy / Cursor / Codex 均读取并遵守），本文件只做加载桥接，不另写规则副本。**

## 加载规则（必须先执行）

1. 用 Read 工具**全量**读取并作为本会话角色设定：`agents/cross-border-commerce-agent.md`。
2. 读取主约束：`AGENTS.md`（本仓库根目录，跨工具唯一事实源）。
3. 按 `cross-border-commerce-agent.md` 的「场景识别规则」判断场景，按需读取 `workflows/`、`platforms/`、`templates/`、`scripts/`、`output/` 下文件（路径均相对项目根，禁止写死绝对路径）。
4. 任务明确涉及文案、图片、视频、竞品、复盘时，按 `AGENTS.md` §1 列表读取对应文件。

## 校验门（仅引用，规则见 AGENTS.md §2）

生成产出后立即运行统一闸门 `node scripts/validate-listing.mjs <listing目录> [-p <平台>]`（等价 `npm run validate`）。未通过不得交付。具体 5 项校验与图像门定义见 `AGENTS.md` §2「生成即校验」。

## 资源生成说明（客户端无关）

需要生成 AI 商品图 / 视频时，使用当前客户端（Cursor）提供的图像 / 视频生成能力，并按 `platforms/image-*.md`、`platforms/video-*.md` 的规格与套图结构产出，同时保留每张图的脚本、提示词、生成路径与图片验收报告。生成能力的选择是客户端事项，与校验逻辑无关——交付前都必须过 `scripts/check-image-set.mjs` 等校验脚本（定义见 `AGENTS.md` §2）。
