package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hiscaler/conductor/assets"
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
	AbsPath string `json:"absPath,omitempty"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	Content string `json:"content,omitempty"`
	RawURL  string `json:"rawUrl,omitempty"`
}

type server struct {
	root string
}

// openFolderImpl 在文件管理器中打开目录；测试可替换以避免真正拉起系统程序。
var openFolderImpl = openInFileManager

// NewHandler 创建物料浏览器的 HTTP 处理器，可供桌面 WebView 或纯 HTTP 模式复用。
func NewHandler(root string) http.Handler {
	mux := http.NewServeMux()
	s := &server{root: root}
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/api/tree", s.tree)
	mux.HandleFunc("/api/file", s.file)
	mux.HandleFunc("/api/open-folder", s.openFolder)
	mux.HandleFunc("/api/readme", s.readme)
	mux.HandleFunc("/raw", s.raw)
	mux.HandleFunc("/doc-asset", s.docAsset)
	mux.HandleFunc("/assets/", s.asset)
	return mux
}

// index 返回单页物料浏览器界面。
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

// file 按方法分发：GET 预览元数据，DELETE 删除图片文件或子目录。
func (s *server) file(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, "":
		s.getFile(w, r)
	case http.MethodDelete:
		s.deleteFile(w, r)
	default:
		w.Header().Set("Allow", "GET, DELETE")
		writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
}

// getFile 返回所选文件或目录的元信息和预览内容。
func (s *server) getFile(w http.ResponseWriter, r *http.Request) {
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
			"absPath":  full,
			"type":     "dir",
			"children": children,
		})
		return
	}

	kind := kindFor(full)
	out := fileInfo{
		Name:    info.Name(),
		Path:    filepath.ToSlash(rel),
		AbsPath: full,
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

// deleteFile 删除 output 根目录内的图片文件，或 output 下的非根子目录（含子内容）。
// 必须显式确认（请求头 X-Confirm-Delete: 1 或 query confirm=1），防止误删。
func (s *server) deleteFile(w http.ResponseWriter, r *http.Request) {
	if !deleteConfirmed(r) {
		writeError(w, errors.New("missing delete confirmation"), http.StatusBadRequest)
		return
	}
	rel := r.URL.Query().Get("path")
	full, err := s.clean(rel)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	info, err := os.Stat(full)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, errors.New("file not found"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	slashRel := filepath.ToSlash(strings.Trim(filepath.ToSlash(rel), "/"))
	if info.IsDir() {
		if slashRel == "" || full == s.root {
			writeError(w, errors.New("refusing to delete output root"), http.StatusBadRequest)
			return
		}
		if err := os.RemoveAll(full); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		if kindFor(full) != "image" {
			writeError(w, errors.New("only image files or directories can be deleted"), http.StatusBadRequest)
			return
		}
		if err := os.Remove(full); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	deletedType := "file"
	if info.IsDir() {
		deletedType = "dir"
	}
	parent := path.Dir(slashRel)
	if parent == "." {
		parent = ""
	}
	writeJSON(w, map[string]any{
		"ok":     true,
		"path":   slashRel,
		"parent": parent,
		"type":   deletedType,
	})
}

// openFolder 在操作系统文件管理器中打开路径对应目录（文件则打开其所在目录）。
// 仅允许 output 根目录内的相对路径，防止逃逸。
func (s *server) openFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	rel := r.URL.Query().Get("path")
	full, err := s.clean(rel)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	folder, err := revealFolderPath(full)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, errors.New("path not found"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusBadRequest)
		return
	}
	// 先回响应再拉起文件管理器，避免 explorer/open 启动慢拖住按钮反馈。
	go func(target string) {
		_ = openFolderImpl(target)
	}(folder)
	writeJSON(w, map[string]any{
		"ok":     true,
		"folder": folder,
	})
}

// deleteConfirmed 检查删除请求是否已通过前端确认门禁。
// 同时接受 header 与 query，避免部分 WebView 丢弃自定义请求头。
func deleteConfirmed(r *http.Request) bool {
	if r.Header.Get("X-Confirm-Delete") == "1" {
		return true
	}
	return r.URL.Query().Get("confirm") == "1"
}

// readme 返回项目 README.md，方便用户在浏览器内查看使用说明。
func (s *server) readme(w http.ResponseWriter, r *http.Request) {
	readmePath, err := findReadme(s.root)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	info, err := os.Stat(readmePath)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	data, err := os.ReadFile(readmePath)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, fileInfo{
		Name:    "README.md",
		Path:    "README.md",
		AbsPath: readmePath,
		Type:    "markdown",
		Size:    info.Size(),
		ModTime: info.ModTime().Format(time.RFC3339),
		Content: string(data),
	})
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

// docAsset 返回 README 所在目录下的相对资源（如图片），供使用说明预览引用。
func (s *server) docAsset(w http.ResponseWriter, r *http.Request) {
	readmePath, err := findReadme(s.root)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	docRoot, err := filepath.Abs(filepath.Dir(readmePath))
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	rel := filepath.Clean(filepath.FromSlash(r.URL.Query().Get("path")))
	if rel == "." {
		rel = ""
	}
	if rel == "" || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		writeError(w, errors.New("invalid path"), http.StatusBadRequest)
		return
	}
	full := filepath.Join(docRoot, rel)
	abs, err := filepath.Abs(full)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	rootWithSep := docRoot + string(filepath.Separator)
	if abs != docRoot && !strings.HasPrefix(abs, rootWithSep) {
		writeError(w, errors.New("path escapes doc root"), http.StatusBadRequest)
		return
	}
	info, err := os.Stat(abs)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}
	if ct := mime.TypeByExtension(strings.ToLower(filepath.Ext(abs))); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Content-Disposition", "inline")
	http.ServeFile(w, r, abs)
}

// asset 返回内嵌静态资源，并强制正确的 SVG MIME，避免浏览器把 logo 当作下载文件。
func (s *server) asset(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/assets/")
	name = path.Clean("/" + name)
	name = strings.TrimPrefix(name, "/")
	if name == "" || name == "." {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(assets.FS, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ext := strings.ToLower(path.Ext(name))
	if ext == ".svg" {
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	} else if ct := mime.TypeByExtension(ext); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write(data)
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
		if skipEntry(entry.Name()) {
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

// skipEntry 判断目录树中应隐藏的系统或无关文件。
func skipEntry(name string) bool {
	switch name {
	case ".gitignore", ".DS_Store", "Thumbs.db", "desktop.ini":
		return true
	}
	return strings.HasPrefix(name, "._")
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

// findReadme 从 output 根目录向上查找项目 README.md。
func findReadme(start string) (string, error) {
	dir := start
	for {
		candidate := filepath.Join(dir, "README.md")
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
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
