# content-package 结构草案（v0.1）

> 状态：讨论中（R1）  
> 归属：[对接总计划](./对接总计划.md) 第 R1 轮产出  
> 定位：Conductor → 中台的机对机内容契约草案，不代表已可刊登。

## 设计约束（已对齐）

- content-package 是**新增**契约，导出时**不得改变 Agent 已产出的内容**，也**不得触发新的内容生产**；只有用户明确指定时才生成中性基础文案，否则相应字段标记为「未生成」。
- 字段结构**不绑定具体平台**，平台差异一律放进 `platform_overrides`。
- 字段键默认使用**中文**，与 `data/product-catalog.csv` 列名和 `platforms/image-set-rules.md` 图型名一一对应，减少映射层；另附英文对照表供中台落库。
- 中性内容与平台覆盖都由 Conductor 产出；接口字段（类目 ID、属性码、API 报文）由中台在推送时组装。

## 顶层结构

```text
content-package
├── 信封            envelope
├── 商品事实        product_facts   （中性，逐条带状态）
├── 中性基础内容    base_content
│   ├── 语义素材    semantics
│   └── 基础文案    copy            （默认未生成）
├── 平台覆盖[]      platform_overrides
├── 媒体[]          media           （实体文件，按 id 引用）
└── 门禁与风险      gates
```

## 1. 信封 envelope

- `schema_version`：契约版本
- `package_id`：内容包标识
- `generated_at`：生成时间
- `spu`：SPU
- `skus[]`：`{ sku, qty }`，沿用现有 SKU 与数量规则
- `listing_id`：单 SKU 或组合标识（沿用 `workflows/output-structure.md`）
- `sales_unit`：销售单位与组合形态（单件 / 多件 / 套装 / 混色 / 混款）
- `product_name`：商品名称（业务文本）
- `suggested_category`：建议类目（业务文本，**非**平台类目 ID）
- `target_platforms[]`：本包已含覆盖的平台与市场，可为空
- `content_status`：内容层自评（`ready` / `partial` / `blocked`），**不是**刊登判定
- `blocked_reasons[]`：内容层受阻原因

## 2. 商品事实 product_facts

逐条记录，键对齐 `data/product-catalog.csv` 的列：

```json
{ "key": "净重g", "value": "400", "unit": "g", "status": "confirmed", "source": "catalog" }
```

- `status`：`confirmed` / `unconfirmed` / `missing`
- `source`：`catalog` / `user` / `supplier` / `image` / `market_suggestion`
- 建议分组：外观、尺寸重量、材质工艺、已确认功能、定制（内容/位置/工艺）、包装（清单/方式/尺寸/毛重）、认证（信息/状态）、使用护理、禁止或未确认声明、适用对象与场景
- `market_suggestion` 永不等于已确认事实；中台不得直接映射为平台属性值。

## 3. 中性基础内容 base_content

### 3.1 语义素材 semantics

- `selling_points[]`、`use_scenarios[]`、`purchase_motivations[]`、`target_audience[]`
- 每条带 `status`（`confirmed` / `unconfirmed`）

### 3.2 基础文案 copy

- `title`、`description`、`keywords[]`
- `generation_status`：`generated` / `not_generated`
- **默认 `not_generated`**：导出动作不触发生成，也不改写 Agent 已产出内容；仅当用户明确要求「生成中性基础文案」时才填充。

## 4. 平台覆盖 platform_overrides[]

每项对应一个平台 + 市场，只表达**内容差异**：

- `platform` / `market` / `locale`
- `title`
- `bullets[]`（要点结构随平台变化，例如 Amazon 5 条、Temu 规格卖点）
- `long_description`
- `search_keywords[]` / `backend_keywords[]` / `ad_copy[]`（有则给，无则省略，不编造）
- `image_set`：
  - `expected_count`、`actual_count`、`missing_roles[]`、`status`（`已完成` / `未完成`）
  - `images[]`：`{ role, media_id }`
- `video_set`：`{ media_id, duration, has_voiceover, has_burned_subtitles, status }`
- `content_notes[]`：该平台内容规则导致的限制与风险

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

顶层与信封：

- 信封 = `envelope`
- 商品事实 = `product_facts`
- 中性基础内容 = `base_content`
- 语义素材 = `semantics`
- 基础文案 = `copy`
- 平台覆盖 = `platform_overrides`
- 媒体 = `media`
- 门禁与风险 = `gates`

商品事实常用键（对齐 CSV 列）：

- 颜色/款式 = `color_variant`
- 尺码/规格 = `size_spec`
- 长度cm / 宽度cm / 高度cm = `length_cm` / `width_cm` / `height_cm`
- 净重g / 包装毛重g = `net_weight_g` / `gross_weight_g`
- 材质 = `material`
- 结构/表面工艺 = `structure_finish`
- 已确认功能 = `confirmed_features`
- 定制内容 / 定制位置 / 定制工艺 = `custom_content` / `custom_position` / `custom_process`
- 包装清单 / 包装方式 = `packing_list` / `packing_method`
- 认证信息 / 认证状态 = `certification` / `certification_status`
- 使用/护理说明 = `care_instructions`
- 禁止/未确认声明 = `prohibited_or_unconfirmed`
- 适用对象 / 使用场景 = `target_audience` / `use_scenarios`

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

### 9.1 product_facts：全列输出 + 空值显式标记

**建议：主表 37 列全部输出**，空值仍占一行，`status=missing`，`value` 为空字符串。

理由：

- 中台一眼能看出「缺什么」，不会把「字段未传」和「字段值为空」搞混。
- 与 CSV 列一一对应，导出逻辑简单，不需要再维护「关键列白名单」。
- 扩展属性表（`product-attributes.csv`）另开数组 `extended_attributes[]`，按 SKU 挂载，不混进主表 37 列。

可选后续优化（不进 v0.1）：中台若嫌包大，可在接收端过滤 `missing`，Conductor 侧仍全量输出。

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

数量归一化规则沿用现有文档：`3个装`、`三个装`、`3-pack`、`pack of 3` 等一律 `销售数量=3`。信封里的 `skus[].qty` 与 `sales_unit.销售数量` 必须一致；不一致时 `content_status=blocked`，原因写入 `blocked_reasons`。

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

| 值 | 含义 | 中台建议行为 |
| --- | --- | --- |
| `ready` | 内容层自评可提交给中台做推送准备 | 可进入「待推送」队列 |
| `partial` | 有可用内容，但缺非阻断项（如基础文案未生成、某平台覆盖未做） | **保存并锁定推送**，允许补齐后再推 |
| `blocked` | 存在阻断项（事实冲突、必备图未完成且本轮声称已完成、高风险未确认） | **保存并锁定推送**，必须人工或 Conductor 更新后再推 |

**建议：中台一律接收并保存**，用状态锁定推送，而不是拒收。拒收会导致 Conductor 与中台对账困难。

### 10.2 blocked_reasons / needs_human_confirmation 示例

- `SKU数量与销售单位不一致`
- `Temu定制类套图未完成（缺尺寸规格图）`
- `认证信息未确认`
- `存在 market_suggestion 被误写为 confirmed`（若导出校验发现）
- `视频缺口播或烧录字幕`

### 10.3 责任边界

| 检查项 | Conductor | 中台 |
| --- | --- | --- |
| 事实分级、套图完成门禁、内容缺失 | 写入 `gates` / `content_status` | 读取后决定是否允许推送 |
| 平台类目/属性是否填得上、API 必填 | 不负责 | 推送时校验 |
| 账号、限流、媒体上传失败 | 不负责 | 负责并回写结果（回写策略见 R5） |

最终「能不能登上某平台」永远是中台在推送当下裁定；Conductor 的 `content_status` 只是内容层自评输入。

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
