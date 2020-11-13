package hitbtc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMessageRouter(t *testing.T) {

	var mr *MessageRouter
	var err error
	mr, err = NewMessageRouter()
	if err != nil {
		t.Fatal(err)
	}
	if mr == nil {
		t.Fatal("mr is nil")
	}
	mr.AddRoute(WSSTicker, echoTicker)
	mr.SubscribeTicker("ETHBTC")
	time.Sleep(10 * time.Second)

}

func TestMessageStruct(t *testing.T) {
	var m Message
	m.Jsonrpc = "2.0"
	m.Method = string(MthdSubscribeTicker)
	m.Params = Order{ID: 123, Symbol: "BTCUSD", Quantity: decimal.New(0, 0)}

	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
