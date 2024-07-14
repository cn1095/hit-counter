package sentry

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
)

// WrappedHub wraps sentry hub for supporting extra features.
type WrappedHub struct {
	hub *SentryHub
}

// With configures options.
func (h *WrappedHub) With(opts ...Option) *WrappedHub {
	for _, opt := range opts {
		h.hub.ConfigureScope(opt.apply)
	}
	return h
}

// Debug sends an error in debug mode.
func (h *WrappedHub) Debug(err error, opts ...Option) {
	sendError(h.hub, err, sentry.LevelDebug, opts...)
}

// DebugMsg sends a message in debug mode.
func (h *WrappedHub) DebugMsg(msg string, opts ...Option) {
	sendMsg(h.hub, msg, sentry.LevelDebug, opts...)
}

// Info sends an error in info mode.
func (h *WrappedHub) Info(err error, opts ...Option) {
	sendError(h.hub, err, sentry.LevelInfo, opts...)
}

// InfoMsg sends a msg in info mode.
func (h *WrappedHub) InfoMsg(msg string, opts ...Option) {
	sendMsg(h.hub, msg, sentry.LevelInfo, opts...)
}

// Warning sends an error in warning mode.
func (h *WrappedHub) Warning(err error, opts ...Option) {
	sendError(h.hub, err, sentry.LevelWarning, opts...)
}

// WarningMsg sends a message in warning mode.
func (h *WrappedHub) WarningMsg(msg string, opts ...Option) {
	sendMsg(h.hub, msg, sentry.LevelWarning, opts...)
}

// Error sends an error in error mode.
func (h *WrappedHub) Error(err error, opts ...Option) {
	sendError(h.hub, err, sentry.LevelError, opts...)
}

// ErrorMsg sends a message in error mode.
func (h *WrappedHub) ErrorMsg(msg string, opts ...Option) {
	sendMsg(h.hub, msg, sentry.LevelError, opts...)
}

// Fatal sends an error in fatal mode.
func (h *WrappedHub) Fatal(err error, opts ...Option) {
	sendError(h.hub, err, sentry.LevelFatal, opts...)
}

// FatalMsg sends a message in fatal mode.
func (h *WrappedHub) FatalMsg(msg string, opts ...Option) {
	sendMsg(h.hub, msg, sentry.LevelFatal, opts...)
}

// GetOrCreateHubOnContext returns a hub, a context from context.
// It might set a hub to a context and return its context if a context doesn't have hub on context.
func GetOrCreateHubOnContext(ctx context.Context) (*WrappedHub, context.Context, bool, error) {
	if ctx == nil {
		return nil, nil, false, fmt.Errorf("[error] invalid params")
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		// set a hub to context.
		newctx := sentry.SetHubOnContext(ctx, hub)
		return &WrappedHub{hub: hub}, newctx, true, nil
	}

	return &WrappedHub{hub: hub}, ctx, false, nil
}

// MustGetHubOnContext returns a hub.
// Raise panic if a hub doesn't exist.
func MustGetHubOnContext(ctx context.Context) *WrappedHub {
	if ctx == nil {
		panic(fmt.Errorf("[error] invalid param(ctx)"))
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		panic(fmt.Errorf("[error] invalid param(hub)"))
	}

	return &WrappedHub{hub: hub}
}

// GetHub returns a hub.
func GetHub() *WrappedHub {
	hub := sentry.CurrentHub().Clone()
	return &WrappedHub{hub: hub}
}
