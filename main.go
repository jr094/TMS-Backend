package main

import (
	"TMS-Backend/api"
	"syscall/js"
)

func main() {
	js.Global().Set("getDailyPricesGoApi", api.GetDailyPricing())
	<-make(chan bool)
}
