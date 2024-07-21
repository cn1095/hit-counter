package main

import (
	"flag"
	"fmt"
	"log/slog"
	"runtime"

	appserver "github.com/gjbae1212/hit-counter/internal/app/server"
	_ "go.uber.org/automaxprocs"
)

var (
	address = flag.String("addr", ":8080", "address")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()
	// show number of logical cpu and go processor(automaxprocs).
	slog.Info(fmt.Sprintf("logical cpu = %d, go processor = %d", runtime.NumCPU(), runtime.GOMAXPROCS(0)))
	appserver.Run(*address, *tls)
}
