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

func applyOverlay() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return `{"err":""}`
		}
		originalYAML := args[0].String()
		overlayString := args[1].String()
		var orig yaml.Node
		err := yaml.Unmarshal([]byte(originalYAML), &orig)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to parse schema: %s\"", err.Error())
		}
		var target overlay.Overlay
		err = yaml.Unmarshal([]byte(overlayString), &target)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to parse schema: %s\"", err.Error())
		}

		err = target.ApplyTo(&orig)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to apply schema: %s\"", err.Error())
		}
		// Unwrap the document node
		if orig.Kind == yaml.DocumentNode && len(orig.Content) == 1 {
			orig = *orig.Content[0]
		}

		out, err := yaml.Marshal(&orig)
		if err != nil {
			return fmt.Sprintf("{\"err\": \"failed to marshal schema: %s\"", err.Error())
		}

		return string(out)
	})
}
func main() {
	js.Global().Set("CalculateOverlay", calculateOverlay())
	js.Global().Set("ApplyOverlay", applyOverlay())
	<-make(chan bool)
}
