package browser

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Conductor 浏览器</title>
  <link rel="icon" href="/assets/favicon.svg" type="image/svg+xml">
  <script src="https://cdn.tailwindcss.com"></script>
  <style>
    html, body { overflow-x:hidden; }
    .tree, .tree ul { list-style:none; margin:0; padding-left:14px; }
    .tree { padding-left:0; }
    .collapsed > ul { display:none; }
    .markdown { line-height:1.7; min-width:0; overflow-wrap:anywhere; word-break:break-word; }
    .markdown-shell { display:grid; grid-template-columns:minmax(0, 1fr) 220px; gap:18px; align-items:start; min-width:0; }
    .copyable { position:relative; padding-right:42px; }
    .copy-btn {
      position:absolute; right:0; top:0.1em; border:0; background:transparent;
      color:rgb(100 116 139); border-radius:4px; padding:2px 6px; font-size:12px; opacity:0; cursor:pointer;
    }
    .copyable:hover .copy-btn { opacity:1; }
    .copy-btn:hover { color:rgb(125 211 252); background:rgb(148 163 184 / 0.08); }
    .markdown h1, .markdown h2, .markdown h3 { line-height:1.25; scroll-margin-top:18px; }
    .markdown h1 { font-size:24px; font-weight:700; color:rgb(248 250 252); border-bottom:1px solid rgb(51 65 85 / 0.7); padding-bottom:8px; margin:0 0 0.85em; }
    .markdown h2 { font-size:21px; font-weight:700; color:rgb(241 245 249); border-bottom:1px solid rgb(51 65 85 / 0.55); padding-bottom:8px; margin:2.2em 0 0.9em; }
    .markdown h3 { font-size:14px; font-weight:600; color:rgb(125 211 252); letter-spacing:0.02em; margin:1.5em 0 0.55em; padding:0; background:transparent; border:0; }
    .markdown h3.copyable { padding-right:52px; }
    .markdown p, .markdown li { color:rgb(203 213 225); font-size:15px; line-height:1.75; }
    .markdown p { margin:0 0 0.9em; padding:0; padding-right:42px; background:transparent; border:0; }
    .markdown h3 + p, .markdown h3 + ul, .markdown h3 + pre, .markdown h3 + table { margin-top:0; }
    .markdown ul { margin:0 0 0.9em; padding:0 0 0 1.25em; background:transparent; border:0; }
    .markdown li { margin:0.3em 0; padding-right:42px; }
    .markdown code { background:rgb(30 41 59 / 0.7); padding:1px 5px; border-radius:4px; color:rgb(226 232 240); }
    .markdown pre {
      white-space:pre-wrap; word-break:break-word; overflow-wrap:anywhere;
      margin:0 0 0.9em; padding:12px 0; padding-right:42px;
      background:transparent; border:0; border-top:1px solid rgb(51 65 85 / 0.45); border-bottom:1px solid rgb(51 65 85 / 0.45);
      color:rgb(226 232 240);
    }
    .markdown table { border-collapse:collapse; width:100%; margin:12px 0; table-layout:fixed; }
    .markdown th, .markdown td { border:1px solid rgb(51 65 85); padding:8px; vertical-align:top; }
    .markdown th, .markdown td { overflow-wrap:anywhere; word-break:break-word; }
    .markdown th { background:rgb(15 23 42); }
    @media (max-width: 1050px) { .markdown-shell { grid-template-columns:1fr; } .toc { position:static; max-height:none; order:-1; } }
    .header-actions {
      display:flex; align-items:center; gap:2px;
      font-size:13px; letter-spacing:0.01em;
    }
    .header-actions > * + * { margin-left:2px; }
    .header-link {
      display:inline-flex; align-items:center; gap:7px;
      padding:6px 10px; border:0; background:transparent; cursor:pointer;
      color:rgb(148 163 184); font:inherit; line-height:1;
      border-radius:6px; transition:color .15s ease, background .15s ease;
    }
    .header-link:hover { color:rgb(224 242 254); background:rgb(148 163 184 / 0.08); }
    .header-link:focus-visible { outline:1px solid rgb(56 189 248 / 0.5); outline-offset:2px; }
    .header-link svg { width:15px; height:15px; stroke-width:1.6; opacity:.85; }
    .header-link:hover svg { opacity:1; }
    .header-sep {
      width:1px; height:14px; margin:0 8px;
      background:rgb(51 65 85 / 0.9);
    }
    .header-auto {
      display:inline-flex; align-items:center; gap:9px;
      padding:4px 4px 4px 10px; cursor:pointer;
      color:rgb(100 116 139); font:inherit; line-height:1;
      transition:color .15s ease;
    }
    .header-auto:hover { color:rgb(148 163 184); }
    .header-auto:has(input:checked) { color:rgb(186 230 253); }
    .header-auto input {
      appearance:none; width:30px; height:16px; margin:0; flex-shrink:0;
      border-radius:999px; background:rgb(51 65 85);
      box-shadow:inset 0 0 0 1px rgb(71 85 105 / 0.6);
      position:relative; cursor:pointer; transition:background .18s ease, box-shadow .18s ease;
    }
    .header-auto input::after {
      content:""; position:absolute; top:2px; left:2px;
      width:12px; height:12px; border-radius:50%;
      background:rgb(226 232 240);
      box-shadow:0 1px 2px rgb(0 0 0 / 0.35);
      transition:transform .18s ease, background .18s ease;
    }
    .header-auto input:checked {
      background:rgb(14 165 233);
      box-shadow:inset 0 0 0 1px rgb(56 189 248 / 0.35);
    }
    .header-auto input:checked::after {
      transform:translateX(14px); background:white;
    }
    .header-auto #refreshState {
      min-width:2.75rem; color:inherit; opacity:.72;
      font-variant-numeric:tabular-nums; font-size:12px;
    }
    .header-auto:not(:has(input:checked)) #refreshState { opacity:.45; }
    @media (max-width: 640px) {
      header { height:auto; min-height:3.5rem; padding-top:10px; padding-bottom:10px; flex-wrap:wrap; }
      .header-link span, .header-auto > span:first-of-type { display:none; }
      .header-sep { margin:0 4px; }
      .header-link { padding:6px 8px; }
    }
    .gallery-grid {
      display:grid;
      grid-template-columns:repeat(auto-fill, minmax(160px, 1fr));
      gap:12px;
    }
    .gallery-card {
      display:flex; flex-direction:column; gap:8px;
      padding:0; border:0; border-radius:10px;
      background:rgb(15 23 42); color:inherit; font:inherit; text-align:left;
      cursor:pointer; overflow:hidden; transition:background .15s ease, outline-color .15s ease;
      outline:1px solid transparent; outline-offset:0;
    }
    .gallery-card:hover { background:rgb(30 41 59); outline-color:rgb(56 189 248 / 0.45); }
    .gallery-card.is-active { outline-color:rgb(56 189 248); }
    .gallery-card img {
      display:block; width:100%; aspect-ratio:1; object-fit:cover; background:rgb(2 6 23);
    }
    .file-list { width:100%; border-collapse:collapse; table-layout:fixed; }
    .file-list th, .file-list td {
      border:1px solid rgb(51 65 85); padding:9px 10px; vertical-align:middle; text-align:left;
    }
    .file-list th {
      background:rgb(15 23 42); color:rgb(148 163 184); font-size:12px; font-weight:500;
    }
    .file-list td { color:rgb(203 213 225); font-size:14px; }
    .file-list tr.file-row { cursor:pointer; }
    .file-list tr.file-row:hover td { background:rgb(30 41 59 / 0.7); color:rgb(224 242 254); }
    .file-list .name { overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
    .file-list .type, .file-list .size { color:rgb(148 163 184); font-size:12px; width:110px; }
    .file-list .size { width:88px; }
    .path-bar { display:flex; min-width:0; align-items:center; gap:10px; margin-top:6px; }
    .path-bar code {
      min-width:0; flex:1; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
      font-size:12px; color:rgb(100 116 139); background:transparent; padding:0;
    }
    .path-bar button {
      flex-shrink:0; border:0; background:transparent; padding:0;
      color:rgb(56 189 248); font:inherit; font-size:12px; cursor:pointer;
    }
    .path-bar button:hover { color:rgb(125 211 252); }
    .gallery-card .meta {
      padding:0 10px 10px; min-width:0;
    }
    .gallery-card .name {
      display:block; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
      font-size:12px; color:rgb(226 232 240);
    }
    .gallery-card .size { font-size:11px; color:rgb(100 116 139); }
    .lightbox {
      position:fixed; inset:0; z-index:50;
      display:flex; align-items:center; justify-content:center;
      background:rgb(2 6 23 / 0.88); padding:24px;
    }
    .lightbox[hidden] { display:none; }
    .lightbox-inner {
      position:relative; display:flex; flex-direction:column; align-items:center;
      max-width:min(1100px, 100%); max-height:100%;
    }
    .lightbox-toolbar {
      display:flex; align-items:center; justify-content:space-between; gap:16px;
      width:100%; margin-bottom:14px; color:rgb(226 232 240); font-size:13px;
    }
    .lightbox-toolbar .title {
      min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
      font-size:14px; font-weight:500; letter-spacing:0.01em; color:rgb(241 245 249);
    }
    .lightbox-actions { display:flex; align-items:center; gap:10px; flex-shrink:0; }
    .lightbox-counter {
      min-width:3.25rem; text-align:right;
      font-size:12px; font-variant-numeric:tabular-nums;
      color:rgb(148 163 184); letter-spacing:0.02em;
    }
    .lightbox-btn {
      display:inline-flex; align-items:center; gap:6px;
      height:32px; padding:0 12px 0 10px;
      border:1px solid transparent; border-radius:999px;
      background:rgb(255 255 255 / 0.06); color:rgb(226 232 240);
      font:inherit; font-size:12px; font-weight:500; letter-spacing:0.02em;
      cursor:pointer; backdrop-filter:blur(8px);
      transition:background .15s ease, border-color .15s ease, color .15s ease, transform .15s ease;
    }
    .lightbox-btn svg { width:14px; height:14px; stroke-width:1.7; opacity:.9; }
    .lightbox-btn:hover {
      background:rgb(255 255 255 / 0.1); color:rgb(255 255 255);
      border-color:rgb(148 163 184 / 0.25);
    }
    .lightbox-btn:active { transform:translateY(1px); }
    .lightbox-btn.primary {
      background:linear-gradient(180deg, rgb(14 165 233), rgb(2 132 199));
      border-color:rgb(56 189 248 / 0.35); color:white;
      box-shadow:0 1px 2px rgb(0 0 0 / 0.25), inset 0 1px 0 rgb(255 255 255 / 0.18);
    }
    .lightbox-btn.primary:hover {
      background:linear-gradient(180deg, rgb(56 189 248), rgb(14 165 233));
      border-color:rgb(125 211 252 / 0.5); color:white;
    }
    .lightbox-btn.ghost {
      background:transparent; border-color:rgb(71 85 105 / 0.7); color:rgb(203 213 225);
    }
    .lightbox-btn.ghost:hover {
      background:rgb(239 68 68 / 0.12); border-color:rgb(248 113 113 / 0.45); color:rgb(254 202 202);
    }
    .lightbox-stage {
      position:relative; display:flex; align-items:center; justify-content:center;
      max-width:100%; max-height:calc(100vh - 120px);
    }
    .lightbox-stage img {
      max-width:100%; max-height:calc(100vh - 120px);
      border-radius:12px; background:rgb(15 23 42);
      box-shadow:0 20px 50px rgb(0 0 0 / 0.45);
    }
    .lightbox-nav {
      position:absolute; top:50%; transform:translateY(-50%);
      display:inline-flex; align-items:center; justify-content:center;
      width:40px; height:40px; border-radius:999px;
      border:1px solid rgb(255 255 255 / 0.12);
      background:rgb(15 23 42 / 0.72); color:rgb(241 245 249);
      font-size:20px; line-height:1; cursor:pointer; backdrop-filter:blur(8px);
      transition:background .15s ease, border-color .15s ease, color .15s ease;
    }
    .lightbox-nav:hover {
      background:rgb(15 23 42 / 0.92); border-color:rgb(56 189 248 / 0.55); color:rgb(125 211 252);
    }
    .lightbox-nav.prev { left:12px; }
    .lightbox-nav.next { right:12px; }
    .lightbox-nav:disabled { opacity:.28; cursor:default; pointer-events:none; }
  </style>
</head>
<body class="bg-slate-950 text-slate-100">
<header class="flex h-14 items-center justify-between gap-4 border-b border-slate-800 bg-slate-900 px-5">
  <div class="flex min-w-0 items-center gap-3">
    <img class="h-9 w-auto shrink-0" src="/assets/coor-logo.svg" alt="Coor">
    <div class="min-w-0 leading-tight">
      <strong class="block text-[15px] tracking-wide text-slate-100">浏览器</strong>
      <span class="hidden text-xs text-slate-400 sm:block">AI 成果浏览</span>
    </div>
  </div>
  <nav class="header-actions" aria-label="工具">
    <button class="header-link" type="button" onclick="openReadme()">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M8 4.5h6.2L17.5 7.8V19.5H8z"/><path stroke-linecap="round" d="M10.2 11h3.8M10.2 14.2h3.8"/></svg>
      <span>使用说明</span>
    </button>
    <button class="header-link" type="button" onclick="loadTree()">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M4.8 12a7.2 7.2 0 0 1 12.3-5.1M19.2 12a7.2 7.2 0 0 1-12.3 5.1"/><path stroke-linecap="round" stroke-linejoin="round" d="M16.8 4.2V8h-3.8M7.2 19.8V16h3.8"/></svg>
      <span>刷新</span>
    </button>
    <span class="header-sep" aria-hidden="true"></span>
    <label class="header-auto" title="自动刷新目录树">
      <span>自动刷新</span>
      <span id="refreshState">5 秒</span>
      <input id="autoRefresh" type="checkbox" checked onchange="toggleAutoRefresh()">
    </label>
  </nav>
</header>
<main class="grid h-[calc(100vh-3.5rem)] grid-cols-[330px_minmax(0,1fr)] overflow-hidden max-[800px]:grid-cols-1 max-[640px]:h-[calc(100vh-4.5rem)]">
  <aside class="overflow-auto border-r border-slate-800 bg-slate-900 p-3 max-[800px]:h-[38vh] max-[800px]:border-b max-[800px]:border-r-0">
    <ul id="tree" class="tree"></ul>
  </aside>
  <section class="min-w-0 overflow-auto overflow-x-hidden p-6 max-[800px]:h-[calc(62vh-3.5rem)]">
    <div id="content" class="text-slate-500">请选择左侧文件或目录。</div>
  </section>
</main>
<div id="lightbox" class="lightbox" hidden onclick="closeLightboxBackdrop(event)">
  <div class="lightbox-inner" onclick="event.stopPropagation()">
    <div class="lightbox-toolbar">
      <div class="title" id="lightboxTitle"></div>
      <div class="lightbox-actions">
        <span id="lightboxCounter" class="lightbox-counter"></span>
        <button type="button" class="lightbox-btn primary" onclick="copyLightboxImage(event)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><rect x="9" y="9" width="11" height="11" rx="2"/><path stroke-linecap="round" stroke-linejoin="round" d="M5 15V7a2 2 0 0 1 2-2h8"/></svg>
          <span>复制图片</span>
        </button>
        <button type="button" class="lightbox-btn ghost" onclick="closeLightbox()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" d="M7 7l10 10M17 7L7 17"/></svg>
          <span>关闭</span>
        </button>
      </div>
    </div>
    <div class="lightbox-stage">
      <button type="button" class="lightbox-nav prev" id="lightboxPrev" onclick="galleryStep(-1)" aria-label="上一张">‹</button>
      <img id="lightboxImage" alt="">
      <button type="button" class="lightbox-nav next" id="lightboxNext" onclick="galleryStep(1)" aria-label="下一张">›</button>
    </div>
  </div>
</div>
<script>
let activePath = "";
let refreshTimer = null;
let expandedPaths = new Set();
let galleryImages = [];
let galleryIndex = -1;

// openReadme 加载项目 README，作为用户使用说明预览。
async function openReadme() {
  activePath = "__readme__";
  markActive();
  const res = await fetch("/api/readme");
  if (!res.ok) {
    document.getElementById("content").innerHTML = "<div class='text-slate-500'>README.md 读取失败</div>";
    return;
  }
  const data = await res.json();
  renderContent(data);
}

// loadTree 刷新左侧目录树，并保留展开状态和当前选中项。
async function loadTree() {
  const res = await fetch("/api/tree");
  if (!res.ok) return;
  const data = await res.json();
  rememberExpanded();
  document.getElementById("tree").innerHTML = renderChildren(data.children || [], 1);
  if (activePath) markActive();
}

// renderChildren 渲染指定层级下的目录节点列表。
function renderChildren(children, depth) {
  return children.map(n => renderNode(n, depth)).join("");
}

// renderNode 渲染左侧树中的单个文件或目录节点。
function renderNode(n, depth) {
  const icon = n.type === "dir" ? "📁" : iconFor(n.type);
  const meta = n.type === "dir" ? "" : "<span class='ml-auto text-xs text-slate-400'>" + formatSize(n.size || 0) + "</span>";
  const hasChildren = n.type === "dir" && n.children && n.children.length;
  const shouldCollapse = hasChildren && depth >= 2 && !expandedPaths.has(n.path);
  const liClass = shouldCollapse ? " class='collapsed'" : "";
  const twisty = hasChildren ? "<span class='twisty w-4 text-center text-slate-400'>" + (shouldCollapse ? "▶" : "▼") + "</span>" : "<span class='w-4'></span>";
  const child = hasChildren ? "<ul>" + renderChildren(n.children, depth + 1) + "</ul>" : "";
  return "<li" + liClass + "><button class='node flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm text-slate-200 hover:bg-slate-800 hover:text-blue-300' data-path='" + escAttr(n.path) + "' onclick='handleNodeClick(event,\"" + escJS(n.path) + "\"," + (hasChildren ? "true" : "false") + ")'>" + twisty + "<span class='w-5 text-center text-slate-400'>" + icon + "</span><span>" + esc(n.name) + "</span>" + meta + "</button>" + child + "</li>";
}

// iconFor 根据文件类型选择显示图标。
function iconFor(type) {
  if (type === "markdown") return "📝";
  if (type === "image") return "🖼️";
  if (type === "video") return "🎬";
  if (type === "json") return "{}";
  return "📄";
}

// handleNodeClick 处理节点点击，目录会展开/收起并打开内容。
function handleNodeClick(event, path, hasChildren) {
  if (hasChildren) toggleNode(event.currentTarget);
  openPath(path);
}

// toggleNode 展开或收起目录节点。
function toggleNode(button) {
  const item = button.closest("li");
  if (!item) return;
  item.classList.toggle("collapsed");
  const path = button.dataset.path;
  if (item.classList.contains("collapsed")) expandedPaths.delete(path);
  else expandedPaths.add(path);
  const twisty = button.querySelector(".twisty");
  if (twisty) twisty.textContent = item.classList.contains("collapsed") ? "▶" : "▼";
}

// rememberExpanded 在重建目录树前记录已展开的目录。
function rememberExpanded() {
  document.querySelectorAll(".node").forEach(button => {
    const item = button.closest("li");
    if (item && item.querySelector("ul") && !item.classList.contains("collapsed")) expandedPaths.add(button.dataset.path);
  });
}

// markActive 高亮左侧当前选中的文件或目录。
function markActive() {
  document.querySelectorAll(".node").forEach(el => {
    const active = el.dataset.path === activePath;
    el.classList.toggle("bg-slate-800", active);
    el.classList.toggle("text-blue-300", active);
  });
}

// openPath 加载选中路径的元信息和预览内容。
async function openPath(path) {
  activePath = path;
  markActive();
  const res = await fetch("/api/file?path=" + encodeURIComponent(path));
  if (!res.ok) {
    document.getElementById("content").innerHTML = "<div class='text-slate-500'>读取失败</div>";
    return;
  }
  const data = await res.json();
  renderContent(data);
}

// renderContent 根据文件类型选择合适的预览方式。
function renderContent(data) {
  closeLightbox();
  if (data.type === "dir") {
    renderDirContent(data);
    return;
  }
  if (data.type === "image") {
    renderImageContent(data);
    return;
  }
  let body = "";
  if (data.type === "markdown") body = renderMarkdownPreview(data.content || "");
  else if (data.type === "video") body = "<div><video class='max-w-full rounded-lg bg-black' src='" + escAttr(data.rawUrl) + "' controls></video></div>";
  else if (data.type === "text" || data.type === "json") body = "<pre class='overflow-auto whitespace-pre-wrap break-words rounded-lg bg-slate-900/60 p-4 text-slate-200'>" + esc(data.content || "") + "</pre>";
  else body = "<p><a href='" + escAttr(data.rawUrl) + "' target='_blank'>下载或打开文件</a></p>";
  const size = formatSize(data.size || 0);
  const sizeHint = size ? "<span class='text-xs text-slate-500'>" + esc(size) + "</span>" : "";
  document.getElementById("content").innerHTML = contentShell(
    "<div class='mb-5'><div class='mb-1 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold tracking-tight'>" + esc(data.name) + "</h1>" + sizeHint + "</div>" + renderPathBar(data) + "</div>",
    body
  );
}

// contentShell 组装右侧预览外壳，避免多层边框套框。
function contentShell(header, body) {
  return "<div class='max-w-6xl'>" + header + body + "</div>";
}

// renderDirContent 渲染目录内容；含图片时优先展示缩略图网格。
function renderDirContent(data) {
  const children = data.children || [];
  const images = children.filter(n => n.type === "image");
  const others = children.filter(n => n.type !== "image");
  galleryImages = images.map(n => ({
    name: n.name,
    path: n.path,
    size: n.size || 0,
    rawUrl: rawUrlFor(n.path)
  }));
  let body = "";
  if (images.length) {
    body += "<div class='mb-2 text-sm text-slate-500'>共 " + images.length + " 张图片，点击可放大预览</div>";
    body += "<div class='gallery-grid'>" + images.map((n, i) => renderGalleryCard(n, i)).join("") + "</div>";
  }
  if (others.length) {
    body += (images.length ? "<div class='mt-6 mb-2 text-sm text-slate-500'>其他文件</div>" : "")
      + renderFileList(others);
  }
  if (!body) body = "<div class='text-slate-500'>空目录</div>";
  document.getElementById("content").innerHTML = contentShell(
    "<div class='mb-5'><h1 class='m-0 text-xl font-semibold tracking-tight'>" + esc(data.name || "output") + "</h1>" + renderPathBar(data) + "</div>",
    body
  );
}

// renderFileList 以带边框表格展示目录中的非图片项。
function renderFileList(items) {
  const rows = items.map(n =>
    "<tr class='file-row' onclick='openPath(\"" + escJS(n.path) + "\")'>"
    + "<td class='name'>" + iconFor(n.type) + " " + esc(n.name) + "</td>"
    + "<td class='type'>" + esc(n.type) + "</td>"
    + "<td class='size'>" + esc(formatSize(n.size || 0) || "—") + "</td></tr>"
  ).join("");
  return "<table class='file-list'><thead><tr><th>名称</th><th>类型</th><th>大小</th></tr></thead><tbody>"
    + rows + "</tbody></table>";
}

// renderImageContent 渲染单张图片，并加载同目录图片供左右切换。
async function renderImageContent(data) {
  const size = formatSize(data.size || 0);
  const sizeHint = size ? "<span class='text-xs text-slate-500'>" + esc(size) + "</span>" : "";
  document.getElementById("content").innerHTML = contentShell(
    "<div class='mb-5'><div class='mb-1 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold tracking-tight'>" + esc(data.name) + "</h1>" + sizeHint + "</div>" + renderPathBar(data) + "</div>",
    "<div class='relative inline-block max-w-full'><img id='previewImage' class='max-w-full rounded-lg bg-slate-900' src='" + escAttr(data.rawUrl) + "' alt='" + escAttr(data.name) + "'><button type='button' class='absolute right-3 top-3 rounded-md bg-slate-950/80 px-2.5 py-1 text-[12px] text-slate-200 hover:text-sky-300' onclick='copyImage(event)'>复制图片</button></div><div id='siblingGallery' class='mt-5'></div>"
  );
  await loadSiblingGallery(data.path);
}

// loadSiblingGallery 拉取当前图片所在目录的全部图片，渲染缩略图并支持灯箱切换。
async function loadSiblingGallery(imagePath) {
  const parent = parentPath(imagePath);
  const res = await fetch("/api/file?path=" + encodeURIComponent(parent));
  if (!res.ok) return;
  const data = await res.json();
  if (data.type !== "dir") return;
  const images = (data.children || []).filter(n => n.type === "image");
  if (images.length < 2) return;
  galleryImages = images.map(n => ({
    name: n.name,
    path: n.path,
    size: n.size || 0,
    rawUrl: rawUrlFor(n.path)
  }));
  const currentIndex = galleryImages.findIndex(n => n.path === imagePath);
  const shell = document.getElementById("siblingGallery");
  if (!shell) return;
  shell.innerHTML = "<div class='mb-2 text-sm text-slate-400'>同目录共 " + images.length + " 张，点击切换</div><div class='gallery-grid'>"
    + images.map((n, i) => renderGalleryCard(n, i, i === currentIndex)).join("")
    + "</div>";
}

// renderGalleryCard 渲染单张缩略图卡片。
function renderGalleryCard(n, index, active) {
  const raw = rawUrlFor(n.path);
  const activeClass = active ? " is-active" : "";
  return "<button type='button' class='gallery-card" + activeClass + "' onclick='openGallery(" + index + ")'>"
    + "<img src='" + escAttr(raw) + "' alt='" + escAttr(n.name) + "' loading='lazy'>"
    + "<div class='meta'><span class='name' title='" + escAttr(n.name) + "'>" + esc(n.name) + "</span>"
    + "<span class='size'>" + esc(formatSize(n.size || 0)) + "</span></div></button>";
}

// rawUrlFor 根据相对路径生成原始文件 URL。
function rawUrlFor(path) {
  return "/raw?path=" + encodeURIComponent(path || "");
}

// parentPath 返回相对路径的父目录，根目录返回空字符串。
function parentPath(path) {
  const parts = String(path || "").split("/").filter(Boolean);
  parts.pop();
  return parts.join("/");
}

// openGallery 打开灯箱并展示指定索引的图片。
function openGallery(index) {
  if (!galleryImages.length) return;
  galleryIndex = Math.max(0, Math.min(index, galleryImages.length - 1));
  updateLightbox();
  document.getElementById("lightbox").hidden = false;
}

// closeLightbox 关闭图片灯箱。
function closeLightbox() {
  const box = document.getElementById("lightbox");
  if (box) box.hidden = true;
  galleryIndex = -1;
}

// closeLightboxBackdrop 点击遮罩空白处关闭灯箱。
function closeLightboxBackdrop(event) {
  if (event.target === event.currentTarget) closeLightbox();
}

// galleryStep 在灯箱中前后切换图片。
function galleryStep(delta) {
  if (!galleryImages.length || galleryIndex < 0) return;
  const next = galleryIndex + delta;
  if (next < 0 || next >= galleryImages.length) return;
  galleryIndex = next;
  updateLightbox();
}

// updateLightbox 根据当前索引刷新灯箱内容与导航状态。
function updateLightbox() {
  const item = galleryImages[galleryIndex];
  if (!item) return;
  const img = document.getElementById("lightboxImage");
  img.src = item.rawUrl;
  img.alt = item.name;
  document.getElementById("lightboxTitle").textContent = item.name;
  document.getElementById("lightboxCounter").textContent = (galleryIndex + 1) + " / " + galleryImages.length;
  document.getElementById("lightboxPrev").disabled = galleryIndex <= 0;
  document.getElementById("lightboxNext").disabled = galleryIndex >= galleryImages.length - 1;
}

// copyLightboxImage 复制灯箱中当前预览的图片。
async function copyLightboxImage(event) {
  event.stopPropagation();
  const btn = event.currentTarget;
  const img = document.getElementById("lightboxImage");
  if (!img || !img.src) return;
  try {
    btn.disabled = true;
    const blob = await imageToPngBlob(img);
    await navigator.clipboard.write([new ClipboardItem({ "image/png": blob })]);
    flashCopied(btn);
  } catch (err) {
    flashCopied(btn, "复制失败");
    console.error(err);
  } finally {
    btn.disabled = false;
  }
}

document.addEventListener("keydown", event => {
  const box = document.getElementById("lightbox");
  if (!box || box.hidden) return;
  if (event.key === "Escape") closeLightbox();
  else if (event.key === "ArrowLeft") galleryStep(-1);
  else if (event.key === "ArrowRight") galleryStep(1);
});

// renderPathBar 显示可复制的本地绝对路径，便于在访达/资源管理器中定位。
function renderPathBar(data) {
  const abs = data.absPath || "";
  if (!abs) return "";
  const folderPath = data.type === "dir" ? abs : abs.replace(/[/\\][^/\\]+$/, "") || abs;
  const showPath = data.type === "dir" ? abs : folderPath;
  return "<div class='path-bar'>"
    + "<code title='" + escAttr(showPath) + "'>" + esc(showPath) + "</code>"
    + "<button type='button' onclick='copyText(event,\"" + escJS(showPath) + "\")'>复制路径</button>"
    + "</div>";
}

// renderMarkdownPreview 渲染 Markdown 正文和右侧悬浮目录。
function renderMarkdownPreview(src) {
  const rendered = renderMarkdown(src);
  return "<div class='markdown-shell'><div class='markdown'>" + rendered.html + "</div>" + renderToc(rendered.headings) + "</div>";
}

// isTableRow 判断一行是否像 Markdown 表格行。
function isTableRow(line) {
  const trimmed = line.trim();
  return trimmed.includes("|") && (trimmed.startsWith("|") || trimmed.endsWith("|") || trimmed.split("|").length > 2);
}

// isTableSeparator 判断一行是否为 Markdown 表格分隔行。
function isTableSeparator(line) {
  return /^\|?\s*:?-{3,}:?\s*(\|\s*:?-{3,}:?\s*)+\|?$/.test(line.trim());
}

// parseTableRow 将 Markdown 表格行拆分为单元格。
function parseTableRow(line) {
  const trimmed = line.trim();
  let cells = trimmed.split("|").map(c => c.trim());
  if (trimmed.startsWith("|")) cells = cells.slice(1);
  if (trimmed.endsWith("|")) cells = cells.slice(0, -1);
  return cells;
}

// renderTableBlock 将 Markdown 表格块渲染为 HTML table。
function renderTableBlock(tableLines) {
  const header = parseTableRow(tableLines[0]);
  const bodyLines = tableLines.slice(2);
  let html = "<table><thead><tr>";
  for (const cell of header) html += "<th>" + inline(cell) + "</th>";
  html += "</tr></thead><tbody>";
  for (const row of bodyLines) {
    const cells = parseTableRow(row);
    html += "<tr>";
    for (const cell of cells) html += "<td>" + inline(cell) + "</td>";
    html += "</tr>";
  }
  html += "</tbody></table>";
  return html;
}

// renderMarkdown 将常用 Markdown 内容转换为预览 HTML，并收集标题。
function renderMarkdown(src) {
  const lines = esc(src).split(/\r?\n/);
  let out = [];
  let headings = [];
  let inList = false;
  let inCode = false;
  let code = [];
  const tick = String.fromCharCode(96);
  const fence = tick + tick + tick;
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.startsWith(fence)) {
      if (inCode) { out.push(copyBlock("pre", code.join("\n"))); code = []; inCode = false; }
      else { if (inList) { out.push("</ul>"); inList = false; } inCode = true; }
      continue;
    }
    if (inCode) { code.push(line); continue; }
    if (isTableRow(line) && i + 1 < lines.length && isTableSeparator(lines[i + 1])) {
      if (inList) { out.push("</ul>"); inList = false; }
      const tableLines = [line, lines[i + 1]];
      i++;
      while (i + 1 < lines.length && isTableRow(lines[i + 1]) && !isTableSeparator(lines[i + 1])) {
        i++;
        tableLines.push(lines[i]);
      }
      out.push(renderTableBlock(tableLines));
      continue;
    }
    if (line.startsWith("### ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(headingBlock("h3", 3, line.slice(4), headings)); }
    else if (line.startsWith("## ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(headingBlock("h2", 2, line.slice(3), headings)); }
    else if (line.startsWith("# ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(headingBlock("h1", 1, line.slice(2), headings)); }
    else if (line.startsWith("- ")) { if (!inList) { out.push("<ul>"); inList = true; } out.push(copyBlock("li", line.slice(2))); }
    else if (line.trim() === "") { if (inList) { out.push("</ul>"); inList = false; } }
    else if (line.includes("|")) { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("p", line)); }
    else { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("p", line)); }
  }
  if (inList) out.push("</ul>");
  return { html: out.join(""), headings };
}

// headingBlock 渲染标题块，并登记到目录导航。
function headingBlock(tag, level, text, headings) {
  const id = "heading-" + headings.length;
  headings.push({ id, level, text: decodeEntities(text) });
  return "<" + tag + " id='" + id + "' data-level='" + level + "' class='copyable'>" + inline(text) + "<button class='copy-btn' onclick='copySection(event)'>复制</button></" + tag + ">";
}

// renderToc 生成 Markdown 右侧悬浮目录导航。
function renderToc(headings) {
  if (!headings.length) return "";
  const levelClass = h => h.level === 1 ? "font-semibold text-slate-200" : h.level === 2 ? "pl-3" : "pl-6 text-xs";
  const links = headings.map(h => "<a class='block rounded px-1 py-1 text-sm leading-snug text-slate-400 hover:bg-slate-800 hover:text-blue-300 " + levelClass(h) + "' href='#" + h.id + "'>" + esc(h.text) + "</a>").join("");
  return "<nav class='toc sticky top-5 max-h-[calc(100vh-7rem)] overflow-auto p-1'><div class='mb-2 text-sm font-semibold text-slate-200'>目录</div>" + links + "</nav>";
}

// copyBlock 为 Markdown 块包裹复制按钮。
function copyBlock(tag, text) {
  const value = decodeEntities(text);
  return "<" + tag + " class='copyable'>" + inline(text) + "<button class='copy-btn' onclick='copyText(event,\"" + escJS(value) + "\")'>复制</button></" + tag + ">";
}

// copySection 复制标题下直到同级或更高级标题前的全部内容。
function copySection(event) {
  event.stopPropagation();
  const heading = event.currentTarget.closest("h1,h2,h3");
  if (!heading) return;
  const level = parseInt(heading.dataset.level || heading.tagName.slice(1), 10);
  const lines = [];
  let node = heading.nextElementSibling;
  while (node) {
    const tag = node.tagName;
    if (tag === "H1" || tag === "H2" || tag === "H3") {
      const nextLevel = parseInt(node.dataset.level || tag.slice(1), 10);
      if (nextLevel <= level) break;
    }
    const text = blockText(node);
    if (text) lines.push(text);
    node = node.nextElementSibling;
  }
  copyFeedback(event.currentTarget, lines.join("\n\n"));
}

// blockText 提取单个渲染块的可复制纯文本。
function blockText(el) {
  const clone = el.cloneNode(true);
  clone.querySelectorAll(".copy-btn").forEach(btn => btn.remove());
  if (clone.tagName === "UL") {
    return Array.from(clone.querySelectorAll("li")).map(li => li.textContent.trim()).filter(Boolean).join("\n");
  }
  if (clone.tagName === "TABLE") {
    return Array.from(clone.querySelectorAll("tr")).map(row =>
      Array.from(row.querySelectorAll("th,td")).map(cell => cell.textContent.trim()).join("\t")
    ).join("\n");
  }
  if (clone.tagName === "PRE") return clone.textContent.replace(/\s+$/, "");
  return clone.textContent.replace(/\s+/g, " ").trim();
}

// copyFeedback 写入剪贴板并在按钮上显示短暂反馈。
function copyFeedback(btn, text) {
  navigator.clipboard.writeText(text).then(() => {
    flashCopied(btn);
  });
}

// flashCopied 在按钮上短暂显示“已复制”反馈。
function flashCopied(btn, label) {
  const target = btn.querySelector("span") || btn;
  const old = target.textContent;
  target.textContent = label || "已复制";
  setTimeout(() => target.textContent = old, 900);
}

// copyText 复制单个 Markdown 块文本，并显示短暂反馈。
function copyText(event, text) {
  event.stopPropagation();
  copyFeedback(event.currentTarget, text);
}

// copyImage 将当前预览图片以 PNG 写入系统剪贴板。
async function copyImage(event) {
  event.stopPropagation();
  const btn = event.currentTarget;
  const img = document.getElementById("previewImage");
  if (!img || !img.src) return;
  try {
    btn.disabled = true;
    const blob = await imageToPngBlob(img);
    await navigator.clipboard.write([new ClipboardItem({ "image/png": blob })]);
    flashCopied(btn);
  } catch (err) {
    flashCopied(btn, "复制失败");
    console.error(err);
  } finally {
    btn.disabled = false;
  }
}

// imageToPngBlob 将图片元素绘制到 canvas 后导出为 PNG Blob（剪贴板兼容性更好）。
function imageToPngBlob(img) {
  return new Promise((resolve, reject) => {
    const draw = () => {
      try {
        const canvas = document.createElement("canvas");
        canvas.width = img.naturalWidth || img.width;
        canvas.height = img.naturalHeight || img.height;
        if (!canvas.width || !canvas.height) {
          reject(new Error("image not ready"));
          return;
        }
        const ctx = canvas.getContext("2d");
        ctx.drawImage(img, 0, 0);
        canvas.toBlob(b => b ? resolve(b) : reject(new Error("toBlob failed")), "image/png");
      } catch (err) {
        reject(err);
      }
    };
    if (img.complete && img.naturalWidth) draw();
    else {
      img.addEventListener("load", draw, { once: true });
      img.addEventListener("error", () => reject(new Error("image load failed")), { once: true });
    }
  });
}

// decodeEntities 将转义后的 HTML 文本还原为可复制的纯文本。
function decodeEntities(s) {
  const el = document.createElement("textarea");
  el.innerHTML = s;
  return el.value;
}

// inline 渲染生成文档中常见的行内 Markdown 标记。
function inline(s) {
  const tick = String.fromCharCode(96);
  const codePattern = new RegExp(tick + "([^" + tick + "]+)" + tick, "g");
  return s.replace(codePattern, "<code>$1</code>").replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
}

// toggleAutoRefresh 开启或关闭目录树自动刷新。
function toggleAutoRefresh() {
  if (document.getElementById("autoRefresh").checked) startAutoRefresh();
  else stopAutoRefresh();
}

// startAutoRefresh 每 5 秒刷新一次目录树。
function startAutoRefresh() {
  stopAutoRefresh();
  refreshTimer = setInterval(loadTree, 5000);
  document.getElementById("refreshState").textContent = "5 秒";
}

// stopAutoRefresh 停止自动刷新并更新工具栏状态。
function stopAutoRefresh() {
  if (refreshTimer) clearInterval(refreshTimer);
  refreshTimer = null;
  document.getElementById("refreshState").textContent = "已关闭";
}

// formatSize 将字节数转换为简洁的可读大小。
function formatSize(n) {
  if (!n) return "";
  if (n < 1024) return n + " B";
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + " KB";
  return (n / 1024 / 1024).toFixed(1) + " MB";
}

// esc 在写入 HTML 前转义文本。
function esc(s) { return String(s || "").replace(/[&<>"']/g, c => ({ "&":"&amp;", "<":"&lt;", ">":"&gt;", "\"":"&quot;", "'":"&#39;" }[c])); }
// escAttr 转义 HTML 属性中的文本。
function escAttr(s) { return esc(s); }
// escJS 转义嵌入行内 JavaScript 字符串的路径。
function escJS(s) { return String(s || "").replace(/\\/g, "\\\\").replace(/"/g, "\\\""); }

loadTree();
openReadme();
startAutoRefresh();
</script>
</body>
</html>`
