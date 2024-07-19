package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/gjbae1212/hit-counter/internal/app/server/config"
	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/gjbae1212/hit-counter/internal/sentry"
	"github.com/gjbae1212/hit-counter/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	perrors "github.com/pkg/errors"
	_ "go.uber.org/automaxprocs"
)

const (
	healthCheckPath = "/healthcheck"
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

	// set panic middleware.
	hs.echo.Use(middleware.Recover())

	// set sentry middleware.
	hs.echo.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// set secure middleware.
	hs.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		HSTSMaxAge:            2592000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: "default-src 'none'; style-src 'unsafe-inline'",
	}))

	if !hs.env.Debug {
		hs.echo.Use(middleware.HTTPSRedirect())
	}

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

	// set cookie duration 24 hour.
	hs.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const (
				cookieName = "ckid"
			)

			hitCtx := c.(*servercontext.HitCounterContext)
			cookie, err := c.Cookie(cookieName)
			if err != nil {
				v := fmt.Sprintf("%s-%d", c.RealIP(), time.Now().UnixNano())
				b64 := base64.StdEncoding.EncodeToString([]byte(v))
				cookie = &http.Cookie{
					Name:     cookieName,
					Value:    b64,
					Expires:  time.Now().Add(24 * time.Hour),
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteNoneMode,
				}
				hitCtx.SetCookie(cookie)
			}
			hitCtx.Set(cookie.Name, cookie.Value)
			return next(hitCtx)
		}
	})

	// main middleware.
	hs.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitCtx := c.(*servercontext.HitCounterContext)
			start := hitCtx.ValueContext("start_time").(time.Time)
			extraLog := hitCtx.ValueContext("extra_log").(map[string]interface{})

			// process main handler.
			err := next(hitCtx)
			stop := time.Now()
			latency := stop.Sub(start)

			if err != nil {
				code := http.StatusInternalServerError
				if httpErr := (*echo.HTTPError)(nil); errors.As(err, &httpErr) {
					code = httpErr.Code
				} else if hitCtx.Response().Status >= http.StatusBadRequest {
					code = hitCtx.Response().Status
				}

				extraLog["status"] = code
				extraLog["error"] = fmt.Sprintf("%v\n", err)
				if code >= http.StatusInternalServerError {
					sentry.Error(err)
					extraLog["latency"] = strconv.FormatInt(int64(latency), 10)
					extraLog["latency_human"] = latency.String()
				}
				hitCtx.Logger().Errorj(extraLog)
				return err
			}

			if extraLog["uri"] != healthCheckPath {
				extraLog["status"] = hitCtx.Response().Status
				extraLog["latency"] = strconv.FormatInt(int64(latency), 10)
				extraLog["latency_human"] = latency.String()
				hitCtx.Logger().Infoj(extraLog)
			}
			return nil
		}
	})

	return nil
}

func (hs *httpServer) initializeRoutes() error {
	// set public route.
	public, err := fs.Sub(web.Public, "public")
	if err != nil {
		return perrors.WithStack(err)
	}
	hs.echo.GET("/public/*",
		echo.WrapHandler(
			http.StripPrefix("/public/", http.FileServer(http.FS(public))),
		),
	)

	// set health check route.

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
