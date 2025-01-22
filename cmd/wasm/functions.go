package main

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/overlay"
	"syscall/js"
)

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

		return overlay.CalculateOverlay(args[0].String(), args[1].String(), args[2].String())
	}))

	js.Global().Set("ApplyOverlay", promisify(func(args []js.Value) (string, error) {
		if len(args) != 2 {
			return "", fmt.Errorf("ApplyOverlay: expected 2 args, got %v", len(args))
		}

		return overlay.ApplyOverlay(args[0].String(), args[1].String())
	}))

	js.Global().Set("GetInfo", promisify(func(args []js.Value) (string, error) {
		if len(args) != 1 {
			return "", fmt.Errorf("GetInfo: expected 1 arg, got %v", len(args))
		}

		return overlay.GetInfo(args[0].String())
	}))

	<-make(chan bool)
}
