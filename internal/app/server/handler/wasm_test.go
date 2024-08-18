package handler

import (
	"net/http/httptest"
	"testing"

	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Wasm(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T) *handler
		ctx    *servercontext.HitCounterContext
		output error
	}{
		{
			name: "success",
			expect: func(t *testing.T) *handler {
				h := &handler{wasmFile: []byte("test")}
				return h
			},
			ctx: servercontext.NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), httptest.NewRecorder())),
			output: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := tc.expect(t)
			output := h.Wasm(tc.ctx)
			assert.Equal(t, tc.output, output)
		})
	}
}
