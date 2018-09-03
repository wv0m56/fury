package engine

import (
	"time"
)

// Options to be passed into Engine creation.
type Options struct {
	// ExpectedLen is the number of expected (k, v) rows in the cache.
	// It's better to overestimate ExpectedLen than to underestimate it.
	// NewEngine panics if ExpectedLen less than 1024 (pointless).
	ExpectedLen int64

	AccessStatsRelevanceWindow time.Duration
	AccessStatsTickStep        time.Duration
	TTLTickStep                time.Duration
	CacheFillTimeout           time.Duration
	O                          Origin

	// MaxPayloadTotalBytes is the total sum of the length of all value/payload
	// (in bytes) from all rows.
	// It must be greater than 10*1000*1000 bytes.
	MaxPayloadTotalBytes int64
}

var OptionsDefault = Options{
	ExpectedLen:                10 * 1000 * 1000,
	AccessStatsRelevanceWindow: 24 * 3600 * time.Second,
	AccessStatsTickStep:        1 * time.Second,
	TTLTickStep:                250 * time.Millisecond,
	CacheFillTimeout:           250 * time.Millisecond,
}
