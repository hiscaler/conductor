# Conductor Project Instructions

本仓库的核心用途是运行“指挥家 Agent”，为跨境电商产品生成从输入审查、平台策略、文案资产到 AI 商品图、视频脚本和下一步动作的标准产出。

## 默认行为

当用户请求跨境电商、商品上架、Temu/Amazon/Shopify/Etsy/TikTok Shop 等平台内容、竞品分析、标题描述关键词、套图、AI 商品图或视频脚本时，必须按本项目的指挥家规则执行。

必须先读取并遵守：

- `agents/cross-border-commerce-agent.md`
- `templates/production-output.md`
- `platforms/platform-profiles.md`
- `platforms/image-set-rules.md`
- `platforms/image-technical-specs.md`
- `workflows/action-menu.md`

如涉及图片或视频生成，还必须读取：

- `agents/visual-production-agent.md`
- `platforms/video-set-rules.md`
- `platforms/video-technical-specs.md`

## 输出契约

最终产出必须严格使用 `templates/production-output.md`。

- 保留 16 个章节和所有表格。
- 不得自由改名、删节或改成普通聊天格式。
- 缺失信息写“未提供”或“待确认”。
- 未执行章节写“本轮未执行”。
- 下一步执行清单不能为空。

如果输出没有 `## 1. 输入摘要` 到 `## 16. 下一步执行清单`，视为没有完成本项目任务。

## 买家侧语言

目标市场为美国时，买家可见内容默认使用英文：

- 商品标题
- 五点描述
- 长描述
- 搜索关键词
- 后台关键词
- 广告短文案
- 图片内文案
- 视频字幕/口播

中文只用于章节标题、流程说明、内部备注和用户沟通。

## 图片生成规则

“套图”“输出套图”“生成套图”“生成一套商品图”默认表示生成 AI 图片，不只是输出脚本。

每张商品图必须是独立文件。禁止用一张拼图、九宫格、预览图或带多个缩略图的单张图片代替整套图。

Temu 定制类商品默认以 8 张独立图片为完成标准：

1. 最终定制主图
2. 到手内容/包装图
3. 定制操作示意图
4. 尺寸规格图
5. 细节放大图
6. 生活场景图
7. 定制效果示例图
8. 卖点集合图

生成图片后必须在第 10、11 节写明：

- 应生成图片数
- 实际独立图片文件数
- 缺失图片类型
- 整套图状态

未达到完成标准时，必须写“未完成”，并在下一步动作中提供“继续生成缺失图片”。

## 定制类商品

定制类商品必须先把定制内容放到商品表面，再生成主图、包装图、场景图和说明图。

- 主图必须展示最终定制商品。
- 包装图必须展示最终定制商品和包装/到手内容。
- 定制操作示意图也必须以最终定制商品为主体，再用箭头或标注说明定制区域。
- `YOUR NAME`、`YOUR TEXT`、`CUSTOM PHOTO` 只能用于操作说明标签，不能作为最终商品图案。
- 未提供真实定制内容时，使用真实感中性样例，例如 `Emma`、`Mom`、`Best Dad`、`Coffee Time` 或通用风景/花朵/几何图片样例，并标注为示例。

## 尺寸和辅图文案

尺寸规格图必须在图片画面内同时标注 `cm` 和 `inch`。只显示 `cm` 的尺寸图不可用。

场景图、包装图、卖点图、规格图和对比图必须有基于已确认事实的简短英文文案，不能只是无文字图片。

不得编造未确认功能，例如：

- insulated
- keeps hot
- keeps cold
- leakproof
- dishwasher safe
- microwave safe

除非用户明确提供这些功能或证明。
