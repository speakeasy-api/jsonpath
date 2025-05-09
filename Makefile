
SOURCE=$(shell find . -iname "*.go")


web/src/assets/wasm/lib.wasm: $(SOURCE)
	./build.sh
