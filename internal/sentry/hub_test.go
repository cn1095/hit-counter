package sentry

import (
	"context"
	"reflect"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestGetHub(t *testing.T) {
	tests := []struct {
		name    string
		isExist bool
	}{
		{
			name:    "success",
			isExist: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hub := GetHub()
			assert.Equal(t, tc.isExist, hub != nil)
		})
	}
}

func TestGetOrCreateHubOnContext(t *testing.T) {
	tests := []struct {
		name      string
		inputCtx  context.Context
		isErr     bool
		isAlready bool
	}{
		{
			name:      "empty",
			inputCtx:  context.Background(),
			isErr:     false,
			isAlready: false,
		},
		{
			name:      "already exist context",
			inputCtx:  sentry.SetHubOnContext(context.Background(), sentry.CurrentHub().Clone()),
			isErr:     false,
			isAlready: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isAlready, sentry.HasHubOnContext(tc.inputCtx))
			hub, outputCtx, ok, err := GetOrCreateHubOnContext(tc.inputCtx)
			assert.Equal(t, tc.isErr, err != nil)
			assert.Equal(t, !tc.isAlready, ok)
			assert.Equal(t, true, sentry.HasHubOnContext(outputCtx))
			assert.Equal(t, hub.hub, sentry.GetHubFromContext(outputCtx))
		})
	}
}

func TestMustGetHubOnContext(t *testing.T) {
	tests := []struct {
		name     string
		inputCtx context.Context
		isErr    bool
	}{
		{
			name:     "error",
			inputCtx: context.Background(),
			isErr:    true,
		},
		{
			name:     "already exist context",
			inputCtx: sentry.SetHubOnContext(context.Background(), sentry.CurrentHub().Clone()),
			isErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.True(t, tc.isErr)
				}
			}()
			hub := MustGetHubOnContext(tc.inputCtx)
			assert.NotEmpty(t, hub)
		})
	}
}

func TestWrappedHub_With(t *testing.T) {
	hub, ctx, ok, err := GetOrCreateHubOnContext(context.Background())
	assert.False(t, err != nil)
	assert.True(t, ok)
	assert.NotEmpty(t, ctx)

	tests := []struct {
		name             string
		hub              *WrappedHub
		opts             []Option
		checkFingerprint []string
	}{
		{
			name: "user and fingerprint and tags",
			hub:  hub,
			opts: []Option{
				WithFingerprint([]string{"fingerprint"}),
			},
			checkFingerprint: []string{"fingerprint"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hub.With(tc.opts...)
			scope := hub.hub.Scope()
			fingerprint := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("fingerprint")
			assert.Equal(t, reflect.ValueOf(tc.checkFingerprint).String(), fingerprint.String())
		})
	}
}

func TestWrappedHub_Debug(t *testing.T)      {}
func TestWrappedHub_DebugMsg(t *testing.T)   {}
func TestWrappedHub_Info(t *testing.T)       {}
func TestWrappedHub_InfoMsg(t *testing.T)    {}
func TestWrappedHub_Warning(t *testing.T)    {}
func TestWrappedHub_WarningMsg(t *testing.T) {}
func TestWrappedHub_Error(t *testing.T)      {}
func TestWrappedHub_ErrorMsg(t *testing.T)   {}
func TestWrappedHub_Fatal(t *testing.T)      {}
func TestWrappedHub_FatalMsg(t *testing.T)   {}
