package server

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/gjbae1212/hit-counter/internal/app/server/config"
	"github.com/gjbae1212/hit-counter/internal/sentry"
	"github.com/labstack/echo/v4"
	perrors "github.com/pkg/errors"
	_ "go.uber.org/automaxprocs"
)

type httpServer struct {
	echo *echo.Echo
	env  *config.LocalEnvironments
}

func (hs *httpServer) initializeMiddlewares() error {
	// TODO:
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
