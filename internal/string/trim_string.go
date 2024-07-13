package string

import (
	"encoding"
	"strings"
)

type TrimString string

var _ encoding.TextUnmarshaler = (*TrimString)(nil)

// UnmarshalText set string with trim space to obj.
func (s *TrimString) UnmarshalText(text []byte) error {
	*s = TrimString(strings.TrimSpace(string(text)))
	return nil
}

// String returns string.
func (s *TrimString) String() string {
	return string(*s)
}
