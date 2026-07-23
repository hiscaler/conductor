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
      "workflows/start-guide.md",
      "workflows/creative-direction-selection.md",
      "workflows/action-menu.md",
      "workflows/output-structure.md",
      "platforms/commerce-semantic-creative-rules.md",
      "templates/production-output.md",
      "docs/agent-testing.md",
      "agents/cross-border-commerce-agent.json",
    ].map(async (path) => [path, await read(path)]),
  ),
);

const agent = JSON.parse(files["agents/cross-border-commerce-agent.json"]);
assert.deepEqual(agent.inputs.required, [], "核心 Agent 不应预设一组统一必填字段");
assert.equal(agent.product_catalog, "../data/product-catalog.xlsx");
assert.equal(agent.creative_direction_selection, "../workflows/creative-direction-selection.md");
assert.equal(agent.commerce_semantic_creative_rules, "../platforms/commerce-semantic-creative-rules.md");
assert.ok(agent.workflow.includes("creative_direction_selection"), "核心 Agent 工作流缺少创意方向选择阶段");

const start = files["workflows/start-guide.md"];
for (const trigger of ["你好", "hello", "开始", "菜单", "帮助"]) {
  assert.match(start, new RegExp(`- ${trigger}`), `缺少启动词：${trigger}`);
}
assert.match(start, /英文字母不区分大小写/);
assert.match(start, /多个 SKU 始终组成一个 listing/);
assert.match(start, /未写数量默认 1/);
assert.match(start, /相同 SKU 自动合并数量/);
assert.doesNotMatch(start, /分别生成.{0,12}组合销售/);
assert.match(start, /创意方向选择/);
assert.match(start, /用户只需回复一个数字/);

const rules = files["AGENTS.md"];
assert.match(rules, /去除整个输入及每个组合项首尾空格/);
assert.match(rules, /英文字母不区分大小写/);
assert.match(rules, /两个 SKU 去除首尾空格并忽略大小写后相同时/);
assert.match(rules, /用户未回传商品资料示例前，不进入正式/);
assert.match(rules, /workflows\/creative-direction-selection\.md/);

const creative = files["workflows/creative-direction-selection.md"];
assert.match(creative, /不替代 `templates\/production-output\.md` 的 16 章正式产出/);
assert.match(creative, /输出 3-5 个.*创意方向/);
assert.match(creative, /用户只回复数字或“自动选择”|用户回复 `1`/);
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

const actions = files["workflows/action-menu.md"];
assert.match(actions, /用户选择启动菜单或下一步动作即视为同意执行该动作/);
assert.doesNotMatch(actions, /\| 是 \|/);
assert.match(actions, /不得覆盖已有文件，除非用户明确要求/);
assert.match(actions, /output\/\{平台\}-\{市场\}\/\{Listing标识\}/);
assert.doesNotMatch(actions, /output\/\{目标平台\}\/\{产品类目\}\/\{产品名称\}/);

const outputStructure = files["workflows/output-structure.md"];
assert.match(outputStructure, /单个已匹配 SKU/);
assert.match(outputStructure, /多个 SKU/);
assert.match(outputStructure, /未建档-\{简短商品名\}/);
assert.match(outputStructure, /找品-\{简短方向\}/);
assert.match(outputStructure, /商品名称、类目.*不再作为目录层级/);

const output = files["templates/production-output.md"];
for (let i = 1; i <= 16; i += 1) {
  assert.match(output, new RegExp(`## ${i}\\.`), `正式输出模板缺少第 ${i} 章`);
}
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

const allText = Object.values(files).join("\n");
assert.doesNotMatch(allText, /product-catalog\.csv/);
assert.doesNotMatch(allText, /每个启用 SKU/);

const workbookUrl = new URL("../data/product-catalog.xlsx", import.meta.url);
const workbook = await stat(workbookUrl);
assert.ok(workbook.size > 10_000, "商品资料库工作簿不存在或内容异常");
const signature = await readFile(workbookUrl);
assert.equal(signature.subarray(0, 2).toString(), "PK", "商品资料库不是有效的 XLSX/ZIP 文件");

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

console.log("规则一致性检查通过：入口、动态输入、SKU、示例确认、创意方向、动作确认、16 章模板、Temu 套图、XLSX 和本地文档链接均一致。");
