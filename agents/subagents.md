# 核心 Agent 的子 Agent 体系

## 总体结构

核心 Agent（Conductor Agent / 指挥家）负责流程编排，不负责亲自完成所有专业任务。每个阶段由一个或多个子 Agent 完成专业产出，核心 Agent 负责派发、验收、合并和推进。

## 子 Agent 列表

| 子 Agent | 阶段 | 职责 | 关键输出 |
| --- | --- | --- | --- |
| `intake-agent` | 信息接收 | 按任务收集最低必要输入，根据商品名称和图片识别产品，并区分事实、推测、市场建议和待确认项 | 输入摘要、产品识别草案、缺失信息、风险字段 |
| `trend-and-video-discovery-agent` | 趋势与短视频找品 | 当用户只有方向、品类、人群、关键词或视频链接时，先按目标国家选择数据源和本地语言关键词，再分析趋势增长、短视频场景、商品出现频率、购买意图和品牌真空 | 目标市场数据源、候选商品清单、趋势/场景证据、品牌真空判断、初筛评分 |
| `product-selection-agent` | 选品初筛 | 判断需求、竞争、平台适配、供应链和新手友好度 | 选品评分、推进/暂缓建议 |
| `market-demand-agent` | 需求验证 | 判断搜索需求、趋势、人群、季节性和使用频率 | 需求强度、目标人群、购买动机 |
| `platform-strategy-agent` | 平台匹配 | 判断适合 Amazon、Temu、Shopify、Etsy 等哪个平台 | 平台优先级、平台机会和风险 |
| `competitor-research-agent` | 竞品分析 | 研究目标市场热销榜、同类搜索结果、代表性详情页和可访问评论，分析标题、图片、价格、评价、痛点和差异化 | 数据来源记录、竞品矩阵、评论痛点、待验证机会 |
| `profit-agent` | 利润测算 | 测算采购、物流、佣金、广告、退货后的利润空间 | 售价建议、毛利区间、风险线 |
| `compliance-agent` | 合规检查 | 检查商标、专利、认证、敏感类目和禁用表达 | 风险等级、禁用表达、补证清单 |
| `supply-chain-agent` | 供应链确认 | 判断 MOQ、交期、质检、包装、SKU 稳定性 | 供应链问题清单、首批建议 |
| `positioning-agent` | 定位卖点 | 综合自身事实、市场研究和仍可执行的未来节日，自动生成核心功能、希望突出卖点、主要购买动机和定制示例，并结合平台策略补充类目必需属性 | 可复制完整商品示例、一句话定位、卖点矩阵、节日机会 |
| `listing-strategy-agent` | 上架策略 | 设计 SKU、价格、offer、物流和字段结构 | 上架策略、字段 checklist |
| `copywriting-agent` | 文案生产 | 按平台生成标题、要点、描述、关键词、广告文案 | 平台适配文案包 |
| `visual-production-agent` | 图片/视频视觉资产 | 按平台规则规划图片和视频，并在工具可用时生成资产 | 套图脚本、视频脚本、生成提示词、资产路径、验收报告 |
| `listing-qa-agent` | 上架验收 | 检查 listing 完整度、一致性和可发布状态 | 上架 QA 报告、阻塞项 |
| `growth-review-agent` | 测试复盘 | 设计上架后测试指标和优化节奏 | 测试计划、复盘模板 |

### 文案生产语言边界

- `copywriting-agent` 输出的所有卖家侧内容必须使用中文，包括章节标题、栏目名、字段名、表头、流程说明、内部备注、测试策略、事实清单、待确认项、风险提示和禁止声明，包括另存为 `文案资产.md` 的独立文件。
- 买家可见的实际文案正文使用目标市场语言。美国市场使用“中文栏目名 + 英文文案正文”，例如 `## 商品标题` 下方写英文商品标题。
- 只有可直接复制到平台前台、搜索字段或广告投放字段的最终文案值使用目标市场语言；给卖家看的解释、提醒和判断不得整段使用英文。
- 禁止使用 `Product Title`、`Five Bullet Points`、`Long Description`、`Search Keywords`、`Backend Keywords` 等英文栏目名。

## 核心 Agent 调度规则

### 1. 阶段判断

核心 Agent 先判断用户当前处于哪个阶段：

- 只有产品想法：从 `product-selection-agent` 开始。
- 只有方向、品类、人群、关键词或视频链接：先调度 `trend-and-video-discovery-agent`，提取候选商品后再进入 `product-selection-agent`。
- 有产品资料但未选平台：先调度 `platform-strategy-agent`。
- 有平台和竞品：进入 `competitor-research-agent`。
- 已确定要上架：进入 `listing-strategy-agent`、`copywriting-agent`、`visual-production-agent`。
- 已有 listing 草稿：进入 `listing-qa-agent`。
- 已上架：进入 `growth-review-agent`。

用户选择任务后，必须按 `workflows/start-guide.md` 先给出带示例值的最小输入模板。需要具体商品时，系统自动完成商品识别与电商研究，并返回包含四项营销建议和平台必需属性的完整可复制示例；用户回传修改后的示例前不得进入生产。

### 2. 串行和并行

必须串行的阶段：

- 信息接收必须先于所有阶段。
- 合规检查必须在最终上架前完成。
- 上架 QA 必须在文案和套图完成后进行。

可以并行的阶段：

- 需求验证、竞品分析、平台匹配可以并行。
- 文案生产、图片套图和视频脚本可以在定位确认后并行。
- 利润测算和供应链确认可以并行。

### 3. 阻塞条件

出现以下情况时，核心 Agent 必须暂停推进：

- 产品功能、材质、规格无法确认。
- 目标平台或目标市场缺失。
- 认证、功效、适配型号等高风险事实无证据。
- 利润测算缺少采购价或物流成本。
- 竞品分析没有任何竞品资料且用户要求“基于真实竞品”。
- 图片脚本依赖产品外观但没有图片或外观说明。

### 4. 验收标准

每个子 Agent 的输出必须满足：

- 明确引用输入事实或标注假设。
- 给出可执行结论，不只写泛泛建议。
- 输出能被下一阶段直接使用。
- 标出风险、缺失信息和需要人工确认的内容。

## 推荐执行链路

```text
intake-agent
  -> trend-and-video-discovery-agent（仅在只有方向/趋势找品/视频找品时启用）
  -> product-selection-agent
  -> market-demand-agent + platform-strategy-agent + competitor-research-agent
  -> profit-agent + compliance-agent + supply-chain-agent
  -> positioning-agent
  -> listing-strategy-agent
  -> copywriting-agent + visual-production-agent
  -> listing-qa-agent
  -> growth-review-agent
```
