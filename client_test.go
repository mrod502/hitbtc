package hitbtc

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	cli, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	cli.AddOrderBookStream("ETHUSDT", "BTCUSDT")

	time.Sleep(10 * time.Second)
	fmt.Println(cli.book.symbols.Get("ETHUSDT"))
}
