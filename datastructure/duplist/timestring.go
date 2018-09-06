package duplist

import (
	"time"
)

// TimeString is a modified skiplist implementation allowing duplicate
// time keys to exist inside the same list. Elements with duplicate keys are
// adjacent inside TimeString, with a later insert placed left of earlier
// ones.
// Elements with different keys are sorted in ascending order as usual.
// TimeString is required for implementing TTL.
// TimeString does not allow random get or delete by specifying a key and
// instead only allows get or delete on the first element of the list, or delete
// by specifying an element pointer.
type TimeString struct {
	front     []*TimeStringElement
	rh        *randomHeight
	maxHeight int
}

func NewTimeString(maxHeight int) *TimeString {
	ts := &TimeString{}
	ts.Init(maxHeight)
	return ts
}

func (ts *TimeString) Init(maxHeight int) {
	ts.maxHeight = maxHeight
	ts.front = make([]*TimeStringElement, maxHeight)
	if maxHeight < 2 || maxHeight >= 64 {
		panic("maxHeight must be between 2 and 64")
	}
	ts.rh = newRandomHeight(maxHeight, nil)
}

func (ts *TimeString) First() *TimeStringElement {
	return ts.front[0]
}

func (ts *TimeString) DelElement(el *TimeStringElement) {
	if el == nil {
		return
	}
	left, it := ts.iterSearch(el)
	if it == el {
		ts.del(left, it)
	}
}

func (ts *TimeString) iterSearch(el *TimeStringElement) (
	left []*TimeStringElement,
	iter *TimeStringElement,
) {

	left = make([]*TimeStringElement, ts.maxHeight)

	for h := ts.maxHeight - 1; h >= 0; h-- {

		if h == ts.maxHeight-1 || left[h+1] == nil {
			iter = ts.front[h]
		} else {
			left[h] = left[h+1]
			iter = left[h].nexts[h]
		}

		for {
			if iter == nil || iter == el || el.key.Before(iter.key) {
				break
			} else {
				left[h] = iter
				iter = iter.nexts[h]
			}
		}
	}

	return
}

func (ts *TimeString) del(left []*TimeStringElement, el *TimeStringElement) {
	for i := 0; i < len(el.nexts); i++ {
		ts.reassignLeftAtIndex(i, left, el.nexts[i])
	}
}

func (ts *TimeString) Insert(key time.Time, val string) *TimeStringElement {

	el := newTimeStringElement(key, val, ts.rh)

	if ts.front[0] == nil {

		ts.insert(ts.front, el, nil)

	} else {

		ts.searchAndInsert(el)
	}
	return el
}

func (ts *TimeString) searchAndInsert(el *TimeStringElement) {
	left, iter := ts.search(el.key)
	ts.insert(left, el, iter)
}

func (ts *TimeString) search(key time.Time) (left []*TimeStringElement, iter *TimeStringElement) {
	left = make([]*TimeStringElement, ts.maxHeight)

	for h := ts.maxHeight - 1; h >= 0; h-- {

		if h == ts.maxHeight-1 || left[h+1] == nil {
			iter = ts.front[h]
		} else {
			left[h] = left[h+1]
			iter = left[h].nexts[h]
		}

		for {
			if iter == nil || key.Before(iter.key) || key.Equal(iter.key) { // slow comparison
				break
			} else {
				left[h] = iter
				iter = iter.nexts[h]
			}
		}
	}

	return
}

func (ts *TimeString) insert(left []*TimeStringElement, el, right *TimeStringElement) {
	for i := 0; i < len(el.nexts); i++ {
		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			ts.takeNextsFromLeftAtIndex(i, left, el)
		}

		ts.reassignLeftAtIndex(i, left, el)
	}
}

func (ts *TimeString) takeNextsFromLeftAtIndex(i int, left []*TimeStringElement, el *TimeStringElement) {
	if left[i] != nil {
		el.nexts[i] = left[i].nexts[i]
	} else {
		el.nexts[i] = ts.front[i]
	}
}

func (ts *TimeString) reassignLeftAtIndex(i int, left []*TimeStringElement, el *TimeStringElement) {
	if left[i] == nil {
		ts.front[i] = el
	} else {
		left[i].nexts[i] = el
	}
}

func (ts *TimeString) DelFirst() {

	for i := 0; i < ts.maxHeight; i++ {
		if ts.front[i] == nil || ts.front[i] != ts.front[0] {
			continue
		}
		ts.front[i] = ts.front[i].nexts[i]
	}
}
