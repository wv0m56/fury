package engine

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tylertreat/BoomFilters"
	"github.com/wv0m56/fury/datastructure/duplist"
	"github.com/wv0m56/fury/datastructure/linkedlist"
)

// internals
func TestAccessStats(t *testing.T) {

	as := &accessStats{
		sync.Mutex{},
		boom.NewCountMinSketch(0.001, 0.99),
		&linkedlist.TimeString{},
		duplist.NewUint64String(24),
		map[string]relevantTuple{},
		20 * time.Millisecond,
		duplist.NewUint64String(24),
		map[string]*duplist.Uint64StringElement{},
	}

	as.addToWindow("foo")
	as.addToWindow("bar")
	as.addToWindow("bar")
	as.addToWindow("bar")
	as.addToWindow("baz")

	as.Lock()

	relevantDLKeys := ""
	for it := as.relevantDuplist.First(); it != nil; it = it.Next() {
		relevantDLKeys += it.Val()
	}
	assert.Equal(t, "bazfoobar", relevantDLKeys)

	relevantLLKeys := ""
	for it := as.relevantLL.Front(); it != nil; it = it.Next() {
		relevantLLKeys += it.Key()
	}
	assert.Equal(t, "foobarbaz", relevantLLKeys)

	assert.Equal(t, 3, len(as.relevantMap))
	assert.Equal(t, 0, len(as.irrelevantMap))
	for _, v := range as.relevantMap {
		if key := v.dlPtr.Val(); !(key == "foo" || key == "bar" || key == "baz") {
			t.Error("map wrong")
		}
	}

	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(3), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))

	as.Unlock()

	go as.startLoop(time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	as.Lock()

	assert.Equal(t, 0, len(as.relevantMap))
	assert.Equal(t, 3, len(as.irrelevantMap))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(3), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))

	as.Unlock()

	as.addToWindow("a")
	as.addToWindow("a")
	as.addToWindow("b")
	as.addToWindow("bar")
	as.addToWindow("a")
	as.addToWindow("a")

	as.Lock()
	assert.Equal(t, 3, len(as.relevantMap))
	assert.Equal(t, 2, len(as.irrelevantMap))

	relevantLLKeys = ""
	for it := as.relevantLL.Front(); it != nil; it = it.Next() {
		relevantLLKeys += it.Key()
	}
	assert.Equal(t, "bbara", relevantLLKeys)

	irrelevantKeys := ""
	for it := as.irrelevantDuplist.First(); it != nil; it = it.Next() {
		irrelevantKeys += it.Val()
	}
	assert.Equal(t, "bazfoo", irrelevantKeys)
	as.Unlock()

	as.updateDataDeletion("b")
	as.updateDataDeletion("baz")

	as.Lock()
	assert.Equal(t, 2, len(as.relevantMap))
	assert.Equal(t, 1, len(as.irrelevantMap))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(4), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))
	assert.Equal(t, uint64(4), as.cms.Count([]byte("a")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("b")))
	as.Unlock()
}
