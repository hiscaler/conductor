package browser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatTreeModTime(t *testing.T) {
	got := formatTreeModTime(time.Date(2026, 10, 5, 15, 19, 0, 0, time.Local))
	if got != "2026-10-5 15:19" {
		t.Fatalf("got %q", got)
	}
}

func TestTreeSortsByNewestModTime(t *testing.T) {
	root := t.TempDir()
	oldDir := filepath.Join(root, "aaa-old")
	newDir := filepath.Join(root, "zzz-new")
	if err := os.MkdirAll(filepath.Join(oldDir, "inner"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	oldFile := filepath.Join(oldDir, "inner", "note.md")
	newFile := filepath.Join(newDir, "fresh.md")
	if err := os.WriteFile(oldFile, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldT := time.Date(2026, 1, 2, 10, 0, 0, 0, time.Local)
	newT := time.Date(2026, 10, 5, 15, 19, 0, 0, time.Local)
	for _, item := range []struct {
		path string
		at   time.Time
	}{
		{oldFile, oldT},
		{filepath.Join(oldDir, "inner"), oldT},
		{oldDir, oldT},
		{newFile, newT},
		{newDir, newT},
	} {
		if err := os.Chtimes(item.path, item.at, item.at); err != nil {
			t.Fatal(err)
		}
	}

	h := NewHandler(root)
	req := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var rootNode node
	if err := json.Unmarshal(rec.Body.Bytes(), &rootNode); err != nil {
		t.Fatal(err)
	}
	if len(rootNode.Children) < 2 {
		t.Fatalf("children=%d", len(rootNode.Children))
	}
	if rootNode.Children[0].Name != "zzz-new" {
		t.Fatalf("first=%s want zzz-new", rootNode.Children[0].Name)
	}
	if rootNode.Children[1].Name != "aaa-old" {
		t.Fatalf("second=%s want aaa-old", rootNode.Children[1].Name)
	}
	if rootNode.Children[0].ModTime != "2026-10-5 15:19" {
		t.Fatalf("modTime=%q", rootNode.Children[0].ModTime)
	}
}
