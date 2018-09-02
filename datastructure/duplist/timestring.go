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
	maxHeight int
}

func NewTimeString(maxHeight int) *TimeString {
	d := &TimeString{}
	d.Init(maxHeight)
	return d
}

func (d *TimeString) Init(maxHeight int) {
	d.front = make([]*TimeStringElement, maxHeight)
	if !(maxHeight < 2 || maxHeight >= 64) {
		d.maxHeight = maxHeight
	} else {
		panic("maxHeight must be between 2 and 64")
	}
}

func (d *TimeString) First() *TimeStringElement {
	return d.front[0]
}

func (d *TimeString) DelElement(el *TimeStringElement) {
	if el == nil {
		return
	}
	left, it := d.iterSearch(el)
	if it == el {
		d.del(left, it)
	}
}

func (d *TimeString) iterSearch(el *TimeStringElement) (
	left []*TimeStringElement,
	iter *TimeStringElement,
) {

	left = make([]*TimeStringElement, d.maxHeight)

	for h := d.maxHeight - 1; h >= 0; h-- {

		if h == d.maxHeight-1 || left[h+1] == nil {
			iter = d.front[h]
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

func (d *TimeString) del(left []*TimeStringElement, el *TimeStringElement) {
	for i := 0; i < len(el.nexts); i++ {
		d.reassignLeftAtIndex(i, left, el.nexts[i])
	}
}

func (d *TimeString) Insert(key time.Time, val string) *TimeStringElement {

	el := newTimeStringElement(key, val, d.maxHeight)

	if d.front[0] == nil {

		d.insert(d.front, el, nil)

	} else {

		d.searchAndInsert(el)
	}
	return el
}

func (d *TimeString) searchAndInsert(el *TimeStringElement) {
	left, iter := d.search(el.key)
	d.insert(left, el, iter)
}

func (d *TimeString) search(key time.Time) (left []*TimeStringElement, iter *TimeStringElement) {
	left = make([]*TimeStringElement, d.maxHeight)

	for h := d.maxHeight - 1; h >= 0; h-- {

		if h == d.maxHeight-1 || left[h+1] == nil {
			iter = d.front[h]
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

func (d *TimeString) insert(left []*TimeStringElement, el, right *TimeStringElement) {
	for i := 0; i < len(el.nexts); i++ {
		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			d.takeNextsFromLeftAtIndex(i, left, el)
		}

		d.reassignLeftAtIndex(i, left, el)
	}
}

func (d *TimeString) takeNextsFromLeftAtIndex(i int, left []*TimeStringElement, el *TimeStringElement) {
	if left[i] != nil {
		el.nexts[i] = left[i].nexts[i]
	} else {
		el.nexts[i] = d.front[i]
	}
}

func (d *TimeString) reassignLeftAtIndex(i int, left []*TimeStringElement, el *TimeStringElement) {
	if left[i] == nil {
		d.front[i] = el
	} else {
		left[i].nexts[i] = el
	}
}

func (d *TimeString) DelFirst() {

	for i := 0; i < d.maxHeight; i++ {
		if d.front[i] == nil || d.front[i] != d.front[0] {
			continue
		}
		d.front[i] = d.front[i].nexts[i]
	}
}
