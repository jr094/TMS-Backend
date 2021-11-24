# TMS-Backend

To compile WASM binary run:
```
GOOS=js GOARCH=wasm go build -o lib.wasm
```

To run packages as regular Go programs:

- Use `dev.go` as main file and call libraries there
- Run the following to run and test code

```
go run dev.go
```

Add the following to the header of any file that uses the `syscall/js` library to restrict compilation of file only when compiling to `.wasm`
```
//go:build js && wasm
```

Try to separate core logic from `syscall/js` logic so programs can be tested and run locally without compiling to `.wasm` each time

## Environment variables

windows

```bash
set GOOS=windows
set GOARCH=amd64

set GOOS=js
set GOARCH=wasm
```
