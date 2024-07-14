package sentry

import (
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

type ClientOptions = sentry.ClientOptions
type User = sentry.User
type Context = sentry.Context
type SentryHub = sentry.Hub

// Initialize initializes sentry sdk.
func Initialize(clientOpt ClientOptions, opts ...Option) error {
	if clientOpt.BeforeSend == nil {
		clientOpt.BeforeSend = beforeSend
	}

	if err := sentry.Init(clientOpt); err != nil {
		return err
	}

	for _, opt := range opts {
		sentry.ConfigureScope(opt.apply)
	}

	return nil
}

// Flush sends to events and waits for duration.
func Flush(d time.Duration) bool {
	return sentry.Flush(d)
}

// Debug sends an error in debug mode.
func Debug(err error, opts ...Option) {
	sendError(sentry.CurrentHub().Clone(), err, sentry.LevelDebug, opts...)
}

// DebugMsg sends a message in debug mode.
func DebugMsg(msg string, opts ...Option) {
	sendMsg(sentry.CurrentHub().Clone(), msg, sentry.LevelDebug, opts...)
}

// Info sends an error in info mode.
func Info(err error, opts ...Option) {
	sendError(sentry.CurrentHub().Clone(), err, sentry.LevelInfo, opts...)
}

// InfoMsg sends a msg in info mode.
func InfoMsg(msg string, opts ...Option) {
	sendMsg(sentry.CurrentHub().Clone(), msg, sentry.LevelInfo, opts...)
}

// Warning sends an error in warning mode.
func Warning(err error, opts ...Option) {
	sendError(sentry.CurrentHub().Clone(), err, sentry.LevelWarning, opts...)
}

// WarningMsg sends a message in warning mode.
func WarningMsg(msg string, opts ...Option) {
	sendMsg(sentry.CurrentHub().Clone(), msg, sentry.LevelWarning, opts...)
}

// Error sends an error in error mode.
func Error(err error, opts ...Option) {
	sendError(sentry.CurrentHub().Clone(), err, sentry.LevelError, opts...)
}

// ErrorMsg sends a message in error mode.
func ErrorMsg(msg string, opts ...Option) {
	sendMsg(sentry.CurrentHub().Clone(), msg, sentry.LevelError, opts...)
}

// Fatal sends an error in fatal mode.
func Fatal(err error, opts ...Option) {
	sendError(sentry.CurrentHub().Clone(), err, sentry.LevelFatal, opts...)
}

// FatalMsg sends a message in fatal mode.
func FatalMsg(msg string, opts ...Option) {
	sendMsg(sentry.CurrentHub().Clone(), msg, sentry.LevelFatal, opts...)
}

func sendMsg(sentryHub *SentryHub, msg string, level sentry.Level, opts ...Option) {
	sentryHub.WithScope(func(scope *sentry.Scope) {
		for _, opt := range opts {
			opt.apply(scope)
		}
		scope.SetLevel(level)
		sentryHub.CaptureMessage(msg)
	})
}

func sendError(sentryHub *SentryHub, err error, level sentry.Level, opts ...Option) {
	sentryHub.WithScope(func(scope *sentry.Scope) {
		for _, opt := range opts {
			opt.apply(scope)
		}
		scope.SetLevel(level)
		sentryHub.CaptureException(err)
	})
}

func beforeSend(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	exceptions := make([]sentry.Exception, 0, len(event.Exception))

	for _, exception := range event.Exception {
		if !strings.HasPrefix(exception.Type, "*error.ErrorWithStackTrace") &&
			!strings.HasPrefix(exception.Type, "*errors.withStack") &&
			!strings.HasPrefix(exception.Type, "*fmt.") {
			exceptions = append(exceptions, exception)
			continue
		}

		if exception.Stacktrace != nil {
			frameSize := len(exception.Stacktrace.Frames)
			if frameSize > 0 {
				topFrame := exception.Stacktrace.Frames[frameSize-1]
				exception.Type = topFrame.Module + "." + topFrame.Function
			}
		}
		exceptions = append(exceptions, exception)
	}
	event.Exception = exceptions
	return event
}
