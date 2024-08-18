//go:build docker

package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustNewHandler(t *testing.T) {
	tests := []struct {
		name    string
		phase   string
		addr    string
		cluster bool
	}{
		{
			name:    "success",
			phase:   "local",
			addr:    "localhost:6978",
			cluster: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var h Handler
			assert.NotPanics(t, func() {
				h = MustNewHandler(tc.phase, tc.addr, tc.cluster)
			})
			assert.NotEmpty(t, h.(*handler).websocketBreaker)
			assert.NotEmpty(t, h.(*handler).counter)
			assert.NotEmpty(t, h.(*handler).limiter)
			assert.NotEmpty(t, h.(*handler).phase)
			assert.NotEmpty(t, h.(*handler).wasmFile)
			assert.NotEmpty(t, h.(*handler).indexTemplate)
		})
	}
}
