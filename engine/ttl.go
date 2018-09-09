package engine

import (
	"time"

	"github.com/wv0m56/fury/datastructure/duplist"
)

type ttlControl struct {
	duplist.TimeString
	m map[string]*duplist.TimeStringElement
	e *Engine
}

// to be invoked as a goroutine e.g. go startLoop()
func (tc *ttlControl) startLoop(step time.Duration) {

	for range time.Tick(step) {

		var somethingExpired bool
		now := time.Now()

		tc.e.rwm.RLock()
		if f := tc.First(); f != nil && now.After(f.Key()) {
			somethingExpired = true
		}
		tc.e.rwm.RUnlock()

		if somethingExpired {
			tc.e.rwm.Lock()
			for f := tc.First(); f != nil && now.After(f.Key()); f = tc.First() {
				tc.e.delDataTTLStats(f.Val())
			}
			tc.e.rwm.Unlock()
		}
	}
}

func (tc *ttlControl) delTTLEntry(key string) {
	if el, ok := tc.m[key]; ok {
		tc.DelElement(el)
		delete(tc.m, key)
	}
}

func (e *Engine) setExpiry(key string, expiry time.Time) {

	if de, ok := e.ttl.m[key]; ok {
		e.ttl.DelElement(de)
	}
	insertedTTL := e.ttl.Insert(expiry, key)
	e.ttl.m[key] = insertedTTL
}

// GetTTL returns the number of seconds left until expiry for the given keys, in
// the order in which keys are passed into args.
// Keys without TTL yields negative values.
func (e *Engine) GetTTL(keys ...string) []float64 {

	var t []float64
	now := time.Now()
	for _, k := range keys {
		d, ok := e.ttl.m[k]
		if ok {
			t = append(t, d.Key().Sub(now).Seconds())
		} else {
			t = append(t, -1)
		}
	}
	return t
}
