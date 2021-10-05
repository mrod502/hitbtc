package hitbtc

import (
	"fmt"

	gocache "github.com/mrod502/go-cache"
)

type OrderBook struct {
	symbols *gocache.InterfaceCache
	msg     chan OrderbookMessage
	subs    *gocache.InterfaceCache
}

func (o *OrderBook) Subscribe(symbol string) (c chan *MarketDepth) {
	c = make(chan *MarketDepth, 128)
	o.subs.Set(symbol, c)
	return
}

func (o *OrderBook) getSub(s string) (c chan *MarketDepth, err error) {
	v, ok := o.subs.Get(s).(chan *MarketDepth)
	if !ok {
		return nil, gocache.ErrInterfaceAssertion
	}
	return v, nil
}

func (o *OrderBook) updater() {
	for {
		msg := <-o.msg
		fmt.Println(len(msg.Update), len(msg.Snapshot))

		for k, v := range msg.Update {

			if o.symbols.Exists(k) {
				md := o.symbols.Get(k).(*MarketDepth)
				md.Update(v)
				bd, ak := md.Values()
				fmt.Printf("%+v\t%+v\n", bd, ak)
			} else {
				mdOther := NewMarketDepth()
				mdOther.Update(v)
				o.symbols.Set(k, mdOther)
				bd, ak := mdOther.Values()
				fmt.Printf("%+v\t%+v\n", bd, ak)
			}
		}

	}
}

func (o *OrderBook) Snapshot(m OrderbookMessage) {
	for k, v := range m.Snapshot {
		md := NewMarketDepth()
		md.Update(v)
		o.symbols.Set(k, md)
	}
}

func (o *OrderBook) Update(m OrderbookMessage) {
	fmt.Println("Call Update", len(o.msg))
	o.msg <- m
}

func NewOrderBook() *OrderBook {
	var o = &OrderBook{
		symbols: gocache.NewInterfaceCache(),
		msg:     make(chan OrderbookMessage, 1024),
		subs:    gocache.NewInterfaceCache(),
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
