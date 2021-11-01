package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func getDailyPricing(stockSymbol string) DailyPrices {
	response, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&apikey=demo")
	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject DailyPrices
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func main() {
	js.Global().Set("getDailyPrices", getDailyPricingWrapper())
	<-make(chan bool)
}

func getDailyPricingWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			// Run this code asynchronously
			go func() {
				// Make the HTTP request
				res, err := http.DefaultClient.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&apikey=demo")
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

				// Resolve the Promise
				resolve.Invoke(response)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func wrap(encoded string, err string) map[string]interface{} {
	return map[string]interface{}{
		"error":   err,
		"encoded": encoded,
	}
}
