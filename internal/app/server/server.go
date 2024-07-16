package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/gjbae1212/hit-counter/internal/app/server/config"
	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/gjbae1212/hit-counter/internal/sentry"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	perrors "github.com/pkg/errors"
	_ "go.uber.org/automaxprocs"
)

type httpServer struct {
	echo *echo.Echo
	env  *config.LocalEnvironments
}

func (hs *httpServer) initializeMiddlewares() error {
	hs.echo.Debug = hs.env.Debug
	hs.echo.HideBanner = true
	hs.echo.HidePort = true
	hs.echo.Server.ReadTimeout = 10 * time.Second
	hs.echo.Server.WriteTimeout = 10 * time.Second

	// set sentry middleware.
	hs.echo.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// set secure middleware.
	hs.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		HSTSMaxAge:            2592000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: "default-src 'none'; style-src 'unsafe-inline'",
	}))

	hs.echo.Use(middleware.HTTPSRedirect())

	hs.echo.Use(middleware.RemoveTrailingSlash())

	hs.echo.Use(middleware.NonWWWRedirect())

	hs.echo.Use(middleware.Rewrite(map[string]string{
		"/static/*": "/public/$1",
	}))

	// set custom context.
	hs.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitCtx := &servercontext.HitCounterContext{Context: c}

			// set start time.
			hitCtx.WithContext("start_time", time.Now())

			// set deadline.
			timeout := 15 * time.Second

			ctx, cancel := context.WithTimeout(hitCtx.GetContext(), timeout)
			defer cancel()
			hitCtx.SetContext(ctx)

			// set extra log.
			extraLog := hitCtx.ExtraLog()
			hitCtx.WithContext("extra_log", extraLog)
			return next(hitCtx)
		}
	})
	// TODO: cookie middleware

	// TODO: main middleware

	return nil
}

func (hs *httpServer) initializeRoutes() error {
	// TODO:
	return nil
}

// Run is a function to run the server.
func Run(addr string, tls bool) {
	if addr == "" {
		panic(perrors.WithStack(fmt.Errorf("[err] invalid address")))
	}

	// initialize and get local environments.
	env := config.MustGetLocalEnvironments()

	// initialize default slog logger.
	logLevel := slog.LevelInfo
	if env.Debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	logger.With(
		slog.String("phase", env.Phase.String()),
	)
	slog.SetDefault(logger)

	// initialize sentry.
	if env.SentryDSN.String() != "" {
		if err := sentry.Initialize(
			sentry.ClientOptions{
				Dsn:              env.SentryDSN.String(),
				Environment:      env.Phase.String(),
				Debug:            env.Debug,
				AttachStacktrace: true,
			},
			sentry.WithTags(
				map[string]string{
					"app": "hit-counter",
				},
			),
		); err != nil {
			panic(perrors.WithStack(err))
		}
	}

	// show number of logical cpu and go processor(automaxprocs).
	slog.Info(fmt.Sprintf("logical cpu = %d, go processor = %d", runtime.NumCPU(), runtime.GOMAXPROCS(0)))

	// create server.
	server := &httpServer{
		echo: echo.New(),
		env:  env,
	}

	// initialize middlewares.
	if err := server.initializeMiddlewares(); err != nil {
		panic(perrors.WithStack(err))
	}

	// initialize routes.
	if err := server.initializeRoutes(); err != nil {
		panic(perrors.WithStack(err))
	}

	if tls {
		if err := server.echo.StartAutoTLS(addr); err != nil {
			slog.Error(perrors.WithStack(err).Error())
		}
	} else {
		if err := server.echo.Start(addr); err != nil {
			slog.Error(perrors.WithStack(err).Error())
		}
	}
}
