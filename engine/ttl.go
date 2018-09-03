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
				tc.DelFirst()
				delete(tc.m, f.Val())
				if _, ok := tc.e.data[f.Val()]; ok {
					delete(tc.e.data, f.Val())
					go tc.e.stats.updateDataDeletion(f.Val())
				}
			}
			tc.e.rwm.Unlock()
		}
	}
}

func (e *Engine) setExpiry(key string, expiry time.Time) {

	if de, ok := e.ttl.m[key]; ok {
		e.ttl.DelElement(de)
	}
	insertedTTL := e.ttl.Insert(expiry, key)
	e.ttl.m[key] = insertedTTL
}
