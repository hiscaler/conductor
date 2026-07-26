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
      "agents/visual-production-agent.md",
      "workflows/start-guide.md",
      "workflows/creative-direction-selection.md",
      "workflows/action-menu.md",
      "workflows/output-structure.md",
      "platforms/market-data-sources.md",
      "platforms/platform-profiles.md",
      "platforms/commerce-semantic-creative-rules.md",
      "platforms/image-set-rules.md",
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
assert.match(readme, /选择创意（需要时） → 正式生产/);
assert.match(readme, /请输入创意编号，例如：1；如需系统自动选择，请输入：0/);
assert.match(readme, /`0` 只在最近一次显示的是创意方向菜单时表示自动选择/);
assert.match(readme, /Temu 定制类商品默认以 8 张独立图片为完成标准/);
assert.match(readme, /同时使用 `cm` 和 `inch`/);
assert.match(readme, /供应链\//);
assert.match(readme, /利润\//);
assert.match(readme, /聊天结果.*默认只显示本轮相关内容/);
assert.match(readme, /上架\/完整生产报告\.md/);
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
assert.match(start, /创意确定后，再按所选任务正式生成文案、图片或视频/);
assert.match(start, /不得笼统写成与当前任务无关的全部产物/);
assert.match(start, /你修改并回传商品资料后，我会先生成具体创意方向/);
assert.match(start, /收到确认后的商品资料后，我会先生成具体创意方向并形成创意策略/);
assert.match(start, /存在多个可执行方向时，提供 3-5 个方向供你选择/);
assert.match(start, /非内容任务不显示创意菜单/);

const rules = files["AGENTS.md"];
assert.match(rules, /去除整个输入及每个组合项首尾空格/);
assert.match(rules, /英文字母不区分大小写/);
assert.match(rules, /两个 SKU 去除首尾空格并忽略大小写后相同时/);
assert.match(rules, /用户未回传商品资料示例前，不进入正式/);
assert.match(rules, /workflows\/creative-direction-selection\.md/);
assert.match(rules, /商品示例回传后先生成具体创意方向/);
assert.match(rules, /node scripts\/output-versioning\.mjs/);
assert.match(rules, /本批次全部产物统一写入脚本返回的下一个 `\{Listing标识\}-vN\/` 目录/);
assert.match(rules, /禁止给文案、图片、视频、脚本、验收报告或完整报告文件名添加 `-vN`/);
assert.match(rules, /同一批次后续.*必须复用已分配目录/);
assert.match(rules, /不得创建对应空目录/);
assert.match(rules, /node scripts\/check-output-layout\.mjs/);
assert.match(rules, /菜单 2“文案 \+ AI 商品图”只是示例/);
assert.match(rules, /node scripts\/next-action-scope\.mjs/);
assert.match(rules, /任一菜单都不得推荐其他菜单范围的动作/);

const coreAgent = files["agents/cross-border-commerce-agent.md"];
assert.match(coreAgent, /node scripts\/output-versioning\.mjs/);
assert.match(coreAgent, /Listing 版本目录/);
assert.match(coreAgent, /禁止给文件名添加 `-vN`/);
assert.match(coreAgent, /菜单 2“文案 \+ AI 商品图”未完成时/);
assert.doesNotMatch(coreAgent, /已生成图片时，默认至少包含/);

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

const semantic = files["platforms/commerce-semantic-creative-rules.md"];
assert.match(semantic, /创意方向选择/);
assert.match(semantic, /创意方向选择流程/);
assert.match(semantic, /销售单位\/组合关系/);
assert.match(semantic, /变体与定制边界/);
assert.match(semantic, /图片只能直接确认可见的外观/);
assert.match(semantic, /跳过菜单不代表正式产出可以缺少创意策略/);
assert.match(semantic, /AI 生成图片是内容表达和场景模拟，不是商品事实/);

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
assert.match(output, /尺寸参考物按条件启用/);
assert.match(output, /`Size Reference`/);
assert.match(output, /`Per Item` 与实际套装数量/);
assert.match(output, /只有 `必选` 和 `可选` 才进入定制母版/);
assert.match(output, /服务为 `不提供` 时不得登记定制母版或生成定制操作示意图/);
assert.match(coreAgent, /不得固定套用某一组示例词/);
assert.match(coreAgent, /示例只能在建立本套“定制母版”时选择一次/);
assert.match(coreAgent, /定制类组合商品还必须记录组合内定制关系/);
assert.match(coreAgent, /商品是否具备定制能力，以及本 Listing 是否由卖家提供定制服务/);
assert.match(coreAgent, /只询问 `1\. 不提供定制`、`2\. 必须定制`、`3\. 定制可选`/);
assert.match(files["agents/visual-production-agent.md"], /支持任意文字时用 `Add Your Text`/);
assert.match(files["agents/visual-production-agent.md"], /接受照片、插画或图案等广义图片时用 `Upload Your Image`/);
assert.match(files["agents/visual-production-agent.md"], /必须先生成并验收最终定制主图/);
assert.match(files["agents/visual-production-agent.md"], /默认一套图只展示一个定制母版/);
assert.match(files["agents/visual-production-agent.md"], /同组组件只选择一个视觉清晰的代表商品演示一次/);
assert.match(files["agents/visual-production-agent.md"], /不得使用 `Upload Your Image` 指向没有图片内容的空白区域/);
assert.match(files["agents/visual-production-agent.md"], /生成尺寸图前判断是否需要参考物/);
assert.match(files["agents/visual-production-agent.md"], /比例失真或可能被误认作到手内容时判定为不可用/);
assert.match(files["agents/visual-production-agent.md"], /商品具备定制能力只描述物理或生产能力/);
assert.match(files["agents/visual-production-agent.md"], /卖家定制服务为 `待确认`：暂停定制资产生产/);
assert.match(files["platforms/image-set-rules.md"], /姓名专用字段可写 `Add Your Name`，自由文字字段写 `Add Your Text`/);
assert.match(files["platforms/image-set-rules.md"], /商品具备定制能力不等于本 Listing 启用卖家定制/);

const productInput = files["templates/product-input.md"];
for (const field of ["商品具备定制能力", "卖家定制服务", "非卖家定制定位", "默认到手状态"]) {
  assert.match(productInput, new RegExp(field), `产品输入模板缺少定制流程字段：${field}`);
}
assert.match(files["workflows/start-guide.md"], /不得把普通成品和 DIY 空白基底拆成同级选项|不再放入同级菜单/);
assert.match(files["package.json"], /check-customization-flow\.mjs/);

assert.match(rules, /node scripts\/temu-description-limit\.mjs/);
assert.match(rules, /姓名定制、自由文字、照片和广义图片应分别按上下文选择/);
assert.match(rules, /先生成并验收主图，再把主图中的最终定制效果登记为本套唯一“定制母版”/);
assert.match(rules, /不得逐张独立随机生成新的定制方案/);
assert.match(rules, /组合内定制关系/);
assert.match(rules, /同一类型.*只选择其中一个代表商品演示一次/);
assert.match(rules, /只有商品大小难以直观判断、其他场景图未提供可靠尺度感且存在稳定参照时/);
assert.match(rules, /`Dimensions Shown Per Item`/);
assert.match(rules, /商品“具备定制能力”和“本 Listing 是否启用卖家定制服务”必须分开记录/);
assert.match(rules, /不得把普通成品和 DIY 空白基底拆成同级选项/);
assert.match(rules, /生成标题前必须.*建立标题关键词数据链/);
assert.match(coreAgent, /每个段落最多 `500` 个字符/);
assert.match(coreAgent, /先建立关键词候选池，再确定词序/);
assert.match(coreAgent, /所有平台都必须执行标题关键词数据链/);
assert.match(coreAgent, /多平台任务分别研究、分别组词并分别输出/);
assert.match(readme, /Temu 商品描述每段最多 `500` 个字符/);
assert.match(readme, /最终采用的主要关键词必须记录来源和使用位置/);
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

console.log("规则一致性检查通过：入口、动态输入、SKU、CSV 商品资料、创意方向、跨平台标题数据链、双层输出、16 章完整报告、动作确认、Temu 描述与套图、本地文档链接均一致。");
