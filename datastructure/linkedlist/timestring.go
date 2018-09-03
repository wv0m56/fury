package linkedlist

import "time"

type TimeString struct {
	front *TimeStringElement
	back  *TimeStringElement
}

type TimeStringElement struct {
	lastReadTime time.Time // last accessed time
	val          string
	prev         *TimeStringElement
	next         *TimeStringElement
}

// approximately sorted
func (ll *TimeString) addToBack(val string) *TimeStringElement {

	e := &TimeStringElement{time.Now(), val, nil, nil}
	e.prev = ll.back

	if ll.back != nil {

		ll.back.next = e
		ll.back = e

	} else {

		ll.front = e
		ll.back = e
	}

	return e
}

func (ll *TimeString) delFront() {

	if ll.front == nil && ll.back == nil {
		return
	}

	if ll.front == ll.back { // 1 element
		ll.back = nil
	}

	ll.front = ll.front.next
	if ll.front != nil {
		ll.front.prev = nil
	}
}

func (ll *TimeString) delByPtr(e *TimeStringElement) {

	if e.prev == nil {

		ll.delFront()
		return
	}

	e.prev.next = e.next
	if e.next != nil {
		e.next.prev = e.prev
	}
}
