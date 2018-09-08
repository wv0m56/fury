package testdummies

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"time"
)

// An origin which returns data with random expiry and random
// 1000-2000 bytes long payload.
type RandomOrigin struct{}

func (_ *RandomOrigin) Fetch(_ string, timeout time.Duration) (io.ReadCloser, *time.Time) {
	expiry := 10*time.Millisecond + time.Duration(rand.Int63n(40*int64(time.Millisecond)))
	t := time.Now().Add(expiry)
	b := make([]byte, 1000+rand.Intn(1000))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &randomReadCloser{bytes.NewReader(b), false, ctx, cancel}, &t
}

type randomReadCloser struct {
	br          *bytes.Reader
	hasBeenRead bool
	ctx         context.Context
	cancel      context.CancelFunc
}

func (rrc *randomReadCloser) Close() error {
	rrc.cancel()
	return rrc.ctx.Err()
}

func (rrc *randomReadCloser) Read(p []byte) (int, error) {
	if !rrc.hasBeenRead {
		// first time Read is called for key
		select {

		case <-rrc.ctx.Done():
			return 0, rrc.ctx.Err()

		case <-time.After(time.Duration(rand.Int63n(20)) * time.Millisecond):
			rrc.hasBeenRead = true
			return rrc.br.Read(p)
		}

	} else {

		select {

		case <-rrc.ctx.Done():
			return 0, rrc.ctx.Err()

		default:
			return rrc.br.Read(p)
		}
	}
}
