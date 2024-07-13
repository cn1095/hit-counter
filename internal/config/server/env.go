package server

import (
	"sync"

	env10 "github.com/caarlos0/env/v10"
	internalstring "github.com/gjbae1212/hit-counter/internal/string"
	perrors "github.com/pkg/errors"
)

var (
	onceLocalEnvs = &sync.Once{}
	localEnvs     *LocalEnvironments
)

type LocalEnvironments struct {
	Debug      bool                        `env:"DEBUG,notEmpty"`
	Phase      internalstring.TrimString   `env:"PHASE,notEmpty"`
	SentryDSN  internalstring.TrimString   `env:"SENTRY_DSN,notEmpty"`
	ForceHttps bool                        `env:"FORCE_HTTPS,notEmpty"`
	RedisAddr  []internalstring.TrimString `env:"REDIS_ADDRS,notEmpty" envSeparator:","`
}

// MustInitializeLocalEnvironments must initialize a local environments obj.
func MustInitializeLocalEnvironments() {
	onceLocalEnvs.Do(func() {
		localEnvs = &LocalEnvironments{}
		if err := env10.Parse(localEnvs); err != nil {
			panic(perrors.WithStack(err))
		}
	})
}

// MustGetLocalEnvironments must get a local environments obj with initialized.
func MustGetLocalEnvironments() *LocalEnvironments {
	if localEnvs == nil {
		MustInitializeLocalEnvironments()
	}
	return localEnvs
}
