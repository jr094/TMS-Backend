//go:build !(js || wasm)

package main

import (
	"fmt"
)

/*
This file will be used to test code and run as a regular GO program without compiling to .wasm
This file is ignored when compiling a .wasm binary

Run this file by running `go run dev.go`
*/
func main() {
	fmt.Println("Hello World")
}
