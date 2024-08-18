package handler

import (
	"net/http"

	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/labstack/echo/v4"
)

// Wasm is API for serving wasm file.
func (h *handler) Wasm(c echo.Context) error {
	hitCtx := c.(*servercontext.HitCounterContext)
	hitCtx.Response().Header().Set("Content-Encoding", "gzip")
	return hitCtx.Blob(http.StatusOK, "application/octet-stream", h.wasmFile)
}
