package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeStringLinkedList(t *testing.T) {

	ll := &TimeString{}
	ll.delFront()

	ll.addToBack("one")
	assert.NotNil(t, ll.front)
	assert.NotNil(t, ll.back)
	assert.Equal(t, ll.front, ll.back)

	ll.delFront()
	assert.Nil(t, ll.front)
	assert.Nil(t, ll.back)

	ptr1 := ll.addToBack("one")
	ll.addToBack("two")
	ptr2 := ll.addToBack("3")
	ll.addToBack("4")
	assert.Equal(t, "one", ll.front.val)
	assert.Equal(t, "4", ll.back.val)

	var vals string
	for it := ll.front; it != nil; it = it.next {
		vals += it.val
	}
	assert.Equal(t, "onetwo34", vals)

	ll.delByPtr(ptr2)
	vals = ""
	for it := ll.front; it != nil; it = it.next {
		vals += it.val
	}
	assert.Equal(t, "onetwo4", vals)

	ll.delByPtr(ptr1)
	vals = ""
	for it := ll.front; it != nil; it = it.next {
		vals += it.val
	}
	assert.Equal(t, "two4", vals)

	ptr3 := ll.addToBack("back")
	vals = ""
	for it := ll.front; it != nil; it = it.next {
		vals += it.val
	}
	assert.Equal(t, "two4back", vals)

	ll.delByPtr(ptr3)
	vals = ""
	for it := ll.front; it != nil; it = it.next {
		vals += it.val
	}
	assert.Equal(t, "two4", vals)
}
