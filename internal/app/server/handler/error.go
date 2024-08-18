package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Error is API for error.
func (h *handler) Error(err error, c echo.Context) {
	code := http.StatusInternalServerError

	if httpErr := (*echo.HTTPError)(nil); errors.As(err, &httpErr) {
		code = httpErr.Code
	}
	_ = c.NoContent(code)
}
