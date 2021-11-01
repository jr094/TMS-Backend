// simple local dev server to serve wasm file, I think node server should work too

package main

import (
	"net/http"
)

func main() {
	panic(http.ListenAndServe(":9090", http.FileServer(http.Dir("."))))
}
