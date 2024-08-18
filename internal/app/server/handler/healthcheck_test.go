package handler

import (
	"net/http/httptest"
	"testing"

	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *servercontext.HitCounterContext
		output error
	}{
		{
			name: "success",
			ctx: servercontext.NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), httptest.NewRecorder())),
			output: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &handler{}
			output := h.HealthCheck(tc.ctx)
			assert.Equal(t, tc.output, output)
		})
	}
}
