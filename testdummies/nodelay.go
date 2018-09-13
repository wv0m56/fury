package testdummies

import (
	"bytes"
	"errors"
	"io"
	"time"
)

type NoDelayOrigin struct{}

func (_ *NoDelayOrigin) Fetch(key string, _ time.Duration) (
	io.ReadCloser, *time.Time, error) {

	return &nodelayReadCloser{bytes.NewReader([]byte(key)), key}, nil, nil
}

type nodelayReadCloser struct {
	br  *bytes.Reader
	key string
}

func (_ *nodelayReadCloser) Close() error {
	return nil
}

func (brc *nodelayReadCloser) Read(p []byte) (int, error) {
	if brc.key == "bench error" {
		return 0, errors.New("fake bench error")
	}
	return brc.br.Read(p)
}
