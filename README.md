## To Use API 
- Import WASM glue code 
```
<script src="wasm_exec.js"></script>
```
- Download `alphaVantageApi.wasm`

- Init WebAssembly in head
```
            const go = new Go();
                WebAssembly.instantiateStreaming(fetch("alphaVantageApi.wasm"), go.importObject)
                    .then(result => {
                    go.run(result.instance);
                });
```
- Call Go function in `async` JavaScript function
```
                try {
                    const response = await getDailyPricesGoApi("API_KEY_HERE", "AAPL")
                    const message = await response.json()

                    // Logs Go api output to console
                    console.log(message)
                } catch (err) {
                    console.error('Caught exception', err)
                }
```