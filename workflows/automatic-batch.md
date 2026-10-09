# 批量生产流程

本流程由核心 Agent 在用户选择“批量”后执行，复用现有内容生产步骤。批次脚本负责输入校验、固定分配、目录预留、状态和关联检查；研究、创意、文案、AI 图片和视频脚本由驱动本项目的 Agent 实际执行。创建批次只表示计划已建立，不代表产物已生成。

## 入口与任务范围

hello 等入口先展示 `1. 单个 / 2. 批量`。单个模式进入现有 13 项菜单。批量模式只显示：

```text
1. 生成标题、描述、关键词
2. 生成文案 + AI 商品图
3. 生成文案 + AI 商品图 + 视频脚本
4. 只生成 AI 商品图
```

批量菜单 1/2/3/4 映射单个菜单 1/2/3/8。批量菜单第 4 项不能解释为找品，第三项只交付视频脚本，不自动生成视频成片。下一步动作及 `上架/next-action.json` 一律使用映射后的单个菜单编号。模式菜单、单个任务菜单、批量任务菜单、服务菜单、类型菜单、创意菜单、下一步菜单分别保存上下文；数字只按最近一次显示的菜单解释。明确的自然语言单项任务沿用现有处理，不强制重新选择模式。

## 最低输入与整批选择

仅支持 CSV 中已有的 SPU。提供带示例值的最小模板：

```text
SPU：MUG
示例图案目录：D:/designs/mug-patterns
目标平台：Temu
目标市场：美国
```

目录可使用 Windows、macOS 或 Linux 的真实路径。示例目录中的图片是卖家指定用于本批次的图案/创意素材；商品外观来源是 `data/products/{SPU}/{SKU}/`。不要要求文件名包含 SKU，也不要把示例图案当作商品白底图或买家已提交的订单素材。批次开始前须确定定制类型：`图文`=图片和文字都可定制，`图`=文字固定、图片定制，`文`=图片固定、文字定制。另须确定每张分配图案的角色：买家需上传的定制图片 (`customization_image`)、商品固定印花 (`fixed_product_artwork`) 或仅供创意参考 (`creative_reference_only`)。若只给了图片目录而没有说明图案角色，且该角色会改变买家下单字段或商品印花内容，先只询问这个角色，不得由定制类型或文件夹本身推断。

平台、市场、任务范围、卖家定制服务、适用定制类型整批确认一次，已明确项直接复用。定制能力不代替 Listing 服务，缺少服务时沿用原数字门禁，缺少类型时再询问类型；不能把用户未回复当作选择。图片任务一次确定预设：默认标准 5 张，不等待默认回复；用户明确选择扩展 9 张才使用 extended_9。跨 SKU 定制能力冲突单独阻塞对应任务，不对所有 SKU 重复提问。

选择批量模式并提交本批次输入，授权逐 Listing 自动选择创意、整理商品档案及补齐缺少的通用展示文字。`图文`或`图`类型且角色为 `customization_image` 时，分配图案是买家图片字段的展示样例；`图文`缺少文字时可按 Listing 生成原创通用示例。`文`类型且角色为 `fixed_product_artwork` 时，分配图案作为 SKU 固定印花，买家只提交文字；`文`类型且角色为 `creative_reference_only` 时，图案只作创意参考，不印到商品上。不提供定制时图案默认只作创意参考。不同角色不得混记；不改变已提供的固定文字或固定图案，不把 AI 展示样例记为买家素材；指定照片、肖像、商标/IP 或授权不明素材仍按已有规则检查，不能适用的图案阻塞原配对，不自动换给另一个 SKU。

## 校验、分配与数量差异

1. 严格校验两个 CSV 的 UTF-8、表头、列数、引号、必填值、SKU 唯一性、数值、日期、枚举、属性外键及本地证据路径；错误包含文件、行、字段、原值和修正方式。SKU 的商品图库单独校验 main、角色冲突；图片任务缺图仅阻塞对应 SKU，文案任务可继续并记录未加载主图。
2. 标准 SKU 去重排序，递归读取示例目录里的 PNG/JPG/JPEG/WebP。隐藏文件、符号链接、非图像、损坏图像和解码超限图像列入排除记录。按解码后的像素去重，避免不同文件名/编码的同一图案重复分配。
3. 批次必须明确 `mapping_mode`。用户选择菜单统一为：

   ```text
   1. 一张图案分配给一个 SKU 变体（one_to_one）
   2. 一张图案应用到全部 SKU 变体（one_to_all_skus）
   ```

   “一张图案”表示对图案目录中的每张有效去重图案分别执行该规则。`one_to_one` 为默认值，每张图案只对应一个 SKU 变体，按稳定排序配对，任务数取图案数和 SKU 数量较小的一方；`one_to_all_skus` 将每张图案应用到该 SPU 的全部 SKU 变体，任务数为图案数 × SKU 数，不按数量截断。已确认分配方式直接复用。
4. `one_to_one` 的多余图案和未配对 SKU 排除在本批次生产范围外，报告逐项列出；`one_to_all_skus` 不产生这类数量截断，但必须记录每个图案组的全部变体 SKU 和展开任务数。两种模式都保留原 CSV 和原图片，不静默重新配对。
5. `one_to_one` 每个 `SKU + 图案` 独立 Listing；`one_to_all_skus` 视为一个可选变体 Listing 的共享图案组，每张图案都为全部 SKU 生成相同图案、文案和定制逻辑，只替换 SKU 商品外观/颜色。原始商品包装和固定包含内容仍来自 CSV，不能为了批量处理擅改事实。
6. 将配对、商品事实快照、扩展属性、原始路径与 SHA-256 保存到批次状态。图案分配固定后不再重新洗牌、增删或补配；原有输入内容改变时暂停确认并建立新批次。新增缺失的商品 main/角度图可以补入原任务，但不能替换已登记输入。

## 创建与执行命令

Agent 将已确认值写入工作区临时 JSON，例如 `.tmp/自动批次输入.json`，不让用户编辑机器字段或运行命令：

```json
{
  "spu": "MUG",
  "sample_directory": "D:/designs/mug-patterns",
  "platform": "Temu",
  "market": "US",
  "auto_task": 2,
  "seller_service": "必选",
  "customization_type": "文",
  "sample_usage": "fixed_product_artwork",
  "mapping_mode": "one_to_one"
}
```

`sample_usage` 需要根据用户对图案角色的确认填写三个枚举之一；未确认则先澄清，不创建批次。Temu 用户选扩展图时添加 `image_preset: extended_9`；其他平台的图数由 Agent 按对应平台规则确定并写入 `image_count`，预设为 `platform_default`，不套用 Temu 的 5/9 张。提供文字时添加 `sample_text`。图片默认值无需用户再选择。脚本支持枚举中文值或 required/optional/not_offered、image_text/image/text；市场统一为两位代码。只创建一次批次：

```text
node scripts/batch-production.mjs create .tmp/自动批次输入.json
node scripts/batch-production.mjs status output/{批次号}-批次状态.json
node scripts/batch-production.mjs claim output/{批次号}-批次状态.json
```

claim 返回单个任务：商品资料、属性、对应图库、固定图案、Listing 目录和 lease_token。先检查实际图案和商品图；没有阻塞时依次执行：

`one_to_all_skus` 还返回 `pattern_group`，包含同组全部 SKU、任务 ID 和共享 `customization_master_id`。同组保持图案、展示文字、创意和文案一致，只替换各 SKU 的真实商品外观与颜色；复用图案组的设计内容，不逐 SKU 随机生成新内容。各 SKU 分别生产并验收，整组对应同一变体商品，不表示买家一次收到全部 SKU。

```text
商品事实与素材核对
→ 复用已有研究步骤（必须真实外部检索，访问限制如实记录）
→ 按本 SKU + 图案形成创意候选、评分并自动选优
→ 整理完整商品资料并核对事实
→ 实际生成本轮文案 / 独立套图 / 视频脚本
→ 单项与统一校验
→ 标记该 Listing 完成，再领取下一任务
```

批量模式不展示创意选择菜单，不等待复制商品示例回传；分别记录“系统自动选定”和“自动整理并核对、未经用户逐项回传”，不能伪造用户确认。这是批量模式对单个模式交互步骤的明确例外，研究证据、事实、类型、视觉验收和交付质量要求继续适用。可以复用同 SPU 的有效市场证据，必须按每个 SKU 的差异与分配图案分别组织创意。

每个阶段把实际结果写入 Listing 文档，同时用步骤 JSON 更新状态：

```json
{"stage":"research","status":"running","message":"已完成实际检索，证据和访问限制已写入商品资料"}
```

```text
node scripts/batch-production.mjs record output/{批次号}-批次状态.json {job_id} {lease_token} .tmp/步骤.json
```

stage 可用 input/research/creative/product_profile/production/validation；阻塞或失败必须提供 blockers 数组。完成请求使用 `status: completed`，脚本会实际检查必需文件、批次关联、图片数量和母版，并运行原统一交付闸门；失败写回失败记录，不能声称完成。没有可用 AI 生成工具时记录未完成，不能用脚本、占位图或文件路径冒充图片。

## 批次目录、双向引用和恢复

批次号使用上海时区分钟时间戳 `批202610071122`（可指定时区）；以原子创建的状态文件预留，冲突加 `-02`、`-03`。跨平台只依赖 Node 文件操作和现有版本分配器：

```text
output/{批次号}-批次总报告.md
output/{批次号}-批次状态.json
output/Temu-US/{批次号}-MUG0110PK/批次关联.json
output/Temu-US/{批次号}-MUG0110PK/商品资料.md
output/Temu-US/{批次号}-MUG0110PK/上架/完整生产报告.md
```

以上为 `one_to_one` 的目录。`one_to_all_skus` 使用 `output/{平台}-{市场}/{批次号}-{pattern_group_id}-{SKU}[-vN]/`，例如 `{批次号}-pattern-001-MUG0110PK`，避免同一 SKU 的不同图案互相覆盖。状态中的 `pattern_groups` 固定每组全部 SKU、图案来源和任务 ID，恢复时检查完整“图案 × SKU”矩阵。

claim 可省略选择器领取下一个就绪任务；指定任务时使用 `job_id`。record/retry 也使用 `job_id`。只有一个任务对应某 SKU 时才允许用 SKU 代替；同一 SKU 对应多个图案时会明确拒绝并列出可选任务 ID。

一次 claim 只为该 Listing 预留一次目录，其余任务暂不创建目录。后续步骤和普通修改复用已保存 output_dir；明确要求新版本才使用现有 -vN 分配规则。批次状态采用排他文件锁与原子替换；重复领取执行中的 SKU 被拒绝，原执行者用 lease_token 继续，阻塞或失败的任务补齐资料后 `retry` 再领取；配对和输出目录保持不变。进程中断在目录预留步骤且目录还未记录时必须人工核对，禁止猜测目录或静默再分配；遗留文件锁必须确认没有执行者后才能清理。

每个 Listing 的所有卖家侧 Markdown 文档开头包含批次号和正确的相对报告链接，放在最终买家文案值之外。图片的套图配置与每张 `.verify.json` 必须包含 batch_id、sku、source_pattern_sha256、source_pattern_usage、batch_report；启用定制还要以同一个 customization_master_id 关联该 Listing 的所有图片。next-action.json 包含 batch_id、batch_report 和映射后的菜单。引用相对于该记录所在目录，不使用机器绝对报告路径，以便整体移动 output 后文档链接仍有效（执行恢复须使用原始状态路径）。批次脚本 `check-listing` 会检查这些关联，不替代人工打开最终图片核对图案和 SKU。

`one_to_all_skus` 的 `批次关联.json` 额外保存 mapping_mode、pattern_group_id、variant_skus、共享 customization_master_id。套图配置及逐图侧车额外保存 job_id、pattern_group_id，并使用 claim 返回的图案组共享 customization_master_id；不同图案组不得共用母版。每个 SKU 子任务仍使用自己的图库和套图配置，最终组装变体 Listing 时按图案组汇集各 SKU 产物。

```text
node scripts/batch-production.mjs retry output/{批次号}-批次状态.json {job_id}
node scripts/batch-production.mjs check-listing {claim返回的output_dir}
```

批次总报告随每次状态更新自动刷新，显示任务范围、整批选择、原始/有效/配对数量、排除原因、逐 SKU 状态、实际已有产物链接、执行记录及下一步。某个 SKU 阻塞时继续其他已就绪 SKU，最后集中询问实际缺失的信息及受影响项。最终聊天按聊天模板给出数量汇总、未完成原因和可点击总报告链接，不逐个展开完整报告。不得因为有其他 SKU 完成就把整批说成完成。
