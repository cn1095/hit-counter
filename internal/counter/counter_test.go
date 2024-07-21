//go:build docker

package counter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockRedisAddr = "localhost:6379"
)

var (
	mockCounter Counter
)

func TestMustCounter(t *testing.T) {
	tests := []struct {
		name  string
		addr  string
		isErr bool
	}{
		{
			name:  "error",
			addr:  "",
			isErr: true,
		},
		{
			name:  "success",
			addr:  mockRedisAddr,
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isErr {
				assert.Panics(t, func() {
					MustCounter(tc.addr, false)
				})
			} else {
				assert.NotPanics(t, func() {
					MustCounter(tc.addr, false)
				})
			}
		})
	}
}

func TestMain(m *testing.M) {
	mockCounter = MustCounter(mockRedisAddr, false)
	os.Exit(m.Run())
}
