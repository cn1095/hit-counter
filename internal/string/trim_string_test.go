package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimString_UnmarshalText(t *testing.T) {
	sample := TrimString("")

	tests := []struct {
		name       string
		input      *TrimString
		inputValue []byte
		output     string
	}{
		{
			name:       "trim space",
			input:      &sample,
			inputValue: []byte(" allan "),
			output:     "allan",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.UnmarshalText(tc.inputValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, string(*tc.input))
		})
	}
}

func TestTrimString_String(t *testing.T) {
	tests := []struct {
		name   string
		input  TrimString
		output string
	}{
		{
			name:   "success",
			input:  TrimString("allan"),
			output: "allan",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, tc.input.String())
		})
	}
}
