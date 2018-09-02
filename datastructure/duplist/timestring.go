package duplist

import (
	"time"
)

// TimeStringDuplist is a modified skiplist implementation allowing duplicate
// time keys to exist inside the same list. Elements with duplicate keys are
// adjacent inside TimeStringDuplist, with a later insert placed left of earlier
// ones.
// Elements with different keys are sorted in ascending order as usual.
// TimeStringDuplist is required for implementing TTL.
// TimeStringDuplist does not allow random get or delete by specifying a key and
// instead only allows get or delete on the first element of the list, or delete
// by specifying an element pointer.
type TimeStringDuplist struct {
	front     []*TimeStringDuplistElement
	maxHeight int
}

func NewTimeStringDuplist(maxHeight int) *TimeStringDuplist {
	d := &TimeStringDuplist{}
	d.Init(maxHeight)
	return d
}

func (d *TimeStringDuplist) Init(maxHeight int) {
	d.front = make([]*TimeStringDuplistElement, maxHeight)
	if !(maxHeight < 2 || maxHeight >= 64) {
		d.maxHeight = maxHeight
	} else {
		panic("maxHeight must be between 2 and 64")
	}
}

func (d *TimeStringDuplist) First() *TimeStringDuplistElement {
	return d.front[0]
}

func (d *TimeStringDuplist) DelElement(el *TimeStringDuplistElement) {
	if el == nil {
		return
	}
	left, it := d.iterSearch(el)
	if it == el {
		d.del(left, it)
	}
}

func (d *TimeStringDuplist) iterSearch(el *TimeStringDuplistElement) (
	left []*TimeStringDuplistElement,
	iter *TimeStringDuplistElement,
) {

	left = make([]*TimeStringDuplistElement, d.maxHeight)

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

func (d *TimeStringDuplist) del(left []*TimeStringDuplistElement, el *TimeStringDuplistElement) {
	for i := 0; i < len(el.nexts); i++ {
		d.reassignLeftAtIndex(i, left, el.nexts[i])
	}
}

func (d *TimeStringDuplist) Insert(key time.Time, val string) *TimeStringDuplistElement {

	el := newDupElem(key, val, d.maxHeight)

	if d.front[0] == nil {

		d.insert(d.front, el, nil)

	} else {

		d.searchAndInsert(el)
	}
	return el
}

func (d *TimeStringDuplist) searchAndInsert(el *TimeStringDuplistElement) {
	left, iter := d.search(el.key)
	d.insert(left, el, iter)
}

func (d *TimeStringDuplist) search(key time.Time) (left []*TimeStringDuplistElement, iter *TimeStringDuplistElement) {
	left = make([]*TimeStringDuplistElement, d.maxHeight)

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

func (d *TimeStringDuplist) insert(left []*TimeStringDuplistElement, el, right *TimeStringDuplistElement) {
	for i := 0; i < len(el.nexts); i++ {
		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			d.takeNextsFromLeftAtIndex(i, left, el)
		}

		d.reassignLeftAtIndex(i, left, el)
	}
}

func (d *TimeStringDuplist) takeNextsFromLeftAtIndex(i int, left []*TimeStringDuplistElement, el *TimeStringDuplistElement) {
	if left[i] != nil {
		el.nexts[i] = left[i].nexts[i]
	} else {
		el.nexts[i] = d.front[i]
	}
}

func (d *TimeStringDuplist) reassignLeftAtIndex(i int, left []*TimeStringDuplistElement, el *TimeStringDuplistElement) {
	if left[i] == nil {
		d.front[i] = el
	} else {
		left[i].nexts[i] = el
	}
}

func (d *TimeStringDuplist) DelFirst() {

	for i := 0; i < d.maxHeight; i++ {
		if d.front[i] == nil || d.front[i] != d.front[0] {
			continue
		}
		d.front[i] = d.front[i].nexts[i]
	}
}
