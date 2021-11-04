document.onload = function(e) {
    const go = new Go();
            WebAssembly.instantiateStreaming(fetch("alphaVantageApi.wasm"), go.importObject)
                .then(result => {
                go.run(result.instance);
            });
}

async function getDailyPrice() {
    try {
        const response = await getDailyPricesGoApi("Z7CGZVE7BGU3U777", "AAPL")
        const message = await response.json()

        // Logs Go api output to console
        console.log(message)
    } catch (err) {
        console.error('Caught exception', err)
    }
}