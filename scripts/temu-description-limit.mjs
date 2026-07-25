import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";

export const TEMU_DESCRIPTION_PARAGRAPH_LIMIT = 500;

/**
 * 按 Unicode 码点统计字符，空格和标点符号都计入长度。
 */
export function countTemuDescriptionCharacters(text) {
  return Array.from(text).length;
}

/**
 * 将同一 Markdown 段落内的换行整理为空格，模拟复制到平台后的段落文本。
 */
export function normalizeTemuDescriptionParagraph(paragraph) {
  return paragraph
    .split(/\r?\n/u)
    .map((line) => line.trim())
    .filter(Boolean)
    .join(" ")
    .replace(/\s+/gu, " ")
    .trim();
}

/**
 * 按空行拆分商品描述，返回需要逐段校验的实际文本。
 */
export function splitTemuDescriptionParagraphs(description) {
  return description
    .trim()
    .split(/\r?\n\s*\r?\n/u)
    .map(normalizeTemuDescriptionParagraph)
    .filter(Boolean);
}

/**
 * 从文案 Markdown 中提取“长描述”栏目，不包含后续验收小节。
 */
export function extractLongDescription(markdown) {
  const lines = markdown.replace(/\r\n/gu, "\n").split("\n");
  const headingIndex = lines.findIndex((line) => /^#{2,4}\s+长描述\s*$/u.test(line.trim()));

  if (headingIndex === -1) {
    throw new Error("未找到“长描述”栏目");
  }

  const content = [];
  for (let index = headingIndex + 1; index < lines.length; index += 1) {
    if (/^#{2,4}\s+\S/u.test(lines[index].trim())) {
      break;
    }
    content.push(lines[index]);
  }

  return content.join("\n").trim();
}

/**
 * 校验 Temu 商品描述的每个段落是否不超过 500 个字符。
 */
export function validateTemuDescription(description) {
  const paragraphs = splitTemuDescriptionParagraphs(description);
  const results = paragraphs.map((paragraph, index) => {
    const characterCount = countTemuDescriptionCharacters(paragraph);
    return {
      paragraph: index + 1,
      characterCount,
      limit: TEMU_DESCRIPTION_PARAGRAPH_LIMIT,
      passed: characterCount <= TEMU_DESCRIPTION_PARAGRAPH_LIMIT,
    };
  });

  return {
    passed: paragraphs.length > 0 && results.every((result) => result.passed),
    paragraphs,
    results,
  };
}

/**
 * 命令行入口：读取文案 Markdown，输出逐段字符统计并用退出码表示结果。
 */
async function main() {
  const filePath = process.argv[2];
  if (!filePath) {
    throw new Error("用法：node scripts/temu-description-limit.mjs <文案 Markdown 路径>");
  }

  const markdown = await readFile(resolve(filePath), "utf8");
  const description = extractLongDescription(markdown);
  const validation = validateTemuDescription(description);

  for (const result of validation.results) {
    const status = result.passed ? "通过" : "超限";
    console.log(`第 ${result.paragraph} 段：${result.characterCount}/${result.limit}，${status}`);
  }

  if (!validation.passed) {
    process.exitCode = 1;
  }
}

const isDirectRun =
  process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url);

if (isDirectRun) {
  main().catch((error) => {
    console.error(error.message);
    process.exitCode = 1;
  });
}
