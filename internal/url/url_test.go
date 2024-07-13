package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		schema   string
		domain   string
		port     string
		path     string
		query    string
		fragment string
	}{
		{
			name:     "sample1",
			input:    "http://naver.com/aa/bb?cc=dd&ee=ff#fragment",
			schema:   "http",
			domain:   "naver.com",
			port:     "80",
			path:     "/aa/bb",
			query:    "cc=dd&ee=ff",
			fragment: "fragment",
		},
		{
			name:     "sample2",
			input:    "cc.com:8080/aa/bb",
			schema:   "http",
			domain:   "cc.com",
			port:     "8080",
			path:     "/aa/bb",
			query:    "",
			fragment: "",
		},
		{
			name:     "sample3",
			input:    "https://naver.com/aa/bb?cc=dd&ee=ff#fragment",
			schema:   "https",
			domain:   "naver.com",
			port:     "443",
			path:     "/aa/bb",
			query:    "cc=dd&ee=ff",
			fragment: "fragment",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schema, domain, port, path, query, fragment, err := ParseURL(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.schema, schema)
			assert.Equal(t, tc.domain, domain)
			assert.Equal(t, tc.port, port)
			assert.Equal(t, tc.path, path)
			assert.Equal(t, tc.query, query)
			assert.Equal(t, tc.fragment, fragment)
		})
	}
}
