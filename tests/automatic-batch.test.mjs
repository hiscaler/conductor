import test from "node:test";
import assert from "node:assert/strict";
import { mkdir, mkdtemp, writeFile, readFile, readdir, rm, copyFile, unlink } from "node:fs/promises";
import { resolve, relative, dirname, sep } from "node:path";
import { fileURLToPath } from "node:url";
import { execFile } from "node:child_process";
import { promisify } from "node:util";
import sharp from "sharp";
import { CATALOG_HEADERS, ATTRIBUTE_HEADERS, parseCsv, loadProductCatalog } from "../scripts/product-catalog.mjs";
import { AUTO_TASKS, MAPPING_LABELS, createBatch, claimJob, recordJob, retryJob, loadBatch, normalizeBatchConfig, validateBatchListing } from "../scripts/batch-production.mjs";

const root = fileURLToPath(new URL("..", import.meta.url));
const execute = promisify(execFile);
test("菜单、规则和报告统一图案分配方式名称并保留机器枚举", async () => {
  assert.deepEqual(MAPPING_LABELS, { one_to_one: "一张图案分配给一个 SKU 变体", one_to_all_skus: "一张图案应用到全部 SKU 变体" });
  for (const file of ["AGENTS.md", "README.md", "workflows/start-guide.md", "workflows/automatic-batch.md", "workflows/output-structure.md", "templates/production-output.md"]) {
    const content = await readFile(resolve(root, file), "utf8");
    for (const [mode, label] of Object.entries(MAPPING_LABELS)) {
      assert.ok(content.includes(label), `${file} 缺少名称 ${label}`);
      assert.ok(content.includes(mode), `${file} 缺少机器枚举 ${mode}`);
    }
  }
});
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
  const scripts = ["batch-production.mjs", "product-catalog.mjs", "output-versioning.mjs", "validate-listing.mjs", "check-report-anchors.mjs", "check-output-layout.mjs", "check-copy-language.mjs", "check-image-set.mjs", "next-action-scope.mjs", "temu-description-limit.mjs", "image-sets.json", "pattern-group-images.mjs"];
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
  assert.throws(() => normalizeBatchConfig({ ...input, mapping_mode: "unknown" }), /mapping_mode/);
  assert.equal(normalizeBatchConfig(input).mapping_mode, "one_to_one");
  assert.throws(() => normalizeBatchConfig({ ...input, customization_type: "文" }), /sample_usage/);
  assert.equal(normalizeBatchConfig({ ...input, customization_type: "仅文字", sample_usage: "fixed_product_artwork" }).image_count, 5);
  assert.throws(() => normalizeBatchConfig({ ...input, platform: "Shopify" }), /image_count/);
  assert.equal(normalizeBatchConfig({ ...input, platform: "Shopify", image_count: 8 }).image_preset, "platform_default");
});

test("每张图案应用到全部 SKU 变体：完整矩阵、任务选择及恢复互不覆盖", async (t) => {
  const f = await fixture(t, { skus: ["SKU-C", "SKU-B", "SKU-A"] });
  const state = await f.create({ ...f.config, mapping_mode: "one_to_all_skus" });
  assert.equal(state.jobs.length, 6);
  assert.equal(state.exclusions.length, 0);
  assert.deepEqual(state.jobs.map((job) => job.sku), ["SKU-A", "SKU-B", "SKU-C", "SKU-A", "SKU-B", "SKU-C"]);
  assert.equal(state.pattern_groups.length, 2);
  for (const group of state.pattern_groups) {
    assert.deepEqual(group.skus, ["SKU-A", "SKU-B", "SKU-C"]);
    assert.equal(group.job_ids.length, 3);
  }
  assert.equal((await loadBatch(state.state_path)).jobs.length, 6);
  assert.match(await readFile(state.report_path, "utf8"), /一张图案应用到全部 SKU 变体/);
  await assert.rejects(claimJob(state.state_path, "SKU-A"), /多个图案任务/);
  const first = await claimJob(state.state_path, "job-001");
  const second = await claimJob(state.state_path, "job-004");
  assert.notEqual(first.job.output_dir, second.job.output_dir);
  assert.notEqual(first.pattern_group.customization_master_id, second.pattern_group.customization_master_id);
  const assoc = JSON.parse(await readFile(resolve(first.job.output_dir, "批次关联.json"), "utf8"));
  assert.equal(assoc.pattern_group_id, first.pattern_group.id);
  assert.deepEqual(assoc.variant_skus, ["SKU-A", "SKU-B", "SKU-C"]);
  await assert.rejects(recordJob(state.state_path, "job-004", first.job.lease_token, { message: "错误令牌" }), /令牌不匹配/);
  await assert.rejects(recordJob(state.state_path, "SKU-A", first.job.lease_token, { message: "歧义选择" }), /多个图案任务/);
  await recordJob(state.state_path, "job-001", first.job.lease_token, { status: "blocked", blockers: ["等待素材"], message: "暂停" });
  await assert.rejects(retryJob(state.state_path, "SKU-A"), /多个图案任务/);
  await retryJob(state.state_path, "job-001");
  const restored = await claimJob(state.state_path, "job-001");
  assert.equal(restored.job.output_dir, first.job.output_dir);
  assert.deepEqual(restored.job.sample, first.job.sample);
  assert.equal((await claimJob(state.state_path)).job.id, "job-002");
});

test("共享图案矩阵损坏不能恢复；旧一对一状态不需要 mapping_mode", async (t) => {
  const f = await fixture(t);
  const legacy = await f.create();
  delete legacy.config.mapping_mode;
  await put(legacy.state_path, JSON.stringify(legacy));
  assert.equal((await loadBatch(legacy.state_path)).jobs.length, 2);
  const state = await f.create({ ...f.config, mapping_mode: "one_to_all_skus" });
  const corruptions = [
    (copy) => copy.jobs.pop(),
    (copy) => { copy.jobs[2].sample = copy.jobs[0].sample; },
    (copy) => { copy.pattern_groups[0].skus.pop(); },
    (copy) => { copy.jobs[1].id = copy.jobs[0].id; },
    (copy) => { copy.pattern_groups[0].customization_master_id = "wrong"; },
    (copy) => { copy.jobs[0].output_dir = copy.jobs[1].output_dir = resolve(copy.output_root, "Temu-US/shared"); },
  ];
  for (const mutate of corruptions) {
    const copy = structuredClone(state);
    mutate(copy);
    await put(state.state_path, JSON.stringify(copy));
    await assert.rejects(loadBatch(state.state_path), /损坏|不一致|共用输出目录/);
  }
  await put(state.state_path, JSON.stringify(state));
  assert.equal((await loadBatch(state.state_path)).jobs.length, 4);
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
  const state = await f.create({ ...f.config, customization_type: "仅文字", sample_usage: "creative_reference_only" });
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

for (const mapping_mode of ["one_to_one", "one_to_all_skus"]) test(`套图完成闸门检查实际数量、SKU、图案来源和母版：${mapping_mode}`, async (t) => {
  const f = await fixture(t);
  const state = await f.create({ ...f.config, auto_task: 4, customization_type: "仅图片", mapping_mode });
  // 此既有闸门用例保留历史逐 SKU 布局，检查旧批次兼容性。
  delete state.config.output_layout;
  await put(state.state_path, JSON.stringify(state));
  const { job, pattern_group } = await claimJob(state.state_path);
  const directory = job.output_dir;
  const writeDoc = async (file, body) => {
    const path = resolve(directory, file);
    await put(path, `批次：${state.batch_id}\n\n${docLink(path, state.report_path)}\n\n${body}`);
  };
  await writeDoc("商品资料.md", "测试商品资料");
  await writeDoc("图片/套图脚本.md", "测试图型及独立图案");
  await writeDoc("图片/图片验收报告.md", "测试技术验收");
  await writeDoc("上架/完整生产报告.md", '<a id="issue-summary"></a>\n\n## 问题速览\n\n无\n\n' + Array.from({ length: 16 }, (_, i) => `## ${i + 1}. 测试章节\n\n测试事实`).join("\n\n"));
  await put(resolve(directory, "上架/next-action.json"), JSON.stringify({ menu: 8, actions: ["view_current_assets"], batch_id: state.batch_id, batch_report: relative(resolve(directory, "上架"), state.report_path) }));
  const names = JSON.parse(await readFile(resolve(root, "scripts/image-sets.json"), "utf8")).Temu.presets.standard_5.required;
  const metadata = { batch_id: state.batch_id, sku: job.sku, source_pattern_sha256: job.sample.sha256, source_pattern_usage: "customization_image", batch_report: relative(resolve(directory, "图片"), state.report_path), customization_master_id: pattern_group?.customization_master_id ?? `${state.batch_id}-${job.sku}-master`, ...(pattern_group ? { job_id: job.id, pattern_group_id: pattern_group.id } : {}) };
  await put(resolve(directory, "图片/套图配置.json"), JSON.stringify({ ...metadata, platform: "Temu", preset: "standard_5", selected_count: 5, selection_source: "default", supplemental_images: [], additional_required: [] }));
  for (const name of names.slice(0, 4)) {
    await sharp({ create: { width: 800, height: 800, channels: 3, background: "red" } }).png().toFile(resolve(directory, "图片", name));
    await put(resolve(directory, "图片", name.replace(/\.png$/, ".verify.json")), JSON.stringify({ ...metadata, watermark: "none", pattern_consistent_with_master: true }));
  }
  await assert.rejects(validateBatchListing(directory), /图片数量/);
  await sharp({ create: { width: 800, height: 800, channels: 3, background: "red" } }).png().toFile(resolve(directory, "图片", names[4]));
  const sidecar = resolve(directory, "图片", names[4].replace(/\.png$/, ".verify.json"));
  await put(sidecar, JSON.stringify({ ...metadata, watermark: "none", pattern_consistent_with_master: true }));
  assert.equal((await validateBatchListing(directory, { runListingGate: true })).checked, true);
  if (pattern_group) {
    const assocPath = resolve(directory, "批次关联.json");
    const assoc = JSON.parse(await readFile(assocPath, "utf8"));
    await put(assocPath, JSON.stringify({ ...assoc, pattern_group_id: "wrong" }));
    await assert.rejects(validateBatchListing(directory), /共享图案组关联不一致/);
    await put(assocPath, JSON.stringify(assoc));
    await put(sidecar, JSON.stringify({ ...metadata, watermark: "none", pattern_consistent_with_master: true, job_id: "wrong" }));
    await assert.rejects(validateBatchListing(directory), /共享图案组 \/ 母版不一致/);
  }
  await put(sidecar, JSON.stringify({ ...metadata, customization_master_id: "another-master" }));
  await assert.rejects(validateBatchListing(directory), /母版不一致/);
  await put(sidecar, JSON.stringify({ ...metadata, sku: "SKU-B" }));
  await assert.rejects(validateBatchListing(directory), /来源不一致/);
});
