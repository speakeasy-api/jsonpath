
SOURCE=$(shell find . -iname "*.go")



dist/lib.wasm: $(SOURCE)
	mkdir -p dist
	rm -f dist/*
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" dist/wasm_exec.js
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ./dist/lib.wasm cmd/wasm/main.go
