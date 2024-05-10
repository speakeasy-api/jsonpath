//go:build js && wasm

package pkg

import (
	"github.com/speakeasy-api/jsonpath/pkg"
	"syscall/js"
)

func jsFuncWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return `{"err":""}`
		}
		originalYAML := args[0].String()
		jsonpath := args[1].String()
		returnValue := pkg.Entry(originalYAML, jsonpath)
		return returnValue
	})
}

func main() {
	js.Global().Set("JSONPath", jsFuncWrapper())
	<-make(chan bool)
}
