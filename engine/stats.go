package engine

import (
	"sync"
	"time"

	boom "github.com/tylertreat/BoomFilters"
	"github.com/wv0m56/fury/datastructure/duplist"
	"github.com/wv0m56/fury/datastructure/linkedlist"
)

// accessStats approximates the access statistics of all keys not yet evicted
// (even this is approximate, i.e. eventually consistent with the cache's state).
type accessStats struct {
	sync.Mutex
	cms               *boom.CountMinSketch
	relevantLL        *linkedlist.TimeString
	relevantMap       map[string]*linkedlist.TimeStringElement
	relevanceWindow   time.Duration
	irrelevantDuplist *duplist.Uint64String
	irrelevantMap     map[string]*duplist.Uint64StringElement
}

func (as *accessStats) isRelevant(key string) bool {
	_, ok := as.relevantMap[key]
	return ok
}

// lock because called from goroutine by engine
func (as *accessStats) addToWindow(key string) {
	as.Lock()
	defer as.Unlock()

	as.cms.Add([]byte(key))

	if ptr, ok := as.relevantMap[key]; ok {
		as.relevantLL.Del(ptr)
	}

	ptr := as.relevantLL.AddToBack(key)
	as.relevantMap[key] = ptr

	as.delIrrelevant(key)
}

// lock because called from goroutine by engine
func (as *accessStats) updateDataDeletion(key string) {
	as.Lock()
	defer as.Unlock()

	_ = as.cms.TestAndRemoveAll([]byte(key))
	as.delRelevant(key)
	as.delIrrelevant(key)
}

func (as *accessStats) delRelevant(key string) {
	if ptr, ok := as.relevantMap[key]; ok {
		as.relevantLL.Del(ptr)
		delete(as.relevantMap, key)
	}
}

func (as *accessStats) delIrrelevant(key string) {
	if ptr, ok := as.irrelevantMap[key]; ok {
		as.irrelevantDuplist.DelElement(ptr)
		delete(as.irrelevantMap, key)
	}
}

func (as *accessStats) startLoop(step time.Duration) {
	for range time.Tick(step) {
		as.Lock()
		for it := as.relevantLL.Front(); it != nil &&
			it.LastAccessed().Add(as.relevanceWindow).Before(time.Now()); it = it.Next() {

			// relevant -> irrelevant
			as.delRelevant(it.Key())
			as.upsertIrrelevant(it.Key())
		}
		as.Unlock()
	}
}

func (as *accessStats) upsertIrrelevant(key string) {
	as.delIrrelevant(key)
	as.irrelevantMap[key] = as.irrelevantDuplist.Insert(as.cms.Count([]byte(key)), key)
}
