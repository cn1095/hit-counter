package counter

import (
	"context"
	"errors"
	"fmt"
	"time"

	intime "github.com/gjbae1212/hit-counter/internal/time"
	perrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	rankDailyFormat = "rank:daily:%s:%s"
	rankTotalFormat = "rank:total:%s"
)

type RankMethod interface {
	IncreaseRankOfDaily(ctx context.Context, group, id string, t time.Time, ttl time.Duration) (*Score, error)
	IncreaseRankOfTotal(ctx context.Context, group, id string) (*Score, error)
	GetRankDailyByLimit(ctx context.Context, group string, t time.Time, limit int) ([]*Score, error)
	GetRankTotalByLimit(ctx context.Context, group string, limit int) ([]*Score, error)
}

// IncreaseRankOfDaily increases count of rank daily.
func (c *counter) IncreaseRankOfDaily(ctx context.Context, group, id string, t time.Time, ttl time.Duration) (*Score, error) {
	if group == "" || id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily empty params")
	}

	daily := intime.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	pipe := c.redisClient.Pipeline()
	incrResult := pipe.ZIncrBy(ctx, key, 1, id)
	// set expired.
	pipe.Expire(ctx, key, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, perrors.WithStack(err)
	}
	incr, err := incrResult.Result()
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &Score{Name: id, Value: int64(incr)}, nil
}

// IncreaseRankOfTotal increases count of rank total.
func (c *counter) IncreaseRankOfTotal(ctx context.Context, group, id string) (*Score, error) {
	if group == "" || id == "" {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal empty params")
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	value, err := c.redisClient.ZIncrBy(ctx, key, 1, id).Result()
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &Score{Name: id, Value: int64(value)}, nil
}

// GetRankDailyByLimit gets count of rank daily by limit.
func (c *counter) GetRankDailyByLimit(ctx context.Context, group string, t time.Time, limit int) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit empty params")
	}

	daily := intime.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	scores := make([]*Score, 0, limit)
	values, err := c.redisClient.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if errors.Is(err, redis.Nil) {
		return scores, nil
	}
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	for _, value := range values {
		name := value.Member.(string)
		scores = append(scores, &Score{Name: name, Value: int64(value.Score)})
	}
	return scores, nil
}

// GetRankTotalByLimit gets count of rank total by limit.
func (c *counter) GetRankTotalByLimit(ctx context.Context, group string, limit int) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit empty params")
	}

	key := fmt.Sprintf(rankTotalFormat, group)

	scores := make([]*Score, 0, limit)
	values, err := c.redisClient.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if errors.Is(err, redis.Nil) {
		return scores, nil
	}
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	for _, value := range values {
		name := value.Member.(string)
		scores = append(scores, &Score{Name: name, Value: int64(value.Score)})
	}

	return scores, nil
}
