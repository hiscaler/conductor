import { mkdir, writeFile } from "node:fs/promises";
import sharp from "sharp";

const src = "assets/favicon.svg";

await mkdir("build/windows", { recursive: true });
await sharp(src).resize(512, 512).png().toFile("build/appicon.png");

// Multi-size ICO for Windows exe
const sizes = [16, 24, 32, 48, 64, 128, 256];
const images = await Promise.all(
  sizes.map((size) => sharp(src).resize(size, size).png().toBuffer())
);

// Minimal ICO writer: ICONDIR + ICONDIRENTRY[] + image data
const entries = [];
let offset = 6 + sizes.length * 16;
const payloads = [];
for (let i = 0; i < images.length; i++) {
  const data = images[i];
  const size = sizes[i];
  entries.push({ size, offset, bytes: data.length });
  payloads.push(data);
  offset += data.length;
}

const header = Buffer.alloc(6 + entries.length * 16);
header.writeUInt16LE(0, 0);
header.writeUInt16LE(1, 2);
header.writeUInt16LE(entries.length, 4);
entries.forEach((e, i) => {
  const o = 6 + i * 16;
  header[o] = e.size >= 256 ? 0 : e.size;
  header[o + 1] = e.size >= 256 ? 0 : e.size;
  header[o + 2] = 0;
  header[o + 3] = 0;
  header.writeUInt16LE(1, o + 4);
  header.writeUInt16LE(32, o + 6);
  header.writeUInt32LE(e.bytes, o + 8);
  header.writeUInt32LE(e.offset, o + 12);
});
await writeFile("build/windows/icon.ico", Buffer.concat([header, ...payloads]));
console.log("wrote build/appicon.png and build/windows/icon.ico");
