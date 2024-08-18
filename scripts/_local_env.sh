#!/bin/bash

# golang
export PATH=$(go env GOPATH)/bin:$PATH

export DEBUG=true
export PHASE=debug
export SENTRY_DSN="http://local@localhost/local"
export REDIS_ADDR="localhost:6379"
export REDIS_CLUSTER=false
