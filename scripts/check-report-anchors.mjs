import { readFile } from "node:fs/promises";
import { relative, resolve } from "node:path";

const reportPaths = process.argv.slice(2);

if (reportPaths.length === 0) {
  console.error("用法：node scripts/check-report-anchors.mjs <报告.md> [更多报告.md]");
  process.exit(1);
}

let hasError = false;

for (const inputPath of reportPaths) {
  const absolutePath = resolve(inputPath);
  const displayPath = relative(process.cwd(), absolutePath) || absolutePath;
  const markdown = await readFile(absolutePath, "utf8");
  const anchorIds = [...markdown.matchAll(/<a\s+id="([^"]+)"\s*><\/a>/g)].map((match) => match[1]);
  const issueLinks = [...markdown.matchAll(/\]\(#(issue-[a-z0-9-]+)\)/g)].map((match) => match[1]);
  const seen = new Set();

  if (!markdown.includes("## 问题速览")) {
    console.error(`${displayPath}：缺少“## 问题速览”`);
    hasError = true;
  }

  for (const anchorId of anchorIds) {
    if (seen.has(anchorId)) {
      console.error(`${displayPath}：问题锚点重复：${anchorId}`);
      hasError = true;
    }
    seen.add(anchorId);
  }

  if (!seen.has("issue-summary")) {
    console.error(`${displayPath}：缺少问题速览锚点 issue-summary`);
    hasError = true;
  }

  for (const targetId of issueLinks) {
    if (!seen.has(targetId)) {
      console.error(`${displayPath}：问题链接没有对应锚点：#${targetId}`);
      hasError = true;
    }
  }
}

if (hasError) {
  process.exit(1);
}

console.log(`报告问题锚点检查通过：${reportPaths.length} 个文件`);
