package browser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestRevealFolderPath(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "Temu-US", "listing")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "note.md")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := revealFolderPath(dir)
	if err != nil || got != dir {
		t.Fatalf("dir: got %q err=%v", got, err)
	}
	got, err = revealFolderPath(file)
	if err != nil || got != dir {
		t.Fatalf("file: got %q want %q err=%v", got, dir, err)
	}
}

func TestFileManagerCommand(t *testing.T) {
	folder := filepath.Join("D:", "wwwroot", "output")
	cases := []struct {
		goos string
		bin  string
	}{
		{"windows", "cmd"},
		{"darwin", "open"},
		{"linux", "xdg-open"},
	}
	for _, tc := range cases {
		cmd, err := fileManagerCommand(tc.goos, folder)
		if err != nil {
			t.Fatalf("%s: %v", tc.goos, err)
		}
		if len(cmd.Args) == 0 || cmd.Args[0] != tc.bin {
			t.Fatalf("%s: args=%v", tc.goos, cmd.Args)
		}
		if tc.goos == "windows" {
			if len(cmd.Args) < 5 || cmd.Args[1] != "/C" || cmd.Args[2] != "start" {
				t.Fatalf("windows args=%v", cmd.Args)
			}
		} else if len(cmd.Args) < 2 {
			t.Fatalf("%s: missing folder arg: %v", tc.goos, cmd.Args)
		}
	}
	if _, err := fileManagerCommand("plan9", folder); err == nil {
		t.Fatal("expected unsupported OS error")
	}
}

func TestOpenFolderAPI(t *testing.T) {
	root := t.TempDir()
	listing := filepath.Join(root, "Temu-US", "listing")
	if err := os.MkdirAll(listing, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(listing, "a.png")
	if err := os.WriteFile(file, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var opened string
	var mu sync.Mutex
	done := make(chan struct{}, 2)
	prev := openFolderImpl
	openFolderImpl = func(folder string) error {
		mu.Lock()
		opened = folder
		mu.Unlock()
		done <- struct{}{}
		return nil
	}
	t.Cleanup(func() { openFolderImpl = prev })

	h := NewHandler(root)
	waitOpened := func(t *testing.T, want string) {
		t.Helper()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for openFolderImpl")
		}
		mu.Lock()
		got := opened
		mu.Unlock()
		if got != want {
			t.Fatalf("opened=%q want %q", got, want)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/api/open-folder?path=Temu-US/listing", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("open dir: status=%d body=%s", rec.Code, rec.Body.String())
	}
	waitOpened(t, listing)
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["ok"] != true {
		t.Fatalf("payload=%#v", payload)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/open-folder?path=Temu-US/listing/a.png", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("open file: status=%d body=%s", rec.Code, rec.Body.String())
	}
	waitOpened(t, listing)

	req = httptest.NewRequest(http.MethodPost, "/api/open-folder?path=../outside", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("escape: status=%d body=%s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/open-folder?path=Temu-US/listing", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET: status=%d", rec.Code)
	}

	if _, err := fileManagerCommand(runtime.GOOS, listing); err != nil {
		t.Fatalf("current OS command: %v", err)
	}
}
