package duplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64StringDuplist(t *testing.T) {

	d := NewUint64String(24)
	assert.Nil(t, d.First())

	fooEl := d.Insert(500000, "foo")
	barEl := d.Insert(500000, "bar")
	bazEl := d.Insert(500000, "baz")
	quxEl := d.Insert(700000, "qux")
	d.Insert(300000, "first")
	lastEl := d.Insert(999999, "last")
	assert.NotNil(t, fooEl)
	assert.Equal(t, "foo", fooEl.Val())

	first := d.First()
	assert.Equal(t, uint64(300000), first.Key())
	assert.Equal(t, "first", first.Val())

	var vals string
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "firstbazbarfooquxlast", vals)

	d.DelElement(fooEl)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "firstbazbarquxlast", vals)

	d.DelElement(bazEl)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "firstbarquxlast", vals)

	d.DelElement(quxEl)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "firstbarlast", vals)

	d.DelElement(first)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "barlast", vals)

	d.DelElement(barEl)
	vals = ""
	for it := d.First(); it != nil; it = it.Next() {
		vals += it.Val()
	}
	assert.Equal(t, "last", vals)

	d.DelElement(lastEl)
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

func BenchmarkUint64StringDuplistDelete(b *testing.B) {

	N := 1000 * 10
	dup := NewUint64String(24)
	var ptr *Uint64StringElement
	for i := 0; i < N; i++ {
		if i == 8888 {
			ptr = dup.Insert(uint64(i), "abc")
		} else {
			dup.Insert(uint64(i), "abc")
		}
		dup.Insert(uint64(i), "abc")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dup.DelElement(ptr)
	}
}
