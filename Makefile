
SOURCE=$(shell find . -iname "*.go")



web/src/assets/wasm/lib.wasm: $(SOURCE)
	mkdir -p dist
	rm -f dist/*
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" web/src/assets/wasm/wasm_exec.js
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ./web/src/assets/wasm/lib.wasm cmd/wasm/main.go
