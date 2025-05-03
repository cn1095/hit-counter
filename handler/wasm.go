package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"embed"
)
//go:embed view/hits.wasm
var embeddedFile embed.FS // 嵌入 hits.wasm 文件
// Wasm is API for serving wasm file.
func (h *Handler) Wasm(c echo.Context) error {
	hctx := c.(*HitCounterContext)
	hctx.Response().Header().Set("Content-Encoding", "gzip")

	// 从嵌入的文件系统中读取 hits.wasm
	file, err := embeddedFile.Open("view/hits.wasm")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "无法读取 wasm 文件")
	}
	defer file.Close()

	// 将嵌入的文件内容直接返回
	return c.Stream(http.StatusOK, "application/wasm", file)
}
