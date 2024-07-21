package counter

import (
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	redisMaxRetries         = 1
	redisMaxConnIdleTimeout = 5 * time.Minute
	redisMaxConnAge         = 10 * time.Minute
	redisPoolSizeRatio      = 10
)

// Score presents result for response.
type Score struct {
	Name  string
	Value int64
}

type Counter interface {
	HitMethod
	RankMethod
}

type counter struct {
	redisClient redis.UniversalClient
}

// MustCounter creates counter.
func MustCounter(redisAddr string, redisCluster bool) Counter {
	if redisAddr == "" {
		panic(fmt.Errorf("[err] empty redis addr"))
	}

	c := &counter{}
	if redisCluster {
		c.redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           []string{redisAddr},
			MaxRetries:      redisMaxRetries,
			ConnMaxIdleTime: redisMaxConnIdleTimeout,
			ConnMaxLifetime: redisMaxConnAge,
			PoolSize:        runtime.GOMAXPROCS(0) * redisPoolSizeRatio,
		})
	} else {
		c.redisClient = redis.NewClient(&redis.Options{
			Addr:            redisAddr,
			MaxRetries:      redisMaxRetries,
			ConnMaxIdleTime: redisMaxConnIdleTimeout,
			ConnMaxLifetime: redisMaxConnAge,
			PoolSize:        runtime.GOMAXPROCS(0) * redisPoolSizeRatio,
		})
	}

	return c
}
