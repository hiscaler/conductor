import { readFile, readdir, stat } from "node:fs/promises";
import { resolve, extname, isAbsolute } from "node:path";
import { createHash } from "node:crypto";

export const CATALOG_HEADERS = [
  "SPU", "SKU", "商品名称", "品牌", "产品类目", "产品子类目", "产品类型", "包含内容", "型号", "颜色/款式",
  "尺码/规格", "长度cm", "宽度cm", "高度cm", "净重g", "材质", "结构/表面工艺", "已确认功能", "商品特点", "适用对象",
  "使用场景", "使用/护理说明", "定制内容", "定制位置", "定制工艺", "包装清单", "包装方式", "包装长度cm", "包装宽度cm", "包装高度cm",
  "包装毛重g", "风险/禁用声明", "资料更新时间", "备注",
];
export const ATTRIBUTE_HEADERS = ["SKU", "属性组", "属性名称", "属性值", "单位", "值类型", "是否平台必需", "适用平台", "资料来源", "资料更新时间", "备注"];
export const IMAGE_EXTENSIONS = new Set([".png", ".jpg", ".jpeg", ".webp"]);
export const sha256 = (bytes) => createHash("sha256").update(bytes).digest("hex");
export const normalizeKey = (value) => String(value).trim().toUpperCase();

// Strict quoted CSV, including multiline fields and UTF-8 BOM. Preserve facts verbatim.
export function parseCsv(text, fileName) {
  const rows = [];
  let values = [], field = "", quoted = false, closed = false, line = 1, rowLine = 1;
  const finish = () => {
    values.push(field);
    if (values.some((value) => value !== "")) rows.push({ line: rowLine, values });
    values = []; field = ""; closed = false;
  };
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    if (quoted) {
      if (c === '"' && text[i + 1] === '"') { field += '"'; i++; }
      else if (c === '"') { quoted = false; closed = true; }
      else { field += c; if (c === "\n") line++; }
    } else if (c === ",") { values.push(field); field = ""; closed = false; }
    else if (c === "\n" || c === "\r") {
      if (c === "\r" && text[i + 1] === "\n") i++;
      finish(); line++; rowLine = line;
    } else if (c === '"' && field === "" && !closed) quoted = true;
    else if (c === '"' || closed) throw new Error(`${fileName} 第 ${line} 行：引号结构错误，请用双引号包围整个字段`);
    else field += c;
  }
  if (quoted) throw new Error(`${fileName} 第 ${rowLine} 行：双引号未闭合`);
  if (field || values.length || closed) finish();
  return rows;
}

async function readTable(path, headers) {
  const bytes = await readFile(path);
  let text;
  try { text = new TextDecoder("utf-8", { fatal: true }).decode(bytes).replace(/^\uFEFF/, ""); }
  catch { throw new Error(`${path}：编码不是有效 UTF-8，请另存为 UTF-8 CSV`); }
  const rows = parseCsv(text, path);
  if (JSON.stringify(rows[0]?.values) !== JSON.stringify(headers)) {
    throw new Error(`${path} 第 1 行：表头不符合固定结构，请按 data/README.md 恢复表头`);
  }
  return { hash: sha256(bytes), rows: rows.slice(1).map((row) => {
    if (row.values.length !== headers.length) throw new Error(`${path} 第 ${row.line} 行：列数为 ${row.values.length}，应为 ${headers.length}，请检查逗号和引号`);
    return { line: row.line, data: Object.fromEntries(headers.map((header, i) => [header, row.values[i]])) };
  }) };
}

const isNumber = (value) => value.trim() !== "" && Number.isFinite(Number(value)) && Number(value) >= 0;
const isDate = (value) => !value || /^\d{4}-\d{2}-\d{2}$/.test(value)
  && !Number.isNaN(Date.parse(`${value}T00:00:00Z`))
  && new Date(`${value}T00:00:00Z`).toISOString().slice(0, 10) === value;
const fail = (file, row, field, fix) => { throw new Error(`${file} 第 ${row.line} 行，${field}=${JSON.stringify(row.data[field])}；${fix}`); };

export async function loadProductCatalog(projectRoot) {
  const catalogPath = resolve(projectRoot, "data/product-catalog.csv");
  const attributesPath = resolve(projectRoot, "data/product-attributes.csv");
  const [catalog, attributes] = await Promise.all([
    readTable(catalogPath, CATALOG_HEADERS), readTable(attributesPath, ATTRIBUTE_HEADERS),
  ]);
  const bySku = new Map();
  for (const row of catalog.rows) {
    for (const field of ["SPU", "SKU", "商品名称", "产品类目", "产品类型"]) {
      if (!row.data[field].trim()) fail(catalogPath, row, field, "请补充必填值");
    }
    for (const field of ["SPU", "SKU"]) {
      if (/[\\/\x00-\x1F]/.test(row.data[field]) || [".", ".."].includes(row.data[field].trim())) fail(catalogPath, row, field, "请去除路径字符");
    }
    const key = normalizeKey(row.data.SKU);
    if (bySku.has(key)) fail(catalogPath, row, "SKU", `与第 ${bySku.get(key).line} 行忽略大小写后重复，请修正 SKU`);
    if (!["普货", "定制类"].includes(row.data.产品类型)) fail(catalogPath, row, "产品类型", "只允许普货、定制类");
    for (const field of ["长度cm", "宽度cm", "高度cm", "净重g", "包装长度cm", "包装宽度cm", "包装高度cm", "包装毛重g"]) {
      if (row.data[field].trim() && !isNumber(row.data[field].trim())) fail(catalogPath, row, field, "请填写不带单位的非负数字");
    }
    if (!isDate(row.data.资料更新时间.trim())) fail(catalogPath, row, "资料更新时间", "请使用有效 yyyy-mm-dd 日期或留空");
    bySku.set(key, { ...row, attributes: [] });
  }
  const keys = new Set();
  for (const row of attributes.rows) {
    for (const field of ["SKU", "属性组", "属性名称", "属性值", "值类型"]) if (!row.data[field].trim()) fail(attributesPath, row, field, "请补充必填值");
    const product = bySku.get(normalizeKey(row.data.SKU));
    if (!product) fail(attributesPath, row, "SKU", "SKU 不存在于商品主表，请修正外键");
    if (!["数字", "文本", "布尔值", "列表"].includes(row.data.值类型)) fail(attributesPath, row, "值类型", "请使用数字、文本、布尔值、列表");
    if (!["", "是", "否"].includes(row.data.是否平台必需)) fail(attributesPath, row, "是否平台必需", "请使用是、否或留空");
    if (row.data.值类型 === "数字" && !isNumber(row.data.属性值.trim())) fail(attributesPath, row, "属性值", "请填写不带单位的非负数字");
    if (!isDate(row.data.资料更新时间.trim())) fail(attributesPath, row, "资料更新时间", "请使用有效 yyyy-mm-dd 日期或留空");
    const source = row.data.资料来源.trim();
    if (source && (isAbsolute(source) || /^(?:data[\\/]|\.[\\/])/.test(source))) {
      await stat(resolve(projectRoot, source)).catch(() => fail(attributesPath, row, "资料来源", "本地证据不存在，请提供真实资料路径"));
    }
    const key = JSON.stringify([normalizeKey(row.data.SKU), row.data.属性组, row.data.属性名称, row.data.适用平台]);
    if (keys.has(key)) fail(attributesPath, row, "属性名称", "同 SKU、属性组、属性名称与平台存在重复项，请修正");
    keys.add(key); product.attributes.push(row.data);
  }
  return { products: [...bySku.values()], sources: [{ path: catalogPath, sha256: catalog.hash }, { path: attributesPath, sha256: attributes.hash }] };
}

export async function loadGallery(projectRoot, product) {
  const directory = resolve(projectRoot, "data/products", product.data.SPU.trim(), product.data.SKU.trim());
  const issues = [], files = [], roles = new Map();
  const entries = await readdir(directory, { withFileTypes: true }).catch((error) => {
    if (error.code !== "ENOENT") throw error;
    issues.push(`第 ${product.line} 行 SKU=${product.data.SKU}：图库目录不存在 ${directory}`); return [];
  });
  for (const entry of entries.sort((a, b) => a.name < b.name ? -1 : 1)) {
    if (!entry.isFile() || !IMAGE_EXTENSIONS.has(extname(entry.name).toLowerCase())) continue;
    const role = entry.name.slice(0, -extname(entry.name).length).toLowerCase();
    if (roles.has(role)) issues.push(`SKU=${product.data.SKU} 图片角色 ${role} 扩展名冲突：${roles.get(role)}、${entry.name}`);
    roles.set(role, entry.name);
    const path = resolve(directory, entry.name);
    files.push({ path, role, sha256: sha256(await readFile(path)) });
  }
  const mains = files.filter((file) => file.role === "main");
  if (mains.length !== 1) issues.push(`第 ${product.line} 行 SKU=${product.data.SKU}：main 主图应为 1 个，实际 ${mains.length} 个，请补图或修正角色冲突`);
  return { directory, files, main: mains.length === 1 ? mains[0].path : null, issues };
}
