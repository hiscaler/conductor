//go:build windows

package browser

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

// windowsFileManagerCommand 用 start 异步打开目录，并隐藏 cmd 黑窗。
func windowsFileManagerCommand(folder string) *exec.Cmd {
	cmd := exec.Command("cmd", "/C", "start", "", filepath.Clean(folder))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
