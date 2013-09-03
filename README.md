# mtgox

An implementation of a Mt. Gox client in Go. It uses the [websocket interface](https://en.bitcoin.it/wiki/MtGox/API/Streaming).

_Note_: This API is experimental.

## Documentation

Example below. [API Documentation on Godoc](http://godoc.org/github.com/yanatan16/golang-mtgox).

## Example

```go
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
```

## License

MIT found in LICENSE file.
