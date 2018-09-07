package duplist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeStringDuplist(t *testing.T) {

	d := NewTimeString(24)
	assert.Nil(t, d.First())

	now := time.Now()
	fooEl := d.Insert(now.Add(50*time.Millisecond), "foo")
	barEl := d.Insert(now.Add(50*time.Millisecond), "bar")
	bazEl := d.Insert(now.Add(50*time.Millisecond), "baz")
	quxEl := d.Insert(now.Add(70*time.Millisecond), "qux")
	d.Insert(now.Add(30*time.Millisecond), "first")
	lastEl := d.Insert(now.Add(99*time.Millisecond), "last")
	assert.NotNil(t, fooEl)
	assert.Equal(t, "foo", fooEl.Val())

	first := d.First()
	assert.Equal(t, now.Add(30*time.Millisecond), first.Key())
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

func BenchmarkTimeStringDuplistInsert(b *testing.B) {

	N := 1000 * 10
	dup := NewTimeString(22)
	for i := 0; i < N; i++ {
		dup.Insert(time.Now(), time.Now().String())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dup.Insert(time.Now(), time.Now().String())
	}
}
