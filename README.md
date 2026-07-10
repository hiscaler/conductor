# Conductor

Conductor 是一个面向跨境电商的 AI 商品增长编排框架。

它的目标不是只生成一段商品文案，而是把一个产品想法、产品资料或产品链接，拆解成从选品判断到平台上架、内容生产和测试复盘的一整套执行流程。

当前核心 agent：

- **指挥家 Agent**

指挥家负责流程处理和任务编排，不负责亲自完成所有专业工作。它会判断当前阶段，分派给对应子 agent，汇总结果，检查风险，并给出下一步动作。

## 适合谁使用

Conductor 适合以下场景：

- 有一个产品想法，想判断是否适合做跨境电商。
- 已经找到供应商，想判断适合 Amazon、Temu、Shopify、Etsy、TikTok Shop 还是其他平台。
- 已经有产品资料，想生成上架文案、关键词和图片套图脚本。
- 已经有竞品链接，想做竞品分析、差评痛点提炼和差异化定位。
- 已经准备上架，想检查 listing 字段、图片、合规和风险。
- 已经上架，想制定测试、投放和复盘计划。

## 核心工作流

Conductor 把跨境电商从选品到上架拆成 12 个节点：

1. 选品来源记录
2. 需求验证
3. 竞争分析
4. 平台匹配
5. 利润测算
6. 合规和侵权检查
7. 供应链确认
8. SKU 和包装设计
9. 定价和 offer 设计
10. 内容资产生产
11. 上架检查
12. 测试和复盘

完整说明见 [商品出海流程](./workflows/product-to-listing.md)。

## Agent 架构

Conductor 采用“指挥家 + 阶段子 agent”的结构。

指挥家 Agent 负责：

- 判断用户当前处于哪个阶段。
- 决定应该调用哪些子 agent。
- 把上一个阶段的结论传给下一个阶段。
- 检查子 agent 的输出是否可执行。
- 在信息不足或风险过高时暂停推进。
- 汇总最终结果，给出下一步执行清单。

阶段子 agent 负责：

- 选品初筛
- 需求验证
- 平台匹配
- 竞品分析
- 利润测算
- 合规检查
- 供应链确认
- 定位卖点
- 上架策略
- 文案生产
- 图片套图和图片生成
- 上架验收
- 测试复盘

完整子 agent 列表见 [指挥家子 Agent 体系](./agents/subagents.md)。
图片生成能力见 [Visual Production Agent](./agents/visual-production-agent.md)。

## 支持的平台

当前内置平台画像覆盖：

- Amazon
- Temu
- Shopify
- Etsy
- TikTok Shop
- eBay
- AliExpress
- Walmart Marketplace

不同平台会使用不同的标题策略、图片策略、关键词策略、合规检查和交付物结构。

例如：

- Amazon 更重视搜索关键词、五点描述、白底主图、评价痛点和 A+ 页面。
- Temu 更重视低价转化、规格清晰、套装展示和高信息密度图片。
- Shopify 更重视品牌信任、落地页结构、场景图、FAQ 和广告承接。
- Etsy 更重视手作感、定制说明、礼品场景、材质细节和 tags。
- TikTok Shop 更重视短视频种草、直播话术、封面点击和移动端理解效率。

完整平台规则见 [平台画像](./platforms/platform-profiles.md)。
平台套图规则见 [平台套图适配规则](./platforms/image-set-rules.md)。
图片尺寸、比例和质量标准见 [平台图片技术规格](./platforms/image-technical-specs.md)。

## 如何使用

### 1. 准备产品信息

优先使用 [产品输入模板](./templates/product-input.md) 填写信息。

最少建议提供：

- 产品名称
- 产品类目
- 目标平台
- 目标市场
- 核心功能
- 规格参数
- 材质或结构说明
- 目标售价或成本区间
- 供货方式或采购状态
- 产品图片、链接或外观描述

如果资料不完整也可以开始，指挥家会先整理缺失信息，并标注哪些内容只能作为假设。

### 2. 选择当前阶段

你可以按当前业务状态发起任务：

```text
我只有一个产品想法，请先帮我判断是否值得做。
```

```text
我已经有产品资料，但还没确定平台，请帮我判断适合哪个平台。
```

```text
我有 3 个竞品链接，请做竞品分析并提炼差异化定位。
```

```text
我已经决定上 Amazon，请输出标题、五点描述、关键词和 7 张图脚本。
```

```text
我已经有 listing 草稿，请帮我做上架前检查。
```

### 3. 指挥家分派子 agent

指挥家会根据阶段自动决定执行链路。

推荐链路：

```text
intake-agent
  -> product-selection-agent
  -> market-demand-agent + platform-strategy-agent + competitor-research-agent
  -> profit-agent + compliance-agent + supply-chain-agent
  -> positioning-agent
  -> listing-strategy-agent
  -> copywriting-agent + visual-production-agent
  -> listing-qa-agent
  -> growth-review-agent
```

### 4. 查看输出

标准输出结构见 [产出模板](./templates/production-output.md)。

`production-output.md` 是强制输出契约，不是参考模板。最终结果必须保留模板里的所有章节、表格和字段；缺失信息写“未提供”或“待确认”，不适用内容写“不适用”。

一次完整产出通常包含：

- 输入摘要
- 缺失信息与假设
- 选品初筛
- 平台策略
- 竞品分析
- 产品定位
- 上架策略
- 文案资产
- 套图脚本
- 生成图片路径
- 图片验收报告
- 合规检查
- 下一步执行清单

## 示例输入

```text
产品名称：可折叠宠物饮水杯
产品类目：宠物用品
目标平台：Amazon、TikTok Shop
目标市场：美国
核心功能：外出遛狗时给宠物饮水
规格参数：容量 350ml，可折叠，带挂扣
材质：食品级硅胶和塑料，具体认证待确认
目标售价：9.99-14.99 美元
供货状态：已有 1688 供应商，支持小批量采购
竞品链接：暂未提供
素材：有产品白底图，可生成场景图
要求：先判断是否值得做，再分别给 Amazon 和 TikTok Shop 的内容策略
```

## 示例输出方向

指挥家会先判断：

- 是否需要补充认证、材质证明、包装尺寸和采购价。
- 这个产品在 Amazon 是否面临同质化和评价门槛。
- 这个产品在 TikTok Shop 是否适合短视频演示。
- 是否应该做套装、颜色组合或便携场景差异化。
- 是否进入文案和图片生产阶段。

如果信息足够，它会继续输出：

- 平台优先级
- 竞品调研清单
- 定位和卖点
- Amazon listing 文案
- TikTok Shop 短视频脚本
- 商品套图脚本
- 上架前检查清单

## 目录说明

```text
agents/
  cross-border-commerce-agent.md     指挥家 Agent 角色、流程和提示词
  cross-border-commerce-agent.json   指挥家 Agent 结构化配置
  subagents.md                       阶段子 agent 和调度规则
  visual-production-agent.md         图片套图和图片生成 agent

platforms/
  platform-profiles.md               不同电商平台的内容和上架规则
  image-set-rules.md                 不同平台的套图结构和图片内容规则
  image-technical-specs.md           不同平台的图片尺寸、比例、格式和质量建议

templates/
  product-input.md                   产品输入模板
  production-output.md               标准产出模板

workflows/
  product-to-listing.md              从选品到上架的完整流程

output/
  images/                            生成后的商品图片建议保存目录
```

## 使用原则

- 不编造认证、材质、功效、适配型号、专利、销量、排名和评价。
- 没有真实竞品资料时，只输出竞品调研清单，不假装已经调研。
- 多个平台同时推进时，必须分别输出平台版本。
- 文案和图片脚本必须服务于具体产品，不输出泛泛模板。
- 上架前必须经过合规检查和 listing QA。
