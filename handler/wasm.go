package handler

import (
	"github.com/cn1095/hit-counter/internal"
	"github.com/labstack/echo/v4"
	"io/fs"
	"net/http"
	"strings"
)

// Wasm is API for serving wasm file.
func (h *Handler) Wasm(c echo.Context) error {
	hctx := c.(*HitCounterContext)
	hctx.Response().Header().Set("Content-Encoding", "gzip")

	// 从嵌入的文件系统中读取 hits.wasm
	file, err := embeddedFiles.Open("view/hits.wasm")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "无法读取 wasm 文件")
	}
	defer file.Close()

	// 将嵌入的文件内容直接返回
	return c.Stream(http.StatusOK, "application/wasm", file)
}
