package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall/js"
)

var c chan bool

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
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject DailyPrices
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func test(this js.Value, inputs []js.Value) interface{} {
	symbol := inputs[0].String()
	var priceData = getDailyPricing(symbol)

	return js.ValueOf(priceData.Metadata.Symbol)
}

func main() {
	window := js.Global()
	doc := window.Get("document")
	body := doc.Get("body")
	div := doc.Call("createElement", "div")
	div.Set("textContent", "hello!!")
	body.Call("appendChild", div)
	body.Set("onclick",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			div := doc.Call("createElement", "div")
			div.Set("textContent", "click!!")
			body.Call("appendChild", div)
			return nil
		}))
	<-make(chan struct{})
}
