import test from "node:test";
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("..", import.meta.url));
const checker = new URL("../scripts/check-data-freshness.mjs", import.meta.url).href;
const fixture = fileURLToPath(new URL("./fixtures/data-freshness-valid.md", import.meta.url));

// Freeze only the fixture process. Production checks continue to use Date.now().
function checkAt(timestamp) {
  const source = `Date.now = () => ${Date.parse(timestamp)};
    process.argv = [process.execPath, ${JSON.stringify(fileURLToPath(new URL(checker)))}, ${JSON.stringify(fixture)}];
    await import(${JSON.stringify(checker)});`;
  const result = spawnSync(process.execPath, ["--input-type=module", "-e", source], {
    cwd: root, encoding: "utf8", env: { ...process.env, TZ: "UTC" },
  });
  assert.ifError(result.error);
  return result;
}

test("固定参考时间内的数据样例通过时效检查", () => {
  const result = checkAt("2026-10-06T00:00:00Z");
  assert.equal(result.status, 0, result.stderr);
  assert.match(result.stdout, /数据有效性校验通过：检查 5 行/);
});

test("同一数据超过 SLA 后仍被真实校验器拒绝", () => {
  const result = checkAt("2026-10-09T00:00:00Z");
  assert.equal(result.status, 1);
  assert.match(result.stderr, /数据较旧（96h > SLA 48h）/);
  assert.doesNotMatch(result.stdout, /数据有效性校验通过/);
});
