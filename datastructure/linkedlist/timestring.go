package linkedlist

import "time"

type TimeString struct {
	front *TimeStringElement
	back  *TimeStringElement
}

type TimeStringElement struct {
	lastAccessed time.Time
	key          string
	prev         *TimeStringElement
	next         *TimeStringElement
}

// approximately sorted
func (ll *TimeString) AddToBack(key string) *TimeStringElement {
	e := &TimeStringElement{time.Now(), key, nil, nil}
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

func (ll *TimeString) Del(ptr *TimeStringElement) {
	if ptr.prev == nil {
		ll.delFront()
		return
	}
	ptr.prev.next = ptr.next
	if ptr.next != nil {
		ptr.next.prev = ptr.prev
	}
}

func (ll *TimeString) Front() *TimeStringElement {
	return ll.front
}

func (el *TimeStringElement) Next() *TimeStringElement {
	return el.next
}

func (el *TimeStringElement) LastAccessed() time.Time {
	return el.lastAccessed
}

func (el *TimeStringElement) Key() string {
	return el.key
}
