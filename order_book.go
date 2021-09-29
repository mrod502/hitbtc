package hitbtc

import (
	gocache "github.com/mrod502/go-cache"
)

type OrderBook struct {
	symbols *gocache.InterfaceCache
	msg     chan OrderbookMessage
}

func (o *OrderBook) updater() {
	for {
		msg := <-o.msg

		for k, v := range msg.Update {
			o.symbols.Set(k, v)
		}
	}
}

func (o *OrderBook) Snapshot(m OrderbookMessage) {

}

func (o *OrderBook) Update(m OrderbookMessage) {
	o.msg <- m
}

func NewOrderBook() *OrderBook {
	var o = &OrderBook{
		symbols: gocache.NewInterfaceCache(),
		msg:     make(chan OrderbookMessage, 1024),
	}
	go o.updater()
	return o
}

type OrderbookMessage struct {
	Update   map[string]BookPage `json:"update,omitempty"`
	Snapshot map[string]BookPage `json:"snapshot,omitempty"`
}

func (o OrderbookMessage) IsSnapshot() bool {
	return len(o.Snapshot) > 0
}

type BookPage struct {
	Timestamp int64   `json:"t"`
	SeqNum    uint64  `json:"s"`
	Ask       []Quote `json:"a"`
	Bid       []Quote `json:"b"`
}
