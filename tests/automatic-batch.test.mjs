import test from "node:test";
import assert from "node:assert/strict";
import { mkdir, mkdtemp, writeFile, readFile, readdir, rm, copyFile, unlink } from "node:fs/promises";
import { resolve, relative, dirname, sep } from "node:path";
import { fileURLToPath } from "node:url";
import { execFile } from "node:child_process";
import { promisify } from "node:util";
import sharp from "sharp";
import { CATALOG_HEADERS, ATTRIBUTE_HEADERS, parseCsv, loadProductCatalog } from "../scripts/product-catalog.mjs";
import { AUTO_TASKS, createBatch, claimJob, recordJob, retryJob, loadBatch, normalizeBatchConfig, validateBatchListing } from "../scripts/batch-production.mjs";

const root = fileURLToPath(new URL("..", import.meta.url));
const execute = promisify(execFile);
const docLink = (from, report) => `[批次总报告](${relative(dirname(from), report).split(sep).map(encodeURIComponent).join("/")})`;
async function put(path, text) { await mkdir(dirname(path), { recursive: true }); await writeFile(path, text); }
async function image(path, colour, compressionLevel = 9) {
  await mkdir(dirname(path), { recursive: true });
  await sharp({ create: { width: 32, height: 32, channels: 3, background: colour } }).png({ compressionLevel }).toFile(path);
}
async function fixture(t, { skus = ["SKU-B", "SKU-A"], patterns = ["red", "blue"], main = true } = {}) {
  await mkdir(resolve(root, ".tmp"), { recursive: true });
  const project = await mkdtemp(resolve(root, ".tmp/automatic-batch-tests-"));
  // Cleanup is limited to this test's resolved, freshly allocated workspace directory.
  t.after(async () => {
    assert.ok(project.startsWith(resolve(root, ".tmp") + sep));
    await rm(project, { recursive: true, force: true });
  });
  const catalogRows = skus.map((sku) => CATALOG_HEADERS.map((field) => ({ SPU: "SPU", SKU: sku, 商品名称: "测试商品", 产品类目: "测试", 产品类型: "定制类", 资料更新时间: "2026-10-07" }[field] ?? "")).join(","));
  await put(resolve(project, "data/product-catalog.csv"), [CATALOG_HEADERS.join(","), ...catalogRows].join("\n") + "\n");
  await put(resolve(project, "data/product-attributes.csv"), ATTRIBUTE_HEADERS.join(",") + "\n");
  for (const sku of skus) if (main) await image(resolve(project, "data/products/SPU", sku, "main.png"), "white");
  const samples = resolve(project, "samples");
  await mkdir(samples);
  for (let i = 0; i < patterns.length; i++) await image(resolve(samples, `pattern-${i}.png`), patterns[i]);
  const scripts = ["batch-production.mjs", "product-catalog.mjs", "output-versioning.mjs", "validate-listing.mjs", "check-report-anchors.mjs", "check-output-layout.mjs", "check-image-set.mjs", "next-action-scope.mjs", "temu-description-limit.mjs", "image-sets.json"];
  await mkdir(resolve(project, "scripts"));
  for (const file of scripts) await copyFile(resolve(root, "scripts", file), resolve(project, "scripts", file));
  const config = { spu: "spu", sample_directory: samples, platform: "Temu", market: "US", auto_task: 2, seller_service: "必选", customization_type: "图文" };
  return { project, samples, config, create: (input = config) => createBatch(input, { projectRoot: project, now: new Date("2026-10-07T03:22:00Z") }) };
}

test("自动第四项映射图片范围；确认门禁及非 Temu 默认数量", () => {
  assert.equal(AUTO_TASKS[4].menu, 8);
  const input = { auto_task: 4, spu: "SPU", platform: "Temu", market: "US", sample_directory: "samples", seller_service: "必选", customization_type: "仅图片" };
  assert.throws(() => normalizeBatchConfig({ ...input, seller_service: "待确认" }), /尚未确认/);
  assert.throws(() => normalizeBatchConfig({ ...input, customization_type: "待确认" }), /尚未确认/);
  assert.equal(normalizeBatchConfig({ ...input, customization_type: "仅文字" }).image_count, 5);
  assert.throws(() => normalizeBatchConfig({ ...input, platform: "Shopify" }), /image_count/);
  assert.equal(normalizeBatchConfig({ ...input, platform: "Shopify", image_count: 8 }).image_preset, "platform_default");
});

test("SKU 和图案固定一对一：像素重复、损坏文件与多余图案被报告", async (t) => {
  const f = await fixture(t, { patterns: ["red", "blue", "green"] });
  await image(resolve(f.samples, "duplicate.png"), "red", 0);
  await put(resolve(f.samples, "broken.png"), "not an image");
  await put(resolve(f.samples, "readme.txt"), "not a pattern");
  const state = await f.create();
  assert.deepEqual(state.jobs.map((job) => job.sku), ["SKU-A", "SKU-B"]);
  assert.equal(new Set(state.jobs.map((job) => job.sample.visual_sha256)).size, 2);
  assert.equal(state.inventory.valid_image_count, 3);
  assert.equal(state.exclusions.length, 4);
  assert.ok(state.exclusions.some((item) => /像素相同/.test(item.reason)));
  const report = await readFile(state.report_path, "utf8");
  assert.match(report, /图片数量多于 SKU/);
  assert.equal((await readdir(state.output_root)).length, 2, "未领取时只存在总报告与状态，不预创建 Listing 目录");
  assert.deepEqual((await loadBatch(state.state_path)).jobs.map((job) => job.sample.path), state.jobs.map((job) => job.sample.path));
});

test("图案少于 SKU 时报告未参与项，仅文字只将图案作为创意参考", async (t) => {
  const f = await fixture(t, { patterns: ["red"] });
  const state = await f.create({ ...f.config, customization_type: "仅文字" });
  assert.equal(state.jobs.length, 1);
  assert.equal(state.jobs[0].sample_usage, "creative_reference_only");
  assert.equal(state.exclusions[0].sku, "SKU-B");
  assert.match(await readFile(state.report_path, "utf8"), /未参与本批次/);
});

test("跨进程创建同分钟批次及重复领取不会共享目录", async (t) => {
  const f = await fixture(t);
  await put(resolve(f.project, "config.json"), JSON.stringify(f.config));
  const results = await Promise.all(Array.from({ length: 3 }, () => execute(process.execPath, [resolve(f.project, "scripts/batch-production.mjs"), "create", resolve(f.project, "config.json")], { cwd: f.project })));
  const batches = results.map((result) => JSON.parse(result.stdout));
  assert.equal(new Set(batches.map((batch) => batch.batch_id)).size, 3);
  const claims = await Promise.allSettled([claimJob(batches[0].state_path, "SKU-A"), claimJob(batches[0].state_path, "SKU-A")]);
  assert.equal(claims.filter((result) => result.status === "fulfilled").length, 1);
  const claimed = claims.find((result) => result.status === "fulfilled").value;
  const other = await claimJob(batches[1].state_path, "SKU-A");
  assert.notEqual(claimed.job.output_dir, other.job.output_dir);
  await assert.rejects(claimJob(batches[0].state_path, "SKU-A"), /不能重复领取/);
  await assert.rejects(recordJob(batches[0].state_path, "SKU-A", "wrong-token", { message: "重复更新" }), /令牌不匹配/);
});

test("单项阻塞不阻止其他 SKU；补图恢复复用原配对", async (t) => {
  const f = await fixture(t, { main: false });
  await image(resolve(f.project, "data/products/SPU/SKU-B/main.png"), "white");
  const state = await f.create();
  assert.equal(state.jobs[0].status, "blocked");
  assert.equal((await claimJob(state.state_path)).job.sku, "SKU-B");
  await image(resolve(f.project, "data/products/SPU/SKU-A/main.png"), "white");
  const restored = await retryJob(state.state_path, "SKU-A");
  assert.deepEqual(restored.sample, state.jobs[0].sample);
  const claimed = await claimJob(state.state_path, "SKU-A");
  const directory = claimed.job.output_dir;
  await recordJob(state.state_path, "SKU-A", claimed.job.lease_token, { status: "blocked", blockers: ["尺寸资料缺失"], message: "单项暂停" });
  await retryJob(state.state_path, "SKU-A");
  assert.equal((await claimJob(state.state_path, "SKU-A")).job.output_dir, directory);
});

test("输入改变时阻止续做；不能仅凭记录请求标记完成", async (t) => {
  const f = await fixture(t);
  const state = await f.create();
  const claim = await claimJob(state.state_path);
  const failed = await recordJob(state.state_path, claim.job.sku, claim.job.lease_token, { status: "completed", stage: "validation", message: "申请验收" });
  assert.equal(failed.status, "failed");
  assert.equal(failed.validation.passed, false);
  await image(claim.job.sample.path, "black");
  await assert.rejects(retryJob(state.state_path, claim.job.sku), /发生变化/);
});

test("真实文案产物通过闸门，错误回链与缺失文件不能通过批次检查", async (t) => {
  const f = await fixture(t);
  const state = await f.create({ ...f.config, auto_task: 1 });
  const claim = await claimJob(state.state_path);
  const directory = claim.job.output_dir;
  const writeDoc = async (file, body) => {
    const path = resolve(directory, file);
    await put(path, `批次：${state.batch_id}\n\n${docLink(path, state.report_path)}\n\n${body}\n`);
  };
  await writeDoc("商品资料.md", "# 商品资料\n\n测试事实");
  await writeDoc("文案/文案资产.md", "## 商品标题\n\nTest product\n\n## 长描述\n\nA sample product description for validation.");
  await writeDoc("文案/关键词清单.md", "## 关键词\n\nsample product");
  await writeDoc("上架/完整生产报告.md", '<a id="issue-summary"></a>\n\n## 问题速览\n\n无\n\n' + Array.from({ length: 16 }, (_, i) => `## ${i + 1}. 测试章节\n\n测试资料`).join("\n\n"));
  await put(resolve(directory, "上架/next-action.json"), JSON.stringify({ menu: 1, actions: ["view_current_assets"], batch_id: state.batch_id, batch_report: relative(resolve(directory, "上架"), state.report_path) }));
  assert.equal((await validateBatchListing(directory)).checked, true);
  const finished = await recordJob(state.state_path, claim.job.sku, claim.job.lease_token, { status: "completed", stage: "validation", message: "实际文件与交付闸门通过" });
  assert.equal(finished.status, "completed");
  assert.match(await readFile(state.report_path, "utf8"), /已完成并校验/);
  await put(resolve(directory, "文案/关键词清单.md"), `批次：${state.batch_id}\n\n[批次总报告](wrong.md)`);
  await assert.rejects(validateBatchListing(directory), /正确的相对报告链接/);
  await unlink(resolve(directory, "文案/文案资产.md"));
  await assert.rejects(validateBatchListing(directory), /ENOENT/);
});

test("CSV 拒绝孤立属性、重复 SKU 和非法引号，并保留合法多行字段", async (t) => {
  assert.equal(parseCsv('A,B\n"x\ny","a,b"\n', "fixture")[1].values[0], "x\ny");
  assert.throws(() => parseCsv('A\n"x"bad', "fixture"), /引号结构错误/);
  const f = await fixture(t, { skus: ["SKU-A", " sku-a "] });
  await assert.rejects(loadProductCatalog(f.project), /忽略大小写后重复/);
  const f2 = await fixture(t);
  await put(resolve(f2.project, "data/product-attributes.csv"), ATTRIBUTE_HEADERS.join(",") + "\nUNKNOWN,规格,容量,1,ml,数字,是,Temu,,2026-10-07,\n");
  await assert.rejects(loadProductCatalog(f2.project), /不存在于商品主表/);
});
