package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type node struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Size     int64  `json:"size,omitempty"`
	ModTime  string `json:"modTime,omitempty"`
	Children []node `json:"children,omitempty"`
}

type fileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	Content string `json:"content,omitempty"`
	RawURL  string `json:"rawUrl,omitempty"`
}

// main 启动本地 HTTP 服务，用于浏览生成的 output 产物。
func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	root := flag.String("root", "output", "directory to browse")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(absRoot, 0755); err != nil {
		log.Fatal(err)
	}

	app := &server{root: absRoot}
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/api/tree", app.tree)
	mux.HandleFunc("/api/file", app.file)
	mux.HandleFunc("/raw", app.raw)

	log.Printf("Output browser serving %s", absRoot)
	log.Printf("Open http://%s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	root string
}

// index 返回单页文件浏览器界面。
func (s *server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

// tree 返回 output 目录树，供左侧栏渲染。
func (s *server) tree(w http.ResponseWriter, r *http.Request) {
	root := node{Name: filepath.Base(s.root), Path: "", Type: "dir"}
	children, err := s.readDir("")
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	root.Children = children
	writeJSON(w, root)
}

// file 返回所选文件或目录的元信息和预览内容。
func (s *server) file(w http.ResponseWriter, r *http.Request) {
	rel := r.URL.Query().Get("path")
	full, err := s.clean(rel)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	info, err := os.Stat(full)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	if info.IsDir() {
		children, err := s.readDir(rel)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]any{
			"name":     info.Name(),
			"path":     filepath.ToSlash(rel),
			"type":     "dir",
			"children": children,
		})
		return
	}

	kind := kindFor(full)
	out := fileInfo{
		Name:    info.Name(),
		Path:    filepath.ToSlash(rel),
		Type:    kind,
		Size:    info.Size(),
		ModTime: info.ModTime().Format(time.RFC3339),
	}
	if kind == "markdown" || kind == "text" || kind == "json" {
		data, err := os.ReadFile(full)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		out.Content = string(data)
	} else {
		out.RawURL = "/raw?path=" + queryEscapePath(rel)
	}
	writeJSON(w, out)
}

// raw 从 output 根目录流式返回图片、视频等二进制文件。
func (s *server) raw(w http.ResponseWriter, r *http.Request) {
	rel := r.URL.Query().Get("path")
	full, err := s.clean(rel)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}
	if ct := mime.TypeByExtension(strings.ToLower(filepath.Ext(full))); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	http.ServeFile(w, r, full)
}

// readDir 递归读取相对目录，并生成目录树节点。
func (s *server) readDir(rel string) ([]node, error) {
	full, err := s.clean(rel)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(full)
	if err != nil {
		return nil, err
	}
	out := make([]node, 0, len(entries))
	for _, entry := range entries {
		if entry.Name() == ".gitignore" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		childRel := filepath.Join(rel, entry.Name())
		n := node{
			Name:    entry.Name(),
			Path:    filepath.ToSlash(childRel),
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
		}
		if entry.IsDir() {
			n.Type = "dir"
			children, err := s.readDir(childRel)
			if err == nil {
				n.Children = children
			}
		} else {
			n.Type = kindFor(entry.Name())
			n.Size = info.Size()
		}
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type == "dir" && out[j].Type != "dir" {
			return true
		}
		if out[i].Type != "dir" && out[j].Type == "dir" {
			return false
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// clean 解析用户传入的相对路径，并防止访问 output 之外的文件。
func (s *server) clean(rel string) (string, error) {
	rel = filepath.Clean(filepath.FromSlash(rel))
	if rel == "." {
		rel = ""
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", errors.New("invalid path")
	}
	full := filepath.Join(s.root, rel)
	abs, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	rootWithSep := s.root + string(filepath.Separator)
	if abs != s.root && !strings.HasPrefix(abs, rootWithSep) {
		return "", errors.New("path escapes output root")
	}
	return abs, nil
}

// kindFor 根据文件扩展名判断前端可用的预览类型。
func kindFor(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg":
		return "image"
	case ".mp4", ".webm", ".mov", ".m4v":
		return "video"
	case ".json":
		return "json"
	case ".txt", ".csv", ".log", ".yaml", ".yml":
		return "text"
	default:
		return "binary"
	}
}

// queryEscapePath 将相对路径编码为 raw 文件 URL 参数。
func queryEscapePath(path string) string {
	return url.QueryEscape(filepath.ToSlash(path))
}

// writeJSON 输出带 UTF-8 头的 JSON 响应。
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// writeError 输出简洁的 JSON 风格错误响应。
func writeError(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), status)
}

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Conductor 输出浏览器</title>
  <style>
    :root { color-scheme: light; --line:#e5e7eb; --muted:#6b7280; --bg:#f8fafc; --panel:#ffffff; --accent:#2563eb; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: system-ui, -apple-system, "Segoe UI", sans-serif; background:var(--bg); color:#111827; }
    header { height:52px; display:flex; align-items:center; gap:12px; padding:0 18px; border-bottom:1px solid var(--line); background:var(--panel); }
    header strong { font-size:16px; }
    header span { color:var(--muted); font-size:13px; }
    main { display:grid; grid-template-columns: 330px 1fr; height:calc(100vh - 52px); }
    aside { border-right:1px solid var(--line); background:var(--panel); overflow:auto; padding:12px; }
    section { overflow:auto; padding:22px; }
    .toolbar { display:flex; gap:8px; align-items:center; margin-bottom:12px; flex-wrap:wrap; }
    .toolbar label { color:var(--muted); font-size:13px; display:flex; align-items:center; gap:5px; }
    button { border:1px solid var(--line); background:#fff; padding:7px 10px; border-radius:6px; cursor:pointer; }
    button:hover { border-color:#bfdbfe; color:var(--accent); }
    .tree, .tree ul { list-style:none; margin:0; padding-left:14px; }
    .tree { padding-left:0; }
    .node { display:flex; align-items:center; gap:7px; width:100%; border:0; background:transparent; text-align:left; padding:5px 6px; border-radius:6px; font-size:14px; }
    .node:hover, .node.active { background:#eff6ff; color:#1d4ed8; }
    .twisty { width:16px; color:var(--muted); text-align:center; }
    .icon { width:18px; text-align:center; color:var(--muted); }
    .collapsed > ul { display:none; }
    .meta { color:var(--muted); font-size:12px; margin-left:auto; }
    .empty { color:var(--muted); padding:24px; border:1px dashed var(--line); border-radius:8px; background:#fff; }
    .card { background:var(--panel); border:1px solid var(--line); border-radius:8px; padding:18px; max-width:1100px; }
    .file-title { display:flex; align-items:baseline; gap:12px; margin:0 0 14px; }
    .file-title h1 { font-size:20px; margin:0; }
    .file-title small { color:var(--muted); }
    .preview img { max-width:100%; height:auto; border:1px solid var(--line); border-radius:8px; background:#fff; }
    .preview video { max-width:100%; border:1px solid var(--line); border-radius:8px; background:#000; }
    pre { white-space:pre-wrap; word-break:break-word; background:#0f172a; color:#e5e7eb; padding:16px; border-radius:8px; overflow:auto; }
    .markdown { line-height:1.7; }
    .copyable { position:relative; padding-right:42px; }
    .copy-btn { position:absolute; right:0; top:0.15em; border:1px solid var(--line); background:#fff; color:var(--muted); border-radius:5px; padding:2px 6px; font-size:12px; opacity:0; }
    .copyable:hover .copy-btn { opacity:1; }
    .copy-btn:hover { color:var(--accent); border-color:#bfdbfe; }
    .markdown h1, .markdown h2, .markdown h3 { line-height:1.25; margin-top:1.2em; }
    .markdown h1 { font-size:24px; border-bottom:1px solid var(--line); padding-bottom:8px; }
    .markdown h2 { font-size:20px; border-bottom:1px solid var(--line); padding-bottom:6px; }
    .markdown h3 { font-size:16px; }
    .markdown code { background:#f1f5f9; padding:2px 4px; border-radius:4px; }
    .markdown table { border-collapse:collapse; width:100%; margin:12px 0; }
    .markdown th, .markdown td { border:1px solid var(--line); padding:8px; vertical-align:top; }
    .markdown th { background:#f8fafc; }
    @media (max-width: 800px) { main { grid-template-columns: 1fr; } aside { height:38vh; border-right:0; border-bottom:1px solid var(--line); } section { height:calc(62vh - 52px); } }
  </style>
</head>
<body>
<header>
  <strong>Conductor 浏览器</strong>
</header>
<main>
  <aside>
    <div class="toolbar">
      <button onclick="loadTree()">刷新</button>
      <label><input id="autoRefresh" type="checkbox" checked onchange="toggleAutoRefresh()"> 自动刷新</label>
      <span id="refreshState">5 秒</span>
    </div>
    <ul id="tree" class="tree"></ul>
  </aside>
  <section>
    <div id="content" class="empty">请选择左侧文件或目录。</div>
  </section>
</main>
<script>
let activePath = "";
let refreshTimer = null;
let expandedPaths = new Set();

// loadTree 刷新左侧目录树，并保留展开状态和当前选中项。
async function loadTree() {
  const res = await fetch("/api/tree");
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
  const meta = n.type === "dir" ? "" : "<span class='meta'>" + formatSize(n.size || 0) + "</span>";
  const hasChildren = n.type === "dir" && n.children && n.children.length;
  const shouldCollapse = hasChildren && depth >= 2 && !expandedPaths.has(n.path);
  const liClass = shouldCollapse ? " class='collapsed'" : "";
  const twisty = hasChildren ? "<span class='twisty'>" + (shouldCollapse ? "▶" : "▼") + "</span>" : "<span class='twisty'></span>";
  const child = hasChildren ? "<ul>" + renderChildren(n.children, depth + 1) + "</ul>" : "";
  return "<li" + liClass + "><button class='node' data-path='" + escAttr(n.path) + "' onclick='handleNodeClick(event,\"" + escJS(n.path) + "\"," + (hasChildren ? "true" : "false") + ")'>" + twisty + "<span class='icon'>" + icon + "</span><span>" + esc(n.name) + "</span>" + meta + "</button>" + child + "</li>";
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
  document.querySelectorAll(".node").forEach(el => el.classList.toggle("active", el.dataset.path === activePath));
}

// openPath 加载选中路径的元信息和预览内容。
async function openPath(path) {
  activePath = path;
  markActive();
  const res = await fetch("/api/file?path=" + encodeURIComponent(path));
  if (!res.ok) {
    document.getElementById("content").innerHTML = "<div class='empty'>读取失败</div>";
    return;
  }
  const data = await res.json();
  renderContent(data);
}

// renderContent 根据文件类型选择合适的预览方式。
function renderContent(data) {
  if (data.type === "dir") {
    const rows = (data.children || []).map(n => "<tr><td>" + iconFor(n.type) + " " + esc(n.name) + "</td><td>" + esc(n.type) + "</td><td>" + formatSize(n.size || 0) + "</td></tr>").join("");
    document.getElementById("content").innerHTML = "<div class='card'><div class='file-title'><h1>" + esc(data.name || "output") + "</h1><small>" + esc(data.path || "") + "</small></div><table><thead><tr><th>名称</th><th>类型</th><th>大小</th></tr></thead><tbody>" + rows + "</tbody></table></div>";
    return;
  }
  let body = "";
  if (data.type === "markdown") body = "<div class='markdown'>" + renderMarkdown(data.content || "") + "</div>";
  else if (data.type === "image") body = "<div class='preview'><img src='" + escAttr(data.rawUrl) + "' alt='" + escAttr(data.name) + "'></div>";
  else if (data.type === "video") body = "<div class='preview'><video src='" + escAttr(data.rawUrl) + "' controls></video></div>";
  else if (data.type === "text" || data.type === "json") body = "<pre>" + esc(data.content || "") + "</pre>";
  else body = "<p><a href='" + escAttr(data.rawUrl) + "' target='_blank'>下载或打开文件</a></p>";
  document.getElementById("content").innerHTML = "<div class='card'><div class='file-title'><h1>" + esc(data.name) + "</h1><small>" + esc(data.path) + " · " + formatSize(data.size || 0) + "</small></div>" + body + "</div>";
}

// renderMarkdown 将常用 Markdown 内容转换为预览 HTML。
function renderMarkdown(src) {
  const lines = esc(src).split(/\r?\n/);
  let out = [];
  let inList = false;
  let inCode = false;
  let code = [];
  const tick = String.fromCharCode(96);
  const fence = tick + tick + tick;
  for (const line of lines) {
    if (line.startsWith(fence)) {
      if (inCode) { out.push(copyBlock("pre", code.join("\n"))); code = []; inCode = false; }
      else { if (inList) { out.push("</ul>"); inList = false; } inCode = true; }
      continue;
    }
    if (inCode) { code.push(line); continue; }
    if (line.startsWith("### ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("h3", line.slice(4))); }
    else if (line.startsWith("## ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("h2", line.slice(3))); }
    else if (line.startsWith("# ")) { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("h1", line.slice(2))); }
    else if (line.startsWith("- ")) { if (!inList) { out.push("<ul>"); inList = true; } out.push(copyBlock("li", line.slice(2))); }
    else if (line.trim() === "") { if (inList) { out.push("</ul>"); inList = false; } }
    else if (line.includes("|")) { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("p", line)); }
    else { if (inList) { out.push("</ul>"); inList = false; } out.push(copyBlock("p", line)); }
  }
  if (inList) out.push("</ul>");
  return out.join("");
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
startAutoRefresh();
</script>
</body>
</html>`
