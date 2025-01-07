
SOURCE=$(shell find . -iname "*.go")



web/src/assets/wasm/lib.wasm: $(SOURCE)
	mkdir -p dist
	rm -f dist/*
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" web/src/assets/wasm/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o ./web/src/assets/wasm/lib.wasm cmd/main/functions.go

.PHONY: tinygo web/src/assets/wasm/lib.tinygo.wasm web/src/assets/wasm/lib.wasm
tinygo:
	brew tap tinygo-org/tools
	brew install tinygo


web/src/assets/wasm/lib.tinygo.wasm: $(SOURCE)
	mkdir -p dist
	rm -f dist/*
	cp "${shell brew --prefix tinygo}/targets/wasm_exec.js" web/src/assets/wasm/wasm_exec.js
	GOOS=js GOARCH=wasm tinygo build -target=wasm -o ./web/src/assets/wasm/lib.tinygo.wasm cmd/main/functions.go
