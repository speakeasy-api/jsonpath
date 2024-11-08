//go:build js && wasm

package main

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/overlay"
	"gopkg.in/yaml.v3"
	"syscall/js"
)

func calculateOverlay() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return `{"err":""}`
		}
		originalYAML := args[0].String()
		targetYAML := args[1].String()
		var orig yaml.Node
		err := yaml.Unmarshal([]byte(originalYAML), &orig)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to parse schema: %s\"", err.Error())
		}
		var target yaml.Node
		err = yaml.Unmarshal([]byte(targetYAML), &target)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to parse schema: %s\"", err.Error())
		}

		overlay, err := overlay.Compare("example overlay", &orig, target)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to compare schemas: %s\"", err.Error())
		}
		out, err := yaml.Marshal(overlay)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to marshal schema: %s\"", err.Error())
		}

		return string(out)
	})
}
func main() {
	js.Global().Set("CalculateOverlay", calculateOverlay())
	<-make(chan bool)
}
