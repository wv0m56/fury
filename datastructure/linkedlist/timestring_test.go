package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeStringLinkedList(t *testing.T) {

	ll := &TimeString{}

	ll.AddToBack("one")
	assert.NotNil(t, ll.front)
	assert.NotNil(t, ll.back)
	assert.Equal(t, ll.front, ll.back)

	ll.delFront()
	assert.Nil(t, ll.front)
	assert.Nil(t, ll.back)

	ptr1 := ll.AddToBack("one")
	ptr2 := ll.AddToBack("two")
	ptr3 := ll.AddToBack("3")
	ptr4 := ll.AddToBack("4")
	assert.Equal(t, "one", ll.front.key)
	assert.Equal(t, "4", ll.back.key)

	var keys string
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "onetwo34", keys)

	ll.Del(ptr3)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "onetwo4", keys)

	ll.Del(ptr4)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "onetwo", keys)

	ll.Del(ptr1)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "two", keys)

	ptr5 := ll.AddToBack("back")
	assert.Equal(t, "back", ll.back.key)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "twoback", keys)

	ll.Del(ptr5)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "two", keys)

	ll.Del(ptr2)
	keys = ""
	for it := ll.front; it != nil; it = it.next {
		keys += it.key
	}
	assert.Equal(t, "", keys)
}
