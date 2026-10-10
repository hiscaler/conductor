import { readdir, readFile, writeFile } from "node:fs/promises";
import { resolve, basename, dirname } from "node:path";
import { randomInt, createHash } from "node:crypto";
import { fileURLToPath } from "node:url";

// 每个 Listing / 共享图案组只选择一次；续做读取原记录，不再次随机。
// spu：CSV 标准 SPU；type：图文/图/文/不提供；directory：已分配的生产目录。
export async function selectReferenceAssets({ spu, type, directory, object, template, root = resolve(dirname(fileURLToPath(import.meta.url)), "..") }) {
  if (!spu || /[\\/]/.test(spu) || spu === "." || spu === "..") throw new Error("SPU 必须是标准目录名称");
  const source = resolve(root, "data/products", spu);
  const lock = resolve(directory, "参考资产选择.json");
  const fingerprint = async (name) => {
    if (basename(name) !== name) throw new Error("参考文件必须位于 SPU 目录根层");
    return createHash("sha256").update(await readFile(resolve(source, name))).digest("hex");
  };
  const readLock = async () => {
    const saved = JSON.parse(await readFile(lock, "utf8"));
    if (saved.spu !== spu || saved.customization_type !== type) throw new Error("原参考选择与 SPU / 定制类型不一致，须明确修改后再生产");
    for (const [kind, requested] of [["object", object], ["template", template]]) {
      const selected = saved[kind];
      if (requested && selected?.metadata !== requested) throw new Error("用户指定与已锁定参考不一致，须明确替换原选择");
      if (selected && (await fingerprint(selected.image) !== selected.sha256 || await fingerprint(selected.metadata) !== selected.metadata_sha256)) throw new Error("已锁定参考资产发生变化，须重新核对");
    }
    return saved;
  };
  try { return await readLock(); } catch (error) { if (error.code !== "ENOENT") throw error; }
  if (!["图文", "图", "文", "不提供"].includes(type)) throw new Error("定制类型无效");
  const names = (await readdir(source)).sort();
  const choose = async (pattern, requested, suitable) => {
    const candidates = [];
    for (const name of names.filter((name) => pattern.test(name))) {
      const meta = JSON.parse(await readFile(resolve(source, name), "utf8"));
      if (meta.spu !== spu || !suitable(meta)) continue;
      const hash = await fingerprint(meta.image);
      if (hash !== meta.sha256) throw new Error(`图片哈希与资料不一致：${name}`);
      candidates.push({ id: meta.id, metadata: name, image: meta.image, sha256: hash, metadata_sha256: await fingerprint(name) });
    }
    const selected = requested ? candidates.find((item) => item.metadata === requested) : candidates.length ? candidates[randomInt(candidates.length)] : null;
    if (requested && !selected) throw new Error(`指定参考不存在或尚不适用：${requested}`);
    return selected;
  };
  const selectedObject = await choose(/^ref-object-\d+\.json$/, object, (meta) => meta.dimensions_verified === true && meta.dimensions_cm && Object.values(meta.dimensions_cm).every((value) => Number.isFinite(value) && value > 0));
  const selectedTemplate = type === "不提供" ? null : await choose(/^ref-custom-(text|image|image-text)-\d+\.json$/, template, (meta) => meta.customization_type === type && ["user_approved", "available_reference"].includes(meta.approval_status));
  const selected = { schema_version: 1, spu, customization_type: type, selected_at: new Date().toISOString(), object: selectedObject, template: selectedTemplate, missing: [!selectedObject ? "缺少已核验尺寸的适用实物参考" : null, type !== "不提供" && !selectedTemplate ? "缺少适用定制模板" : null].filter(Boolean) };
  try { await writeFile(lock, JSON.stringify(selected, null, 2) + "\n", { flag: "wx" }); }
  catch (error) { if (error.code !== "EEXIST") throw error; return readLock(); }
  return selected;
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const [spu, type, directory, object, template] = process.argv.slice(2);
  if (!spu || !type || !directory) { console.error("用法：node scripts/select-reference-assets.mjs <SPU> <图文|图|文|不提供> <生产目录> [ref-object-N.json] [ref-custom-text-N.json]"); process.exitCode = 1; }
  else try { console.log(JSON.stringify(await selectReferenceAssets({ spu, type, directory, object, template }), null, 2)); } catch (error) { console.error(error.message); process.exitCode = 1; }
}
