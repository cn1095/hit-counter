package url

import (
	"fmt"

	"github.com/gjbae1212/hit-counter/pkg/urlx"
	perrors "github.com/pkg/errors"
)

// ParseURL parses url.
func ParseURL(s string) (schema, host, port, path, query, fragment string, err error) {
	if s == "" {
		err = perrors.WithStack(fmt.Errorf("[err] ParseURL empty string"))
		return
	}

	url, err := urlx.Parse(s)
	if err != nil {
		err = perrors.WithStack(err)
		return
	}

	schema = url.Scheme

	host, port, err = urlx.SplitHostPort(url)
	if err != nil {
		err = perrors.WithStack(err)
		return
	}
	if schema == "http" && port == "" {
		port = "80"
	} else if schema == "https" && port == "" {
		port = "443"
	}

	path = url.Path
	query = url.RawQuery
	fragment = url.Fragment
	return
}
