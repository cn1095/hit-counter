package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Error(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *servercontext.HitCounterContext
		err    error
		output int
	}{
		{
			name: "no-error",
			ctx: servercontext.NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), httptest.NewRecorder())),
			err:    nil,
			output: http.StatusInternalServerError,
		},
		{
			name: "error",
			ctx: servercontext.NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), httptest.NewRecorder())),
			err:    fmt.Errorf("[err]"),
			output: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &handler{}
			h.Error(tc.err, tc.ctx)
			assert.Equal(t, tc.output, tc.ctx.Response().Status)
		})
	}
}
