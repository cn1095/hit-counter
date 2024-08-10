package limiter

import (
	"context"
	"errors"

	perrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	blacklistKey = "hits:blacklist"
)

type BlacklistMethod interface {
	IsBlackList(ctx context.Context, id string) (bool, error)
	AddBlackList(ctx context.Context, id string) error
}

var _ BlacklistMethod = (*limiter)(nil)

// IsBlackList checks whether the id is in the blacklist or not.
func (l *limiter) IsBlackList(ctx context.Context, id string) (bool, error) {
	if ctx == nil || id == "" {
		return false, perrors.WithStack(perrors.New("[err] empty params"))
	}

	result, err := l.redisClient.SIsMember(ctx, blacklistKey, id).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, perrors.WithStack(err)
	}

	return result, nil
}

// AddBlackList adds the id to the blacklist.
func (l *limiter) AddBlackList(ctx context.Context, id string) error {
	if ctx == nil || id == "" {
		return perrors.WithStack(perrors.New("[err] empty params"))
	}

	err := l.redisClient.SAdd(ctx, blacklistKey, id).Err()
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}
