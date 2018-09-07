package duplist

import "time"

type TimeStringElement struct {
	key   time.Time
	val   string
	prevs []*TimeStringElement
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

func newTimeStringElement(key time.Time, val string, rh *randomHeightGenerator) *TimeStringElement {
	height := rh.height()
	return &TimeStringElement{
		key, val,
		make([]*TimeStringElement, height),
		make([]*TimeStringElement, height),
	}
}
