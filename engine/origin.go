package engine

import (
	"io"
	"time"
)

// Origin is to be implemented by objects which fetches data from the cache
// engine's backend.
// Fetch fetches the data associated with key (usually over the network)
// and returns it as a closable reader stream.
// Engine will remember the returned expiry and invalidates key when it's time.
type Origin interface {
	Fetch(key string, timeout time.Duration) (rc io.ReadCloser, expiry *time.Time)
}
