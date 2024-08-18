package handler

import (
	websocket "github.com/gjbae1212/go-ws-broadcast"
	"github.com/labstack/echo/v4"
	perrors "github.com/pkg/errors"
)

// WebSocket is API for websocket.
func (h *handler) WebSocket(c echo.Context) error {
	ws, err := websocket.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return perrors.WithStack(err)
	}

	// register websocket to breaker.
	if _, err := h.websocketBreaker.Register(ws); err != nil {
		return perrors.WithStack(err)
	}
	return nil
}
