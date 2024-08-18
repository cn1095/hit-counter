set -e -E -u -o pipefail
trap '_cleanup "$?" "${LINENO}" "${FUNCNAME:-unknown}" ${CMD}' EXIT

cd `dirname $0` && cd ..
CURRENT=`pwd`

source $CURRENT/scripts/_color.sh
source $CURRENT/scripts/_local_env.sh

function test
{
  go test -v $(go list ./... | grep -v vendor)  --count 1 -race -covermode=atomic -timeout 300s
}

function test_with_docker
{
  run_containers # cleanup
  go test -v $(go list ./... | grep -v vendor) --tags=docker --count 1 -race -covermode=atomic -timeout 300s
}

function lint
{
  go list -f '{{.Dir}}' -m  | xargs -I{} golangci-lint run -v {}/...
}

function mock
{
  mockery
}

function tools
{
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1
  go install github.com/vektra/mockery/v2@v2.43.2
}

function run_containers {
  docker-compose -f "${CURRENT}/scripts/docker-compose.yaml" up -d
}

function stop_containers {
  docker-compose -f "${CURRENT}/scripts/docker-compose.yaml" down
}

function _cleanup
{
  local ret=$1
  local no=$2
  local func=$3
  local cmd=$4

  if [[ $cmd == "test_with_docker" ]]; then
      stop_containers
  fi

  if [[ $ret == "0" ]]; then
      echo -e "${Green}[SUCCESS][$cmd]${Color_Off}"
  else
      echo -e "${Red}[ERROR][$ret][$cmd] func: $func, line:$no${Color_Off}"
  fi
}

CMD=$1
shift
$CMD $*
