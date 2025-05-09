#!/bin/bash
readonly SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
readonly REPO_ROOT="$(git rev-parse --show-toplevel)"

# ANSI color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Validate Go version higher or equal to than this.
REQUIRED_GO_VERSION="1.24.1"
CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')

if ! echo "$CURRENT_GO_VERSION $REQUIRED_GO_VERSION" | awk '{
    split($1, current, ".")
    split($2, required, ".")
    if ((current[1] > required[1]) || \
        (current[1] == required[1] && current[2] > required[2]) || \
        (current[1] == required[1] && current[2] == required[2] && current[3] >= required[3])) {
        exit 0;
    }
    exit 1
}'; then
    echo -e "${RED}Error: Go version must be $REQUIRED_GO_VERSION or higher, but found $CURRENT_GO_VERSION${NC}"
    exit 1
fi

get_file_size() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        BYTES=$(stat -f%z "$1")
    else
        BYTES=$(stat -c%s "$1")
    fi
    MB=$((BYTES / 1024 / 1024))
    echo "$MB MB"
}

# Function to build the WASM binary
build_wasm() {
    printf "Installing dependencies... "
    (cd "$SCRIPT_DIR" && go mod tidy) > /dev/null 2>&1
    printf "Done\n"

    printf "Building WASM binary..."
	  GOOS=js GOARCH=wasm go build -o ./web/src/assets/wasm/lib.wasm cmd/wasm/functions.go
    SIZE="$(get_file_size "./web/src/assets/wasm/lib.wasm")"
    printf " Done (%s)\n" "$SIZE"
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "./web/src/assets/wasm/wasm_exec.js"

}

# Initial build
build_wasm "$VERSION"
