package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck is API for checking server status.
func (h *handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "health check!")
}
