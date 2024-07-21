//go:build docker

package counter

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCounter_IncreaseRankOfDaily(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	now := time.Now()
	tests := []struct {
		name   string
		ctx    context.Context
		group  string
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
			name:  "success",
			ctx:   context.Background(),
			group: "github.com",
			id:    "test",
			tt:    now,
			ttl:   time.Minute,
			output: &Score{
				Name:  "test",
				Value: 1,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := mockCounter.IncreaseRankOfDaily(tc.ctx, tc.group, tc.id, tc.tt, tc.ttl)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_IncreaseRankOfTotal(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	tests := []struct {
		name   string
		ctx    context.Context
		group  string
		id     string
		output *Score
		isErr  bool
	}{
		{
			name:  "error",
			isErr: true,
		},
		{
			name:  "success",
			ctx:   context.Background(),
			group: "github.com",
			id:    "test",
			output: &Score{
				Name:  "test",
				Value: 1,
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := mockCounter.IncreaseRankOfTotal(tc.ctx, tc.group, tc.id)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_GetRankDailyByLimit(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	now := time.Now()
	tests := []struct {
		name   string
		expect func(t *testing.T)
		ctx    context.Context
		group  string
		tt     time.Time
		limit  int
		output []*Score
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
				for i := range 10 {
					for range i {
						_, _ = mockCounter.IncreaseRankOfDaily(context.Background(), "github.com",
							"test"+strconv.Itoa(i), now, time.Minute)
					}
				}
			},
			ctx:   context.Background(),
			group: "github.com",
			tt:    now,
			limit: 3,
			output: []*Score{
				{
					Name:  "test9",
					Value: 9,
				},
				{
					Name:  "test8",
					Value: 8,
				},
				{
					Name:  "test7",
					Value: 7,
				},
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			output, err := mockCounter.GetRankDailyByLimit(tc.ctx, tc.group, tc.tt, tc.limit)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				spew.Dump(output)
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}

func TestCounter_GetRankTotalByLimit(t *testing.T) {
	defer func() {
		_, _ = mockCounter.(*counter).redisClient.FlushAll(context.Background()).Result()
	}()

	tests := []struct {
		name   string
		expect func(t *testing.T)
		ctx    context.Context
		group  string
		limit  int
		output []*Score
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
				for i := range 10 {
					for range i {
						_, _ = mockCounter.IncreaseRankOfTotal(context.Background(), "github.com",
							"test"+strconv.Itoa(i))
					}
				}
			},
			ctx:   context.Background(),
			group: "github.com",
			limit: 3,
			output: []*Score{
				{
					Name:  "test9",
					Value: 9,
				},
				{
					Name:  "test8",
					Value: 8,
				},
				{
					Name:  "test7",
					Value: 7,
				},
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t)
			output, err := mockCounter.GetRankTotalByLimit(tc.ctx, tc.group, tc.limit)
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				spew.Dump(output)
				assert.True(t, cmp.Equal(tc.output, output))
			}
		})
	}
}
