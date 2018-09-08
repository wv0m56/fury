package testdummies

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
)

// CustomLengthOrigin simulates origin with custom payload length. Calling
// Fetch("bla/3000", timeout) will fetch a 3000 bytes long payload of zeroes
// associated with key bla.
type CustomLengthOrigin struct{}

func (clo *CustomLengthOrigin) Fetch(key string, timeout time.Duration) (io.ReadCloser, *time.Time) {

	split := strings.Split(key, "/")
	length, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	payload := make([]byte, length)

	return &nodelayReadCloser{bytes.NewReader(payload), key}, nil
}
