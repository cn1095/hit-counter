#!/bin/bash

# golang
export PATH=$(go env GOPATH)/bin:$PATH

export DEBUG=true
export PHASE=debug
export SENTRY_DSN="http://local@localhost/local"
export FORCE_HTTPS=true
export REDIS_ADDRS="localhost:6379,localhost:6380"
