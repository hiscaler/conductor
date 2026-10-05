// 数据有效性校验：新鲜度 SLA + 时间戳强制 + 实时/API 交叉偏差 + 合理性（币种/价格）
// 对应 platforms/market-data-sources.md 的「数据有效性校验」一节。
// 用法：node scripts/check-data-freshness.mjs <报告.md> [更多报告.md]
// 退出码 1 = 存在硬违规（缺时间戳 / 超窗且未标注 / 交叉偏差超阈 / 币种错配 / 价格异常）；0 = 通过。

import { readFile } from "node:fs/promises";
import { relative, resolve } from "node:path";

// 新鲜度 SLA（小时）。实时波动信号最严，趋势指数最宽。
const SLA_HOURS = {
  realtime: 48, // 价格 / 排名 / 库存 / Buy Box
  review: 168, // 评论数量 / 评分（7 天）
  trend: 720, // 趋势 / 季节性指数（30 天）
};

// 合理性校验：目标市场 -> 期望币种（与 market-data-sources.md 地区映射一致）
const MARKET_CURRENCY = {
  US: "USD", 美国: "USD",
  GB: "GBP", 英国: "GBP",
  DE: "EUR", 德国: "EUR", FR: "EUR", 法国: "EUR",
  JP: "JPY", 日本: "JPY",
  CA: "CAD", 加拿大: "CAD",
  AU: "AUD", 澳大利亚: "AUD",
};
// 合理性上限：跨境消费品单价超过此值视为抓取/解析异常（可按品类放宽）
const PRICE_UPPER_BOUND = 1_000_000;

const DATE_RE = /(\d{4})[-/](\d{1,2})[-/](\d{1,2})/;
const RESTRICTED_RE = /未获取|受访问限制|待确认|较旧/;
const DATE_HEADER_RE = /抓取日期|抓取时间|采样日期|as[_-]?of|最后更新|日期/;
const ROLE_HEADER_RE = /数据角色|数据类型|数据性质/;
const STATUS_HEADER_RE = /获取状态|数据获取状态|状态/;
const REALTIME_PRICE_HEADER_RE = /实时.*价格|抓取价|页面价|当前价/;
const API_PRICE_HEADER_RE = /api.*价格|源端价|接口价|缓存价/i;

function parseDate(cell) {
  const m = cell.match(DATE_RE);
  if (!m) return null;
  const [, y, mo, d] = m;
  const dt = new Date(Number(y), Number(mo) - 1, Number(d));
  return Number.isNaN(dt.getTime()) ? null : dt;
}

function parseAmount(cell) {
  const m = cell.replace(/,/g, "").match(/-?\d+(?:\.\d+)?/);
  return m ? Number(m[0]) : null;
}

// 解析报告的目标市场期望币种：优先取自 output/{平台}-{市场} 路径，其次 §1 输入摘要 目标市场：
function resolveExpectedCurrency(markdown, displayPath) {
  const pm = displayPath.match(/output[/\\][^/\\]+-([A-Za-z]{2})[/\\]/i);
  if (pm) {
    const code = pm[1].toUpperCase();
    if (MARKET_CURRENCY[code]) return MARKET_CURRENCY[code];
  }
  const m = markdown.match(/目标市场[：:]\s*([^\n|]+)/);
  if (m) {
    const v = m[1];
    for (const [k, c] of Object.entries(MARKET_CURRENCY)) {
      if (v.includes(k)) return c;
    }
  }
  return null; // 无法确定目标市场时，跳过币种校验，避免误报
}

// 返回 null=币种匹配；否则返回检测到的错误币种或 "NONE"（无币种标记）
function currencyMismatch(cell, expected) {
  const up = cell.toUpperCase();
  const explicit = up.includes("USD") ? "USD"
    : up.includes("GBP") ? "GBP"
    : up.includes("EUR") ? "EUR"
    : up.includes("JPY") ? "JPY"
    : up.includes("CNY") || up.includes("RMB") ? "CNY"
    : null;
  if (explicit) return explicit === expected ? null : explicit;
  if (cell.includes("$")) return expected === "USD" ? null : "USD";
  if (cell.includes("£")) return expected === "GBP" ? null : "GBP";
  if (cell.includes("€")) return expected === "EUR" ? null : "EUR";
  if (cell.includes("¥")) return expected === "JPY" ? null : "CNY"; // ¥ 歧义：仅期望 JPY 时认可，否则按 CNY 判错配
  return "NONE";
}

function slaForRole(roleCell) {
  if (/趋势|季节性|指数/.test(roleCell)) return SLA_HOURS.trend;
  if (/评论|评分/.test(roleCell)) return SLA_HOURS.review;
  return SLA_HOURS.realtime; // 价格 / 排名 / 库存 / Buy Box / 默认
}

function splitRow(line) {
  return line
    .trim()
    .replace(/^\|/, "")
    .replace(/\|$/, "")
    .split("|")
    .map((c) => c.trim());
}

function isSeparator(line) {
  return /^\s*\|?[\s:|-]+\|?\s*$/.test(line) && line.includes("-");
}

const reportPaths = process.argv.slice(2);
if (reportPaths.length === 0) {
  console.error("用法：node scripts/check-data-freshness.mjs <报告.md> [更多报告.md]");
  process.exit(1);
}

let hasError = false;
let checkedRows = 0;

for (const inputPath of reportPaths) {
  const absolutePath = resolve(inputPath);
  const displayPath = relative(process.cwd(), absolutePath) || absolutePath;
  const markdown = await readFile(absolutePath, "utf8");
  const expectedCurrency = resolveExpectedCurrency(markdown, displayPath);
  const lines = markdown.split(/\r?\n/);

  let block = [];
  const flush = () => {
    if (block.length >= 2) checkTable(block, displayPath, expectedCurrency);
    block = [];
  };

  for (const line of lines) {
    if (line.trim().startsWith("|")) {
      block.push(line);
    } else if (block.length) {
      flush();
    }
  }
  if (block.length) flush();

  function checkTable(rows, file, expectedCurrency) {
    const header = splitRow(rows[0]);
    const dateIdx = header.findIndex((h) => DATE_HEADER_RE.test(h));
    const roleIdx = header.findIndex((h) => ROLE_HEADER_RE.test(h));
    const statusIdx = header.findIndex((h) => STATUS_HEADER_RE.test(h));
    const rtPriceIdx = header.findIndex((h) => REALTIME_PRICE_HEADER_RE.test(h));
    const apiPriceIdx = header.findIndex((h) => API_PRICE_HEADER_RE.test(h));
    if (dateIdx === -1) return; // 非数据证据表，跳过

    for (let i = 1; i < rows.length; i++) {
      if (isSeparator(rows[i])) continue;
      const cells = splitRow(rows[i]);
      const rowNo = i + 1;
      const status = statusIdx >= 0 ? (cells[statusIdx] || "") : "";
      const dateCell = (cells[dateIdx] || "").trim();
      if (RESTRICTED_RE.test(status)) continue; // 未获取 / 受访问限制 / 待确认 / 较旧：合规标注，跳过新鲜度

      if (!dateCell) {
        console.error(`${file}：第 ${rowNo} 行缺少抓取日期（强制时间戳）`);
        hasError = true;
        continue;
      }
      const dt = parseDate(dateCell);
      if (!dt) {
        console.error(`${file}：第 ${rowNo} 行抓取日期无法解析：${dateCell}`);
        hasError = true;
        continue;
      }
      const ageHours = (Date.now() - dt.getTime()) / 36e5;
      const sla = roleIdx >= 0 ? slaForRole(cells[roleIdx] || "") : SLA_HOURS.realtime;
      if (ageHours > sla) {
        console.error(
          `${file}：第 ${rowNo} 行数据较旧（${Math.round(ageHours)}h > SLA ${sla}h），但未标注“较旧/受访问限制”`,
        );
        hasError = true;
      }

      // 实时价 vs API 价 交叉偏差（>5% 报错）
      if (rtPriceIdx >= 0 && apiPriceIdx >= 0) {
        const rt = parseAmount(cells[rtPriceIdx] || "");
        const api = parseAmount(cells[apiPriceIdx] || "");
        if (rt != null && api != null && api !== 0) {
          const dev = Math.abs(rt - api) / Math.abs(api);
          if (dev > 0.05) {
            console.error(
              `${file}：第 ${rowNo} 行实时价与 API 价偏差 ${Math.round(dev * 100)}% > 5%（实时 ${rt} / API ${api}）`,
            );
            hasError = true;
          }
        }
      }

      // 合理性校验：币种匹配目标市场 + 价格非异常值
      if (expectedCurrency && (rtPriceIdx >= 0 || apiPriceIdx >= 0)) {
        for (const pi of [rtPriceIdx, apiPriceIdx].filter((p) => p >= 0 && p !== dateIdx)) {
          const cell = (cells[pi] || "").trim();
          const amount = parseAmount(cell);
          if (amount === null) continue; // 无数值（空 / 文本），跳过
          const mismatch = currencyMismatch(cell, expectedCurrency);
          if (mismatch === "NONE") {
            console.error(`${file}：第 ${rowNo} 行价格含数值但无币种标记，无法校验目标市场匹配（标记待确认）`);
            hasError = true;
          } else if (mismatch) {
            console.error(`${file}：第 ${rowNo} 行币种 ${mismatch} 与目标市场期望 ${expectedCurrency} 不符（价格 ${cell}）`);
            hasError = true;
          }
          if (amount === 0) {
            console.error(`${file}：第 ${rowNo} 行价格为 0（异常值，标记待确认）`);
            hasError = true;
          } else if (amount < 0) {
            console.error(`${file}：第 ${rowNo} 行价格为负（异常值，标记待确认）`);
            hasError = true;
          } else if (amount > PRICE_UPPER_BOUND) {
            console.error(`${file}：第 ${rowNo} 行价格 ${amount} 超过合理性上限 ${PRICE_UPPER_BOUND}（疑似抓取/解析异常）`);
            hasError = true;
          }
        }
      }
      checkedRows++;
    }
  }
}

if (hasError) {
  process.exit(1);
}
console.log(`数据有效性校验通过：检查 ${checkedRows} 行（时间戳 / 新鲜度 / 币种匹配 / 价格合理性）`);
