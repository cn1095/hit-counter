package handler

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHitCounterContext_ExtraLog(t *testing.T) {
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	ctx := &HitCounterContext{e.NewContext(r, nil)}
	extraLog := ctx.ExtraLog()
	assert.Equal(t, "GET", extraLog["method"])
	assert.Equal(t, "localhost", extraLog["host"])
	assert.Len(t, extraLog, 6)
}

func TestHitCounterContext_ValueContext(t *testing.T) {
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	ctx := &HitCounterContext{e.NewContext(r, nil)}
	ctx.WithContext("allan", "hi")
	value := ctx.ValueContext("allan")
	assert.Equal(t, "hi", value.(string))
}

func TestHitCounterContext_WithContext(t *testing.T) {
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	ctx := &HitCounterContext{e.NewContext(r, nil)}
	ctx.WithContext("allan", "hi")
	ctx.WithContext("test", "testhi")
	value := ctx.ValueContext("allan")
	assert.Equal(t, "hi", value.(string))
}

func TestHitCounterContext_SetContext(t *testing.T) {
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	ctx := &HitCounterContext{e.NewContext(r, nil)}
	ctx.SetContext(context.WithValue(context.Background(), "test", "allan"))
	value := ctx.ValueContext("test")
	assert.Equal(t, "allan", value.(string))
}

func TestHitCounterContext_GetContext(t *testing.T) {
	e := echo.New()
	r := httptest.NewRequest("GET", "http://localhost", nil)
	ctx := &HitCounterContext{e.NewContext(r, nil)}
	ctx.SetContext(context.WithValue(context.Background(), "test", "allan"))
	value := ctx.GetContext().Value("test")
	assert.Equal(t, "allan", value)
}
