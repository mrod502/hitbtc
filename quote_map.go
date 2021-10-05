package hitbtc

import "sync"

type quoteMap struct {
	v map[float64]Quote
	m *sync.RWMutex
}

func newQuoteMap() *quoteMap {
	return &quoteMap{
		v: make(map[float64]Quote),
		m: new(sync.RWMutex),
	}
}

func (q *quoteMap) min() (Quote, error) {
	for _, v := range q.v {
		return v, nil
	}
	return Quote{}, ErrNoQuotes
}

func (q *quoteMap) set(k float64, v Quote) {
	q.m.Lock()
	defer q.m.Unlock()
	if v.Volume == 0 {
		delete(q.v, v.Price)
		return
	}
	q.v[v.Price] = v
}

func (q *quoteMap) get(p float64) (v Quote, ok bool) {
	q.m.RLock()
	defer q.m.RUnlock()
	v, ok = q.v[p]
	return v, ok
}

func (q *quoteMap) vals() []Quote {
	var v = make([]Quote, 0, len(q.v))
	for _, val := range q.v {
		v = append(v, val)
	}
	return v
}
