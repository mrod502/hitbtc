package hitbtc

import (
	"errors"
	"fmt"

	"go.uber.org/atomic"
)

var (
	ErrNoQuotes = errors.New("no quotes")
)

type MarketDepth struct {
	lastSeqNum *atomic.Uint64
	bids       *quoteMap
	asks       *quoteMap
	lastChange BookPage
}

func (m MarketDepth) LastChange() BookPage {
	return m.lastChange
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
		lastSeqNum: atomic.NewUint64(0),
	}
}

func (m *MarketDepth) Update(b BookPage) error {
	if b.SeqNum < m.lastSeqNum.Load() {
		return fmt.Errorf("seqNum of incoming update (%d) is lower than last (%d)", b.SeqNum, m.lastSeqNum.Load())
	}
	m.lastChange = b
	for _, v := range b.Ask {
		m.bids.update(v.Price, v)
	}
	for _, v := range b.Bid {
		m.bids.update(-v.Price, v)
	}
	return nil
}
