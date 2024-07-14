package sentry

import (
	"reflect"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestWithUser(t *testing.T) {
	sample := User{ID: "gjbae1212"}
	tests := []struct {
		name   string
		input  User
		output string
	}{
		{
			name:   "success",
			input:  sample,
			output: sample.ID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithUser(tc.input)
			scope := sentry.NewScope()
			opt.apply(scope)
			result := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("user").FieldByName("ID").String()
			assert.Equal(t, tc.output, result)
		})
	}
}

func TestWithContexts(t *testing.T) {
	sample := map[string]Context{"1": nil}
	tests := []struct {
		name   string
		input  map[string]Context
		output string
	}{
		{
			name:  "success",
			input: sample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithContexts(tc.input)
			scope := sentry.NewScope()
			opt.apply(scope)
			result := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("contexts")
			assert.Equal(t, reflect.ValueOf(tc.input).String(), result.String())
		})
	}
}

func TestWithExtras(t *testing.T) {
	sample := map[string]any{"1": nil}
	tests := []struct {
		name   string
		input  map[string]any
		output string
	}{
		{
			name:  "success",
			input: sample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithExtras(tc.input)
			scope := sentry.NewScope()
			opt.apply(scope)
			result := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("extra")
			assert.Equal(t, reflect.ValueOf(tc.input).String(), result.String())
		})
	}
}

func TestWithTags(t *testing.T) {
	sample := map[string]string{"1": ""}
	tests := []struct {
		name   string
		input  map[string]string
		output string
	}{
		{
			name:  "success",
			input: sample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithTags(tc.input)
			scope := sentry.NewScope()
			opt.apply(scope)
			result := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("tags")
			assert.Equal(t, reflect.ValueOf(tc.input).String(), result.String())
		})
	}
}

func TestWithFingerprint(t *testing.T) {
	sample := []string{"hello"}
	tests := []struct {
		name   string
		input  []string
		output []string
	}{
		{
			name:  "success",
			input: sample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithFingerprint(tc.input)
			scope := sentry.NewScope()
			opt.apply(scope)
			result := reflect.Indirect(reflect.ValueOf(scope)).FieldByName("fingerprint")
			assert.Equal(t, reflect.ValueOf(tc.input).String(), result.String())
		})
	}
}
