package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeStringLinkedList(t *testing.T) {

	ll := &TimeString{}
	ll.delFront()

	ll.AddToBack("one")
	assert.NotNil(t, ll.front)
	assert.NotNil(t, ll.back)
	assert.Equal(t, ll.front, ll.back)

	ll.delFront()
	assert.Nil(t, ll.front)
	assert.Nil(t, ll.back)

	ptr1 := ll.AddToBack("one")
	ll.AddToBack("two")
	ptr2 := ll.AddToBack("3")
	ll.AddToBack("4")
	assert.Equal(t, "one", ll.front.key)
	assert.Equal(t, "4", ll.back.key)

	var keys string
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "onetwo34", keys)

	ll.Del(ptr2)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "onetwo4", keys)

	ll.Del(ptr1)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "two4", keys)

	ptr3 := ll.AddToBack("back")
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "two4back", keys)

	ll.Del(ptr3)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "two4", keys)
}
