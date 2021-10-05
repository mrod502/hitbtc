package hitbtc

import (
	"errors"

	"go.uber.org/atomic"
)

var (
	ErrNoQuotes = errors.New("no quotes")
)

type MarketDepth struct {
	lastSeqNum *atomic.Uint32
	bids       *quoteMap
	asks       *quoteMap
}

func (m MarketDepth) HighestBid() (Quote, error) {
	return m.bids.min()
}

func (m MarketDepth) LowestAsk() (Quote, error) {
	return m.bids.min()
}

func (m *MarketDepth) Values() (bids, asks []Quote) {
	return m.bids.vals(), m.asks.vals()
}

func (m *MarketDepth) Mid() (v float64, err error) {
	var v1, v2 Quote
	if v1, err = m.bids.min(); err != nil {
		return 0, err
	}
	if v2, err = m.asks.min(); err != nil {
		return 0, err
	}
	return (v1.Price + v2.Price) / 2, nil
}

func NewMarketDepth() *MarketDepth {
	return &MarketDepth{
		bids:       newQuoteMap(),
		asks:       newQuoteMap(),
		lastSeqNum: atomic.NewUint32(0),
	}
}

func (m *MarketDepth) Update(b BookPage) {
	for _, v := range b.Ask {
		m.bids.update(v.Price, v)
	}
	for _, v := range b.Bid {
		m.bids.update(-v.Price, v)
	}
}
