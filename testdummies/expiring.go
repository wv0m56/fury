package testdummies

import (
	"bytes"
	"io"
	"time"
)

type ExpiringOrigin struct{}

func (_ *ExpiringOrigin) Fetch(key string, _ time.Duration) (io.ReadCloser, *time.Time) {
	t := time.Now().Add(24 * time.Hour)
	return &nodelayReadCloser{bytes.NewReader([]byte(key)), key}, &t
}
