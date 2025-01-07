//go:build js && wasm

package wasm

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/overlay"
	"gopkg.in/yaml.v3"
	"syscall/js"
)

func CalculateOverlay(originalYAML, targetYAML string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse source schema: %w", err)
	}
	var target yaml.Node
	err = yaml.Unmarshal([]byte(targetYAML), &target)
	if err != nil {
		return "", fmt.Errorf("failed to parse target schema: %w", err)
	}

	overlay, err := overlay.Compare("example overlay", &orig, target)
	if err != nil {
		return "", fmt.Errorf("failed to compare schemas: %w", err)
	}
	out, err := yaml.Marshal(overlay)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(out), nil
}

func ApplyOverlay(originalYAML, overlayYAML string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse original schema: %w", err)
	}

	var overlay overlay.Overlay
	err = yaml.Unmarshal([]byte(overlayYAML), &overlay)
	if err != nil {
		return "", fmt.Errorf("failed to parse overlay schema: %w", err)
	}

	err = overlay.ApplyTo(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to apply overlay: %w", err)
	}

	// Unwrap the document node if it exists and has only one content node
	if orig.Kind == yaml.DocumentNode && len(orig.Content) == 1 {
		orig = *orig.Content[0]
	}

	out, err := yaml.Marshal(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(out), nil
}

func promisify(fn func(args []js.Value) (string, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// Handler for the Promise
		handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			// Run this code asynchronously
			go func() {
				result, err := fn(args)
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				resolve.Invoke(result)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func main() {
	js.Global().Set("CalculateOverlay", promisify(func(args []js.Value) (string, error) {
		if len(args) != 2 {
			return "", fmt.Errorf("CalculateOverlay: expected 2 args, got %v", len(args))
		}

		return CalculateOverlay(args[0].String(), args[1].String())
	}))

	js.Global().Set("ApplyOverlay", promisify(func(args []js.Value) (string, error) {
		if len(args) != 2 {
			return "", fmt.Errorf("ApplyOverlay: expected 2 args, got %v", len(args))
		}

		return ApplyOverlay(args[0].String(), args[1].String())
	}))

	<-make(chan bool)
}
