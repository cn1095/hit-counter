package api_handler

import (
	"fmt"

	"github.com/cn1095/hit-counter/handler"
	"github.com/cn1095/hit-counter/internal"
)

type Handler struct {
	*handler.Handler
}

// NewHandler creates api handler object.
func NewHandler(h *handler.Handler) (*Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("[错误] api handler %w", internal.ErrorEmptyParams)
	}
	return &Handler{Handler: h}, nil
}
