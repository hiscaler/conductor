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
    .copy-btn { position:absolute; right:0; top:0.15em; border:1px solid rgb(51 65 85); background:rgb(15 23 42); color:rgb(148 163 184); border-radius:5px; padding:2px 6px; font-size:12px; opacity:0; }
    .copyable:hover .copy-btn { opacity:1; }
    .copy-btn:hover { color:rgb(37 99 235); border-color:rgb(191 219 254); }
    .markdown h1, .markdown h2, .markdown h3 { line-height:1.25; margin-top:1.2em; scroll-margin-top:18px; }
    .markdown h1 { font-size:24px; border-bottom:1px solid rgb(51 65 85); padding-bottom:8px; }
    .markdown h2 { font-size:20px; border-bottom:1px solid rgb(51 65 85); padding-bottom:6px; }
    .markdown h3 { font-size:16px; }
    .markdown code { background:rgb(30 41 59); padding:2px 4px; border-radius:4px; }
    .markdown pre { white-space:pre-wrap; word-break:break-word; overflow-wrap:anywhere; }
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
    <div id="content" class="rounded-lg border border-dashed border-slate-700 bg-slate-900 p-6 text-slate-400">请选择左侧文件或目录。</div>
  </section>
</main>
<script>
let activePath = "";
let refreshTimer = null;
let expandedPaths = new Set();

// openReadme 加载项目 README，作为用户使用说明预览。
async function openReadme() {
  activePath = "__readme__";
  markActive();
  const res = await fetch("/api/readme");
  if (!res.ok) {
    document.getElementById("content").innerHTML = "<div class='rounded-lg border border-dashed border-slate-700 bg-slate-900 p-6 text-slate-400'>README.md 读取失败</div>";
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
    document.getElementById("content").innerHTML = "<div class='rounded-lg border border-dashed border-slate-700 bg-slate-900 p-6 text-slate-400'>读取失败</div>";
    return;
  }
  const data = await res.json();
  renderContent(data);
}

// renderContent 根据文件类型选择合适的预览方式。
function renderContent(data) {
  if (data.type === "dir") {
    const cell = " class='border border-slate-700 p-2 align-top'";
    const rows = (data.children || []).map(n => "<tr><td" + cell + ">" + iconFor(n.type) + " " + esc(n.name) + "</td><td" + cell + ">" + esc(n.type) + "</td><td" + cell + ">" + formatSize(n.size || 0) + "</td></tr>").join("");
    document.getElementById("content").innerHTML = "<div class='max-w-6xl rounded-lg border border-slate-800 bg-slate-900 p-5'><div class='mb-4 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold'>" + esc(data.name || "output") + "</h1><small class='text-slate-400'>" + esc(data.path || "") + "</small></div><table class='w-full table-fixed border-collapse'><thead><tr><th class='border border-slate-700 bg-slate-950 p-2 text-left'>名称</th><th class='border border-slate-700 bg-slate-950 p-2 text-left'>类型</th><th class='border border-slate-700 bg-slate-950 p-2 text-left'>大小</th></tr></thead><tbody>" + rows + "</tbody></table></div>";
    return;
  }
  let body = "";
  if (data.type === "markdown") body = renderMarkdownPreview(data.content || "");
  else if (data.type === "image") body = "<div><img class='max-w-full rounded-lg border border-slate-700 bg-slate-950' src='" + escAttr(data.rawUrl) + "' alt='" + escAttr(data.name) + "'></div>";
  else if (data.type === "video") body = "<div><video class='max-w-full rounded-lg border border-slate-700 bg-black' src='" + escAttr(data.rawUrl) + "' controls></video></div>";
  else if (data.type === "text" || data.type === "json") body = "<pre class='overflow-auto whitespace-pre-wrap break-words rounded-lg bg-slate-950 p-4 text-slate-200'>" + esc(data.content || "") + "</pre>";
  else body = "<p><a href='" + escAttr(data.rawUrl) + "' target='_blank'>下载或打开文件</a></p>";
  document.getElementById("content").innerHTML = "<div class='max-w-6xl rounded-lg border border-slate-800 bg-slate-900 p-5'><div class='mb-4 flex items-baseline gap-3'><h1 class='m-0 text-xl font-semibold'>" + esc(data.name) + "</h1><small class='text-slate-400'>" + esc(data.path) + " · " + formatSize(data.size || 0) + "</small></div>" + body + "</div>";
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
  const value = decodeEntities(text);
  return "<" + tag + " id='" + id + "' class='copyable'>" + inline(text) + "<button class='copy-btn' onclick='copyText(event,\"" + escJS(value) + "\")'>复制</button></" + tag + ">";
}

// renderToc 生成 Markdown 右侧悬浮目录导航。
function renderToc(headings) {
  if (!headings.length) return "";
  const levelClass = h => h.level === 1 ? "font-semibold text-slate-200" : h.level === 2 ? "pl-3" : "pl-6 text-xs";
  const links = headings.map(h => "<a class='block rounded px-1 py-1 text-sm leading-snug text-slate-400 hover:bg-slate-800 hover:text-blue-300 " + levelClass(h) + "' href='#" + h.id + "'>" + esc(h.text) + "</a>").join("");
  return "<nav class='toc sticky top-5 max-h-[calc(100vh-7rem)] overflow-auto rounded-lg border border-slate-800 bg-slate-900 p-3'><div class='mb-2 text-sm font-semibold text-slate-200'>目录</div>" + links + "</nav>";
}

// copyBlock 为 Markdown 块包裹复制按钮。
function copyBlock(tag, text) {
  const value = decodeEntities(text);
  return "<" + tag + " class='copyable'>" + inline(text) + "<button class='copy-btn' onclick='copyText(event,\"" + escJS(value) + "\")'>复制</button></" + tag + ">";
}

// copyText 复制单个 Markdown 块文本，并显示短暂反馈。
function copyText(event, text) {
  event.stopPropagation();
  navigator.clipboard.writeText(text).then(() => {
    const btn = event.currentTarget;
    const old = btn.textContent;
    btn.textContent = "已复制";
    setTimeout(() => btn.textContent = old, 900);
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
