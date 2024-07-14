set -e -E -u -o pipefail
trap '_cleanup "$?" "${LINENO}" "${FUNCNAME:-unknown}" ${CMD}' EXIT

cd `dirname $0` && cd ..
CURRENT=`pwd`

source $CURRENT/scripts/_color.sh

function make_wasm
{
  local phase=${1:-unknown}
  case $phase in
  local|production)
    echo "WASM Phase: $phase"
    rm $CURRENT/web/view/$phase/hits.wasm || true
    ;;
  *)
    echo "Invalid Phase: $phase"
    exit 1
    ;;
  esac

  # copy wasm_exec.js
  cp $(go env GOROOT)/misc/wasm/wasm_exec.js $CURRENT/web/public/

  GOOS=js GOARCH=wasm go build -ldflags="-s -w -X main.phase=$phase" -o $CURRENT/web/view/$phase/hits.wasm $CURRENT/cmd/wasm/main.go
  gzip $CURRENT/web/view/$phase/hits.wasm
  mv $CURRENT/web/view/$phase/hits.wasm.gz $CURRENT/web/view/$phase/hits.wasm
}

function _cleanup
{
  local ret=$1
  local no=$2
  local func=$3
  local cmd=$4

  if [[ $ret == "0" ]]; then
      echo -e "${Green}[SUCCESS][$cmd]${Color_Off}"
  else
      echo -e "${Red}[ERROR][$ret][$cmd] func: $func, line:$no${Color_Off}"
  fi
}


CMD=$1
shift
$CMD $*
