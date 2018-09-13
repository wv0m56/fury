package testdummies

import (
	"bytes"
	"io"
	"time"
)

// Origin whose value/payload is always a 10000 bytes long dummy content for
// all keys.
type ZeroesPayloadOrigin struct{}

func (_ *ZeroesPayloadOrigin) Fetch(_ string, _ time.Duration) (
	io.ReadCloser, *time.Time, error) {

	return &zeroesPayloadReadCloser{bytes.NewReader(make([]byte, 10000))}, nil, nil
}

type zeroesPayloadReadCloser struct{ *bytes.Reader }

func (_ *zeroesPayloadReadCloser) Close() error {
	return nil
}

func (tbrc *zeroesPayloadReadCloser) Read(p []byte) (int, error) {
	return tbrc.Read(p)
}
