// Temu 批量共享图案的文件合同；单个及旧批次继续使用原预设。
const sharedRoles = ["02-日常使用场景图", "03-礼赠场景图", "04-卖点集合图", "05-尺寸规格图"];
const stem = (name) => typeof name === "string" ? name.replace(/\.(png|jpe?g)$/i, "") : "";
const file = (name) => typeof name === "string" && /^[^\\/]+\.(png|jpe?g)$/i.test(name);
const same = (a, b) => JSON.stringify(a) === JSON.stringify(b);

export function patternGroupImages(config) {
  const errors = [], owners = new Map();
  const expected = [], galleries = new Map();
  const fail = (message) => errors.push(`图案组套图配置错误：${message}`);
  if (config.platform !== "Temu" || config.mapping_mode !== "one_to_all_skus" || config.output_layout !== "pattern_group") fail("仅适用 Temu 批量 one_to_all_skus / pattern_group");
  for (const field of ["batch_id", "pattern_group_id", "source_pattern_sha256"]) if (typeof config[field] !== "string" || !config[field].trim()) fail(`缺少 ${field}`);
  const count = config.preset === "standard_5" ? 5 : config.preset === "extended_9" ? 9 : null;
  if (!count || config.selected_count !== count) fail("预设及 selected_count 不一致");
  if (typeof config.customization_enabled !== "boolean") fail("customization_enabled 必须明确记录");
  if (typeof config.customization_master_id !== "string" || !config.customization_master_id.trim()) fail("缺少图案组母版 ID");
  if (!["default", "user"].includes(config.selection_source) || (count === 9 && config.selection_source !== "user")) fail("扩展预设须由用户选择");
  const shared = Array.isArray(config.shared_images) ? config.shared_images : [];
  if (shared.length !== 4 || shared.some((name, i) => !file(name) || stem(name) !== sharedRoles[i])) fail("shared_images 须按顺序登记 02/03/04/05 四张共用副图");
  const supplements = Array.isArray(config.supplemental_images) ? config.supplemental_images : [];
  if (!Array.isArray(config.supplemental_images) || supplements.length !== (count === 9 ? 4 : 0)
    || supplements.some((name, i) => !file(name) || !stem(name).startsWith(`${String(i + 7).padStart(2, "0")}-`))) fail("扩展副图使用 07 至 10，06 保留给定制示意图");
  expected.push(...shared, ...supplements);
  const variants = Array.isArray(config.variants) ? config.variants : [];
  const skus = variants.map((variant) => variant?.sku);
  if (!variants.length || skus.some((sku) => typeof sku !== "string" || !sku.trim())
    || new Set(skus.map((sku) => String(sku).toUpperCase())).size !== skus.length
    || !same(skus, [...skus].sort((a, b) => String(a).toUpperCase() < String(b).toUpperCase() ? -1 : 1)) || !same(skus, config.variant_skus)) fail("variants 须按标准 SKU 排序并与 variant_skus 一致");
  const jobIds = variants.map((variant) => variant?.job_id);
  if (jobIds.some((id) => typeof id !== "string" || !id.trim()) || new Set(jobIds).size !== jobIds.length) fail("每个 SKU 须有唯一的 job_id");
  for (const variant of variants) {
    const sku = variant?.sku;
    const main = variant?.main_image;
    if (!file(main) || stem(main) !== `01-${sku}-主图`) fail(`SKU ${sku} 主图须为 01-${sku}-主图`);
    const extras = Array.isArray(variant?.additional_required) ? variant.additional_required : [];
    if (!Array.isArray(variant?.additional_required) || extras.some((name) => !file(name))) fail(`SKU ${sku} additional_required 须为图片文件名数组`);
    const demos = extras.filter((name) => stem(name) === `06-${sku}-定制示意图`);
    if (config.customization_enabled ? demos.length !== 1 : extras.some((name) => /定制示意图/.test(name))) fail(`SKU ${sku} 定制示意图与服务状态不一致`);
    for (const name of extras) if (stem(name) !== `06-${sku}-定制示意图` && !stem(name).includes(`-${sku}-`)) fail(`SKU ${sku} 额外必传图须带 SKU`);
    const gallery = [main, ...shared, ...supplements];
    if (gallery.length !== count) fail(`SKU ${sku} 轮播图数量不符`);
    if (gallery.length + extras.length > 10) fail(`SKU ${sku} 图数超过平台每变体上限`);
    galleries.set(sku, gallery);
    for (const name of [main, ...extras]) { expected.push(name); owners.set(name, sku); }
  }
  if (new Set(expected).size !== expected.length || new Set(expected.map(stem)).size !== expected.length) fail("图片名或角色重复登记");
  return { errors, expected, owners, galleries, shared: [...shared, ...supplements] };
}
