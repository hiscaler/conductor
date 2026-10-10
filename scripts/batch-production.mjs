import { open, readFile, writeFile, rename, unlink, readdir, mkdir, stat } from "node:fs/promises";
import { resolve, relative, dirname, basename, extname, isAbsolute, sep } from "node:path";
import { pathToFileURL, fileURLToPath } from "node:url";
import { randomUUID } from "node:crypto";
import { spawnSync } from "node:child_process";
import sharp from "sharp";
import { loadProductCatalog, loadGallery, IMAGE_EXTENSIONS, normalizeKey, sha256 } from "./product-catalog.mjs";
import { resolveOutputDirectory } from "./output-versioning.mjs";
import { patternGroupImages } from "./pattern-group-images.mjs";

export const AUTO_TASKS = Object.freeze({
  1: { menu: 1, label: "标题、描述、关键词", copy: true, images: false, video_script: false },
  2: { menu: 2, label: "文案 + AI 商品图", copy: true, images: true, video_script: false },
  3: { menu: 3, label: "文案 + AI 商品图 + 视频脚本", copy: true, images: true, video_script: true },
  4: { menu: 8, label: "AI 商品图", copy: false, images: true, video_script: false },
});
const ROOT = fileURLToPath(new URL("..", import.meta.url));
const slash = (path) => path.split(sep).join("/");
const md = (value) => String(value ?? "").replace(/[\r\n]+/g, " ").replace(/\|/g, "\\|");
const url = (path) => slash(path).split("/").map(encodeURIComponent).join("/");
const link = (from, target, label) => {
  const rel = relative(dirname(from), target);
  const safeLabel = String(label).replace(/[\r\n]/g, " ").replace(/[\[\]|]/g, "\\$&");
  return isAbsolute(rel) ? `[${safeLabel}](<${slash(target)}>)` : `[${safeLabel}](${url(rel)})`;
};
const inside = (root, target) => {
  const rel = relative(root, target);
  return rel && rel !== ".." && !rel.startsWith(`..${sep}`) && !isAbsolute(rel);
};
const segment = (value, label) => {
  if (typeof value !== "string" || !value.trim() || /[<>:"/\\|?*\x00-\x1F]/.test(value)
    || /[. ]$/.test(value) || /^(?:\.|\.\.|CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(?:\.|$)/i.test(value)) {
    throw new Error(`${label} 不能安全用于跨平台目录名：${JSON.stringify(value)}`);
  }
  return value.trim();
};
const SERVICE = { "不提供": "not_offered", "必选": "required", "可选": "optional", not_offered: "not_offered", required: "required", optional: "optional" };
const TYPE = { "图文": "image_text", "图": "image", "文": "text", "仅图片": "image", "仅文字": "text", image_text: "image_text", image: "image", text: "text" };
const allSkus = (state) => state.config.mapping_mode === "one_to_all_skus";
export const MAPPING_LABELS = Object.freeze({ one_to_one: "一张图案分配给一个 SKU 变体", one_to_all_skus: "一张图案应用到全部 SKU 变体" });
const groupFor = (state, job) => state.pattern_groups?.find((group) => group.id === job.pattern_group_id);
const groupedOutput = (state) => state.config.output_layout === "pattern_group";
const patternName = (sample) => {
  let name = basename(sample.path, extname(sample.path)).replace(/[<>:"/\\|?*\x00-\x1F]/g, "-").trim().slice(0, 100).replace(/[. ]+$/g, "");
  if (!name || /^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(?:\.|$)/i.test(name)) name = `图案-${name || "未命名"}`;
  if (/-v[1-9]\d*$/i.test(name)) name += "-图案";
  return name;
};
function selectJob(state, selector) {
  if (!selector) return state.jobs.find((job) => job.status === "pending");
  const byId = state.jobs.find((job) => job.id === selector);
  if (byId) return byId;
  const matches = state.jobs.filter((job) => normalizeKey(job.sku) === normalizeKey(selector));
  if (matches.length > 1) throw new Error(`SKU ${selector} 对应多个图案任务，请使用任务 ID：${matches.map((job) => job.id).join("、")}`);
  return matches[0];
}

export function normalizeBatchConfig(input) {
  const mappingMode = input.mapping_mode ?? "one_to_one";
  if (!["one_to_one", "one_to_all_skus"].includes(mappingMode)) throw new Error("mapping_mode 只允许 one_to_one / one_to_all_skus");
  const task = AUTO_TASKS[input.auto_task];
  if (!task) throw new Error("auto_task 只允许批量菜单 1、2、3、4（4 映射单个菜单 8）");
  const spu = segment(input.spu, "SPU");
  const platformInput = segment(input.platform, "目标平台");
  const platform = ["Temu", "Amazon", "Shopify", "Etsy", "TikTok-Shop", "eBay", "Walmart"].find((name) => name.toLowerCase() === platformInput.toLowerCase()) ?? platformInput;
  const market = segment(input.market, "市场代码").toUpperCase();
  if (!/^[A-Z]{2}$/.test(market) || platform === "待定") throw new Error("批量模式需要确定的平台和两位市场代码，例如 Temu / US");
  if (typeof input.sample_directory !== "string" || !input.sample_directory.trim()) throw new Error("必须提供 sample_directory 示例图案目录");
  const sellerService = SERVICE[input.seller_service];
  if (!sellerService) throw new Error("整批 seller_service 尚未确认，请选择 不提供 / 必选 / 可选");
  const customizationType = sellerService === "not_offered" ? null : TYPE[input.customization_type];
  if (sellerService !== "not_offered" && !customizationType) throw new Error("整批 customization_type 尚未确认，请选择 图文 / 图 / 文");
  const sampleUsage = input.sample_usage ?? (customizationType === "image_text" || customizationType === "image" ? "customization_image" : "creative_reference_only");
  if (!["customization_image", "fixed_product_artwork", "creative_reference_only"].includes(sampleUsage)) throw new Error("sample_usage 只允许 customization_image / fixed_product_artwork / creative_reference_only");
  if (task.images && customizationType === "text" && input.sample_usage == null) throw new Error("文定制且提供图案目录时，必须先确认图案是固定印花还是仅供创意参考，并填写 sample_usage");
  if (customizationType === "text" && sampleUsage === "customization_image") throw new Error("文定制不能把图案标记为买家上传图片字段样例");
  const preset = input.image_preset ?? (platform === "Temu" ? "standard_5" : "platform_default");
  if (platform === "Temu" && !["standard_5", "extended_9"].includes(preset)) throw new Error("Temu image_preset 只允许 standard_5 或 extended_9");
  if (platform !== "Temu" && preset !== "platform_default") throw new Error("其他平台请按其既有规则设置 platform_default 和 image_count，不套用 Temu 预设");
  const imageCount = platform === "Temu" ? preset === "extended_9" ? 9 : 5 : input.image_count;
  if (task.images && (!Number.isInteger(imageCount) || imageCount < 1)) throw new Error("其他平台的 image_count 必须由 Agent 根据目标平台规则填写正整数");
  if (input.sample_text != null && typeof input.sample_text !== "string") throw new Error("sample_text 必须是文字或 null");
  const timeZone = input.time_zone ?? "Asia/Shanghai";
  new Intl.DateTimeFormat("en", { timeZone }).format(new Date());
  return {
    spu, platform, market, mapping_mode: mappingMode, auto_task: Number(input.auto_task), task, seller_service: sellerService,
    customization_type: customizationType, sample_usage: sampleUsage, image_preset: task.images ? preset : null, image_count: task.images ? imageCount : null,
    preset_selection_source: input.image_preset ? "user" : "default", time_zone: timeZone,
    sample_directory: resolve(input.sample_directory), sample_text: input.sample_text ?? null,
  };
}

async function inventoryImages(directory) {
  if (!(await stat(directory)).isDirectory()) throw new Error(`示例图案目录无效：${directory}`);
  const candidates = [], excluded = [];
  const walk = async (current) => {
    const entries = await readdir(current, { withFileTypes: true });
    for (const entry of entries.sort((a, b) => a.name < b.name ? -1 : 1)) {
      const path = resolve(current, entry.name);
      if (entry.isSymbolicLink() || entry.name.startsWith(".")) {
        excluded.push({ path, reason: "隐藏文件或符号链接，不参与分配" }); continue;
      }
      if (entry.isDirectory()) { await walk(path); continue; }
      if (!entry.isFile() || !IMAGE_EXTENSIONS.has(extname(entry.name).toLowerCase())) {
        excluded.push({ path, reason: "不是支持的 PNG/JPG/JPEG/WebP 图案文件" }); continue;
      }
      candidates.push(path);
    }
  };
  await walk(directory);
  const hashes = new Map(), images = [];
  for (const path of candidates.sort()) {
    try {
      const bytes = await readFile(path);
      const { data, info } = await sharp(bytes, { limitInputPixels: 40000000 }).rotate().toColourspace("srgb").ensureAlpha().raw().toBuffer({ resolveWithObject: true });
      const visualHash = sha256(Buffer.concat([Buffer.from(`${info.width}x${info.height}:`), data]));
      if (hashes.has(visualHash)) {
        excluded.push({ path, reason: `与 ${hashes.get(visualHash)} 像素相同，去重后不重复分配` }); continue;
      }
      hashes.set(visualHash, path);
      images.push({ path, sha256: sha256(bytes), visual_sha256: visualHash, width: info.width, height: info.height });
    } catch (error) { excluded.push({ path, reason: `图片无法解码或尺寸超限：${error.message}` }); }
  }
  return { images, excluded, files: candidates };
}

const batchStamp = (date, timeZone) => {
  const parts = new Intl.DateTimeFormat("en-GB", { timeZone, year: "numeric", month: "2-digit", day: "2-digit", hour: "2-digit", minute: "2-digit", hourCycle: "h23" }).formatToParts(date);
  const values = Object.fromEntries(parts.map((part) => [part.type, part.value]));
  return `${values.year}${values.month}${values.day}${values.hour}${values.minute}`;
};
async function atomicJson(path, data) {
  const temp = `${path}.${randomUUID()}.tmp`;
  try { await writeFile(temp, JSON.stringify(data, null, 2) + "\n", { flag: "wx" }); await rename(temp, path); }
  finally { await unlink(temp).catch((error) => { if (error.code !== "ENOENT") throw error; }); }
}

export async function createBatch(input, { projectRoot = ROOT, outputRoot = resolve(projectRoot, "output"), now = new Date() } = {}) {
  projectRoot = resolve(projectRoot); outputRoot = resolve(outputRoot);
  const config = normalizeBatchConfig(input);
  if (config.platform === "Temu" && config.mapping_mode === "one_to_all_skus") config.output_layout = "pattern_group";
  if (inside(config.sample_directory, outputRoot) || config.sample_directory === outputRoot || inside(outputRoot, config.sample_directory)) throw new Error("示例图案目录必须与 output 分离，避免把生成产物重新当作输入");
  const catalog = await loadProductCatalog(projectRoot);
  const products = catalog.products.filter((product) => normalizeKey(product.data.SPU) === normalizeKey(config.spu));
  if (!products.length) throw new Error(`CSV 中没有 SPU=${config.spu}，批量模式只处理已建档商品`);
  products.sort((a, b) => normalizeKey(a.data.SKU) < normalizeKey(b.data.SKU) ? -1 : 1);
  for (const product of products) segment(product.data.SKU.trim(), "标准 SKU");
  const inventory = await inventoryImages(config.sample_directory);
  if (!inventory.images.length) throw new Error("示例目录中没有可解码且去重后的有效图片，请补充图案");
  const shared = config.mapping_mode === "one_to_all_skus";
  const directoryNames = [], usedNames = new Set();
  for (const sample of inventory.images) {
    const index = directoryNames.length;
    const name = patternName(sample), id = `pattern-${String(index + 1).padStart(3, "0")}`;
    const duplicate = inventory.images.filter((image) => normalizeKey(patternName(image)) === normalizeKey(name)).length > 1;
    const base = duplicate ? `${name}-${id}` : name;
    let candidate = base, suffix = 2;
    while (usedNames.has(normalizeKey(candidate))) candidate = `${base}-${suffix++}`;
    usedNames.add(normalizeKey(candidate)); directoryNames.push(candidate);
  }
  const count = shared ? products.length * inventory.images.length : Math.min(products.length, inventory.images.length);
  const galleries = new Map();
  const jobs = [];
  for (let i = 0; i < count; i++) {
    const product = products[shared ? i % products.length : i];
    const patternIndex = shared ? Math.floor(i / products.length) : i;
    const gallery = galleries.get(product) ?? await loadGallery(projectRoot, product);
    galleries.set(product, gallery);
    const blockers = config.task.images ? [...gallery.issues] : [];
    if (config.seller_service !== "not_offered" && product.data.产品类型 !== "定制类") blockers.push("商品资料未确认定制能力，不能套用整批定制服务，请确认该 SKU");
    jobs.push({
      id: `job-${String(i + 1).padStart(3, "0")}`, sku: product.data.SKU.trim(), catalog_line: product.line,
      product: product.data, attributes: product.attributes, gallery, sample: inventory.images[patternIndex],
      ...(shared ? { pattern_group_id: `pattern-${String(patternIndex + 1).padStart(3, "0")}` } : {}),
      sample_usage: config.sample_usage,
      status: blockers.length ? "blocked" : "pending", stage: "input", blockers,
      output_dir: null, version: null, lease_token: null, events: [], validation: null,
    });
  }
  await mkdir(outputRoot, { recursive: true });
  const prefix = `批${batchStamp(now, config.time_zone)}`;
  let batchId, statePath, handle;
  for (let sequence = 1; ; sequence++) {
    batchId = sequence === 1 ? prefix : `${prefix}-${String(sequence).padStart(2, "0")}`;
    statePath = resolve(outputRoot, `${batchId}-批次状态.json`);
    if (await stat(resolve(outputRoot, `${batchId}-批次总报告.md`)).then(() => true, (error) => { if (error.code !== "ENOENT") throw error; return false; })) continue;
    try { handle = await open(statePath, "wx"); break; }
    catch (error) { if (error.code !== "EEXIST") throw error; }
  }
  const state = {
    schema_version: 1, batch_id: batchId, project_root: projectRoot, output_root: outputRoot,
    state_path: statePath, report_path: resolve(outputRoot, `${batchId}-批次总报告.md`),
    created_at: now.toISOString(), updated_at: now.toISOString(), config, sources: catalog.sources,
    mapping_policy: shared ? "有效图案稳定排序、像素去重；每张图案 × 全部标准 SKU，组内共用设计与文案，逐 SKU 使用自己的图库" : "标准 SKU 排序；有效图案按相对路径排序、解码像素去重；按位置一对一分配并固定保存，不按文件名推断 SKU",
    ...(shared ? { pattern_groups: inventory.images.map((sample, index) => {
      const id = `pattern-${String(index + 1).padStart(3, "0")}`;
      return { id, sample, ...(config.output_layout === "pattern_group" ? { directory_name: directoryNames[index], output_dir: null, version: null } : {}), skus: products.map((product) => product.data.SKU.trim()), job_ids: jobs.filter((job) => job.pattern_group_id === id).map((job) => job.id), customization_master_id: `${batchId}-${id}-master` };
    }) } : {}),
    inventory: { sku_count: products.length, valid_image_count: inventory.images.length, paired_count: count },
    exclusions: [...inventory.excluded,
      ...(!shared ? inventory.images.slice(count).map((image) => ({ path: image.path, reason: "图片数量多于 SKU，未分配" })) : []),
      ...(!shared ? products.slice(count).map((product) => ({ sku: product.data.SKU.trim(), reason: "SKU 数量多于有效图案，未参与本批次" })) : []),
    ], jobs,
  };
  const creationLock = await open(`${statePath}.lock`, "wx");
  try {
    await creationLock.writeFile(JSON.stringify({ pid: process.pid, started_at: new Date().toISOString() }));
    await handle.writeFile(JSON.stringify(state, null, 2) + "\n");
    await renderBatchReport(state);
  } finally { await handle.close(); await creationLock.close(); await unlink(`${statePath}.lock`); }
  return state;
}

export async function loadBatch(path) {
  const state = JSON.parse(await readFile(resolve(path), "utf8"));
  if (state.schema_version !== 1 || state.state_path !== resolve(path) || !Array.isArray(state.jobs)) throw new Error("批次状态格式错误或已被移动，请使用原批次状态文件");
  if (!inside(state.output_root, state.report_path) || !inside(state.output_root, state.state_path)) throw new Error("批次文件必须位于输出根目录");
  if (![undefined, "one_to_one", "one_to_all_skus"].includes(state.config.mapping_mode)) throw new Error("批次分配模式损坏");
  if (state.config.output_layout !== undefined && (!groupedOutput(state) || !allSkus(state) || state.config.platform !== "Temu")) throw new Error("批次目录模式损坏");
  const skus = new Set(), samples = new Set(), pixels = new Set(), ids = new Set(), pairs = new Set(), directories = new Set();
  for (const job of state.jobs) {
    const pair = JSON.stringify([normalizeKey(job.sku), job.sample.visual_sha256]);
    if (ids.has(job.id) || pairs.has(pair)) throw new Error("批次映射损坏：任务 ID 或 SKU 与图案配对重复");
    ids.add(job.id); pairs.add(pair);
    if (!allSkus(state) && (skus.has(normalizeKey(job.sku)) || samples.has(job.sample.path) || pixels.has(job.sample.visual_sha256))) throw new Error("批次映射损坏：SKU 或图案被重复使用");
    skus.add(normalizeKey(job.sku)); samples.add(job.sample.path); pixels.add(job.sample.visual_sha256);
    if (job.output_dir && (!inside(state.output_root, job.output_dir) || relative(state.output_root, job.output_dir).split(sep).length !== 2)) throw new Error("Listing 目录不符合输出层级");
    if (job.output_dir && directories.has(job.output_dir) && !groupedOutput(state)) throw new Error("批次映射损坏：多个任务共用输出目录");
    if (job.output_dir) directories.add(job.output_dir);
  }
  if (allSkus(state)) {
    const groups = state.pattern_groups;
    if (!Array.isArray(groups) || groups.length !== state.inventory.valid_image_count || new Set(groups.map((group) => group.id)).size !== groups.length
      || new Set(groups.map((group) => group.sample.visual_sha256)).size !== groups.length
      || state.jobs.length !== state.inventory.sku_count * groups.length || state.inventory.paired_count !== state.jobs.length) throw new Error("一张图案应用到全部 SKU 变体的图案组或任务矩阵损坏");
    for (const group of groups) {
      const members = state.jobs.filter((job) => job.pattern_group_id === group.id);
      if (groupedOutput(state)) {
        segment(group.directory_name, "图案目录名称");
        if (group.output_dir && (!inside(state.output_root, group.output_dir) || relative(state.output_root, group.output_dir).split(sep).length !== 2)) throw new Error("图案组目录层级损坏");
        if (groups.some((other) => other.id !== group.id && (normalizeKey(other.directory_name) === normalizeKey(group.directory_name) || (group.output_dir && other.output_dir === group.output_dir)))) throw new Error("不同图案组不得共用目录");
        if (members.some((job) => job.output_dir !== group.output_dir || job.version !== group.version)) throw new Error("图案组成员目录或版本不一致");
      }
      if (group.customization_master_id !== `${state.batch_id}-${group.id}-master` || group.skus.length !== state.inventory.sku_count
        || new Set(group.skus).size !== group.skus.length || members.length !== group.skus.length
        || JSON.stringify(group.skus.map(normalizeKey).sort()) !== JSON.stringify([...skus].sort())
        || JSON.stringify(group.job_ids) !== JSON.stringify(members.map((job) => job.id))
        || members.some((job) => !group.skus.includes(job.sku) || JSON.stringify(job.sample) !== JSON.stringify(group.sample))) throw new Error("一张图案应用到全部 SKU 变体的图案组关联不一致");
    }
  }
  return state;
}

async function withBatchLock(path, action) {
  const lockPath = `${resolve(path)}.lock`;
  let lock;
  try { lock = await open(lockPath, "wx"); }
  catch (error) { if (error.code === "EEXIST") throw new Error("批次正由另一个执行者更新，请稍后重试；中断遗留锁须确认没有执行者后再人工清理"); throw error; }
  try {
    await lock.writeFile(JSON.stringify({ pid: process.pid, started_at: new Date().toISOString() }));
    const state = await loadBatch(path);
    const result = await action(state);
    state.updated_at = new Date().toISOString();
    await atomicJson(state.state_path, state);
    await renderBatchReport(state);
    return result ?? state;
  } finally { await lock.close(); await unlink(lockPath); }
}

export async function checkInputs(state, job) {
  for (const source of [...state.sources, job.sample, ...job.gallery.files]) {
    if (sha256(await readFile(source.path)) !== source.sha256) throw new Error(`输入已改变：${source.path}；不得重新配对或沿用旧商品事实，请确认变化并新建批次`);
  }
  if (state.config.task.images && job.gallery.issues.length) throw new Error(job.gallery.issues.join("；"));
}

export async function claimJob(statePath, sku) {
  return withBatchLock(statePath, async (state) => {
    const job = selectJob(state, sku);
    if (!job) throw new Error("没有可领取的任务，请查看批次总报告");
    if (job.status !== "pending") throw new Error(`SKU ${job.sku} 状态为 ${job.status}，不能重复领取；续做使用已保存 lease_token`);
    if (groupedOutput(state) && state.jobs.some((member) => member.pattern_group_id === job.pattern_group_id && member.status === "allocating")) throw new Error("同组目录分配曾中断，须先人工核对恢复，禁止再次分配");
    try { await checkInputs(state, job); }
    catch (error) {
      job.status = "blocked"; job.blockers = [error.message];
      job.events.push({ at: new Date().toISOString(), stage: "input", message: error.message });
      return { batch_id: state.batch_id, config: state.config, job };
    }
    job.lease_token = randomUUID();
    if (!job.output_dir) {
      // Journal before reserving: an interrupted allocation is never silently repeated.
      job.status = "allocating";
      await atomicJson(state.state_path, state);
      const group = groupedOutput(state) ? groupFor(state, job) : null;
      const allocation = group?.output_dir ? { path: group.output_dir, version: group.version }
        : await resolveOutputDirectory(resolve(state.output_root, `${state.config.platform}-${state.config.market}`, `${state.batch_id}-${group ? group.directory_name : `${job.pattern_group_id ? `${job.pattern_group_id}-` : ""}${job.sku}`}`), { root: state.output_root });
      if (group) {
        group.output_dir = allocation.path; group.version = allocation.version;
        for (const member of state.jobs.filter((item) => item.pattern_group_id === group.id)) { member.output_dir = allocation.path; member.version = allocation.version; }
      } else { job.output_dir = allocation.path; job.version = allocation.version; }
      await atomicJson(state.state_path, state);
    }
    await writeListingAssociation(state, job);
    job.status = "running";
    job.events.push({ at: new Date().toISOString(), stage: job.stage, message: "已领取，后续步骤复用本目录和配对" });
    return { batch_id: state.batch_id, config: state.config, job, ...(allSkus(state) ? { pattern_group: groupFor(state, job) } : {}) };
  });
}

async function writeListingAssociation(state, job) {
  const metadata = {
    batch_id: state.batch_id, job_id: job.id, spu: state.config.spu, sku: job.sku,
    batch_report: slash(relative(job.output_dir, state.report_path)), batch_state: slash(relative(job.output_dir, state.state_path)),
    source_pattern: job.sample.path, source_pattern_sha256: job.sample.sha256, sample_usage: job.sample_usage,
    listing_relationship: allSkus(state) ? "shared_pattern_variants" : "independent", version: job.version,
    ...(allSkus(state) ? { mapping_mode: state.config.mapping_mode, pattern_group_id: job.pattern_group_id, variant_skus: groupFor(state, job).skus, customization_master_id: groupFor(state, job).customization_master_id } : {}),
  };
  if (groupedOutput(state)) {
    delete metadata.sku; delete metadata.job_id;
    metadata.output_layout = "pattern_group";
    metadata.job_ids = groupFor(state, job).job_ids;
  }
  await atomicJson(resolve(job.output_dir, "批次关联.json"), metadata);
}

export async function renderBatchReport(state) {
  const report = state.report_path;
  const counts = Object.fromEntries(["pending", "allocating", "running", "blocked", "failed", "completed"].map((status) => [status, state.jobs.filter((job) => job.status === status).length]));
  const labels = { pending: "待执行", allocating: "目录分配中断待核对", running: "执行中", blocked: "阻塞", failed: "失败", completed: "已完成并校验" };
  const rows = state.jobs.map((job) => {
    const directory = job.output_dir ? link(report, job.output_dir, basename(job.output_dir)) : "尚未分配";
    const files = job.output_dir ? ["商品资料.md", "文案/文案资产.md", "图片/图片验收报告.md", "视频/视频脚本.md", "上架/完整生产报告.md"] : [];
    return { job, directory, files };
  });
  const table = [];
  for (const { job, directory, files } of rows) {
    const saved = [];
    for (const file of files) if (await stat(resolve(job.output_dir, file)).then((value) => value.isFile(), () => false)) saved.push(link(report, resolve(job.output_dir, file), file));
    const usageLabel = { customization_image: "买家上传图片字段的展示样例", fixed_product_artwork: "固定印花图案，买家无需上传", creative_reference_only: "仅供创意参考，未作为固定印花或买家输入" }[job.sample_usage] ?? "待确认";
    table.push(`| ${md(job.id)} / ${md(job.pattern_group_id ?? "独立配对")} / ${md(job.sku)} | ${link(report, job.sample.path, basename(job.sample.path))} | ${usageLabel} | ${labels[job.status] ?? job.status} / ${job.stage} | ${directory} | ${saved.join("、") || "尚无产物"} | ${md(job.blockers.join("；")) || "无"} |`);
  }
  const text = `# ${state.batch_id} 批次总报告\n\n- SPU：${md(state.config.spu)}\n- 平台市场：${state.config.platform}-${state.config.market}\n- 批量任务：${state.config.auto_task}. ${state.config.task.label}（单个菜单 ${state.config.task.menu}）\n- 创建时间：${state.created_at}（批次号时区：${state.config.time_zone}）\n- 更新时间：${state.updated_at}\n- 批次状态：${link(report, state.state_path, "批次状态 JSON")}\n- 图案分配方式：${MAPPING_LABELS[state.config.mapping_mode ?? "one_to_one"]}（${state.config.mapping_mode ?? "one_to_one"}）\n- 语义：对目录中每张有效去重图案分别执行；任务数为${allSkus(state) ? "图案数 × SKU 数" : "图案数和 SKU 数较小的一方"}\n- 处理规则：${state.mapping_policy}\n- 整批选择：服务 ${state.config.seller_service}；类型 ${state.config.customization_type ?? "不适用"}；图片预设 ${state.config.image_preset ?? "本轮未执行"}\n- 数量：SKU ${state.inventory.sku_count}；有效去重图案 ${state.inventory.valid_image_count}；配对 ${state.inventory.paired_count}\n- 完成 ${counts.completed}｜待执行 ${counts.pending}｜执行中 ${counts.running}｜阻塞 ${counts.blocked}｜失败 ${counts.failed}｜目录待核对 ${counts.allocating}\n\n${allSkus(state) ? `## 共享图案组\n\n| 图案组 / 源图案名称 | 全部 SKU 变体 | 展开任务数 | 共享母版 | 组目录 |\n| --- | --- | --- | --- | --- |\n${state.pattern_groups.map((group) => `| ${md(group.id)} / ${md(basename(group.sample.path))} | ${md(group.skus.join("、"))} | ${group.job_ids.length} | ${md(group.customization_master_id)} | ${group.output_dir ? link(report, group.output_dir, basename(group.output_dir)) : "尚未分配或旧布局按SKU目录"} |`).join("\n")}\n\n同组图案和文案一致，只替换各 SKU 对应的商品外观与颜色。\n\n` : ""}## SKU 与图案分配及产物\n\n| 任务 ID / 图案组 / SKU | 分配图案 | 使用方式 | 状态 / 阶段 | Listing 目录 | 已保存产物 | 阻塞原因 |\n| --- | --- | --- | --- | --- | --- | --- |\n${table.join("\n")}\n\n## 数量差异和排除项\n\n| SKU / 文件 | 原因 |\n| --- | --- |\n${state.exclusions.length ? state.exclusions.map((item) => `| ${md(item.sku ?? item.path)} | ${md(item.reason)} |`).join("\n") : "| 无 | 无 |"}\n\n## 执行记录\n\n${state.jobs.map((job) => `### ${job.id} / ${job.sku}\n\n${job.events.length ? job.events.map((event) => `- ${event.at} / ${event.stage}：${md(event.message)}`).join("\n") : "- 尚未执行"}`).join("\n\n")}\n\n## 下一步\n\n${counts.pending ? "继续按任务 ID 领取待执行任务，依次完成所选范围。" : "查看产物，或根据阻塞项补充资料后在原目录续做。"} ${counts.allocating ? "目录分配曾中断，须核对已预留目录与批次关联后恢复，禁止再次分配目录。" : ""}\n`;
  const temp = `${report}.${randomUUID()}.tmp`;
  try { await writeFile(temp, text, { flag: "wx" }); await rename(temp, report); }
  finally { await unlink(temp).catch((error) => { if (error.code !== "ENOENT") throw error; }); }
}

export async function recordJob(statePath, sku, token, update) {
  return withBatchLock(statePath, async (state) => {
    const job = selectJob(state, sku);
    if (!job || !token || job.lease_token !== token) throw new Error("任务令牌不匹配，不能更新其他执行者的任务");
    if (!["running", "completed"].includes(job.status)) throw new Error(`任务当前为 ${job.status}，请先恢复任务`);
    const stage = update.stage ?? job.stage;
    if (!["input", "research", "creative", "product_profile", "production", "validation"].includes(stage)) throw new Error("未知生产阶段");
    const status = update.status ?? "running";
    if (!["running", "blocked", "failed", "completed"].includes(status)) throw new Error("未知任务状态");
    if (!update.message?.trim()) throw new Error("记录步骤必须提供 message 说明实际处理结果");
    if (["blocked", "failed"].includes(status) && (!Array.isArray(update.blockers) || !update.blockers.length)) throw new Error("阻塞或失败时必须提供 blockers");
    if (status === "completed") {
      try { await validateBatchListing(job.output_dir, { runListingGate: true }); }
      catch (error) {
        job.status = "failed"; job.stage = "validation"; job.validation = { passed: false, checked_at: new Date().toISOString() };
        job.blockers = [error.message];
        job.events.push({ at: new Date().toISOString(), stage: "validation", message: `交付校验失败：${error.message}` });
        return job;
      }
      job.validation = { passed: true, checked_at: new Date().toISOString() };
    } else job.validation = null;
    job.status = status; job.stage = stage;
    job.blockers = update.blockers ?? [];
    job.events.push({ at: new Date().toISOString(), stage, message: update.message });
    return job;
  });
}

export async function retryJob(statePath, sku) {
  return withBatchLock(statePath, async (state) => {
    const job = selectJob(state, sku);
    if (!job || !["blocked", "failed", "allocating"].includes(job.status)) throw new Error("只能恢复阻塞、失败或目录分配中断的任务；执行中的任务用原令牌续做");
    if (job.status === "allocating" && !job.output_dir) throw new Error("分配记录缺少目录，请先核对磁盘预留目录并人工恢复记录，禁止再次调用版本分配器");
    const gallery = await loadGallery(state.project_root, { data: job.product, line: job.catalog_line });
    // A supplied missing angle/main may be added; recorded source files may never change silently.
    for (const source of [...state.sources, job.sample, ...job.gallery.files]) if (sha256(await readFile(source.path)) !== source.sha256) throw new Error(`已确认输入发生变化：${source.path}，需要新批次`);
    if (state.config.task.images && gallery.issues.length) throw new Error(gallery.issues.join("；"));
    if (state.config.seller_service !== "not_offered" && job.product.产品类型 !== "定制类") throw new Error("定制能力仍未确认，请修正商品事实并建立新批次");
    job.gallery = gallery; job.status = "pending"; job.blockers = []; job.lease_token = null; job.validation = null;
    job.events.push({ at: new Date().toISOString(), stage: job.stage, message: "资料补齐后恢复，沿用原配对和已分配目录" });
    return job;
  });
}

export async function validateBatchListing(listingDirectory, { runListingGate = false } = {}) {
  const directory = resolve(listingDirectory);
  const association = JSON.parse(await readFile(resolve(directory, "批次关联.json"), "utf8"));
  const state = await loadBatch(resolve(directory, association.batch_state));
  const group = groupedOutput(state) ? state.pattern_groups.find((item) => item.id === association.pattern_group_id) : null;
  if (groupedOutput(state) && !group) throw new Error("Listing 共享图案组关联不一致");
  const job = state.jobs.find((item) => item.id === (group ? group.job_ids[0] : association.job_id));
  if (!job || job.output_dir !== directory || (!group && job.sku !== association.sku) || state.batch_id !== association.batch_id
    || association.source_pattern_sha256 !== job.sample.sha256 || association.source_pattern !== job.sample.path
    || resolve(directory, association.batch_report) !== state.report_path) throw new Error("Listing 批次关联与原始分配记录不一致");
  await checkInputs(state, job);
  if (groupedOutput(state)) {
    if (!group || association.output_layout !== "pattern_group" || JSON.stringify(association.job_ids) !== JSON.stringify(group.job_ids)) throw new Error("图案组任务关联不一致");
    for (const member of state.jobs.filter((item) => item.pattern_group_id === group.id)) await checkInputs(state, member);
  }
  if (allSkus(state)) {
    const group = groupFor(state, job);
    if (association.mapping_mode !== state.config.mapping_mode || association.pattern_group_id !== group.id
      || association.listing_relationship !== "shared_pattern_variants" || association.customization_master_id !== group.customization_master_id
      || JSON.stringify(association.variant_skus) !== JSON.stringify(group.skus)) throw new Error("Listing 共享图案组关联不一致");
  }
  await stat(state.report_path);
  const required = ["商品资料.md", "上架/完整生产报告.md", "上架/next-action.json"];
  if (state.config.task.copy) required.push("文案/文案资产.md", "文案/关键词清单.md");
  if (state.config.task.images) required.push("图片/套图脚本.md", "图片/图片验收报告.md", "图片/套图配置.json");
  if (state.config.task.video_script) required.push("视频/视频脚本.md");
  for (const file of required) {
    const info = await stat(resolve(directory, file));
    if (!info.isFile() || !info.size) throw new Error(`本轮必需产物缺失或为空：${file}`);
  }
  const fullReport = await readFile(resolve(directory, "上架/完整生产报告.md"), "utf8");
  for (let section = 1; section <= 16; section++) if (!new RegExp(`^## ${section}\\. `, "m").test(fullReport)) throw new Error(`完整生产报告缺少第 ${section} 章`);
  if (state.config.task.copy && state.config.platform === "Temu"
    && !/^#{2,4}\s+长描述\s*$/mu.test(await readFile(resolve(directory, "文案/文案资产.md"), "utf8"))) throw new Error("Temu 文案资产缺少长描述栏目");
  const next = JSON.parse(await readFile(resolve(directory, "上架/next-action.json"), "utf8"));
  if (String(next.menu) !== String(state.config.task.menu)) throw new Error("批量菜单必须映射到单个菜单编号后校验下一步范围");
  if (next.batch_id !== state.batch_id || resolve(directory, "上架", next.batch_report ?? "") !== state.report_path) throw new Error("下一步动作记录缺少正确批次引用");
  if (state.config.task.images) {
    const imageDirectory = resolve(directory, "图片");
    const imageFiles = (await readdir(imageDirectory)).filter((name) => /\.(?:png|jpe?g)$/i.test(name));
    const config = JSON.parse(await readFile(resolve(imageDirectory, "套图配置.json"), "utf8"));
    const selectedCount = state.config.image_count;
    if (groupedOutput(state)) {
      const contract = patternGroupImages(config);
      if (config.customization_enabled !== (state.config.seller_service !== "not_offered") || config.selected_count !== selectedCount
        || JSON.stringify(config.variant_skus) !== JSON.stringify(group.skus)
        || config.variants?.some((variant) => !state.jobs.some((member) => member.id === variant.job_id && member.sku === variant.sku && member.pattern_group_id === group.id))) throw new Error("图案组 SKU / 服务 / 任务与批次不一致");
      if (contract.errors.length || imageFiles.length !== contract.expected.length || contract.expected.some((name) => !imageFiles.includes(name))) throw new Error(`图片数量与图案组预设和额外必传图不一致：${contract.errors.join("；")}`);
    } else if (config.selected_count !== selectedCount || imageFiles.length !== selectedCount + (config.additional_required?.length ?? 0)) throw new Error("图片数量与整批预设和额外必传图不一致");
    if (state.config.seller_service !== "not_offered") {
      if (typeof config.customization_master_id !== "string" || !config.customization_master_id.trim()) throw new Error("定制 Listing 套图配置缺少独立母版 ID");
      for (const name of imageFiles) {
        const sidecar = JSON.parse(await readFile(resolve(imageDirectory, name.replace(/\.[^.]+$/, ".verify.json")), "utf8"));
        if (sidecar.customization_master_id !== config.customization_master_id) throw new Error(`本 Listing 图片母版不一致：${name}`);
      }
    }
  }
  const walk = async (current) => {
    for (const entry of await readdir(current, { withFileTypes: true })) {
      const path = resolve(current, entry.name);
      if (entry.isSymbolicLink()) throw new Error(`批次产物不允许符号链接：${path}`);
      if (entry.isDirectory()) { await walk(path); continue; }
      if (entry.name.endsWith(".md")) {
        const content = await readFile(path, "utf8");
        const reportLink = link(path, state.report_path, "批次总报告");
        if (!content.includes(state.batch_id) || !content.includes(reportLink)) throw new Error(`文档缺少批次号或正确的相对报告链接：${path}；应引用 ${reportLink}`);
      }
      if (entry.name === "套图配置.json" || entry.name.endsWith(".verify.json")) {
        const meta = JSON.parse(await readFile(path, "utf8"));
        if (groupedOutput(state)) {
          const config = JSON.parse(await readFile(resolve(directory, "图片/套图配置.json"), "utf8"));
          const contract = patternGroupImages(config);
          const owner = entry.name === "套图配置.json" ? null : contract.owners.get(contract.expected.find((name) => name.replace(/\.[^.]+$/, ".verify.json") === entry.name));
          if (owner) {
            const member = state.jobs.find((item) => item.pattern_group_id === group.id && item.sku === owner);
            if (meta.sku !== owner || meta.job_id !== member.id) throw new Error(`SKU 专属图片任务关联不一致：${path}`);
          } else if (JSON.stringify(meta.variant_skus) !== JSON.stringify(group.skus) || meta.sku != null || meta.job_id != null) throw new Error(`共享图片或配置须关联全部 SKU，不得归属单个任务：${path}`);
        }
        if (meta.batch_id !== state.batch_id || (!groupedOutput(state) && meta.sku !== job.sku) || meta.source_pattern_sha256 !== job.sample.sha256) throw new Error(`图片记录的批次 / SKU / 图案来源不一致：${path}`);
        if (allSkus(state) && ((!groupedOutput(state) && meta.job_id !== job.id) || meta.pattern_group_id !== job.pattern_group_id
          || meta.customization_master_id !== groupFor(state, job).customization_master_id)) throw new Error(`图片记录的共享图案组 / 母版不一致：${path}`);
        if (resolve(dirname(path), meta.batch_report ?? "") !== state.report_path) throw new Error(`图片记录缺少正确批次报告引用：${path}`);
        if (meta.source_pattern_usage !== job.sample_usage) throw new Error(`图片记录未区分定制图案与创意参考：${path}`);
        if (entry.name === "套图配置.json" && meta.preset !== state.config.image_preset) throw new Error("套图预设与整批确认不一致");
        if (!groupedOutput(state) && entry.name === "套图配置.json" && meta.variants) throw new Error("批次子任务仅验收当前 SKU 的套图，禁止混入其他 SKU 的套图配置");
      }
    }
  };
  await walk(directory);
  if (runListingGate) {
    const result = spawnSync(process.execPath, [resolve(state.project_root, "scripts/validate-listing.mjs"), directory, "-p", state.config.platform], { cwd: state.project_root, encoding: "utf8" });
    if (result.error || result.status !== 0) throw new Error(`Listing 交付闸门未通过：${result.error?.message ?? ""}\n${result.stdout ?? ""}\n${result.stderr ?? ""}`);
  }
  return { batch_id: state.batch_id, sku: job.sku, checked: true };
}

async function runCli() {
  const [command, target, sku, token, updatePath] = process.argv.slice(2);
  try {
    let result;
    if (command === "create" && target) result = await createBatch(JSON.parse(await readFile(resolve(target), "utf8")));
    else if (command === "status" && target) { result = await loadBatch(target); await withBatchLock(target, () => undefined); }
    else if (command === "claim" && target) result = await claimJob(target, sku);
    else if (command === "record" && target && sku && token && updatePath) result = await recordJob(target, sku, token, JSON.parse(await readFile(resolve(updatePath), "utf8")));
    else if (command === "retry" && target && sku) result = await retryJob(target, sku);
    else if (command === "check-listing" && target) result = await validateBatchListing(target);
    else throw new Error("用法：batch-production.mjs create <配置.json> | status <状态.json> | claim <状态.json> [任务ID或唯一SKU] | record <状态.json> <任务ID或唯一SKU> <令牌> <步骤.json> | retry <状态.json> <任务ID或唯一SKU> | check-listing <Listing目录>");
    process.stdout.write(JSON.stringify(result, null, 2) + "\n");
    if (command === "record" && result.status === "failed" || command === "claim" && result.job?.status === "blocked") process.exitCode = 1;
  } catch (error) { process.stderr.write(`批次处理失败：${error.message}\n`); process.exitCode = 1; }
}
if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) await runCli();
