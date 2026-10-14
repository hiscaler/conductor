#!/usr/bin/env node
// 检查报告中“已完成文案研究”的声明是否至少有可追溯详情证据。
// 页面真实性、描述内容和商品事实映射仍需人工核验。
import { readFileSync } from "node:fs";

const reportPath = process.argv[2];
if (!reportPath) {
  console.error("用法: node scripts/check-copy-research.mjs <完整生产报告.md>");
  process.exit(2);
}

let report;
try {
  report = readFileSync(reportPath, "utf8");
} catch (error) {
  console.error(`无法读取报告：${error.message}`);
  process.exit(2);
}

const chapter = report.match(/^## 5\. 竞品分析\s*$([\s\S]*?)(?=^## 6\. |$(?![\s\S]))/mu)?.[1] ?? "";
if (!chapter) {
  console.error("FAIL：缺少第 5 章竞品分析");
  process.exit(1);
}

const field = (name) => chapter.match(new RegExp(`^- ${name}：\\s*(.*)$`, "mu"))?.[1]?.trim() ?? "";
const declaredComplete = /^(是|已完成|完成)(?:[；;，,\s]|$)/u.test(field("是否完成真实竞品调研")) ||
  field("文案研究结论") === "已完成";
const evidence = chapter.match(/^### 文案研究证据\s*$([\s\S]*?)(?=^### |^## |$(?![\s\S]))/mu)?.[1] ?? "";

// Older reports without the new block remain readable. Completion claims do not.
if (!evidence) {
  if (declaredComplete) {
    console.error("FAIL：声明完成竞品/文案研究，但缺少“文案研究证据”记录");
    process.exit(1);
  }
  console.log("PASS：未声明已完成研究；旧报告无文案研究证据块");
  process.exit(0);
}

const readCount = Number(field("实际读到描述的详情页数"));
const openCount = Number(field("实际打开详情页数"));
const reviewCount = Number(field("实际读到的单品评论数"));
const status = field("文案研究结论");
const validStatuses = new Set(["已完成", "有限参考", "仅商品事实", "本轮不适用"]);
const errors = [];
if (!validStatuses.has(status)) errors.push("文案研究结论须使用模板状态");
for (const label of ["实际打开详情页数", "实际读到描述的详情页数", "实际读到的单品评论数"]) {
  if (field(label) === "") errors.push(`${label}不得留空`);
}
for (const [label, value] of [["实际打开详情页数", openCount], ["实际读到描述的详情页数", readCount], ["实际读到的单品评论数", reviewCount]]) {
  if (!Number.isSafeInteger(value) || value < 0) errors.push(`${label}须为非负整数`);
}
if (readCount > openCount) errors.push("读到描述的页数不能多于打开详情页数");

const rows = evidence.split(/\r?\n/u).filter((line) => /^\|/u.test(line) && /https?:\/\//u.test(line));
const describedRows = rows.filter((line) => /描述可读|描述已读|描述\s*\d+|描述：/u.test(line));
if (declaredComplete && readCount === 0) errors.push("声明已完成研究，但未读到任何详情描述");
if (declaredComplete && describedRows.length < readCount) errors.push("已读描述数量缺少对应详情 URL 与可读范围表行");
if (readCount > 0 && describedRows.length < readCount) errors.push("详情描述计数与证据表不符");
if (rows.length < openCount) errors.push("打开详情页计数与证据表不符");

if (errors.length) {
  for (const error of errors) console.error(`FAIL：${error}`);
  process.exit(1);
}
console.log(`PASS：文案研究声明与记录数量一致（详情 ${openCount}，描述 ${readCount}，单品评论 ${reviewCount}）；内容真实性须人工核验`);
