package handler

import (
	"path/filepath"
	"io/fs"
	"github.com/cn1095/hit-counter/internal"
	"github.com/labstack/echo/v4"
)

// Wasm is API for serving wasm file.
func (h *Handler) Wasm(c echo.Context) error {
	hctx := c.(*HitCounterContext)
	hctx.Response().Header().Set("Content-Encoding", "gzip")
	// 从嵌入的文件系统中读取 wasm 文件
	data, err := fs.ReadFile(EmbeddedFiles, "view/hits.wasm")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "无法读取 wasm 文件")
	}

	// 返回文件内容作为响应
	return c.Blob(http.StatusOK, "application/wasm", data)
}
