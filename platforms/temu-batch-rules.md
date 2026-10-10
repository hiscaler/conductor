# Temu 批量共享图案规则

## 适用范围与加载入口

本文件是 Temu 批量 `one_to_all_skus` 的目录、图片命名、数量及配置规则的唯一维护位置。目标平台为 Temu、运行模式为批量且选择“一张图案应用到全部 SKU 变体”时，创建批次或生产前必须全量读取本文件。

单个模式、批量 `one_to_one` 和其他平台沿用原规则。这里定义项目交付组织方式；不代表已经核验 Temu 后台当前类目要求，后台额外限制仍须实际核验。

通用入口、输入校验、任务矩阵、claim/record/retry、锁、并行职责及阻塞处理见 [批量生产流程](../workflows/automatic-batch.md)。图片技术与视觉要求仍按 [平台套图规则](./image-set-rules.md) 和 [图像技术规格](./image-technical-specs.md) 执行。

## 图案目录与共用文案

Temu 批量 `one_to_all_skus` 按图案归档：一张有效去重图案对应一个 `output/Temu-US/{批次号}-{图案名称}[-vN]/` 目录，同组全部 SKU 共用它，不再为每个 SKU 创建目录。图案名称取源文件名去掉扩展名，保留原编号和中文；非法字符替换为短横线、移除末尾点与空格，Windows 保留名加“图案-”前缀，超长名称截取前 100 个字符。同名（含大小写、截断及清理后碰撞）追加稳定的 pattern_group_id 区分，最终再次检查唯一性，仍碰撞则加稳定数字序号；源名称末尾类似-v2时加“图案”尾缀避免被误识别为生产版本，不覆盖另一组。其他平台的批次目录仍由通用批量流程规定。

```text
output/Temu-US/{批次号}-图案名称/
  商品资料.md
  文案/
    文案资产.md
    关键词清单.md
  图片/
    01-SKU123-主图.jpg
    01-SKU456-主图.jpg
    02-日常使用场景图.jpg
    03-礼赠场景图.jpg
    04-卖点集合图.jpg
    05-尺寸规格图.jpg
    06-SKU123-定制示意图.jpg
    06-SKU456-定制示意图.jpg
  上架/
    完整生产报告.md
```

以上展示主要交付文件，实际还须按原规则保存批次关联、套图配置、脚本、逐图侧车、验收报告和 next-action.json；未请求文案或视频时不创建对应目录。允许 PNG/JPG，图片名称的序号和角色保持一致。商品资料、共用文案与完整报告每组各一份，商品资料逐 SKU 列出颜色、属性、图库和任务 ID；买家在变体间择一，不能写成一次收到多件。

## 图片数量与命名

仅适用 Temu 批量 `one_to_all_skus`；单个、批量 `one_to_one` 和其他平台保持原规则。

- 每个 SKU 一张 `01-{标准SKU}-主图`，逐 SKU 使用真实商品图库，不借用其他颜色。
- 每组只生成一套 `02-日常使用场景图`、`03-礼赠场景图`、`04-卖点集合图`、`05-尺寸规格图`。场景图和卖点图可组合展示本组多个真实变体，须明确属于可选款式、不暗示组合售卖；尺寸图只展示一个商品，所有 SKU 规格相同时共用，规格不同时暂停该图并说明差异。
- 启用卖家定制服务（必选/可选）时，强制为每个 SKU 追加 `06-{标准SKU}-定制示意图`，不计入默认五张。服务为不提供时禁止生成该图。
- 每 SKU 的五张轮播图为自己的主图 + 02/03/04/05 四张共用副图；定制示意图单列。每组实际文件数为 SKU 数 + 4 + SKU 数（启用定制时）。两张图案、七个 SKU 应有两个目录、14 张主图、8 张共用副图和14 张定制示意图，共36张，不能按14个SKU任务重复生成共用副图。
- 明确选扩展九张时，保留上述固定名称，扩展副图使用 07/08/09/10（角色按真实事实设计），登记 supplemental_images；06 专供 SKU 定制示意图。每 SKU 九张轮播图 + 一张定制示意图，每组文件数为 SKU 数 + 8 + SKU 数。组目录的去重文件数不是单个 SKU 的轮播图数，平台上限按每个 SKU 使用图集核对，后台额外规则仍须核验。

## 配置、侧车与验收关联

Temu 新共享批次使用 output_layout: pattern_group，套图配置只保存一份：shared_images 登记 02/03/04/05，supplemental_images 登记扩展副图，variants 按标准 SKU 排序登记 sku、job_id、main_image 和 additional_required（其中包含该 SKU 的 06 定制示意图）。selected_count 记录每 SKU 轮播张数 5/9；customization_enabled 明确服务状态，customization_master_id 使用组母版。

组级 `批次关联.json` 保存 batch_id、spu、mapping_mode、output_layout、pattern_group_id、job_ids、variant_skus、来源图案路径与哈希、使用角色、共享母版、版本以及报告/状态相对路径，不登记单个 sku/job_id。共用配置和副图侧车关联全部 variant_skus；专属图片侧车记录该 sku/job_id；所有图片关联同一 pattern_group_id、来源哈希和组母版。专属侧车只能对应自己的 SKU，不能用整组声明代替逐图验收。校验按实际去重文件清单及各 SKU 使用图集检查，不将共用副图重复计数。

组套图配置示例（关联值必须从当前 claim 返回值填写，不能把示例值当作真实记录）：

```json
{
  "platform": "Temu",
  "mapping_mode": "one_to_all_skus",
  "output_layout": "pattern_group",
  "batch_id": "当前批次号",
  "pattern_group_id": "pattern-001",
  "source_pattern_sha256": "当前分配图案的真实SHA-256",
  "source_pattern_usage": "fixed_product_artwork",
  "batch_report": "../../../当前批次号-批次总报告.md",
  "customization_master_id": "当前图案组的母版ID",
  "customization_enabled": true,
  "preset": "standard_5",
  "selected_count": 5,
  "selection_source": "default",
  "variant_skus": ["SKU123", "SKU456"],
  "shared_images": ["02-日常使用场景图.jpg", "03-礼赠场景图.jpg", "04-卖点集合图.jpg", "05-尺寸规格图.jpg"],
  "supplemental_images": [],
  "variants": [
    {"sku": "SKU123", "job_id": "job-001", "main_image": "01-SKU123-主图.jpg", "additional_required": ["06-SKU123-定制示意图.jpg"]},
    {"sku": "SKU456", "job_id": "job-002", "main_image": "01-SKU456-主图.jpg", "additional_required": ["06-SKU456-定制示意图.jpg"]}
  ]
}
```

配置位于 图片/，报告链接须由 Agent 根据实际路径和批次号计算，不能沿用占位值。

## 目录预留与历史批次

Temu 新共享批次记录 `output_layout: pattern_group`。图案组首次 claim 只预留一次目录，将 output_dir 和版本同步给组内全部任务；同组后续领取、续做和普通修改复用该目录，不重新调用版本分配器。每个 SKU 仍保留自己的 job_id、lease_token 和任务状态。

状态中的 pattern_groups 固定图案名称、directory_name、output_dir、版本、全部 SKU、图案来源和任务 ID，恢复时检查完整“图案 × SKU”矩阵。

共用文案、配置和副图由组负责人维护，专属图片按 SKU 分别生产及验收；全部组资产和统一交付闸门通过后才记录完成，通用调度细节见批量流程。

没有 output_layout 字段的历史批次保留原目录、文件命名和验收方式，断点续做不自动搬迁；迁移已有产物需单独授权和核对。
