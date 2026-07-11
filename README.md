# Conductor

Conductor 是一个面向跨境电商的 AI 商品增长编排框架。

它的目标不是只生成一段商品文案，而是把一个产品想法、产品资料或产品链接，拆解成从选品判断到平台上架、内容生产和测试复盘的一整套执行流程。

当前核心 agent：

- **Conductor Agent**

Conductor 负责流程处理和任务编排，不负责亲自完成所有专业工作。它会判断当前阶段，分派给对应子 agent，汇总结果，检查风险，并给出下一步动作。

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

Conductor 采用“主 Agent + 阶段子 Agent”的结构。

主 Agent 负责：

- 判断用户当前处于哪个阶段。
- 决定应该调用哪些子 agent。
- 把上一个阶段的结论传给下一个阶段。
- 检查子 agent 的输出是否可执行。
- 在信息不足或风险过高时暂停推进。
- 汇总最终结果，给出下一步执行清单。

阶段子 Agent 负责：

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
- 图片套图、视频脚本和视觉资产生成
- 上架验收
- 测试复盘

完整子 agent 列表见 [Conductor 子 Agent 体系](./agents/subagents.md)。
图片和视频生成能力见 [Visual Production Agent](./agents/visual-production-agent.md)。

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
平台视频规则见 [平台视频适配规则](./platforms/video-set-rules.md)。
视频尺寸、比例、时长和质量标准见 [平台视频技术规格](./platforms/video-technical-specs.md)。

## 如何使用

### 1. 先判断你处在哪个场景

不是每次都需要从选品开始。你可以按当前业务状态直接发起任务，Conductor 会自动跳过不需要的阶段。

| 场景 | 你已经有什么 | 适合让 Conductor 做什么 | 会跳过什么 |
| --- | --- |-------------------| --- |
| 只有产品想法 | 一个品类或产品概念 | 判断是否值得做、适合哪个平台    | 文案和图片生成可以暂缓 |
| 已经有产品 | 产品名称、规格、材质、供应商 | 直接输出平台文案、关键词和套图脚本 | 可跳过选品来源和部分需求验证 |
| 已经决定平台 | 产品资料 + 目标平台 | 直接按平台生成上架内容       | 可跳过平台匹配 |
| 已经有竞品 | 产品资料 + 竞品链接 | 做竞品分析、差异化定位、文案和套图 | 不需要先做通用市场判断 |
| 只要文案 | 产品资料 + 平台 | 只生成标题、描述、关键词、广告文案 | 图片生成、QA、复盘可暂缓 |
| 只要套图 | 产品资料 + 平台 + 卖点 | 只生成套图脚本或图片        | 选品、利润、文案可简化 |
| 只要视频 | 产品资料 + 平台 + 卖点 | 只生成视频脚本、分镜或视频     | 选品、利润、图片可简化 |
| 已经有 listing 草稿 | 标题、描述、图片或草稿 | 做上架前 QA 和合规检查     | 不重新生成全部内容 |
| 已经上架 | 后台数据或运营反馈 | 做测试计划、复盘和优化建议     | 不重新做选品 |

### 2. 准备产品信息

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

如果资料不完整也可以开始，Conductor 会先整理缺失信息，并标注哪些内容只能作为假设。

### 3. 按场景发起任务

#### 场景 A：只有产品想法，先判断是否值得做

适合在还没找供应商、还没决定平台时使用。

```text
我有一个产品想法：便携宠物饮水杯。
目标市场：美国。
请先判断这个产品是否值得做，适合哪些平台，并列出需要补充的信息。
```

Conductor 通常会调度：

```text
intake-agent
-> product-selection-agent
-> market-demand-agent
-> platform-strategy-agent
-> compliance-agent
```

常见产物：

- 选品初筛
- 平台优先级
- 风险清单
- 下一步需要补充的信息

#### 场景 B：已经有产品，只要输出文案和套图

适合你已经有供应商、规格、目标平台，不需要重新判断“要不要做”。

```text
我已经有产品了，不需要完整选品流程。

产品名称：可折叠宠物饮水杯
产品类目：宠物用品
目标平台：Temu
目标市场：美国
核心功能：外出遛狗时给宠物饮水
规格参数：容量 350ml，可折叠，带挂扣
材质：食品级硅胶和塑料，认证待确认
供货状态：已有 1688 供应商

请直接输出 Temu 标题、描述、关键词和 7 张图脚本。
```

Conductor 通常会调度：

```text
intake-agent
-> platform-strategy-agent
-> positioning-agent
-> copywriting-agent
-> visual-production-agent
-> compliance-agent
```

会简化或跳过：

- 深度选品判断
- 多平台匹配
- 利润测算，除非你提供采购价和物流成本

常见产物：

- 文案资产
- 套图脚本
- 图片生成提示词
- 合规提醒
- 可执行动作菜单

#### 场景 C：已经决定平台，只做单平台上架内容

适合你明确要上 Amazon、Temu、Shopify、Etsy 等其中一个平台。

```text
我已经决定上 Amazon。
请不要做多平台比较，直接按 Amazon 输出 listing 文案、关键词、A+ 页面建议和 7 张图片脚本。

产品资料如下：
...
```

Conductor 会只读取目标平台画像和图片规则，不再做平台优先级比较。

常见产物：

- 单平台策略
- 单平台文案
- 单平台套图脚本
- 单平台合规检查

#### 场景 D：已有竞品链接，做竞品分析后再产出

适合你想基于真实竞品找差异化。

```text
我有 3 个竞品链接，请先做竞品分析，再输出差异化定位、文案和套图脚本。

目标平台：Amazon
目标市场：美国
竞品链接：
1. ...
2. ...
3. ...
产品资料：
...
```

Conductor 通常会调度：

```text
intake-agent
-> competitor-research-agent
-> positioning-agent
-> copywriting-agent
-> visual-production-agent
```

注意：

- 如果没有联网或页面不可访问，Conductor 不能假装调研成功。
- 没有真实竞品数据时，只输出竞品调研清单。

#### 场景 E：只保存文案，不生成图片

适合你只想把文案先落地成文件。

```text
请只生成并保存 Temu 文案资产，不生成图片。

产品资料：
...
```

推荐后续动作：

```text
1. 保存文案资产
```

产物位置：

```text
output/{目标平台}/{产品类目}/{产品名称}/文案/文案资产.md
```

#### 场景 F：只生成主图或整套图片

适合你已经有文案和套图脚本，下一步只要图片。

```text
请根据现有套图脚本生成主图。
目标平台：Temu
产品类目：宠物用品
产品名称：可折叠宠物饮水杯
```

推荐后续动作：

```text
1. 生成主图
2. 生成整套图片
```

产物位置：

```text
output/{目标平台}/{产品类目}/{产品名称}/图片/01-主图.png
output/{目标平台}/{产品类目}/{产品名称}/图片/图片验收报告.md
```

注意：

- 没有供应商实物图时，生成图只能作为概念图。
- 正式上架前需要用实物图校准外观、结构、颜色和配件。

#### 场景 G：只生成视频脚本或短视频

适合你已经有产品定位和卖点，需要短视频脚本、分镜、口播或视频生成提示词。

```text
请只输出 Temu 短视频脚本，不生成图片。

产品名称：可折叠宠物饮水杯
产品类目：宠物用品
目标平台：Temu
目标市场：美国
核心卖点：350ml、可折叠、带挂扣、适合外出遛狗
要求：输出 15 秒竖屏视频脚本、分镜、字幕和视频生成提示词。
```

推荐后续动作：

```text
1. 保存视频脚本
2. 生成主视频
3. 生成视频组
```

产物位置：

```text
output/{目标平台}/{产品类目}/{产品名称}/视频/视频脚本.md
output/{目标平台}/{产品类目}/{产品名称}/视频/01-主视频.mp4
output/{目标平台}/{产品类目}/{产品名称}/视频/视频验收报告.md
```

注意：

- 没有实物视频时，生成视频只能作为概念视频或投放方向。
- TikTok Shop、Temu 等移动端平台优先使用 9:16 竖屏视频。
- 视频必须无声可理解，字幕要避开平台按钮和价格区域。

#### 场景 H：已有 listing 草稿，做上架前 QA

适合你已经有标题、描述、图片或上架草稿，需要检查是否能发布。

```text
我已经有 listing 草稿，请做上架前 QA。

目标平台：Temu
产品资料：
...
标题：
...
描述：
...
图片：
...
```

Conductor 通常会调度：

```text
listing-qa-agent
-> compliance-agent
```

常见产物：

- 上架前 QA 报告
- 合规风险清单
- 缺失字段
- 是否可发布建议

#### 场景 I：已经上架，做测试和复盘

适合你已经有曝光、点击、转化、退款、差评等数据。

```text
这个产品已经上架，请根据以下数据做复盘和优化计划。

平台：Temu
曝光：
点击率：
转化率：
退款率：
主要差评：
广告花费：
```

Conductor 通常会调度：

```text
growth-review-agent
-> copywriting-agent
-> visual-production-agent
```

常见产物：

- 测试复盘
- 图片优化建议
- 标题/价格/卖点 A/B 测试计划
- 是否补货、降价、换图或下架建议

### 4. Conductor 分派子 agent

Conductor 会根据阶段自动决定执行链路。

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

### 5. 查看输出

标准输出结构见 [产出模板](./templates/production-output.md)。

`production-output.md` 是强制输出契约，不是参考模板。最终结果必须保留模板里的所有章节、表格和字段；缺失信息写“未提供”或“待确认”，不适用内容写“不适用”。

最后一节“下一步执行清单”会同时给出文字说明和可执行动作。理想交互是弹出“选择下一步动作”的对话框，用户直接点选 `保存文案资产`、`生成主图`、`生成整套图片`、`保存视频脚本`、`生成主视频` 等动作，Conductor 再调度对应子 agent。只有运行环境不支持弹窗时，才退回输入数字序号的方式，例如输入 `1` 或 `1,2`。

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
- 视频脚本
- 生成视频路径
- 视频验收报告
- 合规检查
- 下一步执行清单

### 6. 选择下一步动作

每次输出最后都会给出“下一步执行清单”。

理想情况下，会弹出“选择下一步动作”的对话框。你可以直接点选：

- 保存文案资产
- 生成主图
- 生成整套图片
- 保存视频脚本
- 生成主视频
- 生成视频组
- 保存套图脚本
- 生成上架包
- 做上架前 QA
- 补充竞品调研
- 利润测算
- 合规风险复核
- 制定测试计划

如果环境不支持弹窗，会退回数字序号：

```text
请选择下一步动作：
1. 保存文案资产
2. 生成主图
3. 做上架前 QA

请输入数字，例如：1 或 1,2
```

### 7. 查看输出目录

所有产物都按同一个产品目录聚合：

```text
output/{目标平台}/{产品类目}/{产品名称}/
```

例如：

```text
output/Temu/宠物用品/可折叠宠物饮水杯/
  文案/文案资产.md
  图片/01-主图.png
  图片/图片验收报告.md
  视频/视频脚本.md
  视频/01-主视频.mp4
  视频/视频验收报告.md
  上架/上架包.md
  上架/上架前QA报告.md
  合规/合规风险报告.md
```

这样同一个产品的文案、图片、QA、合规和测试材料都会放在一起。

## 本地输出浏览器

项目内置 Conductor 桌面成果浏览器，用来查看 `output` 目录下生成的内容。

### 普通用户

1. 仓库根目录下载 `Conductor.exe`（免安装），或从 GitHub Actions 产物获取
2. 双击打开即可
3. 默认浏览与 exe 同目录下的 `output` 文件夹

需要 Windows 10/11，并已安装 WebView2（多数系统已自带）。

推送到 `main` 后，GitHub Actions 会自动构建并把最新 `Conductor.exe` 提交到仓库根目录。

### 开发者

桌面窗口（默认）：

```bash
wails dev
# 或
go run -tags desktop,dev .
```

纯 HTTP 模式（用系统浏览器访问）：

```bash
go run . -web
```

默认地址：`http://127.0.0.1:8080`

指定目录或端口：

```bash
go run . -web -addr 127.0.0.1:8090
go run . -web -root output
```

构建可双击的便携程序：

```bash
npm install
npm run gen-icon   # 用 assets/favicon.svg 生成 exe 图标
wails build
Copy-Item build/bin/Conductor.exe .\Conductor.exe
```

产物在仓库根目录 `Conductor.exe`（同时也会生成在 `build/bin/Conductor.exe`）。

功能：

- 打开后默认预览项目 `README.md`，方便先查看使用说明。
- 左侧显示 `output` 目录树。
- 点击左侧目录，右侧显示目录内容。
- 点击 Markdown 文件，右侧直接预览文案、QA 报告、合规报告等内容。
- 点击图片文件，右侧直接预览图片。
- 点击视频文件，右侧使用浏览器播放器预览视频。
- 访问范围限制在 `output` 目录内，避免读取项目其他文件。

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

Conductor 会先判断：

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
  cross-border-commerce-agent.md     Conductor Agent 角色、流程和提示词
  cross-border-commerce-agent.json   Conductor Agent 结构化配置
  subagents.md                       阶段子 agent 和调度规则
  visual-production-agent.md         图片套图和图片生成 agent

main.go                              Conductor 桌面入口（Wails）
internal/browser/                    成果浏览 HTTP UI/API
wails.json                           桌面应用构建配置
frontend/                            Wails 前端占位（页面由 Go Handler 提供）

assets/
  coor-logo.svg                      Coor 浏览器 logo
  favicon.svg                        浏览器标签页图标

platforms/
  platform-profiles.md               不同电商平台的内容和上架规则
  image-set-rules.md                 不同平台的套图结构和图片内容规则
  image-technical-specs.md           不同平台的图片尺寸、比例、格式和质量建议
  video-set-rules.md                 不同平台的视频结构和内容规则
  video-technical-specs.md           不同平台的视频尺寸、比例、时长和质量建议

templates/
  product-input.md                   产品输入模板
  production-output.md               标准产出模板

workflows/
  product-to-listing.md              从选品到上架的完整流程
  action-menu.md                     下一步可执行动作和 subagent 调度规则
  output-structure.md                输出目录和中文文件命名规范

output/
  {目标平台}/{产品类目}/{产品名称}/     同一产品的所有交付物统一保存目录
```

## 使用原则

- 不编造认证、材质、功效、适配型号、专利、销量、排名和评价。
- 没有真实竞品资料时，只输出竞品调研清单，不假装已经调研。
- 多个平台同时推进时，必须分别输出平台版本。
- 文案和图片脚本必须服务于具体产品，不输出泛泛模板。
- 上架前必须经过合规检查和 listing QA。
