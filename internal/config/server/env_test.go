package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustInitializeLocalEnvironments(t *testing.T) {
	tests := []struct {
		name  string
		isErr bool
	}{
		{
			name:  "success",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.True(t, tc.isErr)
				}
			}()
			MustInitializeLocalEnvironments()
		})
	}
}

func TestMustGetLocalEnvironments(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.NoError(t, err.(error))
				}
			}()
			output := MustGetLocalEnvironments()
			assert.NotEmpty(t, output)
			assert.Equal(t, localEnvs, output)
		})
	}
}
