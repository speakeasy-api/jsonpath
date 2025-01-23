//go:build js && wasm

package main

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/config"
	"github.com/speakeasy-api/jsonpath/pkg/overlay"
	"gopkg.in/yaml.v3"
	"syscall/js"
)

func CalculateOverlay(originalYAML, targetYAML, existingOverlay string) (string, error) {
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

	// we go from the original to a new version, then look at the extra overlays on top
	// of that, then add that to the existing overlay
	var existingOverlayDocument overlay.Overlay
	err = yaml.Unmarshal([]byte(existingOverlay), &existingOverlayDocument)
	if err != nil {
		return "", fmt.Errorf("failed to parse overlay schema: %w", err)
	}
	// now modify the original using the existing overlay
	err = existingOverlayDocument.ApplyTo(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to apply existing overlay: %w", err)
	}

	newOverlay, err := overlay.Compare("example overlay", &orig, target)
	if err != nil {
		return "", fmt.Errorf("failed to compare schemas: %w", err)
	}
	// special case, is there only one action and it targets the same as the last overlayDocument.Actions item entry, we'll just replace it.
	if len(newOverlay.Actions) == 1 && len(existingOverlayDocument.Actions) > 0 && newOverlay.Actions[0].Target == existingOverlayDocument.Actions[len(existingOverlayDocument.Actions)-1].Target {
		existingOverlayDocument.Actions[len(existingOverlayDocument.Actions)-1] = newOverlay.Actions[0]
	} else {
		// Otherwise, we'll just append the new overlay to the existing overlay
		existingOverlayDocument.Actions = append(existingOverlayDocument.Actions, newOverlay.Actions...)
	}

	out, err := yaml.Marshal(existingOverlayDocument)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(out), nil
}

func GetInfo(originalYAML string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse source schema: %w", err)
	}

	titlePath, err := jsonpath.NewPath("$.info.title")
	if err != nil {
		return "", err
	}
	versionPath, err := jsonpath.NewPath("$.info.version")
	if err != nil {
		return "", err
	}
	descriptionPath, err := jsonpath.NewPath("$.info.version")
	if err != nil {
		return "", err
	}
	toString := func(node []*yaml.Node) string {
		if len(node) == 0 {
			return ""
		}
		return node[0].Value
	}

	return `{
    "title": "` + toString(titlePath.Query(&orig)) + `",
    "version": "` + toString(versionPath.Query(&orig)) + `",
    "description": "` + toString(descriptionPath.Query(&orig)) + `"
}`, nil
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

func Query(currentYAML, path string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(currentYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse original schema: %w", err)
	}
	parsed, err := jsonpath.NewPath(path, config.WithPropertyNameExtension())
	if err != nil {
		return "", err
	}
	result := parsed.Query(&orig)
	// Marshal it back out
	out, err := yaml.Marshal(result)
	if err != nil {
		return "", err
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
		if len(args) != 3 {
			return "", fmt.Errorf("CalculateOverlay: expected 3 args, got %v", len(args))
		}

		return CalculateOverlay(args[0].String(), args[1].String(), args[2].String())
	}))

	js.Global().Set("ApplyOverlay", promisify(func(args []js.Value) (string, error) {
		if len(args) != 2 {
			return "", fmt.Errorf("ApplyOverlay: expected 2 args, got %v", len(args))
		}

		return ApplyOverlay(args[0].String(), args[1].String())
	}))

	js.Global().Set("GetInfo", promisify(func(args []js.Value) (string, error) {
		if len(args) != 1 {
			return "", fmt.Errorf("GetInfo: expected 1 arg, got %v", len(args))
		}

		return GetInfo(args[0].String())
	}))
	js.Global().Set("Query", promisify(func(args []js.Value) (string, error) {
		if len(args) != 1 {
			return "", fmt.Errorf("Query: expected 2 args, got %v", len(args))
		}

		return Query(args[0].String(), args[1].String())
	}))

	<-make(chan bool)
}
