package main // import "github.com/cn1095/hit-counter"

import (
	"embed"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/cn1095/hit-counter/internal"

	"path/filepath"

	"github.com/cn1095/hit-counter/env"
	"github.com/labstack/echo/v4"
)

//go:embed public/* view/*
var embeddedFiles embed.FS // 嵌入 public 和 view 目录中的所有文件

var (
	address = flag.String("port", ":8080", "port")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// initialize sentry
	name, _ := os.Hostname()
	if err := internal.InitSentry(env.GetSentryDSN(), env.GetPhase(), env.GetPhase(),
		name, true, env.GetDebug()); err != nil {
		log.Println(err)
	}

	e := echo.New()

	// make options for echo server.
	var opts []Option

	// debug option
	opts = append(opts, WithDebugOption(env.GetDebug()))

	var dir string
	var file string
	if env.GetLogPath() != "" {
		dir, file = filepath.Split(env.GetLogPath())
	}

	// logger option
	logger, err := internal.NewLogger(dir, file)
	if err != nil {
		log.Panic(err)
	}
	opts = append(opts, WithLogger(logger))

	// add middleware
	if err := AddMiddleware(e, opts...); err != nil {
		log.Panic(err)
	}
	
	// 设置静态文件目录为 embed 文件系统（替代 e.Static）
	staticFS := echo.MustSubFS(embeddedFiles, "public")
	e.StaticFS("/", staticFS)
	
	// add route
	if err := AddRoute(e, env.GetRedisAddrs()[0]); err != nil {
		log.Panic(err)
	}

	if *tls {
		// start TLS server with let's encrypt certification.
		e.StartAutoTLS(*address)
	} else {
		e.Start(*address)
	}
}
