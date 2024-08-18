package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "error icons",
			input: "public/unknown.png",
			isErr: true,
		},
		{
			name:  "success icons",
			input: "public/icon.png",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := Public.ReadFile(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.NotEmpty(t, output)
			}
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "wasm",
			input: "view/production/hits.wasm",
			isErr: false,
		},
		{
			name:  "html",
			input: "view/production/index.html",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := View.ReadFile(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.NotEmpty(t, output)
			}
		})
	}
}

func TestGetWasm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "error",
			input: "empty",
			isErr: true,
		},
		{
			name:  "success-local",
			input: "local",
			isErr: false,
		},
		{
			name:  "success-production",
			input: "production",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := GetWasm(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.NotEmpty(t, output)
			}
		})
	}
}

func TestGetIndexHtml(t *testing.T) {
	tests := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "error",
			input: "empty",
			isErr: true,
		},
		{
			name:  "success-local",
			input: "local",
			isErr: false,
		},
		{
			name:  "success-production",
			input: "production",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := GetIndexHtml(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.NotEmpty(t, output)
			}
		})
	}
}
