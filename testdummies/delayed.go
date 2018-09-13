// Package testdummies contains implementations of the Origin interface for tests.
package testdummies

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"
)

type DelayedOrigin struct{}

// Fetch fetches dummy data. "error" as key simulates a network error should
// the returned io.ReadCloser is read. Else returns &bytes.Reader([]byte(key))
// implementing a no-op Close() method with 100ms delay. If timeout has
// elapsed and Fetch has not finished fetching data, it terminates and returns
// (err, nil).
func (do *DelayedOrigin) Fetch(key string, timeout time.Duration) (
	io.ReadCloser, *time.Time, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &delayedReadCloser{
		bytes.NewReader([]byte(key)),
		key,
		false,
		ctx,
		cancel,
	}, nil, nil
}

type delayedReadCloser struct {
	br          *bytes.Reader
	key         string
	hasBeenRead bool
	ctx         context.Context
	cancel      context.CancelFunc
}

func (drc *delayedReadCloser) Close() error {
	drc.cancel()
	return drc.ctx.Err()
}

func (drc *delayedReadCloser) Read(p []byte) (int, error) {

	if !drc.hasBeenRead {
		// first time Read is called for key
		select {

		case <-drc.ctx.Done():
			return 0, drc.ctx.Err()

		case <-time.After(100 * time.Millisecond):
			if drc.key == "error" {
				return 0, errors.New("fake error")
			}
			drc.hasBeenRead = true
			return drc.br.Read(p)
		}

	} else {

		select {

		case <-drc.ctx.Done():
			return 0, drc.ctx.Err()

		default:
			return drc.br.Read(p)
		}
	}
}
