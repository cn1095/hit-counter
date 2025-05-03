package handler

import (
	"fmt"

	websocket "github.com/cn1095/go-ws-broadcast"
	"github.com/labstack/echo/v4"
)

// WebSocket is API for websocket.
func (h *Handler) WebSocket(c echo.Context) error {
	ws, err := websocket.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("[错误] WebSocket API %w", err)
	}

	// register websocket to breaker.
	if _, err := h.WebSocketBreaker.Register(ws); err != nil {
		return fmt.Errorf("[错误] WebSocket API %w", err)
	}
	return nil
}
