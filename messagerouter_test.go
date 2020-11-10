package hitbtc

import (
	"testing"
	"time"
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
