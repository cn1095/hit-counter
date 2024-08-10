//go:build docker

package limiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiter_IsBlackList(t *testing.T) {
	tests := []struct {
		name   string
		ready  func(t *testing.T) Limiter
		id     string
		output bool
		isErr  bool
	}{
		{
			name: "error",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			isErr: true,
		},
		{
			name: "success",
			ready: func(t *testing.T) Limiter {
				mock := MustNewLimiter(mockRedisAddr, false)
				_ = mock.AddBlackList(context.Background(), "test")
				return mock
			},
			id:     "test",
			output: true,
			isErr:  false,
		},
		{
			name: "miss",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			id:     "empty",
			output: false,
			isErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lt := tc.ready(t)
			output, err := lt.IsBlackList(context.Background(), tc.id)
			assert.Equal(t, tc.isErr, err != nil)
			assert.Equal(t, tc.output, output)
		})
	}
}

func TestLimiter_AddBlackList(t *testing.T) {
	tests := []struct {
		name  string
		ready func(t *testing.T) Limiter
		id    string
		isErr bool
	}{
		{
			name: "error",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			isErr: true,
		},
		{
			name: "success",
			ready: func(t *testing.T) Limiter {
				return MustNewLimiter(mockRedisAddr, false)
			},
			id:    "test",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lt := tc.ready(t)
			err := lt.AddBlackList(context.Background(), tc.id)
			assert.Equal(t, tc.isErr, err != nil)
		})
	}
}
