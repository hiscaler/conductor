package browser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteDirectoryAndRefuseRoot(t *testing.T) {
	root := t.TempDir()
	listing := filepath.Join(root, "Temu-US", "组合-MUG01x6")
	if err := os.MkdirAll(filepath.Join(listing, "图片"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(listing, "商品资料.md"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}

	h := NewHandler(root)

	// 拒绝删除 output 根目录
	req := httptest.NewRequest(http.MethodDelete, "/api/file?path=", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("delete root: status = %d, body = %s", rec.Code, rec.Body.String())
	}

	// 允许删除 listing 子目录
	req = httptest.NewRequest(http.MethodDelete, "/api/file?path=Temu-US/%E7%BB%84%E5%90%88-MUG01x6", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete listing: status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["ok"] != true || payload["type"] != "dir" || payload["parent"] != "Temu-US" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if _, err := os.Stat(listing); !os.IsNotExist(err) {
		t.Fatalf("listing still exists: %v", err)
	}
}

func TestDeleteStillAllowsImageOnlyForFiles(t *testing.T) {
	root := t.TempDir()
	md := filepath.Join(root, "note.md")
	if err := os.WriteFile(md, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := NewHandler(root)
	req := httptest.NewRequest(http.MethodDelete, "/api/file?path=note.md", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("delete markdown: status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
