package browser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func deleteReq(path string, confirmMode string) *http.Request {
	url := "/api/file?path=" + path
	if confirmMode == "query" {
		url += "&confirm=1"
	}
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	if confirmMode == "header" {
		req.Header.Set("X-Confirm-Delete", "1")
	}
	if confirmMode == "both" {
		req.Header.Set("X-Confirm-Delete", "1")
		// query already not set; set via rebuild
		req = httptest.NewRequest(http.MethodDelete, "/api/file?path="+path+"&confirm=1", nil)
		req.Header.Set("X-Confirm-Delete", "1")
	}
	return req
}

func TestDeleteRequiresConfirm(t *testing.T) {
	root := t.TempDir()
	listing := filepath.Join(root, "Temu-US", "listing")
	if err := os.MkdirAll(listing, 0o755); err != nil {
		t.Fatal(err)
	}
	h := NewHandler(root)

	req := deleteReq("Temu-US/listing", "")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("without confirm: status=%d body=%s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(listing); err != nil {
		t.Fatalf("listing should still exist: %v", err)
	}
}

func TestDeleteAcceptsHeaderOrQueryConfirm(t *testing.T) {
	root := t.TempDir()
	h := NewHandler(root)

	for _, mode := range []string{"header", "query", "both"} {
		dir := filepath.Join(root, "mode-"+mode)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		req := deleteReq("mode-"+mode, mode)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("mode %s: status=%d body=%s", mode, rec.Code, rec.Body.String())
		}
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Fatalf("mode %s: dir still exists", mode)
		}
	}
}

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

	req := deleteReq("", "query")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("delete root: status = %d, body = %s", rec.Code, rec.Body.String())
	}

	req = deleteReq("Temu-US/%E7%BB%84%E5%90%88-MUG01x6", "both")
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
	req := deleteReq("note.md", "query")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("delete markdown: status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteImageWithConfirm(t *testing.T) {
	root := t.TempDir()
	img := filepath.Join(root, "a.png")
	if err := os.WriteFile(img, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := NewHandler(root)
	req := deleteReq("a.png", "query")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete image: status=%d body=%s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(img); !os.IsNotExist(err) {
		t.Fatal("image still exists")
	}
}

func TestIndexHTMLHasConfirmDeleteWiring(t *testing.T) {
	h := NewHandler(t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	body := rec.Body.String()
	needles := []string{
		`id="confirmDialog"`,
		`id="confirmOkBtn"`,
		`settleConfirm(true)`,
		`function askConfirm`,
		`&confirm=1`,
		`X-Confirm-Delete`,
	}
	for _, n := range needles {
		if !strings.Contains(body, n) {
			t.Fatalf("index HTML missing %q", n)
		}
	}
}
