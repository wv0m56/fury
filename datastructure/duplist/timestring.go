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
	rhg       *randomHeightGenerator
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
	ts.rhg = newRandomHeightGenerator(maxHeight, nil)
}

func (ts *TimeString) First() *TimeStringElement {
	return ts.front[0]
}

func (ts *TimeString) DelElement(el *TimeStringElement) {
	if el == nil {
		return
	}

	for i := 0; i < len(el.nexts); i++ {

		if el.prevs[i] == nil {
			ts.front[i] = el.nexts[i]
		} else {
			el.prevs[i].nexts[i] = el.nexts[i]
		}

		if el.nexts[i] != nil {
			el.nexts[i].prevs[i] = el.prevs[i]
		}
	}
}

func (ts *TimeString) Insert(key time.Time, val string) *TimeStringElement {

	el := newTimeStringElement(key, val, ts.rhg)

	if ts.front[0] == nil {
		ts.insert(ts.front, el, nil)
	} else {
		ts.searchAndInsert(el)
	}
	return el
}

func (ts *TimeString) searchAndInsert(el *TimeStringElement) {
	leftAll, iter := ts.search(el.key)
	ts.insert(leftAll, el, iter)
}

func (ts *TimeString) search(key time.Time) (
	leftAll []*TimeStringElement,
	right *TimeStringElement, // iterator and result
) {

	leftAll = make([]*TimeStringElement, ts.maxHeight)

	for h := ts.maxHeight - 1; h >= 0; h-- {

		if h == ts.maxHeight-1 || leftAll[h+1] == nil {
			right = ts.front[h]
		} else {
			leftAll[h] = leftAll[h+1]
			right = leftAll[h].nexts[h]
		}

		for {
			if right == nil || key.Before(right.key) || key.Equal(right.key) { // slow comparison
				break
			} else {
				leftAll[h] = right
				right = right.nexts[h]
			}
		}
	}

	return
}

func (ts *TimeString) insert(leftAll []*TimeStringElement, el, right *TimeStringElement) {

	for i := 0; i < len(el.nexts); i++ {
		el.prevs[i] = leftAll[i]

		if right != nil && i < len(el.nexts) {

			if i < len(right.nexts) {
				right.prevs[i] = el

			} else {

				if leftAll[i] != nil {
					if leftAll[i].nexts[i] != nil {
						leftAll[i].nexts[i].prevs[i] = el
					}

				} else {
					if ts.front[i] != nil {
						ts.front[i].prevs[i] = el
					}
				}
			}
		}

		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			// intercept leftAll[i].nexts[i]
			if leftAll[i] != nil {
				el.nexts[i] = leftAll[i].nexts[i]
			} else {
				el.nexts[i] = ts.front[i]
			}
		}

		// reassign what leftAll[i].nexts[i] points to
		if leftAll[i] == nil {
			ts.front[i] = el
		} else {
			leftAll[i].nexts[i] = el
		}
	}
}
