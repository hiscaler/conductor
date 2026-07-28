import assert from "node:assert/strict";
import { readFile, readdir, stat } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("..", import.meta.url));

const read = (path) => readFile(new URL(`../${path}`, import.meta.url), "utf8");
const files = Object.fromEntries(
  await Promise.all(
    [
      "AGENTS.md",
      "README.md",
      "agents/cross-border-commerce-agent.md",
      "agents/subagents.md",
      "agents/visual-production-agent.md",
      "workflows/product-to-listing.md",
      "workflows/start-guide.md",
      "workflows/creative-direction-selection.md",
      "workflows/action-menu.md",
      "workflows/output-structure.md",
      "platforms/market-data-sources.md",
      "platforms/platform-profiles.md",
      "platforms/commerce-semantic-creative-rules.md",
      "platforms/image-set-rules.md",
      "platforms/image-technical-specs.md",
      "templates/chat-output.md",
      "templates/product-input.md",
      "templates/production-output.md",
      "docs/agent-testing.md",
      "agents/cross-border-commerce-agent.json",
      "package.json",
    ].map(async (path) => [path, await read(path)]),
  ),
);

const agent = JSON.parse(files["agents/cross-border-commerce-agent.json"]);
assert.deepEqual(agent.inputs.required, [], "核心 Agent 不应预设一组统一必填字段");
assert.deepEqual(agent.product_catalog, {
  products: "../data/product-catalog.csv",
  attributes: "../data/product-attributes.csv",
});
assert.equal(agent.creative_direction_selection, "../workflows/creative-direction-selection.md");
assert.equal(agent.commerce_semantic_creative_rules, "../platforms/commerce-semantic-creative-rules.md");
assert.ok(agent.workflow.includes("creative_direction_selection"), "核心 Agent 工作流缺少创意方向选择阶段");
assert.equal(agent.output_contract.chat_template_file, "../templates/chat-output.md");
assert.equal(agent.output_contract.chat_default, "relevant_only");
assert.equal(
  agent.output_contract.full_report_output_path,
  "output/{平台}-{市场}/{Listing版本目录}/上架/完整生产报告.md",
);

const start = files["workflows/start-guide.md"];
const readme = files["README.md"];
assert.match(readme, /!\[Conductor 多场景协作：指挥家按任务类型分流\]\(\.\/assets\/conductor-workflow\.png\)/);
assert.doesNotMatch(readme, /想法 → 找品 → 平台 → 内容 → 验收 → 复盘/);
assert.doesNotMatch(readme, /## 协作路径/);
assert.match(readme, /## 你可以用它做什么/);
assert.match(readme, /## Conductor 如何工作/);
assert.match(readme, /## 不同任务如何执行/);
for (let menu = 1; menu <= 13; menu += 1) {
  assert.match(readme, new RegExp(`\\| ${menu} \\|`), `README 缺少菜单 ${menu} 的协作路径`);
}
assert.match(readme, /组合选择会合并对应路径、复用已有资料并去除重复步骤/);
assert.match(readme, /## 如何开始/);
assert.match(readme, /所有任务共用以下入口/);
assert.match(readme, /## 内容生产任务示例/);
assert.match(readme, /以下仅说明已有商品生成文案、图片或视频的交互/);
assert.doesNotMatch(readme, /## 使用流程/);
assert.match(readme, /创意菜单输入一个方向编号；输入 `0` 由系统自动选择/);
assert.match(readme, /`0` 只在最近一次显示的是创意菜单时表示自动选择/);
assert.match(readme, /供应链\//);
assert.match(readme, /利润\//);
assert.match(readme, /聊天结果.*默认只显示本轮相关内容/);
assert.match(readme, /上架\/完整生产报告\.md/);
assert.doesNotMatch(readme, /CSV 可以直接通过 Git 查看逐行变化/);
assert.doesNotMatch(readme, /商品主表不填写图片路径/);
assert.doesNotMatch(readme, /Agent 每次读取时自动检查表头/);
const readmeCanDoIndex = readme.indexOf("## 你可以用它做什么");
const readmeHowIndex = readme.indexOf("## Conductor 如何工作");
const readmeImageIndex = readme.indexOf("![Conductor 多场景协作");
const readmeStartIndex = readme.indexOf("## 如何开始");
const readmePathsIndex = readme.indexOf("## 不同任务如何执行");
assert.ok(
  readmeCanDoIndex < readmeHowIndex &&
    readmeHowIndex < readmeImageIndex &&
    readmeImageIndex < readmeStartIndex &&
    readmeStartIndex < readmePathsIndex,
  "README 阅读顺序应为：能做什么 → 如何工作 → 总览图 → 如何开始 → 任务路径",
);
for (const trigger of ["你好", "hello", "开始", "菜单", "帮助"]) {
  assert.match(start, new RegExp(`- ${trigger}`), `缺少启动词：${trigger}`);
}
assert.match(start, /英文字母不区分大小写/);
assert.match(start, /多个 SKU 始终组成一个 listing/);
assert.match(start, /未写数量默认 1/);
assert.match(start, /相同 SKU 自动合并数量/);
assert.doesNotMatch(start, /分别生成.{0,12}组合销售/);
assert.match(start, /创意方向选择/);
assert.match(start, /输入 `0` 表示系统自动选择/);
assert.match(start, /创意确定后，我会生成与该方向一致的完整商品资料示例/);
assert.match(start, /示例确认后，再按所选任务正式生成文案、图片或视频/);
assert.match(start, /不得笼统写成与当前任务无关的全部产物/);
assert.match(start, /随后建立内部事实底稿并研究目标市场、评论和热度，再生成具体创意方向/);
assert.match(start, /创意确定后，我会生成与该方向一致的完整商品资料示例/);
assert.match(start, /存在多个可靠方向时，会提供 3-5 个方向供你输入数字选择/);
assert.match(start, /非内容任务不显示创意菜单/);
assert.doesNotMatch(start, /商品示例回传后先生成具体创意方向/);
assert.match(start, /回复“重新选择创意”重新查看上一轮候选/);
assert.match(start, /回复“重新生成创意”废弃上一批候选并重新研究/);
assert.match(start, /数据复盘必须绑定实际线上 Listing/);
assert.match(start, /商品名称、SKU 或以前生成的商品资料不能替代实际商品详情链接/);
assert.match(start, /缺少有效商品详情链接时，只提供可复制的补充资料模板并暂停/);
assert.match(start, /商品链接：请替换为实际商品详情页链接（必填）/);
assert.match(start, /复盘时间范围：2026-07-01 至 2026-07-24/);

const rules = files["AGENTS.md"];
assert.match(rules, /去除整个输入及每个组合项首尾空格/);
assert.match(rules, /英文字母不区分大小写/);
assert.match(rules, /两个 SKU 去除首尾空格并忽略大小写后相同时/);
assert.match(rules, /用户未回传商品资料示例前，不进入正式/);
assert.match(rules, /workflows\/creative-direction-selection\.md/);
assert.match(rules, /内容生产任务在研究完成后先按 .*形成创意策略/);
assert.match(rules, /用户未回传商品资料示例前，不进入正式.*回传后.*先判断修改是否使已选创意失效/);
assert.match(rules, /`重新选择创意` 表示原样重新展示上一轮候选菜单/);
assert.match(rules, /`重新生成创意` 表示废弃上一轮候选/);
assert.match(rules, /不支持 `改选创意 3` 等额外语法/);
assert.match(rules, /node scripts\/output-versioning\.mjs/);
assert.match(rules, /本批次全部产物统一写入脚本返回的下一个 `\{Listing标识\}-vN\/` 目录/);
assert.match(rules, /禁止给文案、图片、视频、脚本、验收报告或完整报告文件名添加 `-vN`/);
assert.match(rules, /同一批次后续.*必须复用已分配目录/);
assert.match(rules, /不得创建对应空目录/);
assert.match(rules, /node scripts\/check-output-layout\.mjs/);
assert.match(rules, /菜单 2“文案 \+ AI 商品图”只是示例/);
assert.match(rules, /node scripts\/next-action-scope\.mjs/);
assert.match(rules, /任一菜单都不得推荐其他菜单范围的动作/);
assert.match(rules, /图片生产前执行资料充足性门禁/);
assert.match(rules, /只影响单张图片的事实不足时，仅暂停对应图片/);
assert.match(rules, /不得作为买家可见文字、徽章或图标进入商品图/);
assert.match(rules, /暂停图片必须记录图片序号、图片类型、暂停原因、最低补充资料/);
assert.match(rules, /目标图片文件尚不存在时继续写入原 Listing 版本目录/);
assert.match(rules, /16 个编号章节之前生成 `问题速览`/);
assert.match(rules, /issue-\{两位章节号\}-\{两位问题序号\}/);
assert.match(rules, /node scripts\/check-report-anchors\.mjs <报告路径>/);
assert.match(rules, /卖家定制服务确认是创意方向和自动商品资料示例的共同前置门禁/);
assert.match(rules, /当前回复只能显示上述三项数字菜单/);
assert.match(rules, /不得同时生成或展示商品资料示例/);
assert.match(rules, /必须直接归一化并跳过三项菜单/);
assert.match(rules, /只说明商品能力，不能据此推断卖家服务为必选或可选/);
assert.match(rules, /已上架数据复盘必须先取得实际商品详情链接、复盘时间范围和运营数据/);
assert.match(rules, /商品名称只能辅助识别，不能替代链接/);
assert.match(rules, /未通过门禁时只索取缺失资料，不得进入复盘、生成报告或分配输出目录/);
assert.match(rules, /`平台商品-\{平台商品ID\}`/);

const coreAgent = files["agents/cross-border-commerce-agent.md"];
assert.match(coreAgent, /node scripts\/output-versioning\.mjs/);
assert.match(coreAgent, /Listing 版本目录/);
assert.match(coreAgent, /禁止给文件名添加 `-vN`/);
assert.match(coreAgent, /菜单 2“文案 \+ AI 商品图”未完成时/);
assert.doesNotMatch(coreAgent, /已生成图片时，默认至少包含/);
assert.match(coreAgent, /已上架数据复盘要求实际商品详情链接、复盘时间范围和运营数据/);
assert.match(coreAgent, /不得按商品名称或历史产物推测线上状态/);
assert.match(coreAgent, /链接门禁未通过时暂停复盘，不生成正式产物/);

const creative = files["workflows/creative-direction-selection.md"];
assert.match(creative, /不替代 `templates\/production-output\.md` 的 16 章完整生产报告/);
assert.match(creative, /输出 3-5 个.*创意方向/);
assert.match(creative, /请输入创意编号，例如：1；如需系统自动选择，请输入：0/);
assert.match(creative, /用户回复 `0` 时/);
assert.match(creative, /`0` 只在最近一次展示的是创意方向菜单时表示自动选择/);
assert.match(creative, /第 6 节“产品定位”/);
assert.match(creative, /不展示菜单不等于没有创意策略/);
assert.match(creative, /证据锚点/);
assert.match(creative, /未请求图片或视频时|未请求的内容载体/);
assert.match(creative, /多个数字时，不静默合并，也不默认采用第一个/);
assert.match(creative, /## 候选池与差异性门禁/);
assert.match(creative, /不少于计划展示数量两倍的内部候选/);
assert.match(creative, /至少必须在“目标人群、主要购买动机、核心场景、视觉叙事”中有两项实质不同/);
assert.match(creative, /无法得到至少 3 个通过差异性门禁的可靠方向时，不强行凑数/);
assert.doesNotMatch(creative, /每日专属感|送礼表达|照片记忆|Make Every Sip Feel Like Yours/);
assert.match(creative, /## 候选评分模型/);
assert.match(creative, /商品事实匹配度 \| 25%/);
assert.match(creative, /市场热度与内容信号 \| 15%/);
assert.match(creative, /已获得加权分 ÷ 已获得维度权重 × 100/);
assert.match(creative, /不得把“未获取”解释为“没有热度”/);
assert.match(creative, /用户回复 `重新选择创意` 时，原样重新展示最近一轮候选菜单/);
assert.match(creative, /用户回复 `重新生成创意` 时，废弃最近一轮候选及当前选择/);
assert.match(creative, /与上一轮方向做跨批次语义去重/);
assert.match(creative, /任一入口触发后都暂停示例确认和正式生产/);

const semantic = files["platforms/commerce-semantic-creative-rules.md"];
assert.match(semantic, /创意方向选择/);
assert.match(semantic, /创意方向选择流程/);
assert.match(semantic, /销售单位\/组合关系/);
assert.match(semantic, /变体与定制边界/);
assert.match(semantic, /图片只能直接确认可见的外观/);
assert.match(semantic, /跳过菜单不代表正式产出可以缺少创意策略/);
assert.match(semantic, /AI 生成图片是内容表达和场景模拟，不是商品事实/);
assert.match(semantic, /不得从固定的“日常、礼赠、照片、节日”分类直接填充菜单/);
assert.match(semantic, /只换方向名称、颜色词、节日名称或钩子视为语义重复/);
assert.match(semantic, /七项加权评分/);
assert.match(semantic, /按可用权重归一化总分并同时展示证据覆盖率/);

const platformProfiles = files["platforms/platform-profiles.md"];
assert.match(platformProfiles, /商品描述限制：每个段落最多 `500` 个字符/);
assert.match(platformProfiles, /字母、数字、空格和标点符号均计入/);
assert.match(platformProfiles, /完整句子边界拆段/);
assert.match(platformProfiles, /\[条件修饰词\].*\[核心品类词\]/);
assert.match(platformProfiles, /自身商品事实 > Temu 后台搜索词/);
for (const platformTitleRule of [
  "Amazon Brand Analytics",
  "Shopify 站内搜索词",
  "Etsy Shop Stats",
  "TikTok Shop 后台搜索",
  "eBay Product Research/Terapeak",
  "AliExpress 后台搜索词",
  "Walmart Seller Center/Search Insights",
]) {
  assert.match(platformProfiles, new RegExp(platformTitleRule.replace("/", "\\/")), `缺少平台标题数据规则：${platformTitleRule}`);
}
assert.match(platformProfiles, /标题结构按本文件对应平台规则执行/);

const marketSources = files["platforms/market-data-sources.md"];
assert.match(marketSources, /## 标题关键词数据链/);
assert.match(marketSources, /自身已确认商品事实[\s\S]*目标平台后台搜索词/);
assert.match(marketSources, /Google Trends.*不能单独证明目标平台流量/);
assert.match(marketSources, /每个进入最终标题的主要关键词必须记录/);
assert.match(marketSources, /## 创意方向热度证据/);
assert.match(marketSources, /目标平台直接信号、搜索趋势信号、内容热度信号或其他平台补充信号/);
assert.match(marketSources, /不把未获取解释为没有需求/);

const actions = files["workflows/action-menu.md"];
assert.match(actions, /用户选择启动菜单或下一步动作即视为同意执行该动作/);
assert.doesNotMatch(actions, /\| 是 \|/);
assert.match(actions, /不得覆盖已有文件，除非用户明确要求/);
assert.match(actions, /node scripts\/output-versioning\.mjs/);
assert.match(actions, /脚本返回的 Listing 版本目录/);
assert.match(actions, /同一批次的后续动作全部复用该目录/);
assert.match(actions, /不得留下空目录/);
assert.match(actions, /任务范围过滤（强制）/);
assert.match(actions, /菜单 2 禁止默认出现/);
assert.match(actions, /node scripts\/next-action-scope\.mjs/);
assert.match(actions, /不得为了凑数量加入扩展任务/);
assert.match(actions, /只处理记录为“暂停\/未生成”且目标文件尚不存在的图片/);
assert.match(actions, /用户直接回复暂停记录所需资料，等同于选择 `generate_missing_images`/);
assert.match(actions, /不包含覆盖任何已存在图片/);
assert.match(actions, /适用于全部启动菜单和自然语言单项任务/);
assert.match(actions, /output\/\{平台\}-\{市场\}\/\{Listing标识\}/);
assert.doesNotMatch(actions, /output\/\{目标平台\}\/\{产品类目\}\/\{产品名称\}/);

const outputStructure = files["workflows/output-structure.md"];
assert.match(outputStructure, /单个已匹配 SKU 且销售数量为 1/);
assert.match(outputStructure, /销售数量大于 1 或包含多个 SKU/);
assert.match(outputStructure, /未建档-\{简短商品名\}/);
assert.match(outputStructure, /找品-\{简短方向\}/);
assert.match(outputStructure, /商品名称、类目.*不再作为目录层级/);
assert.match(outputStructure, /上架\/完整生产报告\.md/);
assert.match(outputStructure, /强制版本预检/);
assert.match(outputStructure, /版本号只添加到 Listing 目录/);
assert.match(outputStructure, /禁止创建 `完整生产报告-v2\.md`/);
assert.match(outputStructure, /目录内.*文件名.*不得添加版本后缀/);
assert.match(outputStructure, /属于同一批次断点续做/);
assert.match(outputStructure, /不得仅凭“最新目录”猜测/);
assert.match(outputStructure, /已上架复盘无法匹配 SKU/);
assert.match(outputStructure, /平台商品-\{平台商品ID\}/);
assert.match(outputStructure, /不得退回使用 `未建档-\{简短商品名\}` 开始复盘/);
assert.match(outputStructure, /未通过时不得调用版本分配器、创建 Listing 目录/);

const subagents = files["agents/subagents.md"];
assert.match(subagents, /验证实际商品链接与当前线上页面/);
assert.match(subagents, /只有商品名称时先由 `intake-agent` 索取 Listing 身份依据/);
assert.match(subagents, /已上架复盘缺少有效商品详情链接/);

const productWorkflow = files["workflows/product-to-listing.md"];
assert.match(productWorkflow, /复盘必须先绑定实际线上 Listing/);
assert.match(productWorkflow, /缺少有效 Listing 身份时只索取资料，不进入分析，不创建完整报告或输出目录/);
assert.match(productWorkflow, /链接只能支持线上内容核验，不能替代未提供的后台指标/);

const chatOutput = files["templates/chat-output.md"];
for (const heading of ["本轮状态", "本轮产出", "缺失与风险", "已保存文件", "下一步动作"]) {
  assert.match(chatOutput, new RegExp(`## \\d\\. ${heading}`), `聊天结果模板缺少：${heading}`);
}
assert.match(chatOutput, /不得展开与当前任务无关的章节/);
assert.match(chatOutput, /完整生产报告\.md/);
assert.match(chatOutput, /1-5 个本轮明确任务范围内/);

const output = files["templates/production-output.md"];
for (let i = 1; i <= 16; i += 1) {
  assert.match(output, new RegExp(`## ${i}\\.`), `正式输出模板缺少第 ${i} 章`);
}
assert.match(output, /菜单 2“文案 \+ AI 商品图”完成后/);
assert.match(output, /只有 1 个合理动作时只提供 1 个/);
for (const imageType of [
  "最终定制主图",
  "到手内容/包装图",
  "定制操作示意图",
  "尺寸规格图",
  "尺寸对比图",
  "细节放大图",
  "日常使用场景图",
  "情绪/礼赠场景图",
  "卖点集合图",
]) {
  assert.match(output, new RegExp(imageType.replace("/", "\\/")), `缺少 Temu 图型：${imageType}`);
}
assert.match(output, /### 创意策略/);
assert.match(output, /文案主线/);
assert.match(output, /图片主线/);
assert.match(output, /视频主线/);
assert.match(output, /创意策略来源：用户选择 \/ 自动选择 \/ 沿用用户方向 \/ 系统默认 \/ 不适用/);
assert.match(output, /证据锚点/);
assert.match(output, /主差异轴/);
assert.match(output, /视觉叙事/);
assert.match(output, /可靠方向少于 3 个时不强行填表/);
assert.match(output, /创意候选评分方法：七项加权评分/);
assert.match(output, /热度数据获取状态：已获取 \/ 部分获取 \/ 未获取 \/ 不适用/);
assert.match(output, /归一化总分/);
assert.match(output, /市场热度证据另表记录/);
assert.match(output, /#### Temu 描述段落验收/);
assert.match(output, /字符数（含空格和标点）/);
assert.match(output, /### 标题关键词证据/);
assert.match(output, /查询入口\/关键词/);
assert.match(output, /在 `Name` 与 `Text`、`Photo` 与 `Image` 之间按上下文确定/);
assert.match(output, /标签与实际定制字段不符/);
assert.match(output, /定制类套图在表格前必须登记“定制母版”/);
assert.match(output, /必须先生成并验收主图，再用主图作为其余图片的视觉参考/);
assert.match(output, /与定制母版不一致/);
assert.match(output, /按“商品结构 \+ 定制位置 \+ 输入类型 \+ 定制工艺”归组/);
assert.match(output, /只选择一个视觉清晰的代表商品演示一次/);
assert.match(output, /不得把 `Text`、`Image` 等输入类型分别指向两个同类型组件/);
assert.match(output, /尺寸规格图和尺寸对比图必须分开处理/);
assert.match(output, /Temu 两张均为独立必备图/);
assert.match(output, /`Size Reference`/);
assert.match(output, /`Per Item` 与实际套装数量/);
assert.match(output, /图片生产前执行资料充足性门禁/);
assert.match(output, /包装形式未确认但实际到手内容已确认时/);
assert.match(output, /相关单图必须判定为“不可用”/);
assert.match(output, /暂停图片另附断点记录/);
assert.match(output, /恢复目标 Listing 版本目录/);
assert.match(output, /<a id="issue-summary"><\/a>/);
assert.match(output, /## 问题速览/);
assert.match(output, /`🔴【阻塞】`、`🟠【高风险】`、`🟡【待确认】`/);
assert.match(output, /\[查看详情\]\(#issue-02-01\)/);
assert.match(output, /\[返回问题速览\]\(#issue-summary\)/);
assert.match(output, /只有 `必选` 和 `可选` 才进入定制母版/);
assert.match(output, /服务为 `不提供` 时不得登记定制母版或生成定制操作示意图/);
assert.match(coreAgent, /不得固定套用某一组示例词/);
assert.match(coreAgent, /示例只能在建立本套“定制母版”时选择一次/);
assert.match(coreAgent, /定制类组合商品还必须记录组合内定制关系/);
assert.match(coreAgent, /商品是否具备定制能力，以及本 Listing 是否由卖家提供定制服务/);
assert.match(coreAgent, /只询问 `1\. 不提供定制`、`2\. 提供定制，买家必须提交定制内容`、`3\. 提供定制，买家可以选择是否定制`/);
assert.match(coreAgent, /定制服务选择是创意方向和自动商品资料示例的共同前置门禁/);
assert.match(coreAgent, /不得在同一回复中附带自动商品资料示例/);
assert.match(coreAgent, /直接归一化并跳过菜单/);
assert.match(files["agents/visual-production-agent.md"], /支持任意文字时用 `Add Your Text`/);
assert.match(files["agents/visual-production-agent.md"], /接受照片、插画或图案等广义图片时用 `Upload Your Image`/);
assert.match(files["agents/visual-production-agent.md"], /必须先生成并验收最终定制主图/);
assert.match(files["agents/visual-production-agent.md"], /默认一套图只展示一个定制母版/);
assert.match(files["agents/visual-production-agent.md"], /同组组件只选择一个视觉清晰的代表商品演示一次/);
assert.match(files["agents/visual-production-agent.md"], /不得使用 `Upload Your Image` 指向没有图片内容的空白区域/);
assert.match(files["agents/visual-production-agent.md"], /尺寸规格图和尺寸对比图是两个独立图型/);
assert.match(files["agents/visual-production-agent.md"], /Temu 两张均为必备图/);
assert.match(files["agents/visual-production-agent.md"], /比例失真或可能被误认作到手内容时判定为不可用/);
assert.match(files["agents/visual-production-agent.md"], /资料充足性门禁/);
assert.match(files["agents/visual-production-agent.md"], /暂停该图，继续生成其他事实充分的图片/);
assert.match(files["agents/visual-production-agent.md"], /包装形式未确认但实际到手内容已确认时/);
assert.match(files["agents/visual-production-agent.md"], /暂停与恢复/);
assert.match(files["agents/visual-production-agent.md"], /复用当前商品档案、创意策略、套图脚本和已完成图片/);
assert.match(files["agents/visual-production-agent.md"], /issue-image-\{两位图片序号\}/);
assert.match(files["agents/visual-production-agent.md"], /商品具备定制能力只描述物理或生产能力/);
assert.match(files["agents/visual-production-agent.md"], /卖家定制服务为 `待确认`：暂停定制资产生产/);
assert.match(files["platforms/image-set-rules.md"], /姓名专用字段可写 `Add Your Name`，自由文字字段写 `Add Your Text`/);
assert.match(files["platforms/image-set-rules.md"], /商品具备定制能力不等于本 Listing 启用卖家定制/);
assert.match(files["platforms/image-set-rules.md"], /两者都是 Temu 独立必备图型/);
assert.match(files["platforms/image-set-rules.md"], /定制类商品默认 `9 张`/);
assert.match(files["platforms/image-technical-specs.md"], /尺寸规格图和尺寸对比图均为独立必备图/);
assert.match(coreAgent, /Temu 定制类商品默认 `应生成图片数=9`/);
assert.match(coreAgent, /尺寸对比图通过真实比例验收/);
assert.match(files["agents/visual-production-agent.md"], /`应生成图片数` 默认为 9/);
assert.match(files["docs/agent-testing.md"], /Temu 定制类默认 9 张图/);
assert.match(files["docs/agent-testing.md"], /尺寸对比图与尺寸规格图为两个独立文件/);
assert.match(rules, /Temu 另将“尺寸对比图”作为独立必备图型/);

const productInput = files["templates/product-input.md"];
for (const field of ["商品具备定制能力", "卖家定制服务", "非卖家定制定位", "默认到手状态"]) {
  assert.match(productInput, new RegExp(field), `产品输入模板缺少定制流程字段：${field}`);
}
assert.match(files["workflows/start-guide.md"], /不得把普通成品和 DIY 空白基底拆成同级选项|不再放入同级菜单/);
assert.match(files["workflows/start-guide.md"], /当前回复只能显示一次纯数字选择/);
assert.match(files["workflows/start-guide.md"], /不得要求用户把定制选项数字和商品资料一起回传/);
assert.match(files["workflows/start-guide.md"], /收到数字选择并写入当前 Listing 后才继续/);
assert.match(files["workflows/start-guide.md"], /直接归一化并跳过菜单/);
assert.match(files["workflows/start-guide.md"], /只表示商品能力，不能代替 Listing 服务选择/);
assert.match(productInput, /## 已上架数据复盘的最小输入/);
assert.match(productInput, /商品名称不能替代链接/);
assert.match(productInput, /只提供商品名称、SKU 或运营数据时/);
assert.match(readme, /数据复盘不能只按商品名称执行/);
assert.match(readme, /缺少有效 Listing 身份时，Conductor 只提示补充资料/);
assert.ok(agent.workflow.includes("live_listing_validation"), "核心 Agent 工作流缺少线上 Listing 验证阶段");
assert.ok(agent.workflow.includes("growth_review"), "核心 Agent 工作流缺少数据复盘阶段");
assert.ok(agent.inputs.optional.includes("live_listing_url"), "核心 Agent 输入缺少实际 Listing 链接");
assert.ok(agent.inputs.optional.includes("review_date_range"), "核心 Agent 输入缺少复盘时间范围");
assert.ok(agent.inputs.optional.includes("seller_backend_metrics_or_export"), "核心 Agent 输入缺少后台运营数据");
assert.match(files["package.json"], /check-customization-flow\.mjs/);
assert.match(files["package.json"], /check-report-anchors\.mjs tests\/fixtures\/report-anchor-valid\.md/);

assert.match(rules, /node scripts\/temu-description-limit\.mjs/);
assert.match(rules, /姓名定制、自由文字、照片和广义图片应分别按上下文选择/);
assert.match(rules, /先生成并验收主图，再把主图中的最终定制效果登记为本套唯一“定制母版”/);
assert.match(rules, /不得逐张独立随机生成新的定制方案/);
assert.match(rules, /组合内定制关系/);
assert.match(rules, /同一类型.*只选择其中一个代表商品演示一次/);
assert.match(rules, /Temu 另将“尺寸对比图”作为独立必备图型/);
assert.match(rules, /两张图不得合并或互相替代/);
assert.match(rules, /`Dimensions Shown Per Item`/);
assert.match(rules, /商品“具备定制能力”和“本 Listing 是否启用卖家定制服务”必须分开记录/);
assert.match(rules, /不得把普通成品和 DIY 空白基底拆成同级选项/);
assert.match(rules, /生成标题前必须.*建立标题关键词数据链/);
assert.match(coreAgent, /每个段落最多 `500` 个字符/);
assert.match(coreAgent, /先建立关键词候选池，再确定词序/);
assert.match(coreAgent, /所有平台都必须执行标题关键词数据链/);
assert.match(coreAgent, /多平台任务分别研究、分别组词并分别输出/);
assert.match(readme, /商品定制能力与本 Listing 的卖家定制服务分开判断/);

const allText = Object.values(files).join("\n");
assert.doesNotMatch(allText, /product-catalog\.xlsx/);
assert.doesNotMatch(allText, /每个启用 SKU/);

const parseCsv = (text, fileName) => {
  const rows = [];
  let row = [];
  let field = "";
  let quoted = false;
  let line = 1;
  let rowLine = 1;

  for (let index = 0; index < text.length; index += 1) {
    const character = text[index];
    if (quoted) {
      if (character === '"' && text[index + 1] === '"') {
        field += '"';
        index += 1;
      } else if (character === '"') {
        quoted = false;
      } else {
        field += character;
        if (character === "\n") line += 1;
      }
      continue;
    }

    if (character === '"' && field === "") quoted = true;
    else if (character === ",") {
      row.push(field);
      field = "";
    } else if (character === "\n") {
      row.push(field.replace(/\r$/, ""));
      rows.push({ line: rowLine, values: row });
      row = [];
      field = "";
      line += 1;
      rowLine = line;
    } else field += character;
  }

  assert.equal(quoted, false, `${fileName} 第 ${rowLine} 行存在未闭合的双引号`);
  if (field !== "" || row.length > 0) {
    row.push(field.replace(/\r$/, ""));
    rows.push({ line: rowLine, values: row });
  }
  return rows;
};

const catalogHeaders = [
  "SPU", "SKU", "商品名称", "品牌", "产品类目", "产品子类目", "产品类型", "包含内容", "型号", "颜色/款式",
  "尺码/规格", "长度cm", "宽度cm", "高度cm", "净重g", "材质", "结构/表面工艺", "已确认功能", "商品特点", "适用对象",
  "使用场景", "使用/护理说明", "定制内容", "定制位置", "定制工艺", "包装清单", "包装方式", "包装长度cm", "包装宽度cm", "包装高度cm",
  "包装毛重g", "风险/禁用声明", "资料更新时间", "备注",
];
const attributeHeaders = ["SKU", "属性组", "属性名称", "属性值", "单位", "值类型", "是否平台必需", "适用平台", "资料来源", "资料更新时间", "备注"];
const catalogPath = resolve(root, "data/product-catalog.csv");
const attributePath = resolve(root, "data/product-attributes.csv");
await stat(resolve(root, "data/product-catalog.xlsx")).then(
  () => assert.fail("商品资料库应只使用 CSV，不应继续保留 data/product-catalog.xlsx"),
  () => undefined,
);
const catalogText = await readFile(catalogPath, "utf8");
const attributeText = await readFile(attributePath, "utf8");
assert.doesNotMatch(catalogText, /\uFFFD/, "data/product-catalog.csv 包含无效 UTF-8 字符");
assert.doesNotMatch(attributeText, /\uFFFD/, "data/product-attributes.csv 包含无效 UTF-8 字符");
assert.doesNotMatch(catalogText.split(/\r?\n/, 1)[0], /认证信息|认证状态|禁止\/未确认声明/, "data/product-catalog.csv 不应恢复旧认证字段");
assert.match(catalogText.split(/\r?\n/, 1)[0], /风险\/禁用声明/, "data/product-catalog.csv 缺少“风险/禁用声明”字段");
const catalogRows = parseCsv(catalogText, "data/product-catalog.csv");
const attributeRows = parseCsv(attributeText, "data/product-attributes.csv");
assert.deepEqual(catalogRows[0]?.values, catalogHeaders, "data/product-catalog.csv 表头不符合约定");
assert.deepEqual(attributeRows[0]?.values, attributeHeaders, "data/product-attributes.csv 表头不符合约定");

const isNumber = (value) => value.trim() !== "" && Number.isFinite(Number(value)) && Number(value) >= 0;
const isDate = (value) => {
  if (value === "") return true;
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;
  const parsed = new Date(`${value}T00:00:00Z`);
  return !Number.isNaN(parsed.getTime()) && parsed.toISOString().slice(0, 10) === value;
};
const indexByHeader = (headers) => Object.fromEntries(headers.map((header, index) => [header, index]));
const catalogIndex = indexByHeader(catalogHeaders);
const attributeIndex = indexByHeader(attributeHeaders);
const skuRows = new Map();

for (const row of catalogRows.slice(1)) {
  assert.equal(row.values.length, catalogHeaders.length, `data/product-catalog.csv 第 ${row.line} 行列数应为 ${catalogHeaders.length}`);
  for (const field of ["SPU", "SKU", "商品名称", "产品类目", "产品类型"]) {
    assert.notEqual(row.values[catalogIndex[field]].trim(), "", `data/product-catalog.csv 第 ${row.line} 行“${field}”不能为空`);
  }
  const sku = row.values[catalogIndex.SKU].trim();
  const spu = row.values[catalogIndex.SPU].trim();
  assert.doesNotMatch(spu, /(?:^|[\\/])\.\.(?:[\\/]|$)|[\\/]/, `data/product-catalog.csv 第 ${row.line} 行 SPU 不得包含路径字符`);
  assert.doesNotMatch(sku, /(?:^|[\\/])\.\.(?:[\\/]|$)|[\\/]/, `data/product-catalog.csv 第 ${row.line} 行 SKU 不得包含路径字符`);
  const normalizedSku = sku.toLocaleLowerCase("en-US");
  assert.equal(skuRows.has(normalizedSku), false, `data/product-catalog.csv 第 ${row.line} 行 SKU“${sku}”忽略大小写后重复`);
  skuRows.set(normalizedSku, row.line);
  assert.ok(["普货", "定制类"].includes(row.values[catalogIndex.产品类型]), `data/product-catalog.csv 第 ${row.line} 行“产品类型”只允许普货或定制类`);
  for (const field of ["长度cm", "宽度cm", "高度cm", "净重g", "包装长度cm", "包装宽度cm", "包装高度cm", "包装毛重g"]) {
    const value = row.values[catalogIndex[field]].trim();
    assert.ok(value === "" || isNumber(value), `data/product-catalog.csv 第 ${row.line} 行“${field}”只能填写大于或等于零的数字，当前值为“${value}”`);
  }
  const updatedAt = row.values[catalogIndex.资料更新时间].trim();
  assert.ok(isDate(updatedAt), `data/product-catalog.csv 第 ${row.line} 行“资料更新时间”必须使用 yyyy-mm-dd`);
  const productImageDirectory = resolve(root, "data", "products", spu, sku);
  const productImageFiles = await readdir(productImageDirectory).catch(() => []);
  const mainImages = productImageFiles.filter((file) => /^main\.(?:png|jpe?g|webp)$/i.test(file));
  assert.notEqual(mainImages.length, 0, `data/product-catalog.csv 第 ${row.line} 行无法按 SPU/SKU 找到主白底图：data/products/${spu}/${sku}/main.{png|jpg|jpeg|webp}`);
  assert.equal(mainImages.length, 1, `data/product-catalog.csv 第 ${row.line} 行主白底图扩展名冲突：${mainImages.join(", ")}`);
}

const attributeKeys = new Set();
for (const row of attributeRows.slice(1)) {
  assert.equal(row.values.length, attributeHeaders.length, `data/product-attributes.csv 第 ${row.line} 行列数应为 ${attributeHeaders.length}`);
  const sku = row.values[attributeIndex.SKU].trim();
  assert.notEqual(sku, "", `data/product-attributes.csv 第 ${row.line} 行“SKU”不能为空`);
  assert.ok(skuRows.has(sku.toLocaleLowerCase("en-US")), `data/product-attributes.csv 第 ${row.line} 行 SKU“${sku}”不存在于商品主表`);
  for (const field of ["属性组", "属性名称", "属性值", "值类型"]) {
    assert.notEqual(row.values[attributeIndex[field]].trim(), "", `data/product-attributes.csv 第 ${row.line} 行“${field}”不能为空`);
  }
  const valueType = row.values[attributeIndex.值类型].trim();
  assert.ok(["数字", "文本", "布尔值", "列表"].includes(valueType), `data/product-attributes.csv 第 ${row.line} 行“值类型”值无效`);
  const required = row.values[attributeIndex.是否平台必需].trim();
  assert.ok(["", "是", "否"].includes(required), `data/product-attributes.csv 第 ${row.line} 行“是否平台必需”只允许是、否或留空`);
  if (valueType === "数字") {
    const value = row.values[attributeIndex.属性值].trim();
    assert.ok(isNumber(value), `data/product-attributes.csv 第 ${row.line} 行数字属性只能填写大于或等于零的数字，当前值为“${value}”`);
  }
  const updatedAt = row.values[attributeIndex.资料更新时间].trim();
  assert.ok(isDate(updatedAt), `data/product-attributes.csv 第 ${row.line} 行“资料更新时间”必须使用 yyyy-mm-dd`);
  const source = row.values[attributeIndex.资料来源].trim();
  if (source.startsWith("data/")) {
    await stat(resolve(root, source)).catch(() => assert.fail(`data/product-attributes.csv 第 ${row.line} 行“资料来源”文件不存在：${source}`));
  }
  const key = [sku.toLocaleLowerCase("en-US"), row.values[attributeIndex.属性组], row.values[attributeIndex.属性名称], row.values[attributeIndex.适用平台]].join("|");
  assert.equal(attributeKeys.has(key), false, `data/product-attributes.csv 第 ${row.line} 行存在重复扩展属性`);
  attributeKeys.add(key);
}

const markdownFiles = [];
const collectMarkdown = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if ([".git", "node_modules", "output"].includes(entry.name)) continue;
    const path = resolve(directory, entry.name);
    if (entry.isDirectory()) await collectMarkdown(path);
    else if (entry.name.endsWith(".md")) markdownFiles.push(path);
  }
};
await collectMarkdown(root);
for (const markdownPath of markdownFiles) {
  const markdown = await readFile(markdownPath, "utf8");
  for (const match of markdown.matchAll(/\[[^\]]*\]\(([^)]+)\)/g)) {
    const target = match[1].trim().replace(/^<|>$/g, "").split("#", 1)[0];
    if (!target || /^(?:https?:|mailto:|data:)/.test(target) || target.includes("{")) continue;
    const localPath = resolve(dirname(markdownPath), decodeURIComponent(target));
    await stat(localPath).catch(() => assert.fail(`本地文档链接不存在：${markdownPath} -> ${target}`));
  }
}

console.log("规则一致性检查通过：入口、动态输入、已上架复盘门禁、SKU、CSV 商品资料、创意方向、跨平台标题数据链、双层输出、16 章完整报告、动作确认、Temu 描述与套图、本地文档链接均一致。");
