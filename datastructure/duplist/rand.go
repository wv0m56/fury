package duplist

import (
	"math/rand"
	"time"
)

type randomHeight struct {
	maxHeight int
	src       rand.Source
}

func newRandomHeight(maxHeight int, src rand.Source) *randomHeight {
	rh := &randomHeight{maxHeight, nil}
	rh.setRandSource(src)
	return rh
}

// true=heads, false=tails
func (rh *randomHeight) flipCoin() bool {
	if rh.src.Int63()%2 == 0 {
		return true
	}
	return false
}

// SetRandSource sets the random number generator used to perform the coin flips
// to determine an element's "height". It is not thread safe and is meant to be
// called only once before using the package.
// If src is nil, time.Now().UnixNano() is used to seed.
func (rh *randomHeight) setRandSource(src rand.Source) {
	if src != nil {
		rh.src = src
	}
	rh.src = rand.NewSource(time.Now().UnixNano())
}

func (rh *randomHeight) height() int {
	var n int
	for n = 1; n < rh.maxHeight; n++ {
		if !rh.flipCoin() {
			break
		}
	}
	return n
}
