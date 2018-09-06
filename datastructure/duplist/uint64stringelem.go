package duplist

type Uint64StringElement struct {
	key   uint64
	val   string
	nexts []*Uint64StringElement
}

func (el *Uint64StringElement) Key() uint64 {
	return el.key
}

func (el *Uint64StringElement) Val() string {
	return el.val
}

func (el *Uint64StringElement) Next() *Uint64StringElement {
	return el.nexts[0]
}

func newUint64StringElement(key uint64, val string, maxHeight int) *Uint64StringElement {
	return &Uint64StringElement{key, val, make([]*Uint64StringElement, maxHeight)}
}
