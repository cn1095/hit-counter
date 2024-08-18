package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

type HitCounterContext struct {
	echo.Context
}

// GetContext returns a context in request.
func (c *HitCounterContext) GetContext() context.Context {
	return c.Request().Context()
}

// SetContext sets a context to request.
func (c *HitCounterContext) SetContext(ctx context.Context) {
	c.SetRequest(c.Request().WithContext(ctx))
}

// WithContext set a context with new value to request.
func (c *HitCounterContext) WithContext(key, val any) {
	ctx := c.GetContext()
	c.SetContext(context.WithValue(ctx, key, val))
}

// ValueContext returns values in request context.
func (c *HitCounterContext) ValueContext(key any) any {
	return c.GetContext().Value(key)
}

// ExtraLog returns log struct.
func (c *HitCounterContext) ExtraLog() map[string]any {
	return map[string]interface{}{
		"host":       c.Request().Host,
		"ip":         c.RealIP(),
		"uri":        c.Request().RequestURI,
		"method":     c.Request().Method,
		"referer":    c.Request().Referer(),
		"user-agent": c.Request().UserAgent(),
	}
}

// NewHitCounterContext returns a new HitCounterContext.
func NewHitCounterContext(c echo.Context) *HitCounterContext {
	return &HitCounterContext{Context: c}
}
