package limiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithRateLimitCount(t *testing.T) {
	tests := []struct {
		name  string
		count int64
	}{
		{
			name:  "success",
			count: 10,
		},
	}

	for _, tc := range tests {
		cfg := &config{}
		opt := WithRateLimitCount(tc.count)
		opt.apply(cfg)
		assert.Equal(t, tc.count, cfg.rateLimitCount)
	}
}

func TestWithRateLimitWindow(t *testing.T) {
	tests := []struct {
		name   string
		window int64
	}{
		{
			name:   "success",
			window: 20,
		},
	}

	for _, tc := range tests {
		cfg := &config{}
		opt := WithRateLimitWindow(tc.window)
		opt.apply(cfg)
		assert.Equal(t, tc.window, cfg.rateLimitWindow)
	}
}
