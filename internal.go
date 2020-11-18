package hitbtc

import (
	"sync"
)

type orderCache struct {
	v map[uint64]Order
	m *sync.RWMutex
}

func newOrderCache() *orderCache {
	return &orderCache{v: make(map[uint64]Order), m: &sync.RWMutex{}}
}

func (o *orderCache) set(v Order) {
	o.m.Lock()
	o.v[v.ID] = v
	o.m.Unlock()
}

func (o *orderCache) get(k uint64) (out Order) {
	o.m.RLock()
	out = o.v[k]
	o.m.RUnlock()
	return
}

func (o *orderCache) del(k uint64) {
	o.m.RLock()
	delete(o.v, k)
	o.m.RUnlock()
}

type messageIDs struct {
	v map[uint64]string
	m *sync.RWMutex
}

func newOrderIDCache() *messageIDs {
	return &messageIDs{v: make(map[uint64]string), m: &sync.RWMutex{}}
}

func (o *messageIDs) set(k uint64, v string) {
	o.m.Lock()
	o.v[k] = v
	o.m.Unlock()
}

func (o *messageIDs) get(k uint64) (out string) {
	o.m.RLock()
	out = o.v[k]
	o.m.RUnlock()
	return
}

func (o *messageIDs) del(k uint64) {
	o.m.RLock()
	delete(o.v, k)
	o.m.RUnlock()
}
