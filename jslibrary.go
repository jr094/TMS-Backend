//go:build js && wasm

package main

import (
	"TMS-Backend/api"
	"syscall/js"
)

func main() {
	/*
		List all functions to be exposed to Javascript here
		Example:
			Set("Name function will have in JavaScript", Go function)
	*/
	js.Global().Set("getDailyPricesGoApi", api.GetDailyPricing())
	<-make(chan bool)
}
