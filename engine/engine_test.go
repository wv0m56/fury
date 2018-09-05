package engine

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wv0m56/fury/testdummies"
)

var testOptionsDefault Options

func init() {
	testOptionsDefault = OptionsDefault
	testOptionsDefault.MaxPayloadTotalBytes = 4 * 1000 * 1000 * 1000
	testOptionsDefault.O = &testdummies.DelayedOrigin{}
}

func TestOptionValuesSanity(t *testing.T) {
	fmt.Println("// Tests include time sensitive features. It is assumed that")
	fmt.Println("// they are performed on a computer with a modern CPU.")
	fmt.Println("// Failing that, executions might fall behind and produce")
	fmt.Println("// failures when nothing is wrong.")

	opts := testOptionsDefault
	opts.ExpectedLen = 999
	e, err := NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "ExpectedLen must be >= 1024", err.Error())

	opts = testOptionsDefault
	opts.MaxPayloadTotalBytes = 9 * 1000 * 1000
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "MaxPayloadTotalSize must be >= 10*1000*1000 bytes", err.Error())

	opts = testOptionsDefault
	opts.CacheFillTimeout = 9 * time.Millisecond
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "cachefill timeout too small", err.Error())

	opts = testOptionsDefault
	opts.TTLTickStep = 999 * time.Microsecond
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "TTL tick step too small", err.Error())

	opts = testOptionsDefault
	opts.AccessStatsTickStep = 999 * time.Microsecond
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "access stats tick step too small or bigger than relevance window", err.Error())

	opts = testOptionsDefault
	opts.AccessStatsRelevanceWindow = 99 * time.Millisecond
	opts.AccessStatsTickStep = 1 * time.Millisecond
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "access stats relevance window too small", err.Error())

	opts = testOptionsDefault
	opts.AccessStatsRelevanceWindow = 100 * time.Millisecond
	opts.AccessStatsTickStep = 101 * time.Millisecond
	e, err = NewEngine(&opts)
	assert.Nil(t, e)
	assert.Equal(t, "access stats tick step too small or bigger than relevance window", err.Error())
}

func TestSimpleIO(t *testing.T) {

	e, err := NewEngine(&testOptionsDefault)
	assert.Nil(t, err)
	assert.Nil(t, err)

	valR, err := e.Get("water")
	assert.Nil(t, err)
	b, err := ioutil.ReadAll(valR)
	assert.Nil(t, err)
	assert.Equal(t, "water", string(b))
	assert.Nil(t, err)

	// trigger error, see fake.fakeReadCloser
	valR, err = e.Get("error")
	assert.NotNil(t, err)
	assert.Nil(t, valR)
	valR, err = e.Get("error") // make sure the row was not committed above
	assert.NotNil(t, err)
	assert.Nil(t, valR)
}

func TestCachefillTimeout(t *testing.T) {

	opts := testOptionsDefault // origin has 100 ms delay
	opts.CacheFillTimeout = 110 * time.Millisecond
	e, err := NewEngine(&opts)
	assert.Nil(t, err)

	_, err = e.Get("TestCachefillTimeout")
	assert.Nil(t, err)

	opts.CacheFillTimeout = 90 * time.Millisecond
	e2, err := NewEngine(&opts)
	assert.Nil(t, err)
	_, err = e2.Get("TestCachefillTimeout2")
	assert.NotNil(t, err)
	assert.Equal(t, "context deadline exceeded", err.Error())
}

// Test how much time N concurrent calls to CacheFill spend resolving lock
// contention, given 0 network delay.
func BenchmarkHotKey(b *testing.B) {

	N := 10000
	opts := testOptionsDefault
	opts.O = &testdummies.NoDelayOrigin{}
	e, _ := NewEngine(&opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(N)
		for j := 0; j < N; j++ {
			go func() {
				e.Get("hot key")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// Similar to BenchmarkHotKey, except this time origin returns an error.
func BenchmarkErrorKey(b *testing.B) {

	N := 10000
	opts := testOptionsDefault
	opts.O = &testdummies.NoDelayOrigin{}
	e, _ := NewEngine(&opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(N)
		for j := 0; j < N; j++ {
			go func() {
				e.Get("bench error")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}