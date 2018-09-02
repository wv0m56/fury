package duplist

import "time"

type TimeStringElement struct {
	key   time.Time
	val   string
	nexts []*TimeStringElement
}

func (el *TimeStringElement) Key() time.Time {
	return el.key
}

func (el *TimeStringElement) Val() string {
	return el.val
}

func (el *TimeStringElement) Next() *TimeStringElement {
	return el.nexts[0]
}

func newDupElem(key time.Time, val string, maxHeight int) *TimeStringElement {
	lvl := 1 + addHeight(maxHeight)
	return &TimeStringElement{key, val, make([]*TimeStringElement, lvl)}
}

func addHeight(maxHeight int) int {
	var n int
	for n = 0; n < maxHeight-1; n++ {
		if !flipCoin() {
			break
		}
	}
	return n
}
