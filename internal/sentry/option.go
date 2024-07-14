package sentry

import "github.com/getsentry/sentry-go"

type Option interface {
	apply(scope *sentry.Scope)
}

type option func(scope *sentry.Scope)

func (opt option) apply(scope *sentry.Scope) { opt(scope) }

var _ Option = (option)(nil)

// WithUser returns an Option interface setting user to scope.
func WithUser(user User) Option {
	return option(func(scope *sentry.Scope) {
		scope.SetUser(user)
	})
}

// WithTags returns an Option interface setting list of tag to scope.
func WithTags(tags map[string]string) Option {
	return option(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

// WithExtras returns an Option interface setting list of extra to scope.
func WithExtras(extras map[string]any) Option {
	return option(func(scope *sentry.Scope) {
		scope.SetExtras(extras)
	})
}

// WithContexts returns an Option interface setting list of context to scope.
func WithContexts(ctxs map[string]Context) Option {
	return option(func(scope *sentry.Scope) {
		scope.SetContexts(ctxs)
	})
}

// WithFingerprint returns an Option interface setting fingerprint to scope.
func WithFingerprint(fingerprint []string) Option {
	return option(func(scope *sentry.Scope) {
		scope.SetFingerprint(fingerprint)
	})
}
