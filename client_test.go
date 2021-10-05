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
	cli.AddListener("ETHUSDT", func(b *MarketDepth) error {
		bid, ask := b.Values()
		fmt.Printf("ETH:%+v\t%+v\n", bid, ask)
		return nil
	})

	cli.AddListener("BTCUSDT", func(b *MarketDepth) error {
		bid, ask := b.Values()
		fmt.Printf("BTC:%+v,%+v\n", bid, ask)
		return nil
	})
	err = cli.AddOrderBookStream("ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Second)

}
