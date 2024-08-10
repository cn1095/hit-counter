package limiter

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-module/carbon/v2"
	_ "github.com/golang-module/carbon/v2"
	"github.com/maypok86/otter"
	_ "github.com/maypok86/otter"
	perrors "github.com/pkg/errors"
)

type RateMethod interface {
	Allow(ctx context.Context, id string) bool
}

var _ RateMethod = (*limiter)(nil)

// Allow checks whether the id is allowed or not.
func (l *limiter) Allow(ctx context.Context, id string) bool {
	if ctx == nil || id == "" {
		return false
	}
	return l.windowCounter.allow(id)
}

type fixedWindowCounter struct {
	sync.RWMutex

	limitWindow int64
	limitCount  int64
	cache       otter.Cache[int64, *fixedWindow]
	createdAt   carbon.Carbon
}

func (fwc *fixedWindowCounter) allow(id string) bool {
	key := fwc.cacheKey()

	fwc.RLock()
	window, ok := fwc.cache.Get(key)
	fwc.RUnlock()
	if !ok {
		fwc.Lock()
		if window, ok = fwc.cache.Get(key); !ok {
			window = &fixedWindow{syncMap: &sync.Map{}}
			fwc.cache.Set(key, window)
		}
		fwc.Unlock()
	}

	if window.take(id) <= fwc.limitCount {
		return true
	}
	return false
}

func (fwc *fixedWindowCounter) cacheKey() int64 {
	diff := fwc.createdAt.DiffAbsInMinutes(carbon.Now(carbon.UTC))
	return diff / fwc.limitWindow
}

type fixedWindow struct {
	syncMap *sync.Map
}

func (fw *fixedWindow) take(id string) int64 {
	value, _ := fw.syncMap.LoadOrStore(id, &atomic.Int64{})
	i32 := value.(*atomic.Int64)
	return i32.Add(1)
}

func newFixedWindowCounter(window, count int64) (*fixedWindowCounter, error) {
	if window == 0 || count == 0 {
		return nil, perrors.WithStack(fmt.Errorf("[err] invalid params"))
	}
	cache, err := otter.MustBuilder[int64, *fixedWindow](int(10 * window)).
		Cost(func(key int64, value *fixedWindow) uint32 {
			return 1
		}).
		WithTTL(time.Minute*time.Duration(window) + time.Minute).
		Build()
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return &fixedWindowCounter{
		limitWindow: window,
		limitCount:  count,
		cache:       cache,
		createdAt:   carbon.Now(carbon.UTC),
	}, nil
}
