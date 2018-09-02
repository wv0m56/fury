package duplist

import "time"

type TimeStringDuplistElement struct {
	key   time.Time
	val   string
	nexts []*TimeStringDuplistElement
}

func (el *TimeStringDuplistElement) Key() time.Time {
	return el.key
}

func (el *TimeStringDuplistElement) Val() string {
	return el.val
}

func (el *TimeStringDuplistElement) Next() *TimeStringDuplistElement {
	return el.nexts[0]
}

func newDupElem(key time.Time, val string, maxHeight int) *TimeStringDuplistElement {
	lvl := 1 + addHeight(maxHeight)
	return &TimeStringDuplistElement{key, val, make([]*TimeStringDuplistElement, lvl)}
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
