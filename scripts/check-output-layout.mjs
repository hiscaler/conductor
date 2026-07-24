import { readdir } from "node:fs/promises";
import { resolve } from "node:path";
import { pathToFileURL } from "node:url";

const findEmptyDirectories = async (root) => {
  const emptyDirectories = [];

  const walk = async (directory, includeCurrent) => {
    let entries;
    try {
      entries = await readdir(directory, { withFileTypes: true });
    } catch (error) {
      if (error?.code === "ENOENT") return;
      throw error;
    }

    if (includeCurrent && entries.length === 0) {
      emptyDirectories.push(directory);
      return;
    }

    for (const entry of entries) {
      if (!entry.isDirectory()) continue;
      await walk(resolve(directory, entry.name), true);
    }
  };

  await walk(root, false);
  return emptyDirectories;
};

export const checkOutputLayout = async (rootPath = "output") => {
  const root = resolve(rootPath);
  const emptyDirectories = await findEmptyDirectories(root);
  if (emptyDirectories.length > 0) {
    const details = emptyDirectories.map((directory) => `- ${directory}`).join("\n");
    throw new Error(`发现未按需创建的空目录：\n${details}`);
  }
  return { root, emptyDirectories: [] };
};

const runCli = async () => {
  const targets = process.argv.slice(2);
  if (targets.length > 1) {
    process.stderr.write("输出布局检查失败：最多接受一个检查目录\n");
    process.exitCode = 1;
    return;
  }

  try {
    await checkOutputLayout(targets[0] ?? "output");
    process.stdout.write("输出布局检查通过：未发现空目录\n");
  } catch (error) {
    process.stderr.write(`${error.message}\n`);
    process.exitCode = 1;
  }
};

if (process.argv[1] && pathToFileURL(process.argv[1]).href === import.meta.url) {
  await runCli();
}
