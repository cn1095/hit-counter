package sentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name      string
		clientOpt ClientOptions
		opts      []Option
		output    error
	}{
		{
			name:      "success",
			clientOpt: ClientOptions{},
			opts:      []Option{WithTags(map[string]string{"tag": ""}), WithUser(User{ID: "id"}), WithExtras(map[string]any{"extra": nil})},
			output:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Initialize(tc.clientOpt, tc.opts...)
			assert.Equal(t, tc.output, err)
		})
	}
}

func TestFlush(t *testing.T) {
	err := Initialize(ClientOptions{})
	assert.NoError(t, err)
	tests := []struct {
		name   string
		output bool
	}{
		{
			name:   "success",
			output: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Flush(time.Second)
			assert.Equal(t, tc.output, result)
		})
	}
}

func TestDebug(t *testing.T)      {}
func TestDebugMsg(t *testing.T)   {}
func TestInfo(t *testing.T)       {}
func TestInfoMsg(t *testing.T)    {}
func TestWarning(t *testing.T)    {}
func TestWarningMsg(t *testing.T) {}
func TestError(t *testing.T)      {}
func TestErrorMsg(t *testing.T)   {}
func TestFatal(t *testing.T)      {}
func TestFatalMsg(t *testing.T)   {}
