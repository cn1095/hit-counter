package main

import (
	"flag"

	appserver "github.com/gjbae1212/hit-counter/internal/app/server"
)

var (
	address = flag.String("addr", ":8080", "address")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()
	appserver.Run(*address, *tls)
}
