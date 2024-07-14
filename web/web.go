package web

import "embed"

//go:embed public
var Public embed.FS

//go:embed view
var View embed.FS
