package duplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlipCoin(t *testing.T) {
	rh := newRandomHeightGenerator(24, nil)
	var heads, tails float32
	assert.True(t, heads == 0.0)
	for i := 0; i < 100000; i++ {
		if rh.flipCoin() {
			heads++
		} else {
			tails++
		}
	}
	ratio := heads / tails
	assert.True(t, ratio < 1.05 && ratio > 0.95,
		"50-50 probability means ratio is close to 1")
}

func BenchmarkFlipCoin(b *testing.B) {
	// must not take to long to flip a coin
	rh := newRandomHeightGenerator(24, nil)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		rh.flipCoin()
	}
}
