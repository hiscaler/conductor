import { readdir } from "node:fs/promises";
import { basename, dirname, isAbsolute, relative, resolve, sep } from "node:path";
import { pathToFileURL } from "node:url";

const VERSION_SUFFIX = /-v[1-9]\d*$/;

const escapeRegExp = (value) => value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");

const isInside = (root, target) => {
  const rel = relative(root, target);
  return rel !== "" && !rel.startsWith("..") && !isAbsolute(rel);
};

const readDirectory = async (directory) => {
  try {
    return await readdir(directory);
  } catch (error) {
    if (error?.code === "ENOENT") return [];
    throw error;
  }
};

/**
 * 为一次正式生产分配独立的 Listing 版本目录。
 *
 * 首次生产使用基础 Listing 目录；基础目录或历史版本目录存在时，返回
 * `{Listing标识}-vN`。目录内的商品资料、文案、图片、视频和报告始终
 * 使用规范中的稳定文件名，不添加版本后缀。
 */
export async function resolveOutputDirectory(listingPath, options = {}) {
  if (typeof listingPath !== "string" || listingPath.trim() === "") {
    throw new Error("必须提供基础 Listing 输出目录");
  }

  const root = resolve(options.root ?? "output");
  const absolutePath = resolve(listingPath);
  if (!isInside(root, absolutePath)) {
    throw new Error(`Listing 输出目录必须位于 ${root} 内：${listingPath}`);
  }

  const relativeParts = relative(root, absolutePath).split(sep);
  if (relativeParts.length !== 2) {
    throw new Error(
      `必须传入 output/{平台}-{市场}/{Listing标识} 层级的基础目录：${listingPath}`,
    );
  }

  const listingName = basename(absolutePath);
  if (VERSION_SUFFIX.test(listingName)) {
    throw new Error(`必须传入不带版本后缀的基础 Listing 目录：${listingPath}`);
  }

  const parent = dirname(absolutePath);
  const entries = await readDirectory(parent);
  const pattern = new RegExp(`^${escapeRegExp(listingName)}(?:-v([1-9]\\d*))?$`);

  let greatestExistingVersion = 0;
  for (const entry of entries) {
    const match = entry.match(pattern);
    if (!match) continue;
    const version = match[1] ? Number(match[1]) : 1;
    greatestExistingVersion = Math.max(greatestExistingVersion, version);
  }

  const version = greatestExistingVersion === 0 ? 1 : greatestExistingVersion + 1;
  const directoryName = version === 1 ? listingName : `${listingName}-v${version}`;

  return {
    version,
    root,
    requested: listingPath,
    path: resolve(parent, directoryName),
  };
}

const runCli = async () => {
  const targets = process.argv.slice(2);
  if (targets.length !== 1) {
    process.stderr.write(
      "版本分配失败：只接受一个基础 Listing 目录，例如 output/Temu-US/MUG0110RD\n",
    );
    process.exitCode = 1;
    return;
  }

  try {
    const result = await resolveOutputDirectory(targets[0]);
    process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
  } catch (error) {
    process.stderr.write(`版本分配失败：${error.message}\n`);
    process.exitCode = 1;
  }
};

if (process.argv[1] && pathToFileURL(process.argv[1]).href === import.meta.url) {
  await runCli();
}
