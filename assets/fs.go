package assets

import "embed"

// FS 内嵌网站静态资源，避免依赖进程工作目录。
//
//go:embed *.svg
var FS embed.FS
