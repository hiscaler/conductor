package browser

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// revealFolderPath 返回应在文件管理器中打开的目录：目录原样返回，文件则返回其父目录。
func revealFolderPath(full string) (string, error) {
	info, err := os.Stat(full)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return full, nil
	}
	parent := filepath.Dir(full)
	if parent == "" || parent == "." {
		return "", errors.New("cannot resolve parent folder")
	}
	return parent, nil
}

// openInFileManager 在当前操作系统的文件管理器中打开目录。
// 使用 Start + 后台 Wait，避免阻塞 HTTP 请求。
func openInFileManager(folder string) error {
	cmd, err := fileManagerCommand(runtime.GOOS, folder)
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open folder: %w", err)
	}
	go func() { _ = cmd.Wait() }()
	return nil
}

// fileManagerCommand 按操作系统构造打开目录的命令，便于单测。
func fileManagerCommand(goos, folder string) (*exec.Cmd, error) {
	if folder == "" {
		return nil, errors.New("empty folder path")
	}
	switch goos {
	case "windows":
		// cmd start 立即返回；直接调 explorer 往往会卡住等 shell 就绪。
		return windowsFileManagerCommand(folder), nil
	case "darwin":
		return exec.Command("open", folder), nil
	case "linux":
		return exec.Command("xdg-open", folder), nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", goos)
	}
}
