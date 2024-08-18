package handler

import (
	"fmt"
	"html/template"

	websocket "github.com/gjbae1212/go-ws-broadcast"
	"github.com/gjbae1212/hit-counter/internal/counter"
	"github.com/gjbae1212/hit-counter/internal/limiter"
	"github.com/gjbae1212/hit-counter/internal/sentry"
	"github.com/gjbae1212/hit-counter/web"
	"github.com/labstack/echo/v4"
	perrors "github.com/pkg/errors"
)

type Handler interface {
	Error(err error, c echo.Context)
	HealthCheck(c echo.Context) error

	WebSocket(c echo.Context) error
	Wasm(c echo.Context) error
	Index(c echo.Context) error
}

type handler struct {
	phase            string
	counter          counter.Counter
	limiter          limiter.Limiter
	websocketBreaker websocket.Breaker

	wasmFile      []byte
	indexTemplate *template.Template
}

// MustNewHandler returns a new Handler instance.
func MustNewHandler(phase string, redisAddr string, redisCluster bool) Handler {
	if phase == "" || redisAddr == "" {
		panic(fmt.Errorf("[err] empty redis addr"))
	}

	// create websocket breaker.
	websocketBreaker, err := websocket.NewBreaker(websocket.WithMaxReadLimit(1024),
		websocket.WithMaxMessagePoolLength(500),
		websocket.WithErrorHandlerOption(func(err error) {
			sentry.Error(perrors.WithStack(err))
		}))
	if err != nil {
		panic(perrors.WithStack(err))
	}

	// get wasm file.
	wasmFile, err := web.GetWasm(phase)
	if err != nil {
		panic(perrors.WithStack(err))
	}

	// get index html file.
	indexHtml, err := web.GetIndexHtml(phase)
	if err != nil {
		panic(perrors.WithStack(err))
	}
	// parse index template.
	indexTemplate, err := template.New("index").Parse(string(indexHtml))
	if err != nil {
		panic(perrors.WithStack(err))
	}

	return &handler{
		phase:   phase,
		counter: counter.MustNewCounter(redisAddr, redisCluster),
		limiter: limiter.MustNewLimiter(redisAddr, redisCluster,
			limiter.WithRateLimitWindow(1),
			limiter.WithRateLimitCount(50)),
		websocketBreaker: websocketBreaker,
		wasmFile:         wasmFile,
		indexTemplate:    indexTemplate,
	}
}
