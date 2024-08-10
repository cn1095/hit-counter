package limiter

import (
	"fmt"
	"runtime"
	"time"

	perrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	redisMaxRetries         = 1
	redisMaxConnIdleTimeout = 5 * time.Minute
	redisMaxConnAge         = 10 * time.Minute
	redisPoolSizeRatio      = 10
)

var (
	defaultOpts = []Option{
		WithRateLimitWindow(1),
		WithRateLimitCount(50),
	}
)

type Limiter interface {
	RateMethod
	BlacklistMethod
}

type limiter struct {
	cfg           *config
	redisClient   redis.UniversalClient
	windowCounter *fixedWindowCounter
}

// MustNewLimiter returns a new Limiter instance.
func MustNewLimiter(redisAddr string, redisCluster bool, opts ...Option) Limiter {
	if redisAddr == "" {
		panic(fmt.Errorf("[err] empty redis addr"))
	}

	cfg := &config{}
	for _, opt := range append(defaultOpts, opts...) {
		opt.apply(cfg)
	}

	windowCounter, err := newFixedWindowCounter(cfg.rateLimitWindow, cfg.rateLimitCount)
	if err != nil {
		panic(perrors.WithStack(err))
	}

	var redisClient redis.UniversalClient
	if redisCluster {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           []string{redisAddr},
			MaxRetries:      redisMaxRetries,
			ConnMaxIdleTime: redisMaxConnIdleTimeout,
			ConnMaxLifetime: redisMaxConnAge,
			PoolSize:        runtime.GOMAXPROCS(0) * redisPoolSizeRatio,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:            redisAddr,
			MaxRetries:      redisMaxRetries,
			ConnMaxIdleTime: redisMaxConnIdleTimeout,
			ConnMaxLifetime: redisMaxConnAge,
			PoolSize:        runtime.GOMAXPROCS(0) * redisPoolSizeRatio,
		})
	}

	return &limiter{
		cfg:           cfg,
		redisClient:   redisClient,
		windowCounter: windowCounter,
	}
}
