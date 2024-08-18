package handler

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHitCounterContext_ExtraLog(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *HitCounterContext
		output map[string]any
	}{
		{
			name: "success",
			ctx: NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), nil)),
			output: map[string]any{"method": "GET", "host": "localhost", "ip": ""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.ctx.ExtraLog()
			assert.Equal(t, tc.output["method"], output["method"])
			assert.Equal(t, tc.output["localhost"], output["localhost"])
		})
	}
}

func TestHitCounterContext_ValueContext(t *testing.T) {
	ctx := NewHitCounterContext(echo.New().NewContext(
		httptest.NewRequest("GET", "http://localhost", nil), nil))
	ctx.WithContext("allan", "hi")
	tests := []struct {
		name   string
		ctx    *HitCounterContext
		input  any
		output any
	}{
		{
			name:   "success",
			ctx:    ctx,
			input:  "allan",
			output: "hi",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.ctx.ValueContext(tc.input)
			assert.Equal(t, tc.output, output)
		})
	}
}

func TestHitCounterContext_GetContext(t *testing.T) {
	ctx := NewHitCounterContext(echo.New().NewContext(
		httptest.NewRequest("GET", "http://localhost", nil), nil))
	ctx.SetContext(context.WithValue(context.Background(), "allan", "hi"))
	tests := []struct {
		name   string
		ctx    *HitCounterContext
		input  any
		output any
	}{
		{
			name:   "success",
			ctx:    ctx,
			input:  "allan",
			output: "hi",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.ctx.GetContext()
			assert.Equal(t, tc.output, output.Value(tc.input))
		})
	}
}
