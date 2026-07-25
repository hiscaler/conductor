import assert from "node:assert/strict";
import {
  countTemuDescriptionCharacters,
  extractLongDescription,
  validateTemuDescription,
} from "./temu-description-limit.mjs";

assert.equal(countTemuDescriptionCharacters("A, B."), 5, "空格和标点必须计入字符数");
assert.equal(
  countTemuDescriptionCharacters("可定制！"),
  4,
  "中文和中文标点必须按字符统计",
);

const exactlyAtLimit = validateTemuDescription("a".repeat(500));
assert.equal(exactlyAtLimit.passed, true, "500 个字符应通过");
assert.equal(exactlyAtLimit.results[0].characterCount, 500);

const overLimit = validateTemuDescription("a".repeat(501));
assert.equal(overLimit.passed, false, "501 个字符应判定为超限");
assert.equal(overLimit.results[0].characterCount, 501);

const twoParagraphs = validateTemuDescription(`${"a".repeat(450)}\n\n${"b".repeat(420)}`);
assert.equal(twoParagraphs.passed, true, "多个段落必须分别计数");
assert.equal(twoParagraphs.results.length, 2);

const markdown = [
  "## 文案资产",
  "",
  "### 长描述",
  "",
  "First paragraph.",
  "",
  "Second paragraph!",
  "",
  "#### Temu 描述段落验收",
  "",
  "| 段落 | 字符数 |",
].join("\n");
assert.equal(
  extractLongDescription(markdown),
  "First paragraph.\n\nSecond paragraph!",
  "提取长描述时必须排除后续验收内容",
);

console.log("Temu 商品描述字符限制检查通过。");
