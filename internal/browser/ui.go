package browser

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Conductor 物料浏览器</title>
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
      width:1440px; max-width:100%; margin:0 auto;
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
      width:1440px; max-width:100%; margin:0 auto;
      flex:1 1 auto; min-height:0;
      display:flex; flex-direction:column;
      background:rgb(2 6 23);
      border:0;
    }
    .app-frame > main { flex:1; min-height:0; }
    .content-panel { width:100%; }
    .breadcrumb {
      display:flex; flex-wrap:wrap; align-items:center; gap:6px;
      margin:0 0 14px; font-size:13px; line-height:1.45;
    }
    .breadcrumb .crumb {
      border:0; background:transparent; padding:0; font:inherit; max-width:100%;
      overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
    }
    .breadcrumb .crumb.link {
      cursor:pointer; color:rgb(125 211 252);
      transition:color .15s ease;
    }
    .breadcrumb .crumb.link:hover { color:rgb(186 230 253); }
    .breadcrumb .crumb.current { color:rgb(148 163 184); cursor:default; }
    .breadcrumb .crumb-sep {
      color:rgb(71 85 105); user-select:none; flex-shrink:0;
    }
    .app-footer {
      flex:0 0 auto; width:100%;
      border-top:1px solid rgb(51 65 85 / 0.7);
      background:linear-gradient(180deg, rgb(8 15 30) 0%, rgb(15 23 42) 100%);
    }
    .app-footer-inner {
      width:1440px; max-width:100%; margin:0 auto;
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
    .markdown-shell { display:block; min-width:0; }
    /* 宽屏：目录固定在内容区外，不占正文宽度；窄屏：收进内容区内并可折叠 */
    .toc-dock {
      position:fixed; z-index:35;
      top:96px; right:max(12px, calc((100vw - 1440px) / 2 - 236px));
      width:224px; max-height:calc(100vh - 120px);
      display:flex; flex-direction:column; gap:0;
      padding:0; box-sizing:border-box; overflow:hidden;
      border:1px solid rgb(71 85 105 / 0.55);
      border-radius:14px;
      background:linear-gradient(165deg, rgb(22 32 51 / 0.96) 0%, rgb(12 20 36 / 0.96) 100%);
      backdrop-filter:blur(12px);
      box-shadow:
        0 0 0 1px rgb(255 255 255 / 0.03) inset,
        0 12px 32px rgb(0 0 0 / 0.35);
    }
    .toc-dock.is-collapsed {
      width:auto; max-height:none; border-radius:999px;
      background:rgb(15 23 42 / 0.92);
    }
    .toc-dock.is-collapsed .toc-body,
    .toc-dock.is-collapsed .toc-head-meta { display:none; }
    .toc-head {
      display:flex; align-items:center; gap:8px;
      padding:10px 12px;
      border-bottom:1px solid rgb(51 65 85 / 0.55);
      background:rgb(15 23 42 / 0.45);
    }
    .toc-dock.is-collapsed .toc-head {
      border-bottom:0; padding:0; background:transparent;
    }
    .toc-head-meta {
      display:flex; flex-direction:column; gap:1px; min-width:0; flex:1;
    }
    .toc-head-title {
      font-size:12px; font-weight:700; letter-spacing:0.08em;
      text-transform:uppercase; color:rgb(186 230 253);
    }
    .toc-head-count {
      font-size:11px; color:rgb(100 116 139); font-variant-numeric:tabular-nums;
    }
    .toc-toggle {
      display:inline-flex; align-items:center; justify-content:center;
      flex:0 0 auto; width:28px; height:28px; margin:0; padding:0;
      border:1px solid rgb(71 85 105 / 0.55); border-radius:8px;
      background:rgb(30 41 59 / 0.55); color:rgb(148 163 184);
      cursor:pointer; transition:color .15s ease, background .15s ease, border-color .15s ease;
    }
    .toc-toggle:hover {
      color:rgb(186 230 253); background:rgb(51 65 85 / 0.7);
      border-color:rgb(56 189 248 / 0.35);
    }
    .toc-toggle svg { width:14px; height:14px; display:block; }
    .toc-dock.is-collapsed .toc-toggle {
      width:auto; height:auto; gap:8px; padding:8px 12px; border-radius:999px;
      border-color:rgb(71 85 105 / 0.65);
    }
    .toc-dock.is-collapsed .toc-toggle-label { display:inline; font-size:12px; font-weight:600; color:inherit; }
    .toc-toggle-label { display:none; }
    .toc-body {
      overflow:auto; min-height:0; padding:8px 8px 10px;
      display:flex; flex-direction:column; gap:2px;
    }
    .toc-body a {
      position:relative; border-radius:8px;
      padding:7px 10px 7px 12px;
      font-size:12.5px; line-height:1.35; color:rgb(148 163 184);
      text-decoration:none;
      border-left:2px solid transparent;
      transition:color .12s ease, background .12s ease, border-color .12s ease;
      display:flex; align-items:center;
      box-sizing:border-box; min-height:32px;
    }
    .toc-body a .toc-label {
      display:-webkit-box; -webkit-line-clamp:2; -webkit-box-orient:vertical;
      overflow:hidden; text-overflow:ellipsis;
      min-width:0; width:100%;
    }
    .toc-body a:hover {
      background:rgb(56 189 248 / 0.08); color:rgb(226 232 240);
      border-left-color:rgb(56 189 248 / 0.45);
    }
    .toc-body a.is-active {
      background:rgb(56 189 248 / 0.12); color:rgb(186 230 253);
      border-left-color:rgb(56 189 248);
    }
    .toc-body a.level-1 {
      margin-top:6px; padding-top:8px; padding-bottom:8px; min-height:36px;
      font-size:13px; font-weight:600; color:rgb(241 245 249);
    }
    .toc-body a.level-1:first-child { margin-top:0; }
    .toc-body a.level-2 { padding-left:16px; color:rgb(163 174 191); }
    .toc-body a.level-3 {
      padding-left:24px; font-size:12px; color:rgb(120 133 153);
    }
    @media (max-width: 1679px) {
      /* 视口不足以把目录放到 1440 定宽外侧时，仍浮在右侧边缘，不挤占正文 */
      .toc-dock { right:12px; }
    }
    @media (max-width: 1100px) {
      /* 窄屏：目录进入内容区顶部，可折叠 */
      .markdown-shell { display:flex; flex-direction:column; gap:14px; }
      .toc-dock {
        position:static; order:-1; z-index:auto;
        width:100%; max-height:none;
        right:auto; top:auto;
        box-shadow:none;
      }
      .toc-dock.is-collapsed {
        width:100%; border-radius:12px;
      }
      .toc-dock.is-collapsed .toc-head { padding:8px 10px; }
      .toc-dock.is-collapsed .toc-head-meta { display:flex; }
      .toc-dock.is-collapsed .toc-toggle { border-radius:8px; padding:0; width:28px; height:28px; }
      .toc-dock.is-collapsed .toc-toggle-label { display:none; }
      .toc-dock:not(.is-collapsed) .toc-body {
        max-height:min(40vh, 320px); overflow:auto;
      }
    }
    .markdown h1, .markdown h2, .markdown h3 { line-height:1.25; scroll-margin-top:18px; position:relative; }
    .markdown h1 { font-size:24px; font-weight:700; color:rgb(248 250 252); border-bottom:1px solid rgb(51 65 85 / 0.7); padding-bottom:8px; margin:0 0 0.85em; }
    .markdown h2 { font-size:21px; font-weight:700; color:rgb(241 245 249); border-bottom:1px solid rgb(51 65 85 / 0.55); padding-bottom:8px; margin:2.2em 0 0.9em; }
    .markdown h3 { font-size:14px; font-weight:600; color:rgb(125 211 252); letter-spacing:0.02em; margin:1.5em 0 0.55em; padding:0; background:transparent; border:0; }
    .markdown h1.copyable, .markdown h2.copyable, .markdown h3.copyable { padding-right:52px; }
    .markdown .copy-btn {
      position:absolute; right:0; top:0.15em; border:0; background:transparent;
      color:rgb(100 116 139); border-radius:4px; padding:2px 6px; font-size:12px; opacity:0; cursor:pointer;
    }
    .markdown h1 .copy-btn, .markdown h2 .copy-btn { top:0.35em; }
    .markdown h1:hover .copy-btn, .markdown h2:hover .copy-btn, .markdown h3:hover .copy-btn { opacity:1; }
    .markdown .copy-btn:hover { color:rgb(125 211 252); background:rgb(148 163 184 / 0.08); }
    .markdown p, .markdown li { color:rgb(203 213 225); font-size:15px; line-height:1.75; }
    .markdown p { margin:0 0 0.9em; padding:0; background:transparent; border:0; }
    .markdown h3 + p, .markdown h3 + ul, .markdown h3 + pre, .markdown h3 + table { margin-top:0; }
    .markdown ul { margin:0 0 0.9em; padding:0 0 0 1.25em; background:transparent; border:0; }
    .markdown li { margin:0.3em 0; }
    .markdown code { background:rgb(30 41 59 / 0.7); padding:1px 5px; border-radius:4px; color:rgb(226 232 240); }
    .markdown pre {
      white-space:pre-wrap; word-break:break-word; overflow-wrap:anywhere;
      margin:0 0 0.9em; padding:12px 0;
      background:transparent; border:0; border-top:1px solid rgb(51 65 85 / 0.45); border-bottom:1px solid rgb(51 65 85 / 0.45);
      color:rgb(226 232 240);
    }
    .markdown table { border-collapse:collapse; width:100%; margin:12px 0; table-layout:fixed; }
    .markdown col.col-narrow { width:3.25rem; }
    .markdown th.col-narrow, .markdown td.col-narrow {
      width:3.25rem; text-align:center; white-space:nowrap;
      overflow-wrap:normal; word-break:normal;
    }
    .markdown th, .markdown td { border:1px solid rgb(51 65 85); padding:8px; vertical-align:top; }
    .markdown th, .markdown td { overflow-wrap:anywhere; word-break:break-word; }
    .markdown th { background:rgb(15 23 42); }
    .markdown .md-image { margin:1em 0 1.2em; padding:0; border:0; background:transparent; }
    .markdown .md-image img {
      display:block; width:100%; max-width:100%; height:auto;
      border-radius:12px; background:rgb(15 23 42);
    }
    .markdown a[id]:empty { display:block; height:0; overflow:hidden; scroll-margin-top:18px; }
    .markdown a[href^="#"] { color:rgb(125 211 252); text-decoration:underline; text-underline-offset:2px; }
    .markdown a[href^="#"]:hover { color:rgb(186 230 253); }
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
    .file-list col.row-actions { width:5.5rem; }
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
    .path-actions.compact {
      flex-shrink:0; display:inline-flex; align-items:stretch;
      border:1px solid rgb(71 85 105 / 0.7); border-radius:8px; overflow:hidden;
      background:rgb(2 6 23 / 0.82);
    }
    .path-actions.compact > button {
      margin:0; border:0; border-radius:0; border-right:1px solid rgb(71 85 105 / 0.55);
      background:transparent; color:rgb(226 232 240);
      padding:5px 10px; font:inherit; font-size:12px; line-height:1.2; cursor:pointer;
      white-space:nowrap;
      transition:color .15s ease, background .15s ease;
    }
    .path-actions.compact > button:last-child { border-right:0; }
    .path-actions.compact > button:hover {
      color:rgb(186 230 253); background:rgb(30 41 59 / 0.85);
    }
    .path-actions.compact > button.danger { color:rgb(254 202 202); }
    .path-actions.compact > button.danger:hover {
      color:rgb(255 255 255); background:rgb(185 28 28 / 0.7);
    }
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
      background:rgb(148 163 184 / 0.12); border-color:rgb(148 163 184 / 0.45); color:rgb(241 245 249);
    }
    .lightbox-btn.danger {
      background:rgb(127 29 29 / 0.45); border-color:rgb(248 113 113 / 0.35); color:rgb(254 202 202);
    }
    .lightbox-btn.danger:hover {
      background:rgb(185 28 28 / 0.7); border-color:rgb(252 165 165 / 0.45); color:rgb(255 255 255);
    }
    .preview-actions {
      position:absolute; right:12px; top:12px;
      display:flex; align-items:center; gap:8px;
    }
    .preview-actions.in-header {
      position:static;
    }
    .preview-action {
      border:1px solid rgb(71 85 105 / 0.7); border-radius:8px;
      background:rgb(2 6 23 / 0.82); color:rgb(226 232 240);
      padding:6px 10px; font:inherit; font-size:12px; cursor:pointer;
      backdrop-filter:blur(6px);
      transition:color .15s ease, background .15s ease, border-color .15s ease;
    }
    .preview-action:hover { color:rgb(186 230 253); border-color:rgb(56 189 248 / 0.4); }
    .preview-action.danger { color:rgb(254 202 202); border-color:rgb(248 113 113 / 0.35); }
    .preview-action.danger:hover {
      color:rgb(255 255 255); background:rgb(185 28 28 / 0.7);
      border-color:rgb(252 165 165 / 0.45);
    }
    .file-list .row-actions {
      width:5.5rem; text-align:right; white-space:nowrap;
    }
    .file-list .row-actions .preview-action {
      padding:4px 8px; font-size:11px; background:transparent;
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
    .confirm-dialog {
      position:fixed; inset:0; z-index:60;
      display:flex; align-items:center; justify-content:center;
      background:rgb(2 6 23 / 0.72); padding:24px;
    }
    .confirm-dialog[hidden] { display:none; }
    .confirm-card {
      width:min(420px, 100%);
      border:1px solid rgb(71 85 105 / 0.8); border-radius:14px;
      background:rgb(15 23 42); box-shadow:0 24px 60px rgb(0 0 0 / 0.45);
      padding:20px 20px 16px; color:rgb(226 232 240);
    }
    .confirm-card h2 {
      margin:0 0 8px; font-size:16px; font-weight:600; color:rgb(248 250 252);
    }
    .confirm-card p {
      margin:0; font-size:13px; line-height:1.55; color:rgb(148 163 184);
      white-space:pre-wrap; word-break:break-word;
    }
    .confirm-card .confirm-target {
      margin-top:10px; padding:8px 10px; border-radius:8px;
      background:rgb(2 6 23 / 0.65); border:1px solid rgb(51 65 85 / 0.8);
      color:rgb(226 232 240); font-size:13px; font-weight:500;
    }
    .confirm-actions {
      display:flex; justify-content:flex-end; gap:8px; margin-top:18px;
    }
    .confirm-actions button {
      border-radius:8px; padding:7px 12px; font:inherit; font-size:13px; cursor:pointer;
      border:1px solid rgb(71 85 105 / 0.8); background:rgb(30 41 59); color:rgb(226 232 240);
    }
    .confirm-actions button:hover { border-color:rgb(100 116 139); }
    .confirm-actions button.danger {
      background:rgb(185 28 28 / 0.85); border-color:rgb(248 113 113 / 0.45); color:rgb(255 255 255);
    }
    .confirm-actions button.danger:hover { background:rgb(220 38 38 / 0.95); }
  </style>
</head>
<body class="text-slate-100">
<header class="app-header">
  <div class="app-header-inner">
    <div class="brand">
      <img src="/assets/coor-logo.svg" alt="Coor">
      <div class="brand-copy">
        <strong>Conductor</strong>
        <span>物料浏览器 · 浏览与验收 output 产出</span>
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
      <p>面向跨境电商的 AI 选品到上架工作流。本页用于浏览与验收 <code>output/</code> 中的文案、图片与视频等物料。</p>
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
        <button type="button" class="lightbox-btn danger" onclick="deleteLightboxImage(event)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M5 7h14M10 7V5h4v2m-5 3v7m4-7v7M7 7l1 12a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2l1-12"/></svg>
          <span>删除</span>
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
<div id="confirmDialog" class="confirm-dialog" hidden role="dialog" aria-modal="true" aria-labelledby="confirmTitle" onclick="settleConfirm(false)">
  <div class="confirm-card" onclick="event.stopPropagation()">
    <h2 id="confirmTitle">确认删除</h2>
    <p id="confirmMessage"></p>
    <div id="confirmTarget" class="confirm-target" hidden></div>
    <div class="confirm-actions">
      <button type="button" id="confirmCancelBtn" onclick="event.stopPropagation(); settleConfirm(false)">取消</button>
      <button type="button" class="danger" id="confirmOkBtn" onclick="event.stopPropagation(); settleConfirm(true)">确认删除</button>
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
let rootName = "output";
let confirmResolver = null;

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
  rootName = data.name || "output";
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
  if (tocSpyCleanup) {
    tocSpyCleanup();
    tocSpyCleanup = null;
  }
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
    contentHeader(data, "<div class='mb-1 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold tracking-tight'>" + esc(data.name) + "</h1>" + sizeHint + "</div>"),
    body
  );
  applyTocCollapsedState();
}

// contentShell 组装右侧预览外壳，避免多层边框套框。
function contentShell(header, body) {
  return "<div class='content-panel'>" + header + body + "</div>";
}

// contentHeader 组装预览区标题区：面包屑 + 标题 + 可复制路径。
function contentHeader(data, titleHtml) {
  return "<div class='mb-5'>" + renderBreadcrumb(data) + titleHtml + renderPathBar(data) + "</div>";
}

// renderBreadcrumb 根据相对路径生成可点击的面包屑导航。
function renderBreadcrumb(data) {
  if (activePath === "__readme__") {
    return "<nav class='breadcrumb' aria-label='面包屑'><span class='crumb current'>使用说明</span></nav>";
  }
  const rel = (data && data.path != null) ? data.path : (activePath || "");
  const parts = String(rel).split("/").filter(Boolean);
  const items = [{ name: rootName || "全部", path: "" }];
  let acc = "";
  for (const part of parts) {
    acc = acc ? acc + "/" + part : part;
    items.push({ name: part, path: acc });
  }
  const html = items.map((item, i) => {
    const isLast = i === items.length - 1;
    if (isLast) return "<span class='crumb current' title='" + escAttr(item.name) + "'>" + esc(item.name) + "</span>";
    return "<button type='button' class='crumb link' title='" + escAttr(item.name) + "' onclick='openPath(\"" + escJS(item.path) + "\")'>" + esc(item.name) + "</button>"
      + "<span class='crumb-sep' aria-hidden='true'>/</span>";
  }).join("");
  return "<nav class='breadcrumb' aria-label='面包屑'>" + html + "</nav>";
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
  const titleRow = "<div class='mb-1'><h1 class='m-0 text-xl font-semibold tracking-tight'>"
    + esc(data.name || "output") + "</h1></div>";
  document.getElementById("content").innerHTML = contentShell(
    contentHeader(data, titleRow),
    body
  );
}

// renderFileList 以带边框表格展示目录中的非图片项。
function renderFileList(items) {
  const rows = items.map(n => {
    const actions = n.type === "dir"
      ? "<button type='button' class='preview-action danger' onclick='event.stopPropagation(); deleteDirectory(\"" + escJS(n.path) + "\",\"" + escJS(n.name) + "\")'>删除</button>"
      : "—";
    return "<tr class='file-row' onclick='openPath(\"" + escJS(n.path) + "\")'>"
      + "<td class='name'>" + iconFor(n.type) + " " + esc(n.name) + "</td>"
      + "<td class='type'>" + esc(typeLabel(n.type)) + "</td>"
      + "<td class='size'>" + esc(formatSize(n.size || 0) || "—") + "</td>"
      + "<td class='row-actions'>" + actions + "</td></tr>";
  }).join("");
  return "<table class='file-list'><colgroup><col class='name'><col class='type'><col class='size'><col class='row-actions'></colgroup>"
    + "<thead><tr><th class='name'>名称</th><th class='type'>类型</th><th class='size'>大小</th><th class='row-actions'>操作</th></tr></thead><tbody>"
    + rows + "</tbody></table>";
}

// renderImageContent 渲染单张图片，并加载同目录图片供左右切换。
async function renderImageContent(data) {
  const size = formatSize(data.size || 0);
  const sizeHint = size ? "<span class='text-xs text-slate-500'>" + esc(size) + "</span>" : "";
  document.getElementById("content").innerHTML = contentShell(
    contentHeader(data, "<div class='mb-1 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold tracking-tight'>" + esc(data.name) + "</h1>" + sizeHint + "</div>"),
    "<div class='relative inline-block max-w-full'>"
      + "<img id='previewImage' class='max-w-full rounded-lg bg-slate-900' src='" + escAttr(data.rawUrl) + "' alt='" + escAttr(data.name) + "'>"
      + "<div class='preview-actions'>"
      + "<button type='button' class='preview-action' onclick='copyImage(event)'>复制图片</button>"
      + "<button type='button' class='preview-action danger' onclick='deleteImage(\"" + escJS(data.path) + "\",\"" + escJS(data.name) + "\")'>删除</button>"
      + "</div></div><div id='siblingGallery' class='mt-5'></div>"
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
  galleryImages = images.map(n => ({
    name: n.name,
    path: n.path,
    size: n.size || 0,
    rawUrl: rawUrlFor(n.path)
  }));
  const shell = document.getElementById("siblingGallery");
  if (!shell) return;
  if (images.length < 2) {
    shell.innerHTML = "";
    return;
  }
  const currentIndex = galleryImages.findIndex(n => n.path === imagePath);
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

// deleteLightboxImage 删除灯箱中当前图片（先确认再删除）。
async function deleteLightboxImage(event) {
  event.stopPropagation();
  const item = galleryImages[galleryIndex];
  if (!item) return;
  await deleteImage(item.path, item.name, { keepLightbox: true });
}

// deleteDirectory 确认后递归删除 output 内的子目录，并打开其父目录。
async function deleteDirectory(path, name) {
  const rel = String(path || "").replace(/\\/g, "/").replace(/^\/+|\/+$/g, "");
  if (!rel) {
    window.alert("不能删除 output 根目录。");
    return;
  }
  const label = name || rel.split("/").pop() || rel;
  const ok = await askConfirm({
    title: "确认删除目录",
    message: "将永久删除该目录及其全部内容，此操作不可恢复。",
    target: label,
    confirmLabel: "确认删除"
  });
  if (!ok) return;

  const res = await fetch("/api/file?path=" + encodeURIComponent(rel) + "&confirm=1", {
    method: "DELETE",
    headers: { "X-Confirm-Delete": "1" }
  });
  let payload = null;
  try { payload = await res.json(); } catch (_) {}
  if (!res.ok) {
    const msg = (payload && payload.error) ? payload.error : ("删除失败（" + res.status + "）");
    window.alert(msg);
    return;
  }
  const parent = (payload && payload.parent != null) ? payload.parent : parentPath(rel);
  closeLightbox();
  await loadTree();
  await openPath(parent);
}

// sameRelPath 比较相对路径，忽略首尾斜杠与反斜杠差异。
function sameRelPath(a, b) {
  const norm = (p) => String(p || "").replace(/\\/g, "/").replace(/^\/+|\/+$/g, "");
  return norm(a) === norm(b);
}

// deleteImage 确认后删除 output 内的图片；有相邻图则直接打开该图片文件（灯箱内删除则继续在灯箱显示）。
async function deleteImage(path, name, options) {
  const keepLightbox = !!(options && options.keepLightbox);
  const label = name || String(path || "").split("/").pop() || path;
  const ok = await askConfirm({
    title: "确认删除图片",
    message: "将永久删除该图片，此操作不可恢复。",
    target: label,
    confirmLabel: "确认删除"
  });
  if (!ok) return;

  const parent = parentPath(path);
  let siblings;
  if (galleryImages.length && galleryImages.some(g => sameRelPath(g.path, path))) {
    siblings = galleryImages.map(g => ({ path: g.path, name: g.name }));
  } else {
    siblings = await listSiblingImages(parent);
  }
  let idx = siblings.findIndex(n => sameRelPath(n.path, path));
  if (idx < 0 && galleryIndex >= 0 && galleryImages[galleryIndex] && sameRelPath(galleryImages[galleryIndex].path, path)) {
    idx = galleryIndex;
  }
  let nextPath = "";
  if (idx >= 0) {
    if (idx + 1 < siblings.length) nextPath = siblings[idx + 1].path; // 优先下一张
    else if (idx - 1 >= 0) nextPath = siblings[idx - 1].path; // 已是最后一张则前一张
  }

  const res = await fetch("/api/file?path=" + encodeURIComponent(path || "") + "&confirm=1", {
    method: "DELETE",
    headers: { "X-Confirm-Delete": "1" }
  });
  let payload = null;
  try { payload = await res.json(); } catch (_) {}
  if (!res.ok) {
    const msg = (payload && payload.error) ? payload.error : ("删除失败（" + res.status + "）");
    window.alert(msg);
    return;
  }
  const fallbackParent = (payload && payload.parent != null) ? payload.parent : parent;

  await loadTree();

  if (!nextPath) {
    closeLightbox();
    await openPath(fallbackParent);
    return;
  }

  // 始终打开相邻图片文件本身，避免停在目录缩略图预览。
  await openPath(nextPath);

  if (keepLightbox) {
    await refreshGalleryImages(parent);
    let i = galleryImages.findIndex(g => sameRelPath(g.path, nextPath));
    if (i < 0) {
      galleryImages = [{
        name: String(nextPath).split("/").pop() || nextPath,
        path: nextPath,
        size: 0,
        rawUrl: rawUrlFor(nextPath)
      }];
      i = 0;
    }
    openGallery(i);
  } else {
    closeLightbox();
  }
}

// listSiblingImages 列出目录下全部图片（保持服务端排序）。
async function listSiblingImages(parent) {
  const res = await fetch("/api/file?path=" + encodeURIComponent(parent || ""));
  if (!res.ok) return [];
  const data = await res.json();
  if (data.type !== "dir") return [];
  return (data.children || []).filter(n => n.type === "image");
}

// refreshGalleryImages 删除后重建灯箱用的同目录图片列表。
async function refreshGalleryImages(parent) {
  const images = await listSiblingImages(parent);
  galleryImages = images.map(n => ({
    name: n.name,
    path: n.path,
    size: n.size || 0,
    rawUrl: rawUrlFor(n.path)
  }));
}

document.addEventListener("keydown", event => {
  const confirmBox = document.getElementById("confirmDialog");
  if (confirmBox && !confirmBox.hidden) {
    if (event.key === "Escape") {
      event.preventDefault();
      settleConfirm(false);
    }
    return;
  }
  const box = document.getElementById("lightbox");
  if (!box || box.hidden) return;
  if (event.key === "Escape") closeLightbox();
  else if (event.key === "ArrowLeft") galleryStep(-1);
  else if (event.key === "ArrowRight") galleryStep(1);
});

// askConfirm 显示页面内确认框（不依赖浏览器原生 confirm，避免被拦截后静默通过/失败）。
function askConfirm(options) {
  const opts = options || {};
  const box = document.getElementById("confirmDialog");
  const titleEl = document.getElementById("confirmTitle");
  const msgEl = document.getElementById("confirmMessage");
  const targetEl = document.getElementById("confirmTarget");
  const okBtn = document.getElementById("confirmOkBtn");
  const cancelBtn = document.getElementById("confirmCancelBtn");
  if (!box || !titleEl || !msgEl || !okBtn || !cancelBtn) return Promise.resolve(false);

  // 若已有未完成确认，先取消旧的，避免 Promise 永久挂起。
  if (confirmResolver) settleConfirm(false);

  titleEl.textContent = opts.title || "请确认";
  msgEl.textContent = opts.message || "";
  if (opts.target) {
    targetEl.hidden = false;
    targetEl.textContent = opts.target;
  } else {
    targetEl.hidden = true;
    targetEl.textContent = "";
  }
  okBtn.textContent = opts.confirmLabel || "确认";
  box.hidden = false;
  setTimeout(() => { try { okBtn.focus(); } catch (_) {} }, 0);

  return new Promise(resolve => {
    confirmResolver = resolve;
  });
}

// settleConfirm 关闭确认框并返回用户选择。
function settleConfirm(ok) {
  const box = document.getElementById("confirmDialog");
  const resolve = confirmResolver;
  confirmResolver = null;
  if (box) box.hidden = true;
  if (resolve) resolve(!!ok);
}

// renderPathBar 显示本地绝对路径，右侧为 Compact 按钮组（删除/复制/打开）。
function renderPathBar(data) {
  const abs = data.absPath || "";
  if (!abs) return "";
  const folderPath = data.type === "dir" ? abs : abs.replace(/[/\\][^/\\]+$/, "") || abs;
  const showPath = data.type === "dir" ? abs : folderPath;
  const rel = (data && data.path != null) ? data.path : "";
  // README 位于 output 外，打开文件夹 API 只允许 output 内路径。
  const canOpen = activePath !== "__readme__";
  const canDeleteDir = data.type === "dir" && !!rel;
  const actions = "<div class='path-actions compact' role='group' aria-label='路径操作'>"
    + (canDeleteDir
      ? "<button type='button' class='danger' onclick='event.stopPropagation(); deleteDirectory(\"" + escJS(rel) + "\",\"" + escJS(data.name || "") + "\")'>删除目录</button>"
      : "")
    + "<button type='button' onclick='event.stopPropagation(); copyText(event,\"" + escJS(showPath) + "\")'>复制路径</button>"
    + (canOpen
      ? "<button type='button' onclick='event.stopPropagation(); openFolder(event,\"" + escJS(rel) + "\")'>打开文件夹</button>"
      : "")
    + "</div>";
  return "<div class='path-bar'>"
    + "<code title='" + escAttr(showPath) + "'>" + esc(showPath) + "</code>"
    + actions
    + "</div>";
}

// openFolder 请求本地服务在操作系统文件管理器中打开对应目录。
async function openFolder(event, relPath) {
  event.stopPropagation();
  const btn = event.currentTarget;
  try {
    btn.disabled = true;
    const res = await fetch("/api/open-folder?path=" + encodeURIComponent(relPath || ""), { method: "POST" });
    if (!res.ok) {
      flashCopied(btn, "打开失败");
      console.error(await res.text());
      return;
    }
    flashCopied(btn, "已打开");
  } catch (err) {
    flashCopied(btn, "打开失败");
    console.error(err);
  } finally {
    btn.disabled = false;
  }
}

// renderMarkdownPreview 渲染 Markdown 正文；目录默认浮在内容区外，窄屏时进入内容区。
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

// isNarrowTableHeader 判断表头是否应使用窄列（如序号）。
function isNarrowTableHeader(text) {
  return /^(序号|展示序号|#|No\.?|ID)$/i.test(String(text || "").trim());
}

// renderTableBlock 将 Markdown 表格块渲染为 HTML table。
function renderTableBlock(tableLines) {
  const header = parseTableRow(tableLines[0]);
  const bodyLines = tableLines.slice(2);
  const narrowFlags = header.map(cell => isNarrowTableHeader(decodeEntities(cell)));
  let html = "<table><colgroup>";
  for (const narrow of narrowFlags) html += narrow ? "<col class='col-narrow'>" : "<col>";
  html += "</colgroup><thead><tr>";
  header.forEach((cell, i) => {
    html += "<th" + (narrowFlags[i] ? " class='col-narrow'" : "") + ">" + inline(cell) + "</th>";
  });
  html += "</tr></thead><tbody>";
  for (const row of bodyLines) {
    const cells = parseTableRow(row);
    html += "<tr>";
    for (let i = 0; i < cells.length; i++) {
      html += "<td" + (narrowFlags[i] ? " class='col-narrow'" : "") + ">" + inline(cells[i] || "") + "</td>";
    }
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
    else if (isEscapedEmptyAnchorLine(line)) { if (inList) { out.push("</ul>"); inList = false; } out.push(restoreEmptyAnchors(line.trim())); }
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
  return "<" + tag + " id='" + id + "' data-level='" + level + "' class='copyable'>" + inline(text)
    + "<button type='button' class='copy-btn' onclick='copySection(event)'>复制</button></" + tag + ">";
}

// stripMarkdownNoise 去掉标题中的 HTML / 行内标记，供目录显示。
function stripMarkdownNoise(s) {
  const tick = String.fromCharCode(96);
  return String(s || "")
    .replace(/<[^>]+>/g, "")
    .replace(new RegExp("[*_" + tick + "]", "g"), "")
    .replace(/\s+/g, " ")
    .trim();
}

// renderToc 生成可折叠目录；宽屏固定在内容区外，窄屏进入内容区顶部。
function renderToc(headings) {
  if (!headings.length) return "";
  const links = headings.map(h => {
    const level = h.level === 1 ? "level-1" : h.level === 2 ? "level-2" : "level-3";
    const label = stripMarkdownNoise(h.text);
    return "<a class='" + level + "' href='#" + h.id + "' title='" + escAttr(label) + "'><span class='toc-label'>" + esc(label) + "</span></a>";
  }).join("");
  return "<aside class='toc-dock' id='tocDock'>"
    + "<div class='toc-head'>"
    + "<div class='toc-head-meta'>"
    + "<div class='toc-head-title'>目录</div>"
    + "<div class='toc-head-count'>" + headings.length + " 个章节</div>"
    + "</div>"
    + "<button type='button' class='toc-toggle' onclick='toggleToc()' aria-expanded='true' title='折叠目录'>"
    + "<span class='toc-toggle-label'>目录</span>"
    + "<svg class='toc-toggle-icon' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.8' aria-hidden='true'>"
    + "<path stroke-linecap='round' stroke-linejoin='round' d='M6 9l6 6 6-6'/>"
    + "</svg>"
    + "</button>"
    + "</div>"
    + "<nav class='toc-body' aria-label='目录'>" + links + "</nav>"
    + "</aside>";
}

// TOC_COLLAPSED_KEY 记住用户是否收起 Markdown 目录。
const TOC_COLLAPSED_KEY = "conductor.tocCollapsed";
let tocSpyCleanup = null;

// applyTocCollapsedState 根据本地存储恢复目录折叠状态。
function applyTocCollapsedState() {
  const dock = document.getElementById("tocDock");
  if (!dock) return;
  const collapsed = localStorage.getItem(TOC_COLLAPSED_KEY) === "1";
  dock.classList.toggle("is-collapsed", collapsed);
  const btn = dock.querySelector(".toc-toggle");
  if (btn) {
    btn.setAttribute("aria-expanded", collapsed ? "false" : "true");
    btn.title = collapsed ? "展开目录" : "折叠目录";
    const icon = btn.querySelector(".toc-toggle-icon");
    if (icon) {
      icon.innerHTML = collapsed
        ? "<path stroke-linecap='round' stroke-linejoin='round' d='M4 7h16M4 12h10M4 17h14'/>"
        : "<path stroke-linecap='round' stroke-linejoin='round' d='M6 9l6 6 6-6'/>";
    }
  }
  if (!collapsed) initTocScrollSpy();
  else if (tocSpyCleanup) {
    tocSpyCleanup();
    tocSpyCleanup = null;
  }
}

// toggleToc 折叠或展开 Markdown 目录，并写入本地存储。
function toggleToc() {
  const dock = document.getElementById("tocDock");
  if (!dock) return;
  const collapsed = !dock.classList.contains("is-collapsed");
  localStorage.setItem(TOC_COLLAPSED_KEY, collapsed ? "1" : "0");
  applyTocCollapsedState();
}

// initTocScrollSpy 点击立即激活，并随内容区滚动同步高亮章节。
function initTocScrollSpy() {
  if (tocSpyCleanup) {
    tocSpyCleanup();
    tocSpyCleanup = null;
  }
  const root = document.getElementById("layoutContent");
  const dock = document.getElementById("tocDock");
  const links = Array.from(document.querySelectorAll("#tocDock .toc-body a[href^='#']"));
  if (!root || !dock || !links.length) return;
  const items = links.map(a => {
    const id = (a.getAttribute("href") || "").slice(1);
    return { a, id, el: id ? document.getElementById(id) : null };
  }).filter(item => item.el);
  if (!items.length) return;

  const setActive = (id) => {
    links.forEach(a => a.classList.toggle("is-active", a.getAttribute("href") === "#" + id));
  };

  const updateFromScroll = () => {
    const marker = root.getBoundingClientRect().top + Math.min(120, root.clientHeight * 0.2);
    let current = items[0].id;
    for (const item of items) {
      if (item.el.getBoundingClientRect().top <= marker) current = item.id;
      else break;
    }
    setActive(current);
  };

  const onClick = (event) => {
    const a = event.target.closest("a[href^='#']");
    if (!a || !dock.contains(a)) return;
    const id = (a.getAttribute("href") || "").slice(1);
    const el = id ? document.getElementById(id) : null;
    if (!el) return;
    event.preventDefault();
    setActive(id);
    el.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  dock.addEventListener("click", onClick);
  root.addEventListener("scroll", updateFromScroll, { passive: true });
  tocSpyCleanup = () => {
    dock.removeEventListener("click", onClick);
    root.removeEventListener("scroll", updateFromScroll);
  };
  updateFromScroll();
}

// copyBlock 渲染 Markdown 文本块（段落、列表项、代码块）。
function copyBlock(tag, text) {
  return "<" + tag + ">" + inline(text) + "</" + tag + ">";
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
  if (clone.tagName === "FIGURE") return "";
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

// isEscapedEmptyAnchorLine 判断一行是否仅为转义后的空锚点标签。
function isEscapedEmptyAnchorLine(line) {
  return /^&lt;a\s+id=&quot;[a-zA-Z][\w.:-]*&quot;\s*&gt;\s*&lt;\/a&gt;\s*$/.test(String(line || "").trim());
}

// restoreEmptyAnchors 将安全的空锚点从转义文本还原为真实 HTML。
function restoreEmptyAnchors(s) {
  return String(s || "").replace(
    /&lt;a\s+id=&quot;([a-zA-Z][\w.:-]*)&quot;\s*&gt;\s*&lt;\/a&gt;/g,
    '<a id="$1"></a>'
  );
}

// renderHashLink 仅允许页内哈希链接，避免把任意 URL 注入预览。
function renderHashLink(text, href) {
  const target = decodeEntities(href || "").trim();
  if (!/^#[a-zA-Z][\w.:-]*$/.test(target)) return "[" + text + "](" + href + ")";
  return '<a href="' + escAttr(target) + '">' + text + "</a>";
}

// inline 渲染生成文档中常见的行内 Markdown 标记。
function inline(s) {
  const tick = String.fromCharCode(96);
  const codePattern = new RegExp(tick + "([^" + tick + "]+)" + tick, "g");
  return restoreEmptyAnchors(s)
    .replace(/!\[([^\]]*)\]\(([^)]+)\)/g, (_, alt, src) =>
      "<img class='inline-md-image' src='" + escAttr(resolveMarkdownImageSrc(src)) + "' alt='" + alt + "' loading='lazy'>"
    )
    .replace(/\[([^\]]+)\]\((#[^)]+)\)/g, (_, text, href) => renderHashLink(text, href))
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
