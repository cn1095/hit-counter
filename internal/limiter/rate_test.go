//go:build docker

package limiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiter_Allow(t *testing.T) {
	tests := []struct {
		name   string
		ready  func(t *testing.T) Limiter
		id     string
		output bool
	}{
		{
			name: "false(invalid-param)",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			output: false,
		},
		{
			name: "success(allow)",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			id:     "test",
			output: true,
		},
		{
			name: "success(not allow)",
			ready: func(t *testing.T) Limiter {
				mock := MustNewLimiter(mockRedisAddr, false)
				for range 50 {
					mock.Allow(context.Background(), "hello")
				}
				return mock
			},
			id:     "hello",
			output: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lt := tc.ready(t)
			output := lt.Allow(context.Background(), tc.id)
			assert.Equal(t, tc.output, output)
		})
	}
}
