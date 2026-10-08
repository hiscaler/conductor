#!/usr/bin/env node
// 统一交付校验闸门（conductor 级强制，客户端无关）
//
// 所有驱动 conductor 的客户端 / 工具执行同一动作：
//   node scripts/validate-listing.mjs <listing目录> [-p <平台>]
//   （或 npm run validate -- <listing目录> [-p <平台>]）
//
// 该脚本统一跑全部 5 项强制校验，任一项未过即退出非 0。
// 退出码：0 = 全部通过（可交付）；非 0 = 有未通过项（不得交付）。
//
// 设计原则：
// - 脚本位于 conductor/scripts/，纯 Node 标准库、接收路径参数，不绑定任何客户端。
// - 用 process.execPath 调用子脚本，保证与运行本脚本的 Node 一致，不硬编码路径。

import { resolve, join } from "node:path";
import { fileURLToPath } from "node:url";
import { readdirSync, readFileSync, existsSync } from "node:fs";
import { spawnSync } from "node:child_process";

const SCRIPT = fileURLToPath(import.meta.url);
const ROOT = resolve(SCRIPT, "../.."); // conductor 项目根
const NODE = process.execPath;

function usage() {
  console.error("用法: node scripts/validate-listing.mjs <listing目录> [-p <平台>]");
  process.exit(2);
}

const args = process.argv.slice(2);
const listingArg = args[0];
if (!listingArg) usage();
const pIdx = args.indexOf("-p");
const platform = pIdx >= 0 ? args[pIdx + 1] : undefined;
const dryRun = args.includes("--dry-run");

const listingDir = resolve(process.cwd(), listingArg);
if (!existsSync(listingDir)) {
  console.error(`错误：listing 目录不存在：${listingDir}`);
  process.exit(2);
}

// ---- 定位各子文件 ----
const reportPath = join(listingDir, "上架", "完整生产报告.md");
const copyDir = join(listingDir, "文案");
const imgDir = join(listingDir, "图片");
const nextActionMeta = join(listingDir, "上架", "next-action.json");

const copyMds = existsSync(copyDir)
  ? readdirSync(copyDir).filter((f) => f.toLowerCase().endsWith(".md")).map((f) => join(copyDir, f))
  : [];

// ---- 定义校验项 ----
const checks = [];

if (existsSync(reportPath)) {
  checks.push({ name: "报告锚点", script: "check-report-anchors.mjs", args: [reportPath] });
} else {
  checks.push({ name: "报告锚点", script: null, skip: `缺 ${reportPath}` });
}

// Keyword lists and seller-side batch references have no buyer description.
// Only check documents that actually contain the description field.
const descriptionMds = copyMds.filter((f) => /^#{2,4}\s+长描述\s*$/mu.test(readFileSync(f, "utf8")));
for (const f of descriptionMds) {
  checks.push({ name: `Temu 描述长度 (${copyDir.split(/[\\/]/).pop()}/${f.split(/[\\/]/).pop()})`, script: "temu-description-limit.mjs", args: [f] });
}
if (descriptionMds.length === 0) {
  checks.push({ name: "Temu 描述长度", script: null, skip: `缺 ${copyDir}/*.md` });
}

const market = listingDir.split(/[\\/]/).slice(0, -1).reverse()
  .map((part) => part.match(/^[^-]+-([A-Z]{2})$/)?.[1])
  .find(Boolean);
if (market === "US" && copyMds.length > 0) {
  checks.push({ name: "美国站文案语言", script: "check-copy-language.mjs", args: [listingDir, "--market", market] });
} else if (market === "US") {
  checks.push({ name: "美国站文案语言", script: null, skip: `缺 ${copyDir}/*.md` });
}

checks.push({ name: "输出结构", script: "check-output-layout.mjs", args: [listingDir] });

// Automatic batches keep their own identity and required delivery scope.
// The association checker is read-only, so calling it from this gate cannot recurse.
if (existsSync(join(listingDir, "批次关联.json")) || /^批\d{12}(?:-\d+)?-/.test(listingDir.split(/[\\/]/).pop())) {
  checks.push({ name: "批次关联与任务范围", script: "batch-production.mjs", args: ["check-listing", listingDir] });
}

if (existsSync(imgDir)) {
  const imgArgs = [imgDir];
  if (platform) imgArgs.push("-p", platform);
  checks.push({ name: "图像套图", script: "check-image-set.mjs", args: imgArgs });
} else {
  checks.push({ name: "图像套图", script: null, skip: `缺 ${imgDir}` });
}

// 下一步动作范围：依赖启动菜单编号 + 动作 ID，存于 listing/上架/next-action.json
if (existsSync(nextActionMeta)) {
  let meta;
  try {
    meta = JSON.parse(readFileSync(nextActionMeta, "utf8"));
  } catch (e) {
    checks.push({ name: "下一步动作范围", script: null, skip: `next-action.json 解析失败：${e.message}` });
    meta = null;
  }
  if (meta && meta.menu && Array.isArray(meta.actions)) {
    checks.push({ name: "下一步动作范围", script: "next-action-scope.mjs", args: [String(meta.menu), ...meta.actions] });
  } else if (meta) {
    checks.push({ name: "下一步动作范围", script: null, skip: "next-action.json 缺 menu/actions" });
  }
} else {
  checks.push({ name: "下一步动作范围", script: null, skip: `缺 ${nextActionMeta}（动作范围校验为菜单相关项，需提供 next-action.json 或由人工执行）` });
}

// ---- 执行 ----
console.log(`\n=== conductor 统一交付校验：${listingDir} ===`);
if (platform) console.log(`平台：${platform}（来自 -p 参数）\n`);
else console.log(`平台：由路径自动识别\n`);

let pass = 0;
let fail = 0;
let skip = 0;

for (const c of checks) {
  if (c.skip) {
    console.log(`· ${c.name.padEnd(22)} SKIP  ${c.skip}`);
    skip++;
    continue;
  }
  const scriptPath = join(ROOT, "scripts", c.script);
  if (dryRun) {
    console.log(`· ${c.name.padEnd(22)} WOULD-RUN: node ${c.script} ${c.args.join(" ")}`);
    continue;
  }
  const res = spawnSync(NODE, [scriptPath, ...c.args], { cwd: ROOT, encoding: "utf8" });
  const code = res.status ?? 1;
  if (code === 0) {
    console.log(`· ${c.name.padEnd(22)} PASS`);
    if (c.script === "check-image-set.mjs" && res.stdout) {
      const reviewNotes = res.stdout.trim().split("\n").filter((line) =>
        line.includes("仍须视觉核验") || line.includes("机器只检查侧车字段和比例算术"),
      );
      for (const note of reviewNotes) console.log(`    ${note.trim()}`);
    }
    pass++;
  } else {
    console.log(`· ${c.name.padEnd(22)} FAIL  (exit ${code})`);
    if (res.error) console.log(`    ! spawn 错误：${res.error.message}（若报 EBUSY 多为沙箱禁止嵌套 spawn，真实终端环境正常）`);
    if (res.stdout) console.log(res.stdout.trim().split("\n").map((l) => "    " + l).join("\n"));
    if (res.stderr) console.log(res.stderr.trim().split("\n").map((l) => "    ! " + l).join("\n"));
    fail++;
  }
}

console.log(`\n=== 汇总：PASS ${pass} / FAIL ${fail} / SKIP ${skip} ===`);
if (fail > 0) {
  console.error(`\n✗ 未通过校验，不得交付。先修复上列 FAIL 项。`);
  process.exit(1);
}
if (skip > 0) {
  console.log(`\n△ 有 ${skip} 项跳过（多为缺文件或菜单相关项未提供 meta），交付前请确认其不适用或人工补跑。`);
}
console.log(`\n✓ 全部通过，可交付。`);
process.exit(0);
