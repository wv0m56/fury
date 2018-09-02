package duplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64StringDuplist(t *testing.T) {

	d := NewUint64String(24)
	assert.Nil(t, d.First())

	el := d.Insert(500000, "foo")
	d.Insert(500000, "bar")
	d.Insert(500000, "baz")
	d.Insert(700000, "qux")
	d.Insert(300000, "first")
	assert.NotNil(t, el)
	assert.Equal(t, "foo", el.Val())

	first := d.First()
	assert.Equal(t, uint64(300000), first.Key())
	assert.Equal(t, "first", first.Val())

	var vals string
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "firstbazbarfooqux", vals)

	d.DelFirst()
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "bazbarfooqux", vals)

	d.DelElement(el)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "bazbarqux", vals)

	d.DelFirst()
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "barqux", vals)

	d.DelFirst()
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "qux", vals)

	d.DelFirst()
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "", vals)
}

func BenchmarkUint64StringDuplistInsert(b *testing.B) {

	N := 1000 * 10
	dup := NewUint64String(22)
	for i := 0; i < N; i++ {
		dup.Insert(uint64(i), "abc")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dup.Insert(8888, "abc")
	}
}
