package engine

import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync"
	"time"

	boom "github.com/tylertreat/BoomFilters"
	"github.com/wv0m56/fury/datastructure/duplist"
	"github.com/wv0m56/fury/datastructure/linkedlist"
)

type Engine struct {
	rwm             *sync.RWMutex
	data            map[string][]byte
	fillCond        map[string]*condition
	ttl             *ttlControl
	stats           *accessStats
	o               Origin
	timeout         time.Duration
	payloadTotal    int64
	maxPayloadTotal int64
}

// NewEngine creates a new cache engine with a skiplist as the underlying data
// structure.
func NewEngine(opts *Options) (*Engine, error) {

	{ // sanity checks

		if opts.ExpectedLen < 1024 {
			return nil, errors.New("ExpectedLen must be >= 1024")
		}

		if opts.MaxPayloadTotalBytes < 10*1000*1000 {
			return nil, errors.New("MaxPayloadTotalSize must be >= 10*1000*1000 bytes")
		}

		if opts.CacheFillTimeout < 10*time.Millisecond {
			return nil, errors.New("cachefill timeout too small")
		}

		if opts.TTLTickStep < 1*time.Millisecond {
			return nil, errors.New("TTL tick step too small")
		}

		if opts.AccessStatsTickStep < 1*time.Millisecond ||
			opts.AccessStatsTickStep > opts.AccessStatsRelevanceWindow {

			return nil, errors.New("access stats tick step too small or bigger than relevance window")
		}

		if opts.AccessStatsRelevanceWindow < 100*time.Millisecond {
			return nil, errors.New("access stats relevance window too small")
		}
	}

	// log2(ExpectedLen)-1
	n := int(math.Floor(math.Log2(float64(opts.ExpectedLen / 2))))

	e := &Engine{
		&sync.RWMutex{},
		make(map[string][]byte),
		make(map[string]*condition),
		&ttlControl{
			*(duplist.NewTimeString(n)),
			make(map[string]*duplist.TimeStringElement),
			nil,
		},
		&accessStats{
			sync.Mutex{},
			boom.NewCountMinSketch(0.001, 0.99),
			&linkedlist.TimeString{},
			duplist.NewUint64String(n),
			make(map[string]relevantTuple),
			opts.AccessStatsRelevanceWindow,
			duplist.NewUint64String(n - 1),
			make(map[string]*duplist.Uint64StringElement),
		},
		opts.O,
		opts.CacheFillTimeout,
		0,
		opts.MaxPayloadTotalBytes,
	}

	e.ttl.e = e

	go e.ttl.startLoop(opts.TTLTickStep)
	go e.stats.startLoop(opts.AccessStatsTickStep)

	return e, nil
}

type condition struct {
	sync.Cond
	count int
	b     []byte
	err   error
}

func (e *Engine) Get(key string) (r *bytes.Reader, err error) {
	return e.get(key)
}

func (e *Engine) get(key string) (*bytes.Reader, error) {

	go e.stats.addToWindow(key)

	r := e.tryget(key)
	if r != nil { // cache hit
		return r, nil
	}

	// cache miss
	r, err := e.cacheFill(key)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (e *Engine) tryget(key string) *bytes.Reader {
	e.rwm.RLock()
	defer e.rwm.RUnlock()

	if b, ok := e.data[key]; ok {
		return bytes.NewReader(b)
	}

	return nil
}

func (e *Engine) cacheFill(key string) (*bytes.Reader, error) {

	e.rwm.Lock()
	if b, ok := e.data[key]; ok {
		e.rwm.Unlock()
		return bytes.NewReader(b), nil
	}

	// still locked
	if cond, ok := e.fillCond[key]; ok && cond != nil {

		cond.count++
		return e.blockUntilFilled(key)

	} else {

		e.fillCond[key] = &condition{*sync.NewCond(e.rwm), 1, nil, nil}
		go e.firstFill(key)
		return e.blockUntilFilled(key)
	}
}

func (e *Engine) firstFill(key string) {
	var (
		err error
		rw  *rowWriter
		rc  io.ReadCloser
		exp *time.Time
	)

	defer func() {
		if err != nil {

			if rc != nil {
				_ = rc.Close()
			}
			e.rwm.Lock()
			e.fillCond[key].err = err

		} else {

			e.rwm.Lock()

			if rw.b != nil {

				if rowPayloadSize := rw.b.Len(); rowPayloadSize != 0 &&
					e.payloadTotal+int64(rowPayloadSize) > e.maxPayloadTotal {

					if twiceSpace := 2 * int64(rowPayloadSize); int64(twiceSpace) > e.maxPayloadTotal {
						e.evictUntilFree(e.maxPayloadTotal)
					} else {
						e.evictUntilFree(twiceSpace)
					}
				}

				if exp != nil && exp.After(time.Now()) {
					rw.commit()
					e.setExpiry(key, *exp)
				} else if exp == nil {
					rw.commit()
				}

				e.payloadTotal += int64(rw.b.Len())
			}

			if rw.b != nil && rw.b.Bytes() != nil {
				e.fillCond[key].b = rw.b.Bytes()
			} else {
				e.fillCond[key].b = []byte{} // terminates cond.Wait() loop
			}

		}

		e.fillCond[key].Broadcast()
		e.rwm.Unlock()
	}()

	// fetch from remote and fill up buffer
	rc, exp, err = e.o.Fetch(key, e.timeout)
	if err != nil {
		return
	}
	rw = &rowWriter{key, nil, e}

	if rc != nil {
		_, err = io.Copy(rw, rc)
	} else {
		err = errors.New("nil ReadCloser from Fetch")
	}

	return
}

func (e *Engine) blockUntilFilled(key string) (r *bytes.Reader, err error) {

	c := e.fillCond[key]
	for c.b == nil && c.err == nil {
		e.fillCond[key].Wait()
	}

	if c.err != nil {
		err = c.err
	}

	if b := c.b; b != nil {
		r = bytes.NewReader(e.fillCond[key].b)
	}

	e.fillCond[key].count--
	if e.fillCond[key].count == 0 {
		delete(e.fillCond, key)
	}

	e.rwm.Unlock()

	return
}

// still holding top level lock throughout
func (e *Engine) evictUntilFree(wantedFreeSpace int64) {
	if wantedFreeSpace > e.maxPayloadTotal {
		panic("cache-fill candidate is larger than allowed total") // reconsider
	}

	var enoughFreed bool
	e.stats.Lock()
	defer e.stats.Unlock()

	for it := e.stats.irrelevantDuplist.First(); it != nil; it = it.Next() {

		e.delDataTTLStats(it.Val())
		if freeSpace := e.maxPayloadTotal - e.payloadTotal; freeSpace > wantedFreeSpace {
			enoughFreed = true
			break
		}
	}

	if !enoughFreed {
		for it := e.stats.relevantDuplist.First(); it != nil; it = it.Next() {

			e.delDataTTLStats(it.Val())

			if freeSpace := e.maxPayloadTotal - e.payloadTotal; freeSpace > wantedFreeSpace {
				break
			}
		}
	}
}

func (e *Engine) delData(key string) {
	if b, ok := e.data[key]; ok {
		e.payloadTotal -= int64(len(b))
		delete(e.data, key)
	}
}

func (e *Engine) delDataTTLStats(key string) {
	e.delData(key)
	e.ttl.delTTLEntry(key)
	go e.stats.updateDataDeletion(key)
}

type rowWriter struct {
	key string
	b   *bytes.Buffer
	e   *Engine
}

func (rw *rowWriter) Write(p []byte) (n int, err error) {
	if rw.b == nil {
		rw.b = bytes.NewBuffer(nil)
	}
	return rw.b.Write(p)
}

// no locking.
func (rw *rowWriter) commit() {
	if rw.b != nil {
		rw.e.data[rw.key] = rw.b.Bytes()
	} else {
		rw.e.data[rw.key] = []byte{}
	}
}

// Invalidate deletes keys from the data, TTL, and access stats.
// Only invoke Invalidate for manual cluster control (e.g. global purge).
// Normally, control the invalidation process by setting sensible TTL
// values at origin.
func (e *Engine) Invalidate(keys ...string) {
	e.rwm.Lock()
	for _, v := range keys {
		e.delDataTTLStats(v)
	}
	e.rwm.Unlock()
}
