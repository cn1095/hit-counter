package counter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	intime "github.com/gjbae1212/hit-counter/internal/time"
	perrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	hitDailyFormat = "hit:daily:%s:%s"
	hitTotalFormat = "hit:total:%s"
)

type HitMethod interface {
	IncreaseHitOfDaily(ctx context.Context, id string, t time.Time, ttl time.Duration) (*Score, error)
	IncreaseHitOfTotal(ctx context.Context, id string) (*Score, error)
	GetHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error)
	GetHitOfTotal(ctx context.Context, id string) (*Score, error)
	GetHitOfDailyAndTotal(ctx context.Context, id string, t time.Time) (daily *Score, total *Score, err error)
	GetHitOfDailyByRange(ctx context.Context, id string, timeRange []time.Time) ([]*Score, error)
}

// IncreaseHitOfDaily increases count of hit daily.
func (c *counter) IncreaseHitOfDaily(ctx context.Context, id string, t time.Time, ttl time.Duration) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily empty params")
	}

	daily := intime.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	pipe := c.redisClient.Pipeline()
	incrResult := pipe.Incr(ctx, key)
	// set expired.
	pipe.Expire(ctx, key, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, perrors.WithStack(err)
	}

	value, err := incrResult.Result()
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return &Score{Name: key, Value: value}, nil
}

// IncreaseHitOfTotal increases count of hit total.
func (c *counter) IncreaseHitOfTotal(ctx context.Context, id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal empty params")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	value, err := c.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &Score{Name: key, Value: value}, nil
}

// GetHitOfDaily gets count of hit daily.
func (c *counter) GetHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] GetHitOfDaily empty param")
	}

	daily := intime.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	v, err := c.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return &Score{Name: key, Value: value}, nil
}

// GetHitOfTotal gets count of hit total.
func (c *counter) GetHitOfTotal(ctx context.Context, id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] GetHitOfTotal empty param")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := c.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return &Score{Name: key, Value: value}, nil
}

// GetHitOfDailyAndTotal returns daily score and  accumulate score.
func (c *counter) GetHitOfDailyAndTotal(ctx context.Context, id string, t time.Time) (daily *Score, total *Score, retErr error) {
	if id == "" || t.IsZero() {
		retErr = fmt.Errorf("[err] GetHitOfDailyAndTotal empty params")
		return
	}

	dailyKey := fmt.Sprintf(hitDailyFormat, intime.TimeToDailyStringFormat(t), id)
	totalKey := fmt.Sprintf(hitTotalFormat, id)

	v, err := c.redisClient.MGet(ctx, dailyKey, totalKey).Result()
	if errors.Is(err, redis.Nil) {
		return
	}
	if err != nil {
		retErr = perrors.WithStack(err)
		return
	}

	if v[0] != nil {
		dailyValue, err := strconv.ParseInt(v[0].(string), 10, 64)
		if err != nil {
			retErr = perrors.WithStack(err)
			return
		}
		daily = &Score{Name: dailyKey, Value: dailyValue}
	}

	if v[1] != nil {
		totalValue, err := strconv.ParseInt(v[1].(string), 10, 64)
		if err != nil {
			retErr = perrors.WithStack(err)
			return
		}
		total = &Score{Name: totalKey, Value: totalValue}
	}

	return
}

// GetHitOfDailyByRange returns daily scores with range.
func (c *counter) GetHitOfDailyByRange(ctx context.Context, id string, timeRange []time.Time) ([]*Score, error) {
	if id == "" || len(timeRange) == 0 {
		return nil, fmt.Errorf("[err] GetHitOfDailyByRange empty params")
	}

	var keys []string
	for _, t := range timeRange {
		keys = append(keys, fmt.Sprintf(hitDailyFormat, intime.TimeToDailyStringFormat(t), id))
	}

	v, err := c.redisClient.MGet(ctx, keys...).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	scores := make([]*Score, 0, len(keys))
	for i, key := range keys {
		if v[i] != nil {
			dailyValue, err := strconv.ParseInt(v[i].(string), 10, 64)
			if err != nil {
				return nil, perrors.WithStack(err)
			}
			scores = append(scores, &Score{Name: key, Value: dailyValue})
		} else {
			scores = append(scores, nil)
		}
	}

	return scores, nil
}
