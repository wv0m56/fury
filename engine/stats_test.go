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
		10 * time.Millisecond,
		duplist.NewUint64String(24),
		map[string]*duplist.Uint64StringElement{},
	}

	go as.startLoop(time.Millisecond)

	as.addToWindow("foo")
	as.addToWindow("bar")
	as.addToWindow("baz")
	as.addToWindow("bar")

	as.Lock()

	assert.Equal(t, 3, len(as.relevantMap))
	assert.Equal(t, 0, len(as.irrelevantMap))
	for _, v := range as.relevantMap {
		if key := v.dlPtr.Val(); !(key == "foo" || key == "bar" || key == "baz") {
			t.Error("map wrong")
		}
	}

	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(2), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))
	assert.True(t, as.isRelevant("foo"))
	assert.True(t, as.isRelevant("bar"))
	assert.True(t, as.isRelevant("baz"))
	assert.False(t, as.isRelevant("zzz"))

	as.Unlock()

	time.Sleep(15 * time.Millisecond)

	as.Lock()

	assert.False(t, as.isRelevant("foo"))
	assert.Equal(t, 0, len(as.relevantMap))
	assert.Equal(t, 3, len(as.irrelevantMap))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(2), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))

	irrelevantKeys := ""
	for it := as.irrelevantDuplist.First(); it != nil; it = it.Next() {
		irrelevantKeys += it.Val()
	}
	assert.Equal(t, "bazfoobar", irrelevantKeys)

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

	relevantKeys := ""
	for it := as.relevantLL.Front(); it != nil; it = it.Next() {
		relevantKeys += it.Key()
	}
	assert.Equal(t, "bbara", relevantKeys)

	irrelevantKeys = ""
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
	assert.Equal(t, uint64(3), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))
	assert.Equal(t, uint64(4), as.cms.Count([]byte("a")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("b")))
	as.Unlock()
}
