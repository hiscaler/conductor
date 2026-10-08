#!/usr/bin/env node
// Checks buyer-facing copy for Chinese characters in US-market Listings.
// Seller-side headings and research notes are intentionally excluded.

import { existsSync, readFileSync } from "node:fs";
import { basename, join, resolve } from "node:path";

const args = process.argv.slice(2);
const listingArg = args.find((arg) => !arg.startsWith("--"));
if (!listingArg) {
  console.error("用法: node scripts/check-copy-language.mjs <listing目录> [--market US]");
  process.exit(2);
}

const marketIndex = args.indexOf("--market");
const explicitMarket = marketIndex >= 0 ? args[marketIndex + 1] : undefined;
const listingDir = resolve(process.cwd(), listingArg);
if (!existsSync(listingDir)) {
  console.error(`错误：Listing 目录不存在：${listingDir}`);
  process.exit(2);
}

const pathParts = listingDir.split(/[\\/]+/);
const market = explicitMarket ?? pathParts
  .slice(0, -1)
  .reverse()
  .map((part) => part.match(/^[^-]+-([A-Z]{2})$/)?.[1])
  .find(Boolean);

if (!market) {
  console.error("错误：无法从目录识别目标市场；请用 --market <两位市场代码> 指定。");
  process.exit(2);
}
if (market !== "US") {
  console.log(`SKIP：当前语言检查只针对美国站中文字符，目标市场为 ${market}。`);
  process.exit(0);
}

const copyDir = join(listingDir, "文案");
const copyAsset = join(copyDir, "文案资产.md");
const keywordList = join(copyDir, "关键词清单.md");
const buyerHeadings = new Set([
  "商品标题",
  "要点描述",
  "长描述",
  "搜索关键词",
  "后台关键词",
  "广告短文案",
]);
const sectionHeading = /^#{2,4}\s+(.+?)\s*#*$/u;
const errors = [];
const checkedFields = [];
let assetExceptions = [];

function getExceptions(markdown) {
  const exceptions = [];
  for (const line of markdown.split(/\r?\n/u)) {
    const heading = line.match(sectionHeading);
    if (heading && buyerHeadings.has(heading[1])) break;
    const exception = line.match(/^>\s*文案语言例外：(.+?)\s*$/u)?.[1]?.trim();
    if (exception) exceptions.push(exception);
  }
  return exceptions;
}

function checkValue(value, label, exceptions, lineNumber) {
  let checked = value;
  for (const exception of exceptions) checked = checked.replaceAll(exception, "");
  const matches = [...checked.matchAll(/\p{Script=Han}+/gu)];
  for (const match of matches) {
    errors.push(`${label}:${lineNumber} 含中文「${match[0]}」`);
  }
}

if (existsSync(copyAsset)) {
  const markdown = readFileSync(copyAsset, "utf8");
  const lines = markdown.split(/\r?\n/u);
  const exceptions = getExceptions(markdown);
  assetExceptions = exceptions;
  let active = null;

  for (let index = 0; index < lines.length; index++) {
    const heading = lines[index].match(sectionHeading);
    if (heading) {
      if (active) active.end = index;
      active = buyerHeadings.has(heading[1])
        ? { name: heading[1], start: index + 1, end: lines.length }
        : null;
      if (active) checkedFields.push(active.name);
      continue;
    }
    if (active && lines[index].trim()) {
      checkValue(lines[index], `${basename(copyAsset)}「${active.name}」`, exceptions, index + 1);
    }
  }
}

if (existsSync(keywordList)) {
  const markdown = readFileSync(keywordList, "utf8");
  const exceptions = assetExceptions;
  const rows = markdown.split(/\r?\n/u);
  for (let index = 0; index < rows.length; index++) {
    const row = rows[index];
    if (!/^\s*\|/u.test(row)) continue;
    const firstCell = row.split("|")[1]?.trim() ?? "";
    if (!firstCell || /^[-: ]+$/u.test(firstCell) || firstCell.includes("关键词")) continue;
    checkedFields.push("关键词清单:关键词");
    checkValue(firstCell, `${basename(keywordList)}「关键词」`, exceptions, index + 1);
  }
}

if (checkedFields.length === 0) {
  console.error("FAIL：未找到可检查的买家文案字段或关键词值；不能以缺少检查对象判定通过。");
  process.exit(1);
}

if (errors.length) {
  console.error(`FAIL：美国站买家文案发现中文字符（检查字段：${[...new Set(checkedFields)].join("、")}）`);
  for (const error of errors) console.error(`- ${error}`);
  console.error("请改为自然英文表达；只有用户明确提供且必须原样展示的定制原文，才可按文案资产顶部的“文案语言例外”精确豁免。");
  process.exit(1);
}

console.log(`PASS：美国站买家文案未发现中文字符（检查字段：${[...new Set(checkedFields)].join("、")}）。`);
