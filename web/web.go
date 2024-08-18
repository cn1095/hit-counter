package web

import (
	"embed"

	perrors "github.com/pkg/errors"
)

//go:embed public
var Public embed.FS

//go:embed view
var View embed.FS

// GetWasm returns wasm file.
func GetWasm(phase string) ([]byte, error) {
	bys, err := View.ReadFile("view/" + phase + "/hits.wasm")
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return bys, nil
}

// GetIndexHtml returns index.html file.
func GetIndexHtml(phase string) ([]byte, error) {
	bys, err := View.ReadFile("view/" + phase + "/index.html")
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return bys, nil
}
