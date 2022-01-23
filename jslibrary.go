//go:build js && wasm

package main

import (
	"TMS-Backend/api"
	"TMS-Backend/parser"
	"encoding/json"
	"fmt"
	"syscall/js"
)

func main() {
	/*
		List all functions to be exposed to Javascript here
		Example:
			Set("Name function will have in JavaScript", Go function)
	*/
	js.Global().Set("getDailyPricesGoApi", api.GetDailyPricing())
	js.Global().Set("parseDocument", ParseDocument())
	<-make(chan bool)
}

func ParseDocument() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		filepath := args[0].String()
		documentType := args[1].String()
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			go func() {
				res, err := parser.ParseDocument(filepath, documentType)

				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				body, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
					return
				}

				arrayConstructor := js.Global().Get("Uint8Array")
				dataJS := arrayConstructor.New(len(body))
				js.CopyBytesToJS(dataJS, body)

				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(dataJS)
				resolve.Invoke(response)
			}()
			return nil
		})
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
