package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	tests := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "error icons",
			input: "icons/unknown.svg",
			isErr: true,
		},
		{
			name:  "success icons",
			input: "icons/1password.svg",
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := Icons.ReadFile(tc.input)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.NotEmpty(t, output)
			}
		})
	}
}
