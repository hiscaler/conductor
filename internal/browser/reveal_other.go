//go:build !windows

package browser

import (
	"os/exec"
	"path/filepath"
)

// windowsFileManagerCommand 非 Windows 构建占位，供 fileManagerCommand 单测构造命令。
func windowsFileManagerCommand(folder string) *exec.Cmd {
	return exec.Command("cmd", "/C", "start", "", filepath.Clean(folder))
}
