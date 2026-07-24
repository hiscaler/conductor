import assert from "node:assert/strict";
import { mkdtemp, mkdir, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { resolveOutputDirectory } from "./output-versioning.mjs";

const tempRoot = await mkdtemp(join(tmpdir(), "conductor-output-versioning-"));
const listing = join(tempRoot, "Temu-US", "MUG0110RD");

const firstBatch = await resolveOutputDirectory(listing, { root: tempRoot });
assert.equal(firstBatch.version, 1);
assert.equal(firstBatch.path, listing, "首次生产必须使用基础 Listing 目录");

await mkdir(join(firstBatch.path, "文案"), { recursive: true });
await mkdir(join(firstBatch.path, "图片"), { recursive: true });
await mkdir(join(firstBatch.path, "上架"), { recursive: true });
await writeFile(join(firstBatch.path, "文案", "文案资产.md"), "first-copy", { flag: "wx" });
await writeFile(join(firstBatch.path, "图片", "01-最终定制主图.png"), "first-image", {
  flag: "wx",
});
await writeFile(join(firstBatch.path, "上架", "完整生产报告.md"), "first-report", {
  flag: "wx",
});

const secondBatch = await resolveOutputDirectory(listing, { root: tempRoot });
assert.equal(secondBatch.version, 2);
assert.equal(secondBatch.path, `${listing}-v2`, "第二次生产必须使用 Listing-v2 目录");

await mkdir(join(secondBatch.path, "文案"), { recursive: true });
await mkdir(join(secondBatch.path, "图片"), { recursive: true });
await mkdir(join(secondBatch.path, "上架"), { recursive: true });
await writeFile(join(secondBatch.path, "文案", "文案资产.md"), "second-copy", { flag: "wx" });
await writeFile(join(secondBatch.path, "图片", "01-最终定制主图.png"), "second-image", {
  flag: "wx",
});
await writeFile(join(secondBatch.path, "上架", "完整生产报告.md"), "second-report", {
  flag: "wx",
});

assert.equal(
  await readFile(join(firstBatch.path, "文案", "文案资产.md"), "utf8"),
  "first-copy",
  "第二次生产不得修改首次产物",
);
assert.equal(
  await readFile(join(secondBatch.path, "文案", "文案资产.md"), "utf8"),
  "second-copy",
);
assert.ok(
  !secondBatch.path.includes("文案资产-v2"),
  "版本后缀只能添加到 Listing 目录，不能添加到文件名",
);

await mkdir(`${listing}-v3`, { recursive: true });
const fourthBatch = await resolveOutputDirectory(listing, { root: tempRoot });
assert.equal(fourthBatch.version, 4, "必须跳过已存在的不完整 Listing 版本目录");
assert.equal(fourthBatch.path, `${listing}-v4`);

await assert.rejects(
  resolveOutputDirectory(join(tempRoot, "..", "outside"), { root: tempRoot }),
  /必须位于/,
);
await assert.rejects(
  resolveOutputDirectory(join(listing, "图片"), { root: tempRoot }),
  /Listing标识/,
);
await assert.rejects(
  resolveOutputDirectory(`${listing}-v2`, { root: tempRoot }),
  /不带版本后缀/,
);

process.stdout.write("Listing 目录版本行为检查通过\n");
