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
    html, body { height:100%; overflow-x:hidden; }
    body {
      margin:0;
      display:flex; flex-direction:column;
      background:rgb(2 6 23);
    }
    * {
      scrollbar-width:thin;
      scrollbar-color:rgb(71 85 105 / 0.85) transparent;
    }
    *::-webkit-scrollbar {
      width:8px; height:8px;
    }
    *::-webkit-scrollbar-track {
      background:transparent;
    }
    *::-webkit-scrollbar-thumb {
      background:rgb(71 85 105 / 0.75);
      border-radius:999px;
      border:2px solid transparent;
      background-clip:padding-box;
    }
    *::-webkit-scrollbar-thumb:hover {
      background:rgb(100 116 139 / 0.95);
      border:2px solid transparent;
      background-clip:padding-box;
    }
    *::-webkit-scrollbar-corner { background:transparent; }
    .app-header {
      flex:0 0 auto; width:100%;
      border-bottom:1px solid rgb(51 65 85 / 0.7);
      background:
        linear-gradient(180deg, rgb(15 23 42) 0%, rgb(8 15 30) 100%);
      box-shadow:0 8px 24px rgb(0 0 0 / 0.18);
    }
    .app-header-inner {
      width:1280px; max-width:100%; margin:0 auto;
      min-height:72px; padding:14px 28px; box-sizing:border-box;
      display:flex; align-items:center; justify-content:space-between; gap:24px;
    }
    .brand {
      display:flex; align-items:center; gap:14px; min-width:0;
    }
    .brand img {
      height:42px; width:auto; flex-shrink:0;
      filter:drop-shadow(0 2px 8px rgb(56 189 248 / 0.18));
    }
    .brand-copy { min-width:0; display:flex; flex-direction:column; gap:2px; }
    .brand-copy strong {
      display:block; font-size:20px; font-weight:700; letter-spacing:0.04em;
      color:rgb(248 250 252); line-height:1.15;
    }
    .brand-copy span {
      display:block; font-size:12px; color:rgb(148 163 184); letter-spacing:0.02em;
    }
    .app-frame {
      width:1280px; max-width:100%; margin:0 auto;
      flex:1 1 auto; min-height:0;
      display:flex; flex-direction:column;
      background:rgb(2 6 23);
      border:0;
    }
    .app-frame > main { flex:1; min-height:0; }
    .content-panel { width:100%; }
    .app-footer {
      flex:0 0 auto; width:100%;
      border-top:1px solid rgb(51 65 85 / 0.7);
      background:linear-gradient(180deg, rgb(8 15 30) 0%, rgb(15 23 42) 100%);
    }
    .app-footer-inner {
      width:1280px; max-width:100%; margin:0 auto;
      padding:22px 28px 20px; box-sizing:border-box;
      display:grid; grid-template-columns:minmax(0,1.6fr) minmax(0,0.8fr);
      gap:28px; align-items:start;
    }
    .footer-brand {
      display:flex; flex-direction:column; gap:8px; min-width:0;
    }
    .footer-brand strong {
      font-size:15px; font-weight:700; letter-spacing:0.03em; color:rgb(248 250 252);
    }
    .footer-brand p {
      margin:0; font-size:12px; line-height:1.6; color:rgb(148 163 184);
    }
    .footer-brand code {
      color:rgb(186 230 253); font-size:12px;
    }
    .footer-col h3 {
      margin:0 0 10px; font-size:12px; font-weight:600; letter-spacing:0.06em;
      text-transform:uppercase; color:rgb(100 116 139);
    }
    .footer-col ul {
      list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:8px;
    }
    .footer-col a, .footer-col button {
      display:inline-flex; align-items:center; gap:6px;
      border:0; padding:0; background:transparent; cursor:pointer;
      color:rgb(203 213 225); font:inherit; font-size:13px; text-decoration:none;
      transition:color .15s ease;
    }
    .footer-col a:hover, .footer-col button:hover { color:rgb(125 211 252); }
    .footer-meta {
      grid-column:1 / -1;
      margin-top:4px; padding-top:16px;
      border-top:1px solid rgb(51 65 85 / 0.55);
      font-size:12px; color:rgb(100 116 139);
    }
    @media (max-width: 800px) {
      .app-footer-inner { grid-template-columns:1fr; gap:20px; padding:18px 16px 16px; }
    }
    .layout-main {
      display:flex; flex-direction:row; height:100%; min-height:0; overflow:hidden;
    }
    .layout-aside {
      width:var(--aside-width, 300px); flex:0 0 auto; min-width:180px; max-width:60%;
      overflow:auto; border-right:0; background:rgb(15 23 42); padding:12px;
    }
    .layout-splitter {
      flex:0 0 6px; width:6px; cursor:col-resize; position:relative;
      background:transparent; touch-action:none;
    }
    .layout-splitter::after {
      content:""; position:absolute; top:0; bottom:0; left:-4px; right:-4px;
    }
    .layout-splitter:hover,
    .layout-splitter.is-dragging { background:rgb(56 189 248 / 0.55); }
    .layout-content {
      flex:1 1 auto; min-width:0; overflow:auto; overflow-x:hidden; padding:24px;
    }
    body.is-resizing { cursor:col-resize; user-select:none; }
    body.is-resizing iframe, body.is-resizing img, body.is-resizing video { pointer-events:none; }
    @media (max-width: 800px) {
      .layout-main { flex-direction:column; }
      .layout-aside {
        width:100% !important; max-width:none; min-width:0; flex:0 0 38vh;
        border-right:0; border-bottom:1px solid rgb(30 41 59);
      }
      .layout-splitter { display:none; }
      .layout-content { flex:1 1 auto; min-height:0; }
    }
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
    .markdown .md-image { margin:1em 0 1.2em; padding:0; border:0; background:transparent; }
    .markdown .md-image img {
      display:block; width:100%; max-width:100%; height:auto;
      border-radius:12px; background:rgb(15 23 42);
    }
    @media (max-width: 1050px) { .markdown-shell { grid-template-columns:1fr; } .toc { position:static; max-height:none; order:-1; } }
    .header-actions {
      display:flex; align-items:center; gap:4px;
      font-size:14px; letter-spacing:0.01em;
    }
    .header-actions > * + * { margin-left:2px; }
    .header-link {
      display:inline-flex; align-items:center; gap:8px;
      padding:9px 14px; border:0; background:transparent; cursor:pointer;
      color:rgb(203 213 225); font:inherit; line-height:1;
      border-radius:10px; transition:color .15s ease, background .15s ease;
    }
    .header-link:hover { color:rgb(255 255 255); background:rgb(148 163 184 / 0.1); }
    .header-link:focus-visible { outline:1px solid rgb(56 189 248 / 0.5); outline-offset:2px; }
    .header-link svg { width:16px; height:16px; stroke-width:1.6; opacity:.9; }
    .header-link:hover svg { opacity:1; }
    .header-sep {
      width:1px; height:18px; margin:0 10px;
      background:rgb(71 85 105 / 0.9);
    }
    .header-auto {
      display:inline-flex; align-items:center; gap:10px;
      padding:6px 6px 6px 12px; cursor:pointer;
      color:rgb(148 163 184); font:inherit; line-height:1;
      border-radius:999px; transition:color .15s ease, background .15s ease;
    }
    .header-auto:hover { color:rgb(226 232 240); background:rgb(148 163 184 / 0.08); }
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
      .app-header-inner { min-height:64px; padding:12px 16px; }
      .brand img { height:36px; }
      .brand-copy strong { font-size:17px; }
      .brand-copy span { display:none; }
      .header-link span, .header-auto > span:first-of-type { display:none; }
      .header-sep { margin:0 4px; }
      .header-link { padding:8px 10px; }
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
    .file-list col.type { width:4.5rem; }
    .file-list col.size { width:5rem; }
    .file-list th, .file-list td {
      border:1px solid rgb(51 65 85); padding:9px 10px; vertical-align:middle; text-align:left;
    }
    .file-list th {
      background:rgb(15 23 42); color:rgb(148 163 184); font-size:12px; font-weight:500;
    }
    .file-list td { color:rgb(203 213 225); font-size:14px; }
    .file-list tr.file-row { cursor:pointer; }
    .file-list tr.file-row:hover td { background:rgb(30 41 59 / 0.7); color:rgb(224 242 254); }
    .file-list .name { overflow-wrap:anywhere; word-break:break-word; white-space:normal; }
    .file-list .type, .file-list .size {
      color:rgb(148 163 184); font-size:12px; white-space:nowrap; vertical-align:top;
    }
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
<body class="text-slate-100">
<header class="app-header">
  <div class="app-header-inner">
    <div class="brand">
      <img src="/assets/coor-logo.svg" alt="Coor">
      <div class="brand-copy">
        <strong>Conductor</strong>
        <span>跨境电商智能编排与交付预览</span>
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
  </div>
</header>
<div class="app-frame">
<main id="layoutMain" class="layout-main">
  <aside id="layoutAside" class="layout-aside">
    <ul id="tree" class="tree"></ul>
  </aside>
  <div id="layoutSplitter" class="layout-splitter" role="separator" aria-orientation="vertical" aria-label="调整左右栏宽度" title="拖动调整左右宽度"></div>
  <section id="layoutContent" class="layout-content">
    <div id="content" class="text-slate-500">请选择左侧文件或目录。</div>
  </section>
</main>
</div>
<footer class="app-footer">
  <div class="app-footer-inner">
    <div class="footer-brand">
      <strong>Conductor</strong>
      <p>面向跨境电商的 AI 选品到上架工作流。本页用于浏览与验收 <code>output/</code> 中的文案、图片与视频成果。</p>
    </div>
    <div class="footer-col">
      <h3>产品</h3>
      <ul>
        <li><button type="button" onclick="openReadme()">使用说明</button></li>
      </ul>
    </div>
    <div class="footer-meta">
      <span>© <span id="footerYear"></span> Conductor</span>
    </div>
  </div>
</footer>
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
let markdownAssetMode = "output"; // output: /raw ; doc: /doc-asset（README）
let markdownBasePath = "";

// openReadme 加载项目 README，作为用户使用说明预览。
async function openReadme() {
  activePath = "__readme__";
  markdownAssetMode = "doc";
  markdownBasePath = "";
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
  if (type === "dir") return "📁";
  if (type === "markdown") return "📝";
  if (type === "image") return "🖼️";
  if (type === "video") return "🎬";
  if (type === "json") return "{}";
  return "📄";
}

// typeLabel 将内部类型转为界面展示用的中文名称。
function typeLabel(type) {
  if (type === "dir") return "目录";
  if (type === "markdown") return "Markdown";
  if (type === "image") return "图片";
  if (type === "video") return "视频";
  if (type === "json") return "JSON";
  if (type === "text") return "文本";
  return type || "文件";
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
  markdownAssetMode = "output";
  markdownBasePath = parentPath(path);
  markActive();
  const res = await fetch("/api/file?path=" + encodeURIComponent(path));
  if (!res.ok) {
    document.getElementById("content").innerHTML = "<div class='text-slate-500'>读取失败</div>";
    return;
  }
  const data = await res.json();
  // 目录下仅有一个文件时直接打开，减少多余点击。
  if (data.type === "dir") {
    const onlyFile = soleFileChild(data.children || []);
    if (onlyFile) {
      return openPath(onlyFile.path);
    }
  }
  renderContent(data);
}

// soleFileChild 若目录恰好只有一个非目录子项，返回该文件节点，否则返回 null。
function soleFileChild(children) {
  if (!children || children.length !== 1) return null;
  const only = children[0];
  return only && only.type !== "dir" ? only : null;
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
  return "<div class='content-panel'>" + header + body + "</div>";
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
    + "<td class='type'>" + esc(typeLabel(n.type)) + "</td>"
    + "<td class='size'>" + esc(formatSize(n.size || 0) || "—") + "</td></tr>"
  ).join("");
  return "<table class='file-list'><colgroup><col class='name'><col class='type'><col class='size'></colgroup>"
    + "<thead><tr><th class='name'>名称</th><th class='type'>类型</th><th class='size'>大小</th></tr></thead><tbody>"
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
    if (isImageLine(line)) {
      if (inList) { out.push("</ul>"); inList = false; }
      out.push(renderImageBlock(line));
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

// isImageLine 判断一行是否为独立的 Markdown 图片语法。
function isImageLine(line) {
  return /^!\[[^\]]*\]\([^)]+\)\s*$/.test(String(line || "").trim());
}

// renderImageBlock 将 Markdown 图片行渲染为可加载的 figure/img。
function renderImageBlock(line) {
  const m = String(line || "").trim().match(/^!\[([^\]]*)\]\(([^)]+)\)\s*$/);
  if (!m) return copyBlock("p", line);
  const alt = m[1];
  const src = resolveMarkdownImageSrc(m[2]);
  return "<figure class='md-image'><img src='" + escAttr(src) + "' alt='" + alt + "' loading='lazy'></figure>";
}

// resolveMarkdownImageSrc 将 Markdown 图片路径解析为浏览器可访问的 URL。
function resolveMarkdownImageSrc(src) {
  let path = decodeEntities(src || "").trim().replace(/^<|>$/g, "");
  if (!path) return "";
  if (/^(https?:|data:|blob:|\/)/i.test(path)) return path;
  path = path.replace(/^\.\//, "");
  if (markdownAssetMode === "doc") {
    return "/doc-asset?path=" + encodeURIComponent(path);
  }
  const parts = [];
  for (const part of (markdownBasePath + "/" + path).split("/")) {
    if (!part || part === ".") continue;
    if (part === "..") parts.pop();
    else parts.push(part);
  }
  return "/raw?path=" + encodeURIComponent(parts.join("/"));
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
  return s
    .replace(/!\[([^\]]*)\]\(([^)]+)\)/g, (_, alt, src) =>
      "<img class='inline-md-image' src='" + escAttr(resolveMarkdownImageSrc(src)) + "' alt='" + alt + "' loading='lazy'>"
    )
    .replace(codePattern, "<code>$1</code>")
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
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

const ASIDE_WIDTH_KEY = "conductor.asideWidth";
const ASIDE_MIN = 180;
const ASIDE_MAX_RATIO = 0.6;

// applyAsideWidth 设置左侧目录栏宽度，并写入 CSS 变量。
function applyAsideWidth(px) {
  const main = document.getElementById("layoutMain");
  if (!main || window.matchMedia("(max-width: 800px)").matches) return;
  const max = Math.floor(main.clientWidth * ASIDE_MAX_RATIO);
  const width = Math.max(ASIDE_MIN, Math.min(max, Math.round(px)));
  document.documentElement.style.setProperty("--aside-width", width + "px");
  return width;
}

// restoreAsideWidth 从本地存储恢复用户上次调整的左右栏宽度。
function restoreAsideWidth() {
  const saved = parseInt(localStorage.getItem(ASIDE_WIDTH_KEY) || "300", 10);
  applyAsideWidth(Number.isFinite(saved) ? saved : 300);
}

// initAsideResize 启用左右分栏拖拽调整宽度。
function initAsideResize() {
  const splitter = document.getElementById("layoutSplitter");
  const main = document.getElementById("layoutMain");
  if (!splitter || !main) return;
  restoreAsideWidth();
  let dragging = false;
  const onMove = (event) => {
    if (!dragging) return;
    const x = event.touches ? event.touches[0].clientX : event.clientX;
    const width = applyAsideWidth(x - main.getBoundingClientRect().left);
    if (width) localStorage.setItem(ASIDE_WIDTH_KEY, String(width));
  };
  const onUp = () => {
    if (!dragging) return;
    dragging = false;
    splitter.classList.remove("is-dragging");
    document.body.classList.remove("is-resizing");
    window.removeEventListener("pointermove", onMove);
    window.removeEventListener("pointerup", onUp);
    window.removeEventListener("touchmove", onMove);
    window.removeEventListener("touchend", onUp);
  };
  const onDown = (event) => {
    if (window.matchMedia("(max-width: 800px)").matches) return;
    event.preventDefault();
    dragging = true;
    splitter.classList.add("is-dragging");
    document.body.classList.add("is-resizing");
    window.addEventListener("pointermove", onMove);
    window.addEventListener("pointerup", onUp);
    window.addEventListener("touchmove", onMove, { passive: false });
    window.addEventListener("touchend", onUp);
  };
  splitter.addEventListener("pointerdown", onDown);
  splitter.addEventListener("touchstart", onDown, { passive: false });
  window.addEventListener("resize", () => {
    const current = parseInt(getComputedStyle(document.documentElement).getPropertyValue("--aside-width"), 10);
    applyAsideWidth(Number.isFinite(current) ? current : 300);
  });
}

loadTree();
openReadme();
startAutoRefresh();
initAsideResize();
document.getElementById("footerYear").textContent = String(new Date().getFullYear());
</script>
</body>
</html>`
