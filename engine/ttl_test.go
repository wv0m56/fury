package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wv0m56/fury/testdummies"
)

func TestTTL(t *testing.T) {

	opts := testOptionsDefault
	opts.O = &testdummies.NoDelayOrigin{}
	opts.TTLTickStep = 1 * time.Millisecond
	e, err := NewEngine(&opts)
	assert.Nil(t, err)

	e.Get("a")
	e.Get("b")
	e.Get("c")
	e.Get("d")
	e.Get("e")
	e.Get("f")

	setTTL := func(key string, ttl time.Duration) {
		e.setExpiry(key, time.Now().Add(ttl))
	}

	e.rwm.Lock()

	setTTL("c", 19*time.Millisecond)
	setTTL("f", 25*time.Millisecond)
	setTTL("z", 11*time.Millisecond)

	assert.Equal(t, 6, len(e.data))
	b, ok := e.data["c"]
	assert.True(t, ok)
	assert.Equal(t, "c", string(b))

	e.rwm.Unlock()

	// confirm element deletion after expiry
	time.Sleep(20 * time.Millisecond)

	e.rwm.Lock()

	assert.Equal(t, 5, len(e.data))
	b, ok = e.data["c"]
	assert.False(t, ok)

	e.rwm.Unlock()

	time.Sleep(6 * time.Millisecond)

	e.rwm.Lock()

	assert.Equal(t, 4, len(e.data))
	b, ok = e.data["f"]
	assert.False(t, ok)

	// GetTTL
	setTTL("d", 15*time.Second)
	setTTL("a", 24*time.Second)
	secs := e.GetTTL("a", "d", "ff")
	assert.Equal(t, 3, len(secs))

	assert.True(t, roughly(24.0, secs[0]))
	assert.True(t, roughly(15.0, secs[1]))
	assert.Equal(t, -1.0, secs[2])

	e.rwm.Unlock()

	time.Sleep(50 * time.Millisecond)

	e.rwm.Lock()

	secs = e.GetTTL("a", "d", "ff")
	assert.True(t, roughly(23.95, secs[0]))
	assert.True(t, roughly(14.95, secs[1]))

	// Overwrite existing TTL
	setTTL("a", 700*time.Second)
	secs = e.GetTTL("a", "d", "ff")
	assert.True(t, roughly(700, secs[0]))

	e.rwm.Unlock()
}

func roughly(a, b float64) bool {
	if (a/b > 0.999 && a/b <= 1.0) || (a/b >= 1.0 && a/b < 1.001) {
		return true
	}
	return false
}
