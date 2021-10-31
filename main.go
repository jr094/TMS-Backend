package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ApiResponse struct {
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

func main() {
	response, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&apikey=demo")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject ApiResponse
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.Metadata.Symbol)
}
