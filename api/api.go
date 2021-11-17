//go:build js && wasm

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// TODO: Separated JS code from GO code so we can run and test this file without compiling to .wasm
	"syscall/js"
)

type DailyPrices struct {
	Metadata PriceMetadata    `json:"Meta Data"`
	Prices   map[string]Price `json:"Time Series (Daily)"`
}

type PriceMetadata struct {
	Symbol      string `json:"2. Symbol"`
	LastUpdated string `json:"3. Last Refreshed"`
}

type Price struct {
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

func GetDailyPricing() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apiKey := args[0]
		stockSymbol := args[1]
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			go func() {
				res, err := http.DefaultClient.Get(fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", stockSymbol, apiKey))
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}
				defer res.Body.Close()

				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				var responseObject DailyPrices
				json.Unmarshal(data, &responseObject)

				b, err := json.Marshal(responseObject)
				if err != nil {
					fmt.Println(err)
					return
				}

				arrayConstructor := js.Global().Get("Uint8Array")
				dataJS := arrayConstructor.New(len(b))
				js.CopyBytesToJS(dataJS, b)

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
