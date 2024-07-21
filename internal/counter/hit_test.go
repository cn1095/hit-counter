//go:build docker

package counter

import (
	"context"
	"fmt"
	"testing"
	"time"

	intime "github.com/gjbae1212/hit-counter/internal/time"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCounter_IncreaseHitOfDaily(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()
	now := time.Now()
	daily := intime.TimeToDailyStringFormat(now)

	tests := []struct {
		name   string
		ctx    context.Context
		id     string
		tt     time.Time
		ttl    time.Duration
		output *Score
		isErr  bool
	}{
		{
			name:  "error",
			isErr: true,
		},
		{
			name: "success",
			ctx:  context.Background(),
			id:   "test",
			tt:   now,
			ttl:  time.Minute,
			output: &Score{
				Name:  fmt.Sprintf(hitDailyFormat, daily, "test"),
				Value: 1,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := mockCounter.IncreaseHitOfDaily(tc.ctx, tc.id, tc.tt, tc.ttl)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_IncreaseHitOfTotal(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	tests := []struct {
		name   string
		ctx    context.Context
		id     string
		output *Score
		isErr  bool
	}{
		{
			name:  "error",
			isErr: true,
		},
		{
			name: "success",
			ctx:  context.Background(),
			id:   "test",
			output: &Score{
				Name:  fmt.Sprintf(hitTotalFormat, "test"),
				Value: 1,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := mockCounter.IncreaseHitOfTotal(tc.ctx, tc.id)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_GetHitOfDaily(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	now := time.Now()
	daily := intime.TimeToDailyStringFormat(now)

	tests := []struct {
		name   string
		expect func(t *testing.T)
		ctx    context.Context
		id     string
		tt     time.Time
		output *Score
		isErr  bool
	}{
		{
			name:   "error",
			expect: func(t *testing.T) {},
			isErr:  true,
		},
		{
			name: "success",
			expect: func(t *testing.T) {
				for i := 0; i < 10; i++ {
					_, _ = mockCounter.IncreaseHitOfDaily(context.Background(), "test", now, time.Minute)
				}
			},
			ctx: context.Background(),
			id:  "test",
			tt:  now,
			output: &Score{
				Name:  fmt.Sprintf(hitDailyFormat, daily, "test"),
				Value: 10,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			output, err := mockCounter.GetHitOfDaily(tc.ctx, tc.id, tc.tt)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_GetHitOfTotal(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	tests := []struct {
		name   string
		expect func(t *testing.T)
		ctx    context.Context
		id     string
		output *Score
		isErr  bool
	}{
		{
			name:   "error",
			expect: func(t *testing.T) {},
			isErr:  true,
		},
		{
			name: "success",
			expect: func(t *testing.T) {
				for i := 0; i < 10; i++ {
					_, _ = mockCounter.IncreaseHitOfTotal(context.Background(), "test")
				}
			},
			ctx: context.Background(),
			id:  "test",
			output: &Score{
				Name:  fmt.Sprintf(hitTotalFormat, "test"),
				Value: 10,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			output, err := mockCounter.GetHitOfTotal(tc.ctx, tc.id)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_GetHitOfDailyAndTotal(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	now := time.Now()
	daily := intime.TimeToDailyStringFormat(now)

	tests := []struct {
		name        string
		expect      func(t *testing.T)
		ctx         context.Context
		id          string
		tt          time.Time
		outputDaily *Score
		outputTotal *Score
		isErr       bool
	}{
		{
			name:   "error",
			expect: func(t *testing.T) {},
			isErr:  true,
		},
		{
			name: "success",
			expect: func(t *testing.T) {
				for i := 0; i < 10; i++ {
					_, _ = mockCounter.IncreaseHitOfTotal(context.Background(), "test")
					_, _ = mockCounter.IncreaseHitOfDaily(context.Background(), "test", now, time.Minute)
				}
			},
			ctx: context.Background(),
			id:  "test",
			tt:  now,
			outputDaily: &Score{
				Name:  fmt.Sprintf(hitDailyFormat, daily, "test"),
				Value: 10,
			},
			outputTotal: &Score{
				Name:  fmt.Sprintf(hitTotalFormat, "test"),
				Value: 10,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			outputDaily, outputTotal, err := mockCounter.GetHitOfDailyAndTotal(tc.ctx, tc.id, tc.tt)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.outputDaily, outputDaily))
				assert.True(t, cmp.Equal(tc.outputTotal, outputTotal))
			}
		})
	}
}

func TestCounter_GetHitOfDailyByRange(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	now := time.Now()
	daily := intime.TimeToDailyStringFormat(now)

	tests := []struct {
		name      string
		expect    func(t *testing.T)
		ctx       context.Context
		id        string
		timeRange []time.Time
		output    []*Score
		isErr     bool
	}{
		{
			name:   "error",
			expect: func(t *testing.T) {},
			isErr:  true,
		},
		{
			name: "success",
			expect: func(t *testing.T) {
				for i := 0; i < 10; i++ {
					_, _ = mockCounter.IncreaseHitOfDaily(context.Background(), "test", now, time.Minute)
				}
			},
			ctx:       context.Background(),
			id:        "test",
			timeRange: []time.Time{now},
			output: []*Score{
				{
					Name:  fmt.Sprintf(hitDailyFormat, daily, "test"),
					Value: 10,
				},
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			output, err := mockCounter.GetHitOfDailyByRange(tc.ctx, tc.id, tc.timeRange)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}
