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
		map[string]*linkedlist.TimeStringElement{},
		10 * time.Millisecond,
		duplist.NewUint64String(24),
		map[string]*duplist.Uint64StringElement{},
	}

	go as.startLoop(time.Millisecond)

	as.addToWindow("foo")
	as.addToWindow("bar")
	as.addToWindow("baz")

	as.Lock()

	assert.Equal(t, 3, len(as.relevantMap))
	assert.Equal(t, 0, len(as.irrelevantMap))
	for _, v := range as.relevantMap {
		if !(v.Key() == "foo" || v.Key() == "bar" || v.Key() == "baz") {
			t.Error("map wrong")
		}
	}

	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))
	assert.True(t, as.isRelevant("foo"))

	as.Unlock()

	time.Sleep(15 * time.Millisecond)

	as.Lock()

	assert.False(t, as.isRelevant("foo"))
	assert.Equal(t, 0, len(as.relevantMap))
	assert.Equal(t, 3, len(as.irrelevantMap))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))

	as.Unlock()

	as.addToWindow("a")
	as.addToWindow("a")
	as.addToWindow("a")
	as.addToWindow("a")
	as.addToWindow("b")

	as.Lock()
	assert.Equal(t, 2, len(as.relevantMap))
	assert.Equal(t, 3, len(as.irrelevantMap))
	as.Unlock()

	as.updateDataDeletion("b")
	as.updateDataDeletion("baz")

	as.Lock()
	assert.Equal(t, 1, len(as.relevantMap))
	assert.Equal(t, 2, len(as.irrelevantMap))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("foo")))
	assert.Equal(t, uint64(1), as.cms.Count([]byte("bar")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("baz")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("zzz")))
	assert.Equal(t, uint64(4), as.cms.Count([]byte("a")))
	assert.Equal(t, uint64(0), as.cms.Count([]byte("b")))
	as.Unlock()
}
