package limiter

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (opt option) apply(cfg *config) { opt(cfg) }

var _ Option = (option)(nil)

type config struct {
	rateLimitWindow int64
	rateLimitCount  int64
}

// WithRateLimitWindow returns an Option interface setting rate limit window to config.
// 1 window is equal to 1 minute.
func WithRateLimitWindow(window int64) Option {
	return option(func(cfg *config) {
		cfg.rateLimitWindow = window
	})
}

// WithRateLimitCount returns an Option interface setting rate limit count to config.
func WithRateLimitCount(count int64) Option {
	return option(func(cfg *config) {
		cfg.rateLimitCount = count
	})
}
