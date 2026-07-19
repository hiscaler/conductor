package main

import (
	"context"
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hiscaler/conductor/internal/browser"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// App 是 Wails 桌面应用的生命周期载体。
type App struct {
	ctx  context.Context
	root string
}

// NewApp 创建桌面应用实例。
func NewApp(root string) *App {
	return &App{root: root}
}

// startup 在窗口启动时保存上下文。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// main 默认启动 Conductor 桌面窗口；加 -web 则仅提供 HTTP 服务。
func main() {
	// Windows 注册表常把 .svg 标成 image/svg，浏览器无法作为图片渲染。
	_ = mime.AddExtensionType(".svg", "image/svg+xml")

	web := flag.Bool("web", false, "run as local HTTP server instead of desktop window")
	addr := flag.String("addr", "127.0.0.1:8080", "listen address for -web mode")
	rootFlag := flag.String("root", "", "directory to browse (default: output next to executable, or ./output in -web mode)")
	flag.Parse()

	absRoot, err := resolveRoot(*rootFlag, *web)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(absRoot, 0755); err != nil {
		log.Fatal(err)
	}

	handler := browser.NewHandler(absRoot)
	if *web {
		log.Printf("Output browser serving %s", absRoot)
		log.Printf("Open http://%s", *addr)
		if err := http.ListenAndServe(*addr, handler); err != nil {
			log.Fatal(err)
		}
		return
	}

	app := NewApp(absRoot)
	err = wails.Run(&options.App{
		Title:            "Conductor",
		Width:            1440,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
		AssetServer: &assetserver.Options{
			Handler: handler,
		},
		OnStartup: app.startup,
		Bind:      []any{app},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// resolveRoot 解析要浏览的 output 目录；桌面模式默认使用 exe 同目录下的 output。
func resolveRoot(root string, web bool) (string, error) {
	if root == "" {
		if web {
			root = "output"
		} else {
			root = filepath.Join(exeDir(), "output")
		}
	}
	return filepath.Abs(root)
}

// exeDir 返回当前可执行文件所在目录，便于便携分发。
func exeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return filepath.Dir(exe)
}
