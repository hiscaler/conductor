# content-package 结构草案（v0.1）

> 状态：讨论中（R1）  
> 归属：[对接总计划](./对接总计划.md) 第 R1 轮产出  
> 定位：Conductor → 中台的机对机内容契约草案，不代表已可刊登。

## 设计约束（已对齐）

- content-package 是**新增**契约，导出时**不得改变 Agent 已产出的内容**，也**不得触发新的内容生产**；从已有文案资产组装 `copywriting`，不另起空占位当成品。
- 字段结构**不绑定具体平台 API**；默认按单平台包落盘，上品文案放 `copywriting`，平台差异（图集等）放 `platform_overrides`。
- 字段键默认使用**中文**，与 `data/product-catalog.csv` 列名和 `platforms/image-set-rules.md` 图型名一一对应；另附英文对照表供中台落库。上架文案对象英文键为 **`copywriting`**（不用 `copyright`）。
- 中性内容与平台覆盖都由 Conductor 产出；接口字段（类目 ID、属性码、API 报文）由中台在推送时组装。

## 顶层结构

```text
content-package
├── 元数据          metadata
├── 商品事实        product_facts   （中性，逐条带状态）
├── 中性基础内容    base_content
│   ├── 语义素材    semantics
│   └── 上架文案    copywriting     （单平台包的上品正文，必填）
├── 平台覆盖[]      platform_overrides
├── 媒体[]          media           （实体文件，按 id 引用）
└── 门禁与风险      gates
```

## 1. 元数据 metadata

包头元数据（说明这是哪一包、谁的货、面向哪、内容层状态），不放标题/描述等正文。
- `schema_version`：契约版本
- `package_id`：内容包标识
- `generated_at`：生成时间
- `spu`：SPU
- `skus[]`：`{ sku, qty }`，沿用现有 SKU 与数量规则
- `listing_id`：单 SKU 或组合标识（沿用 `workflows/output-structure.md`）
- `sales_unit`：销售单位与组合形态（单件 / 多件 / 套装 / 混色 / 混款）
- `product_name`：商品名称（业务文本）
- `suggested_category`：建议类目（业务文本，**非**平台类目 ID）
- `content_platform` / `content_market`：本包 `copywriting` 与套图实际按哪个平台-市场生产（与落盘目录一致）
- `recommended_platforms[]`：推荐可推送的平台与市场列表（仅建议，**不强制**中台照单全推）
- `content_status`：内容层自评（`ready` / `partial` / `blocked`）；仅 `ready` 可发送到中台
- `blocked_reasons[]`：内容层受阻原因

说明：

- `recommended_platforms` **只是推荐标记**；中台决定实际推送到哪些平台，也可推到列表外的平台。
- Conductor **不因**推荐了多个平台，就必须预生成多套 `platform_overrides.copywriting`。
- 本包正文与图集按 `content_platform`/`content_market` 生产；中台推到其他平台时，自行做内容适配或接口映射（超出本包已提供内容的部分由中台负责）。

## 2. 商品事实 product_facts

逐条记录，**`key` 必须与 `data/product-catalog.csv` 当前表头完全一致**，不得缩写或改名。

当前主表表头（以仓库 CSV 为准，变更时同步本契约）：

```text
SPU, SKU, 商品名称, 品牌, 产品类目, 产品子类目, 产品类型, 包含内容, 型号,
颜色/款式, 尺码/规格, 长度cm, 宽度cm, 高度cm, 净重g, 材质, 结构/表面工艺,
已确认功能, 商品特点, 适用对象, 使用场景, 使用/护理说明,
定制内容, 定制位置, 定制工艺,
包装清单, 包装方式, 包装长度cm, 包装宽度cm, 包装高度cm, 包装毛重g,
风险/禁用声明, 资料更新时间, 备注
```

注意：正确键名是 `颜色/款式`、`材质`、`风险/禁用声明` 等，**不是**口语缩写「颜色」「认证信息」。主表若无某列（例如当前无独立「认证信息」列），不得在 product_facts 里虚构该 key。

```json
{ "key": "颜色/款式", "value": "白色杯身、红色内胆、红色把手", "status": "confirmed", "source": "catalog" }
```

- `status`：`confirmed` / `unconfirmed` / `missing`
- `source`：`catalog` / `user` / `supplier` / `image` / `market_suggestion`
- `unit`：**可选**。键名已含单位时（如 `净重g`、`长度cm`）**不写** `unit`
- 品类专属参数走 `product-attributes.csv` → `extended_attributes[]`，不塞进主表列
- `market_suggestion` 永不等于已确认事实；中台不得直接映射为平台属性值。

## 3. 中性基础内容 base_content

### 3.1 语义素材 semantics

- `selling_points[]`、`use_scenarios[]`、`purchase_motivations[]`、`target_audience[]`
- 每条带 `status`（`confirmed` / `unconfirmed`）

### 3.2 上架文案 copywriting

字段名使用 **`copywriting`**（不要用 `copyright`，后者表示版权）。中文说明为「上架文案」。

- `title`、`bullets[]`、`description`、`keywords[]`（以及有则给的 `search_keywords[]` / `backend_keywords[]` / `ad_copy[]`）
- **单平台包（默认）**：此处为本包目标平台的**成品上品文案**，上品必填；导出时从 Agent 已产出文案组装，不另起一套空占位。
- **多平台同包（例外）**：顶层 `copywriting` 可为共享草稿或省略正文，各平台成品文案改放 `platform_overrides[].copywriting`（见 §4）。
- `generation_status`：`generated` / `not_generated`；单平台上品路径下应为 `generated` 且正文非空，否则 `content_status` 不得为 `ready`。

## 4. 平台覆盖 platform_overrides[]

每项对应一个平台 + 市场。

**单平台包（默认，与 `output/{平台}-{市场}/` 落盘一致）：**

- 只表达**非文案**差异，**不含** title / bullets / description / keywords
- 允许字段：
  - `platform` / `market` / `locale`
  - `image_set`：`expected_count`、`actual_count`、`missing_roles[]`、`status`（`已完成` / `未完成`），`images[]`：`{ role, media_id }`
  - `video_set`：`{ media_id, duration, has_voiceover, has_burned_subtitles, status }`
  - `content_notes[]`：该平台内容规则导致的限制与风险

**多平台同包（例外）：**

- 可增加 `copywriting` 对象承载该平台成品文案；合并时以覆盖为准，顶层 `copywriting` 仅作共享草稿。

覆盖片段**不含**平台类目 ID、属性码、接口字段名或 API 报文。

## 5. 媒体 media[]

实体资产去重清单，供覆盖片段用 `media_id` 引用：

- `media_id`
- `type`：`image` / `video`
- `role`：图型名，使用中文并对齐 `platforms/image-set-rules.md`（如 `尺寸规格图`、`到手内容图`、`定制操作示意图`）
- `path_or_url`：相对路径或约定 URL
- `format` / `width` / `height`
- `spec_ok`：是否满足对应平台技术规格

同一文件被多平台复用时只描述一次，由各覆盖片段引用。

## 6. 门禁与风险 gates

作为中台判断的**输入**，最终提交与否由中台裁定：

- `image_set_complete`（布尔或按平台）
- `video_complete`
- `compliance_risks[]`
- `needs_human_confirmation[]`
- `fact_grading_summary`：事实分级摘要（哪些是市场建议、不可当已确认）

## 7. 明确不提供

- 各平台 OpenAPI 原始请求体
- 平台类目 ID、属性码、接口字段名
- 平台账号、店铺 ID、授权 token
- 仅存在于聊天中的临时说明（以落盘文件为准）

## 8. 中英字段对照（供中台落库）

顶层与元数据：

- 元数据 = `metadata`
- `recommended_platforms` = 推荐平台列表（非强制推送清单）
- 商品事实 = `product_facts`
- 中性基础内容 = `base_content`
- 语义素材 = `semantics`
- 上架文案 = `copywriting`
- 平台覆盖 = `platform_overrides`
- 媒体 = `media`
- 门禁与风险 = `gates`

商品事实常用键（**必须与 CSV 表头一致**，下列仅为英文对照）：

- 颜色/款式 = `color_style`
- 尺码/规格 = `size_spec`
- 长度cm / 宽度cm / 高度cm = `length_cm` / `width_cm` / `height_cm`
- 净重g / 包装毛重g = `net_weight_g` / `gross_weight_g`
- 材质 = `material`
- 结构/表面工艺 = `structure_finish`
- 已确认功能 = `confirmed_features`
- 商品特点 = `product_features`
- 定制内容 / 定制位置 / 定制工艺 = `custom_content` / `custom_position` / `custom_process`
- 包装清单 / 包装方式 = `packing_list` / `packing_method`
- 风险/禁用声明 = `risk_or_prohibited_claims`
- 适用对象 / 使用场景 = `target_audience` / `use_scenarios`
- 资料更新时间 = `data_updated_at`

图型 role 常用值（对齐 image-set-rules）：

- 最终定制主图 = `custom_main`
- 到手内容/包装图 = `whats_included`
- 定制操作示意图 = `customization_guide`
- 尺寸规格图 = `size_spec`
- 细节放大图 = `detail`
- 日常使用场景图 = `daily_scene`
- 情绪/礼赠场景图 = `gift_scene`
- 卖点集合图 = `feature_grid`

对照表仅用于中台落库映射，不改变 Conductor 侧的中文字段与图型名。

## 9. R1 默认建议（供讨论，可推翻）

以下为 Conductor 侧建议默认值，不是最终决议。讨论时直接改本节或打回未决。

### 9.1 product_facts：按当前 CSV 表头全列输出

**建议：`product-catalog.csv` 当前全部列都输出**（现为 34 列，以文件表头为准），空值仍占一行，`status=missing`，`value` 为空字符串。

理由：

- `key` 与表头逐字一致（如 `颜色/款式`，不是「颜色」）。
- 中台能区分「缺列」与「值为空」。
- CSV 增删列时，以仓库主表为唯一真相，本契约跟随更新。

扩展属性表（`product-attributes.csv`）另开数组 `extended_attributes[]`，按 SKU 挂载，不混进主表列。

### 9.2 sales_unit：结构化对象，对齐现有归一化规则

**建议结构：**

```json
{
  "形态": "多件装",
  "销售数量": 3,
  "是否同款": true,
  "混色混款": "同款",
  "包含内容": "…",
  "数量归一化备注": "3-pack → 销售数量 3"
}
```

**形态枚举（中文，对齐 product-to-listing / image-set-rules）：**

| 形态 | 说明 |
| --- | --- |
| `单件` | 默认；未写数量时 |
| `多件装` | N 个装 / N-pack / pack of N / set of N，同款多件 |
| `对装` | 明确成对售卖 |
| `套装` | 多组件组成一套 |
| `礼盒装` / `旅行装` / `入门套装` | 有明确礼盒或套装场景时使用 |
| `混色装` / `混款装` | 必须同时写清组合是固定还是随机 |
| `待确认` | 如「N 件套」无法判断同款多件还是多组件时使用，不得当成可刊登事实 |

数量归一化规则沿用现有文档：`3个装`、`三个装`、`3-pack`、`pack of 3` 等一律 `销售数量=3`。`metadata` 里的 `skus[].qty` 与 `sales_unit.销售数量` 必须一致；不一致时 `content_status=blocked`，原因写入 `blocked_reasons`。

### 9.3 media_id：稳定、可读、与文件解耦

**建议格式：**

```text
{type}-{role_slug}-{index}
```

示例：

- `image-custom_main-01`
- `image-size_spec-01`
- `video-main-01`

规则：

- `type`：`image` | `video`
- `role_slug`：用 §8 英文对照（如 `size_spec`），避免中文进 id
- `index`：两位数字，同 role 多张时递增；单张也用 `01`
- **不把路径或文件名直接当 id**，避免改名导致引用断裂；`path_or_url` 单独字段承载实际位置
- 同一物理文件只分配一个 `media_id`；多平台覆盖通过引用复用

### 9.4 落盘位置与文件名：叠加、不改现有目录语义

**短期建议（不改 `output/{平台}-{市场}/` 人工验收路径）：**

```text
output/{平台}-{市场}/{Listing版本目录}/上架/content-package.json
```

说明：

- 与「完整生产报告」同级，放在 `上架/` 下，方便人工和中台一起发现。
- 文件名固定为 `content-package.json`；版本跟随 Listing 版本目录（`MUG0110RD` / `MUG0110RD-v2`），**不给 json 文件名再加 `-v2`**。
- 当一次生产只面向一个平台-市场时：该目录下的 content-package 里，`base_content` + **一个** `platform_overrides` 项即可；`target_platforms` 填该项。
- 当未来要在同一包里带多平台覆盖时：仍可先按「主任务平台-市场」目录落一份；多平台合并策略留到 R7，短期不强迫改目录树。

**明确不做（本阶段）：**

- 不新增平行的「平台无关根目录」
- 不要求改 Agent 日常输出路径
- 导出 content-package 是可选后续动作，默认讨论阶段只定契约，不强制每次生产自动写出

## 10. R2 预写：content_status 与门禁（草案，供讨论）

### 10.1 content_status 取值

| 值 | 含义 | 是否可发送到中台 |
| --- | --- | --- |
| `ready` | 上品所需内容已齐备（见下方 ready 门禁） | **可以发送** |
| `partial` | 有可用内容但未齐（如图集未完成、copywriting 未齐） | **不可发送**；仅留在 Conductor 本地 |
| `blocked` | 存在阻断项（事实冲突、高风险未确认等） | **不可发送**；仅留在 Conductor 本地 |

**已确认规则：只有全部准备妥当后才可发送到中台。**  
`partial` / `blocked` 包不传中台、不入库草稿；在 Conductor 侧补齐并变为 `ready` 后再导出/推送。

**首期 ready 门禁（已确认）：**

- 必备：`copywriting` 正文齐（title / bullets / description 等按平台规则）
- 必备：约定套图齐（`image_set.status=已完成`，无缺失必出图型）
- 不强制：视频（无视频不挡 ready；若本轮声称已生成视频则须通过视频门禁）
- 仍须：无 `blocked` 级事实冲突；上品必填的 `confirmed` 事实可用

中台若仍收到非 ready 包（异常路径），应拒收并返回明确错误，便于对账；正常路径下 Conductor 不得发出此类包。

### 10.2 blocked_reasons / needs_human_confirmation 示例

- `SKU数量与销售单位不一致`
- `Temu定制类套图未完成（缺尺寸规格图）`
- `风险/禁用声明中含未确认项且未在文案中规避`
- `存在 market_suggestion 被误写为 confirmed`（若导出校验发现）
- `视频缺口播或烧录字幕`
- `copywriting 未生成或正文为空`

### 10.3 责任边界

| 检查项 | Conductor | 中台 |
| --- | --- | --- |
| 事实分级、套图完成门禁、文案是否齐备 | 写入 `gates` / `content_status`；**非 ready 不发送** | 可做二次校验；收到非 ready 则拒收 |
| 平台类目/属性是否填得上、API 必填 | 不负责 | 推送平台时校验 |
| 账号、限流、媒体上传失败 | 不负责 | 负责并回写结果（回写策略见 R5） |

最终「能不能登上某平台」仍由中台在推送当下裁定；Conductor 的 `ready` 只表示**内容包已齐、允许交给中台**。

## 11. R4 预写：图片 / 视频如何传到中台（草案，供讨论）

### 11.1 原则

1. **JSON 不内嵌二进制**。`content-package.json` 只描述媒体元数据与引用；图片、视频始终是独立文件。
2. **不改变 Conductor 现有落盘**。媒体继续写在 `output/.../图片/`、`视频/`；传递协议叠加在现有目录之上。
3. **中台入库时收齐媒体副本或稳定 URL**；推送到跨境平台时，由中台再上传到平台/CDN，不要求 Conductor 懂平台媒体接口。
4. **引用完整性**：每个 `platform_overrides` 里引用的 `media_id`，必须能在本包 `media[]` 中解析到可获取的文件；解析失败则该平台覆盖视为不可推送。

### 11.2 `media[]` 建议补强字段

在现有 `path_or_url` / `role` / 规格字段之外，建议增加：

| 字段 | 说明 |
| --- | --- |
| `sha256` | 文件校验和，防传损、便于中台去重 |
| `byte_size` | 字节数 |
| `mime_type` | 如 `image/png`、`video/mp4` |
| `delivery` | 本条媒体的交付方式，见下 |
| `local_relative_path` | 相对 Listing 版本目录的路径，如 `图片/01-主图.png`（Conductor 侧真相） |
| `remote_url` | 若已上传对象存储，填可拉取 URL；未上传则为空 |

`path_or_url` 可逐步拆成 `local_relative_path` + `remote_url`，避免一个字段语义混杂。

### 11.3 交付方式（分阶段）

| 阶段 | 方式 | 适用 | Conductor 做什么 | 中台做什么 |
| --- | --- | --- | --- | --- |
| **P0 试点** | **打包交付**：`content-package.zip` = `content-package.json` + 被引用的图片/视频（保持相对目录） | 人工导出、联调、无对象存储时 | 按 `media[]` 收集已引用文件打成 zip；路径与 `local_relative_path` 一致 | 解压入库，按 `media_id` 建素材库 |
| **P1** | **共享目录拉取**：中台有权读 Conductor 机器或挂载盘上的 `output/` | 同机房/内网部署 | 只写 content-package；`delivery=shared_fs` | 按路径拉取文件并复制到中台存储 |
| **P2 推荐稳态** | **先上传对象存储，再交 JSON**：媒体进 OSS/S3，`remote_url` 填入（可带过期签名） | 正式对接 | 导出前上传媒体（或由旁路上传任务完成），写 `remote_url` + `sha256`；JSON 仍可不改 Agent 文案 | 拉取 URL 入库；推平台时再转存平台媒体接口 |

**默认推进顺序：P0 → P2**；P1 仅当双方明确同网共享磁盘时采用，不作为通用方案。

### 11.4 P0 打包约定（建议）

```text
content-package.zip
├── content-package.json
├── 图片/
│   ├── 01-主图.png
│   └── ...
└── 视频/          （若有被引用的视频）
    └── 01-主视频.mp4
```

规则：

- zip 内只包含 **被本包 `media[]` 引用** 的文件，不整包塞入未引用的草稿或验收报告。
- 相对路径必须与 `local_relative_path` 一致，便于校验。
- 完整生产报告、竞品、供应链等 Markdown **默认不进 zip**（除非中台明确要归档）；它们继续留在 `output/` 给人看。
- 视频体积大：P0 可允许「JSON + 图片进 zip，视频单独传或只给 URL」；若视频未随包，对应 `media` 项 `delivery=pending`，`content_status` 至少为 `partial`。

### 11.5 责任边界

| 事项 | Conductor | 中台 |
| --- | --- | --- |
| 生成并本地落盘图片/视频 | 负责 | 不负责 |
| 计算 `sha256` / 写入 `media[]` 元数据 | 负责（导出时） | 入库时复验 |
| 上传到对象存储（P2） | 可由导出任务或旁路上传完成 | 也可约定由中台从共享盘拉完再自己存——但推荐 Conductor 侧先给出稳定 URL |
| 上传到 Temu/Amazon 媒体接口 | 不负责 | 推送时负责 |
| 媒体失败重试、平台媒体 ID 回写 | 不负责 | 负责（回写见 R5） |

### 11.6 与 R3（内容包传输）的关系

- **内容包元数据**（JSON）与 **媒体二进制** 可以同一通道（zip），也可以拆开（JSON API + 媒体 URL）。
- 建议契约上把二者解耦：中台可以先收 JSON 建草稿，媒体未齐时状态为 `partial`；媒体齐套后再允许推送。
- R3 定「谁推谁拉」时，必须同时回答：媒体走同通道还是 URL 通道。

### 11.7 明确不做

- 不在聊天或 Markdown 里塞 base64 大图
- 不要求中台直接读 Cursor 会话或临时生成缓存
- 不让 Conductor 调用各平台的图片上传 API

---

## 相关文档

- [对接总计划](./对接总计划.md)
- [轮次默认建议（R3–R8）](./轮次默认建议.md)：传输 API、回写、试点、多平台合并、工程清单、错误码与安全
