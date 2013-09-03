package mtgox

import (
	"testing"
	"time"
	"os"
    "fmt"
)

func Init () (gox *StreamingApi){
    gox, err := NewFromConfig(os.ExpandEnv("$HOME/.mtgox"))
	if err != nil {
		panic(err)
	}
    gox.Start()
    gox.HandleErrors(func(err error){ fmt.Println(err) })
    return
}

func TestDepthReceive(t *testing.T) {
    gox := Init()
	depth := gox.Depth

	select {
    case <-depth:
        return
	case <-time.After(1000 * time.Second):
		t.Error("No depth received after ten seconds.")
	}
}
func TestTickerReceive(t *testing.T) {
    gox := Init()
	ticker := gox.Ticker

	select {
    case <-ticker:
        return
	case <-time.After(1000 * time.Second):
		t.Error("No ticker received after ten seconds.")
	}
}

func TestInfoCall(t *testing.T) {
    gox := Init()
	infoc := gox.RequestInfo()
	select {
	case _ = <-infoc:
	case <-time.After(10 * time.Second):
		t.Error("No info after ten seconds.")
	}
}
func TestOrders(t *testing.T) {
    gox := Init()
	orders := gox.RequestOrders()
	select {
	case _ = <-orders:
	case <-time.After(10 * time.Second):
		t.Error("No orders after ten seconds.")
	}
}
