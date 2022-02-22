//go:build js && wasm

package main

import (
	"TMS-Backend/api"
	"TMS-Backend/parser"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

func main() {
	/*
		List all functions to be exposed to Javascript here
		Example:
			Set("Name function will have in JavaScript", Go function)
	*/
	js.Global().Set("getDailyPricesGoApi", api.GetDailyPricing())
	// Parse CSV string in by calling parseDocument(csvString, brokerType "tdameritrade, interactivebrokers or questrade")
	js.Global().Set("parseDocument", ParseDocument())
	<-make(chan bool)
}

func GetBrokerEnum(broker string) parser.Broker {
	lowercase := strings.ToLower(broker)
	switch lowercase {
	case "tdameritrade":
		return parser.TDAmeritrade
	case "interactivebrokers":
		return parser.InteraciveBrokers
	case "questrade":
		return parser.Questrade
	default:
		return parser.UndefinedBroker
	}
}

func ParseDocument() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		filepath := args[0].String()
		documentType := args[1].String()
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			go func() {
				res, err := parser.ParseDocument(filepath, GetBrokerEnum(documentType))

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
