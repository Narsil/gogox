package mtgox

import (
	"fmt"
	"os"
)

func ExampleStreamingApi() {
	gox, err := NewFromConfig(os.ExpandEnv("$MTGOX_CONFIG"))
    gox.Start()
	if err != nil {
		panic(err)
	}

	tickers := gox.Ticker
	if err != nil {
		panic(err)
	}

	go func() {
		for ticker := range tickers {
			fmt.Println("Got ticker", ticker)
		}
	}()


	orderchan := gox.SubmitOrder("bid", 100000000, 10000) // Both are in _int notation

	order := <-orderchan
    fmt.Println("Yay submitted an order!", order)
}
