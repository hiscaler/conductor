import { readdir, readFile, stat } from "node:fs/promises";
import { readFileSync } from "node:fs";
import { resolve, extname, join } from "node:path";
import { pathToFileURL } from "node:url";
import { patternGroupImages } from "./pattern-group-images.mjs";

// 平台图集预设：单一事实来源在 scripts/image-sets.json（与 output/{平台}-{市场} 路径首段一致）。
// Temu 的所选预设和类目额外必传图记录在 图片/套图配置.json。未知平台跳过结构门。
const IMAGE_SETS = (() => {
  try {
    return JSON.parse(readFileSync(new URL("./image-sets.json", import.meta.url), "utf8"));
  } catch {
    return {};
  }
})();

// 从图片目录路径推断平台：output/<Platform>-<Market>/... 的首段即平台标识。
function detectPlatform(dirPath, override) {
  if (override) return override;
  const parts = resolve(dirPath).split(/[\\/]/);
  for (const seg of parts) {
    const m = seg.match(/^([A-Za-z][A-Za-z0-9]*)-[A-Za-z][A-Za-z0-9]*$/);
    if (m && IMAGE_SETS[m[1]]) return m[1];
  }
  return null;
}

// 统一允许上传格式：PNG/JPG 被所有平台接受，故输出统一限定为这两种（用户决策 2026-10-05）。
// 其他图像格式（webp/svg/bmp/gif 等）一律判违规，上传前转 PNG/JPG。不再按平台区分。
const ALLOWED_FORMATS = [".png", ".jpg", ".jpeg"];
// 图像类扩展名（含矢量/非常规）：格式门只校验“长得像图”的文件；文档/侧车(.md/.json/.verify.json)忽略。
const IMAGE_EXT = new Set([".png", ".jpg", ".jpeg", ".webp", ".svg", ".bmp", ".gif", ".tiff", ".tif", ".heic", ".avif"]);

function getDimensions(buf, ext) {
  if (ext === ".png") {
    if (buf.length < 24) return null;
    // PNG: IHDR width @16, height @20 (big-endian uint32)
    return { width: buf.readUInt32BE(16), height: buf.readUInt32BE(20) };
  }
  if (ext === ".jpg" || ext === ".jpeg") {
    // JPEG: scan for SOF0..SOF15 (exclude DHT/0xC4, JPG/0xC8, DAC/0xCC)
    let i = 2;
    while (i + 9 < buf.length) {
      if (buf[i] !== 0xff) {
        i += 1;
        continue;
      }
      const marker = buf[i + 1];
      if (marker === 0xc0 || (marker >= 0xc1 && marker <= 0xcf && marker !== 0xc4 && marker !== 0xc8 && marker !== 0xcc)) {
        const height = buf.readUInt16BE(i + 5);
        const width = buf.readUInt16BE(i + 7);
        return { width, height };
      }
      if (marker === 0xd9 || marker === 0xda) break;
      const len = buf.readUInt16BE(i + 2);
      i += 2 + len;
    }
    return null;
  }
  return null;
}

const SIZE_COMPARISON_RE = /(尺寸对比图|size[-_\s]?(comparison|reference))/i;
const MAX_SCALE_ERROR_PCT = 3;

function isPositiveNumber(value) {
  return typeof value === "number" && Number.isFinite(value) && value > 0;
}

function isHttpUrl(value) {
  if (typeof value !== "string" || !value.trim()) return false;
  try {
    const url = new URL(value);
    return url.protocol === "https:" || url.protocol === "http:";
  } catch {
    return false;
  }
}

function validateSizeComparisonEvidence(name, data) {
  const errors = [];
  const reference = data?.reference_object;
  const scale = data?.scale_check;
  const prefix = `尺寸对比证据缺失/无效：${name}`;

  if (!reference || typeof reference !== "object") {
    errors.push(`${prefix} 缺少 reference_object`);
  } else {
    for (const field of ["name", "model", "source_title"]) {
      if (typeof reference[field] !== "string" || !reference[field].trim()) {
        errors.push(`${prefix} reference_object.${field} 必须为非空字符串`);
      }
    }
    for (const field of ["source_url", "visual_source_url"]) {
      if (!isHttpUrl(reference[field])) {
        errors.push(`${prefix} reference_object.${field} 必须是可打开核验的 HTTP(S) 来源链接`);
      }
    }
    if (typeof reference.source_checked_at !== "string" || !/^\d{4}-\d{2}-\d{2}$/.test(reference.source_checked_at)) {
      errors.push(`${prefix} reference_object.source_checked_at 必须为 YYYY-MM-DD`);
    }
    for (const field of ["identity_verified", "dimensions_verified", "appearance_verified"]) {
      if (reference[field] !== true) {
        errors.push(`${prefix} reference_object.${field} 必须在人工核验后设为 true`);
      }
    }
    const dims = reference.dimensions_cm;
    const hasWidth = isPositiveNumber(dims?.width_cm);
    const hasDiameter = isPositiveNumber(dims?.diameter_cm);
    if (!isPositiveNumber(dims?.height_cm) || (!hasWidth && !hasDiameter)) {
      errors.push(`${prefix} reference_object.dimensions_cm 须含正数 height_cm 及 width_cm 或 diameter_cm`);
    }
  }

  if (!scale || typeof scale !== "object") {
    errors.push(`${prefix} 缺少 scale_check`);
    return errors;
  }

  if (!["height", "width", "diameter"].includes(scale.axis)) {
    errors.push(`${prefix} scale_check.axis 仅允许 height、width、diameter`);
  }
  for (const field of ["product_cm", "reference_cm", "product_px", "reference_px"]) {
    if (!isPositiveNumber(scale[field])) errors.push(`${prefix} scale_check.${field} 必须为正数`);
  }
  if (scale.measured_in_final_image !== true) errors.push(`${prefix} 必须记录在最终落盘图片中测量像素边界`);
  if (scale.labels_verified !== true) errors.push(`${prefix} 必须人工核实双方名称、尺寸标签`);
  if (scale.dual_units_verified !== true) errors.push(`${prefix} 必须人工核实尺寸图中的双单位`);
  if (!isPositiveNumber(scale.tolerance_pct) || scale.tolerance_pct > MAX_SCALE_ERROR_PCT) {
    errors.push(`${prefix} scale_check.tolerance_pct 必须大于 0 且不超过 ${MAX_SCALE_ERROR_PCT}%`);
  }

  if (reference && scale.axis && isPositiveNumber(scale.reference_cm)) {
    const dims = reference.dimensions_cm || {};
    const sourceSize = scale.axis === "height"
      ? dims.height_cm
      : scale.axis === "diameter"
        ? (dims.diameter_cm ?? dims.width_cm)
        : dims.width_cm;
    if (["height", "width", "diameter"].includes(scale.axis) && !isPositiveNumber(sourceSize)) {
      errors.push(`${prefix} reference_object.dimensions_cm 缺少与 scale_check.axis=${scale.axis} 对应的真实尺寸`);
    } else if (isPositiveNumber(sourceSize) && Math.abs(scale.reference_cm - sourceSize) / sourceSize > 0.001) {
      errors.push(`${prefix} scale_check.reference_cm 与来源记录的 reference_object 尺寸不一致`);
    }
  }

  if ([scale.product_cm, scale.reference_cm, scale.product_px, scale.reference_px, scale.tolerance_pct].every(isPositiveNumber)) {
    const expectedRatio = scale.product_cm / scale.reference_cm;
    const observedRatio = scale.product_px / scale.reference_px;
    const errorPct = Math.abs(observedRatio / expectedRatio - 1) * 100;
    if (errorPct > scale.tolerance_pct) {
      errors.push(`${prefix} 比例误差 ${errorPct.toFixed(2)}% 超过侧车容差 ${scale.tolerance_pct}%（机器重算；默认最大容差 ${MAX_SCALE_ERROR_PCT}%）`);
    }
  }

  return errors;
}

const checkImageSet = async (dirPath, platform) => {
  const root = resolve(dirPath);
  const entries = await readdir(root, { withFileTypes: true });
  const files = entries.filter((e) => e.isFile()).map((e) => e.name);

  const failures = [];
  const warnings = [];
  const skuByImage = new Map();
  let sharedCustomizationMasterId;
  let customizationEnabled = false;
  let patternGroupConfig;
  let sharedGroupImages = [];

  const resolvedPlatform = detectPlatform(root, platform);
  const platformCfg = resolvedPlatform ? IMAGE_SETS[resolvedPlatform] : undefined;
  const platformLabel = resolvedPlatform || "默认(未知平台)";
  if (resolvedPlatform) {
    warnings.push(`平台识别：${resolvedPlatform}（来自路径或 -p 参数；结构门按该平台必备图集校验）`);
  } else {
    warnings.push(`平台未识别：路径中无已知平台段且未用 -p 指定；结构门跳过`);
  }

  // 1) 格式门（统一规则，不区分平台）：所有平台输出必须为 PNG/JPG。
  //    仅校验“图像类扩展名”文件；文档/侧车(.md/.json 等)忽略。
  //    图像扩展名但不属于 ALLOWED_FORMATS → 违规（含 webp / svg / bmp / gif 等），上传前转 PNG/JPG。
  const rasterImages = [];
  for (const name of files) {
    const ext = extname(name).toLowerCase();
    if (!IMAGE_EXT.has(ext)) continue; // 非图像文件（文档/侧车/元数据）不参与图像校验
    if (!ALLOWED_FORMATS.includes(ext)) {
      const why = ext === ".svg" ? "（矢量线框禁止，image-set-rules L75）" : "";
      failures.push(`格式违规：${name} 格式 ${ext} 不在允许上传格式 ${ALLOWED_FORMATS.join("/")} 内${why}；上传前转 PNG/JPG`);
      continue;
    }
    rasterImages.push(name);
  }
  if (rasterImages.length > 0) {
    warnings.push("必须人工逐张打开最终图，对照本 Listing 的分配图案、主图母版、样例短句和 SKU 颜色；本脚本只检查侧车字段，不读取图像像素，机器 PASS 不能证明图案一致。");
  }

  // 2) 几何门：1:1、宽高均 ≥ 800px、单张 < 2MB
  for (const name of rasterImages) {
    const buf = await readFile(join(root, name));
    const dim = getDimensions(buf, extname(name).toLowerCase());
    if (!dim) {
      warnings.push(`几何跳过：无法解析 ${name} 的尺寸（非标准 PNG/JPEG 头）`);
    } else {
      if (dim.width !== dim.height) {
        failures.push(`几何违规：${name} 非 1:1（${dim.width}x${dim.height}），电商平台硬性 1:1`);
      }
      if (Math.min(dim.width, dim.height) < 800) {
        failures.push(`几何违规：${name} 短边 ${Math.min(dim.width, dim.height)}px < 800px 下限`);
      }
    }
    const st = await stat(join(root, name));
    if (st.size > 2 * 1024 * 1024) {
      warnings.push(`体积告警：${name} ${(st.size / 1024 / 1024).toFixed(2)}MB > 2MB，上传可能失败`);
    }
  }

  // 3) 结构门：按平台预设及本 Listing 的显式套图配置校验。
  if (resolvedPlatform && platformCfg?.presets) {
    const configName = "套图配置.json";
    const configPath = join(root, configName);
    if (!files.includes(configName)) {
      failures.push(`结构配置缺失：平台 ${resolvedPlatform} 必须提供 图片/${configName}，记录 standard_5 / extended_9 选择及类目额外必传图`);
    } else {
      let config;
      try {
        config = JSON.parse(await readFile(configPath, "utf8"));
      } catch {
        failures.push(`结构配置损坏：${configName} 不是合法 JSON`);
      }
      if (config) {
        const preset = platformCfg.presets[config.preset];
        if (config.platform !== resolvedPlatform) {
          failures.push(`结构配置错误：${configName} platform=${JSON.stringify(config.platform)}，应为 ${resolvedPlatform}`);
        }
        if (Object.hasOwn(config, "variants") && !Array.isArray(config.variants)) {
          failures.push(`结构配置错误：${configName} variants 必须是数组`);
        }
        if (!preset) {
          failures.push(`结构配置错误：${configName} preset=${JSON.stringify(config.preset)} 不存在`);
        } else if (config.selected_count !== preset.count) {
          failures.push(`结构配置错误：${configName} selected_count 必须等于 ${config.preset} 的 ${preset.count}`);
        }
        if (!["default", "user"].includes(config.selection_source)) {
          failures.push(`结构配置错误：${configName} selection_source 仅允许 default 或 user`);
        }
        if (config.preset === "extended_9" && config.selection_source !== "user") {
          failures.push(`结构配置错误：extended_9 必须是用户明确选择，不能作为默认预设`);
        }
        const isFileNameList = (list) => Array.isArray(list)
          && list.every((name) => typeof name === "string" && name.trim() && name === name.split(/[\\/]/).pop());
        const validateExpectedFiles = (expectedUnfiltered, label) => {
          if (new Set(expectedUnfiltered).size !== expectedUnfiltered.length) {
            failures.push(`结构配置错误：${configName} ${label}中的预设图、扩展图和额外必传图不能重复登记`);
          }
          const expected = [...new Set(expectedUnfiltered)];
          const missing = expected.filter((name) => !files.includes(name));
          if (missing.length > 0) failures.push(`结构缺失：${label}缺少图型 ${missing.join("、")}`);
          return expected;
        };

        if (!preset) {
          // The invalid preset was reported above; avoid dereferencing it below.
        } else if (config.output_layout === "pattern_group") {
          patternGroupConfig = config;
          const contract = patternGroupImages(config);
          failures.push(...contract.errors);
          sharedCustomizationMasterId = config.customization_master_id;
          customizationEnabled = config.customization_enabled === true;
          sharedGroupImages = contract.shared;
          for (const [name, sku] of contract.owners) skuByImage.set(name, sku);
          const expected = validateExpectedFiles(contract.expected, "图案组 ");
          const unexpected = rasterImages.filter((name) => !expected.includes(name));
          if (unexpected.length) failures.push(`结构超出图案组套图：${unexpected.join("、")}`);
          if (rasterImages.length !== expected.length) failures.push(`图案组数量不符：应为 ${expected.length} 张，实际 ${rasterImages.length} 张`);
        } else if (Array.isArray(config.variants)) {
          if (config.variants.length < 2) {
            failures.push(`结构配置错误：${configName} variants 至少需要两个 SKU；单 SKU Listing 不使用变体分组`);
          }
          sharedCustomizationMasterId = config.shared_customization_master_id;
          customizationEnabled = config.customization_enabled === true;
          if (typeof config.customization_enabled !== "boolean") {
            failures.push(`结构配置错误：多 SKU 变体必须用 customization_enabled 明确记录是否启用定制`);
          }
          if (config.customization_enabled === true
            && (typeof sharedCustomizationMasterId !== "string" || !sharedCustomizationMasterId.trim())) {
            failures.push(`结构配置错误：启用定制的多 SKU 变体必须提供非空 shared_customization_master_id`);
          }

          const allExpected = [];
          const seenSkus = new Set();
          const seenNames = new Set();
          const variantNames = [];
          for (const variant of config.variants) {
            const sku = typeof variant?.sku === "string" ? variant.sku.trim() : "";
            if (!sku) {
              failures.push(`结构配置错误：${configName} 每个 variants 项都必须有非空 sku`);
              continue;
            }
            const skuKey = sku.toUpperCase();
            if (seenSkus.has(skuKey)) failures.push(`结构配置错误：variants 中 SKU 重复：${sku}`);
            seenSkus.add(skuKey);
            variantNames.push(sku);

            const supplements = variant.supplemental_images;
            const extras = variant.additional_required;
            if (!isFileNameList(supplements)) {
              failures.push(`结构配置错误：SKU ${sku} 的 supplemental_images 必须是仅含文件名的数组`);
              continue;
            }
            if (!isFileNameList(extras)) {
              failures.push(`结构配置错误：SKU ${sku} 的 additional_required 必须是仅含文件名的数组`);
              continue;
            }
            const requiredSupplementCount = preset.supplemental_count ?? 0;
            if (supplements.length !== requiredSupplementCount) {
              failures.push(`结构配置错误：SKU ${sku} 的 ${config.preset} 需要登记 ${requiredSupplementCount} 张 supplemental_images，当前 ${supplements.length} 张`);
            }
            const prefix = `${sku}_`;
            const names = [...supplements, ...extras];
            for (const name of names) {
              if (!name.startsWith(prefix)) {
                failures.push(`结构配置错误：SKU ${sku} 的图片名必须以 ${prefix} 开头：${name}`);
              }
            }
            const required = preset.required.map((name) => `${prefix}${name}`);
            const expected = validateExpectedFiles([...required, ...supplements, ...extras], `SKU ${sku} `);
            for (const name of expected) {
              if (seenNames.has(name)) failures.push(`结构配置错误：多个变体重复登记图片 ${name}`);
              seenNames.add(name);
              skuByImage.set(name, sku);
            }
            allExpected.push(...expected);
            if (expected.length < (platformCfg.minimum_count ?? 0) || expected.length > (platformCfg.maximum_count ?? Infinity)) {
              failures.push(`数量超出平台范围：SKU ${sku} 配置 ${expected.length} 张，${resolvedPlatform} 允许 ${platformCfg.minimum_count}-${platformCfg.maximum_count} 张/变体`);
            }
          }

          const compareSku = (a, b) => {
            const left = a.toUpperCase();
            const right = b.toUpperCase();
            return left < right ? -1 : left > right ? 1 : 0;
          };
          const sortedNames = [...variantNames].sort(compareSku);
          if (variantNames.some((sku, index) => sku !== sortedNames[index])) {
            failures.push(`结构配置错误：${configName} variants 必须按标准 SKU 排序`);
          }
          const unexpected = rasterImages.filter((name) => !allExpected.includes(name));
          if (unexpected.length > 0) failures.push(`结构超出已登记变体套图：含未登记图片 ${unexpected.join("、")}`);
          if (rasterImages.length !== allExpected.length) {
            failures.push(`数量不符：各 SKU 变体套图应为 ${allExpected.length} 张，实际 ${rasterImages.length} 张`);
          }
          if (allExpected.length > (platformCfg.maximum_count ?? Infinity)) {
            failures.push(`平台总图数冲突：${config.variants.length} 个变体共需 ${allExpected.length} 张，${resolvedPlatform} 单 Listing 上限为 ${platformCfg.maximum_count} 张；暂停生成，请卖家调整变体范围或确认改为独立 Listing`);
          }
        } else {
          const supplements = config.supplemental_images;
          const extras = config.additional_required;
          if (!isFileNameList(supplements)) {
            failures.push(`结构配置错误：${configName} supplemental_images 必须是仅含文件名的数组`);
          }
          if (!isFileNameList(extras)) {
            failures.push(`结构配置错误：${configName} additional_required 必须是仅含文件名的数组`);
          } else if (preset && isFileNameList(supplements)) {
            const requiredSupplementCount = preset.supplemental_count ?? 0;
            if (supplements.length !== requiredSupplementCount) {
              failures.push(`结构配置错误：${configName} 的 ${config.preset} 需要登记 ${requiredSupplementCount} 张 supplemental_images，当前 ${supplements.length} 张`);
            }
            const expected = validateExpectedFiles([...preset.required, ...supplements, ...extras], "");
            const unexpected = rasterImages.filter((name) => !expected.includes(name));
            if (unexpected.length > 0) failures.push(`结构超出所选套图：${config.preset} 含未登记图片 ${unexpected.join("、")}；新增必传图须列入 additional_required`);
            if (rasterImages.length !== expected.length) {
              failures.push(`数量不符：所选 ${config.preset} 与类目额外必传图应为 ${expected.length} 张，实际 ${rasterImages.length} 张`);
            }
            if (expected.length < (platformCfg.minimum_count ?? 0) || expected.length > (platformCfg.maximum_count ?? Infinity)) {
              failures.push(`数量超出平台范围：配置应为 ${expected.length} 张，Temu 允许 ${platformCfg.minimum_count}-${platformCfg.maximum_count} 张`);
            }
          }
        }
      }
    }
  } else if (resolvedPlatform && platformCfg?.required) {
    const missing = platformCfg.required.filter((n) => !files.includes(n));
    if (missing.length > 0) {
      failures.push(`结构缺失：平台 ${resolvedPlatform} 缺少必备图型 ${missing.join("、")}`);
    }
  }

  // 4) 内容门（水印 + 图案一致性）：逐图须有 .verify.json 侧车，且明确通过
  //    脚本不读像素，无法自动判定水印/图案一致性；要求生产者显式确认，杜绝「假验收」。
  for (const name of rasterImages) {
    const base = name.replace(/\.[^.]+$/, "");
    const sidecar = `${base}.verify.json`;
    if (!files.includes(sidecar)) {
      failures.push(`内容未确认：${name} 缺少侧车 ${sidecar}（水印/图案一致性未经确认，验收不通过）`);
      continue;
    }
    let data;
    try {
      data = JSON.parse(await readFile(join(root, sidecar), "utf8"));
    } catch {
      failures.push(`侧车损坏：${sidecar} 不是合法 JSON`);
      continue;
    }
    const expectedSku = skuByImage.get(name);
    if (patternGroupConfig && (data?.pattern_group_id !== patternGroupConfig.pattern_group_id
      || data?.customization_master_id !== patternGroupConfig.customization_master_id)) failures.push(`图案组或母版不一致：${name}`);
    if (patternGroupConfig && sharedGroupImages.includes(name)
      && (JSON.stringify(data?.variant_skus) !== JSON.stringify(patternGroupConfig.variant_skus) || data?.sku != null || data?.job_id != null)) failures.push(`共享副图须记录全部 variant_skus：${name}`);
    if (patternGroupConfig && expectedSku
      && data?.job_id !== patternGroupConfig.variants.find((variant) => variant.sku === expectedSku)?.job_id) failures.push(`SKU 专属图片 job_id 不一致：${name}`);
    if (expectedSku && data?.sku !== expectedSku) {
      failures.push(`SKU 对应关系未确认：${name} 侧车 sku=${JSON.stringify(data?.sku)}，应为 ${expectedSku}`);
    }
    if (expectedSku && customizationEnabled && data?.customization_master_id !== sharedCustomizationMasterId) {
      failures.push(`定制母版不一致：${name} 侧车 customization_master_id=${JSON.stringify(data?.customization_master_id)}，应与同 Listing 变体共享母版 ${JSON.stringify(sharedCustomizationMasterId)}`);
    }
    const wm = data?.watermark;
    if (wm !== "clear" && wm !== "none") {
      failures.push(`水印未清：${name} 侧车 watermark=${JSON.stringify(wm)}（须为 clear/none，image-set-rules L47）`);
    }
    if (data?.pattern_consistent_with_master !== true) {
      failures.push(`图案不一致：${name} 侧车 pattern_consistent_with_master=${JSON.stringify(data?.pattern_consistent_with_master)}（须为 true，L46 母版一致性）`);
    }
    if (SIZE_COMPARISON_RE.test(name)) {
      failures.push(...validateSizeComparisonEvidence(name, data));
      warnings.push(`尺寸对比图 ${name}：机器只检查侧车字段和比例算术；人工仍须打开 reference_object 来源页、来源图片及最终文件核验实物身份、外观和像素测量真实性`);
    }
  }

  return { root, resolvedPlatform, rasterImages, failures, warnings };
};

const runCli = async () => {
  const args = process.argv.slice(2);
  let target = null;
  let platform = null;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "-p" || args[i] === "--platform") {
      platform = args[i + 1];
      i += 1;
    } else if (!target) {
      target = args[i];
    }
  }
  if (!target) {
    process.stderr.write("图像套图检查失败：缺少图片目录参数\n用法：node scripts/check-image-set.mjs <图片目录> [-p <平台>]\n");
    process.exitCode = 1;
    return;
  }
  try {
    const { root, resolvedPlatform, rasterImages, failures, warnings } = await checkImageSet(target, platform);
    for (const w of warnings) process.stdout.write(`⚠ ${w}\n`);
    if (failures.length > 0) {
      process.stdout.write(`\n图像套图检查未通过（${root}），共 ${failures.length} 项：\n`);
      for (const f of failures) process.stdout.write(`- ${f}\n`);
      process.exitCode = 1;
      return;
    }
    process.stdout.write(`图像套图检查通过：${root}（平台=${resolvedPlatform || "未知/跳过结构门"}，${rasterImages.length} 张栅格图；格式/几何/结构/侧车字段及尺寸对比比例算术符合。实物身份和人工测量真实性仍须视觉核验）\n`);
  } catch (error) {
    process.stderr.write(`图像套图检查失败：${error.message}\n`);
    process.exitCode = 1;
  }
};

if (process.argv[1] && pathToFileURL(process.argv[1]).href === import.meta.url) {
  await runCli();
}
