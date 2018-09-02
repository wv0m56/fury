package duplist

// Uint64String is repetition....cuz generics.
// Uint64String is required by EvictPolicy implementation to pool keys which are
// no longer relevant.
type Uint64String struct {
	front     []*Uint64StringElement
	maxHeight int
}

func NewUint64String(maxHeight int) *Uint64String {
	d := &Uint64String{}
	d.Init(maxHeight)
	return d
}

func (d *Uint64String) Init(maxHeight int) {
	d.front = make([]*Uint64StringElement, maxHeight)
	if !(maxHeight < 2 || maxHeight >= 64) {
		d.maxHeight = maxHeight
	} else {
		panic("maxHeight must be between 2 and 64")
	}
}

func (d *Uint64String) First() *Uint64StringElement {
	return d.front[0]
}

func (d *Uint64String) DelElement(el *Uint64StringElement) {
	if el == nil {
		return
	}
	left, it := d.iterSearch(el)
	if it == el {
		d.del(left, it)
	}
}

func (d *Uint64String) iterSearch(el *Uint64StringElement) (
	left []*Uint64StringElement,
	iter *Uint64StringElement,
) {

	left = make([]*Uint64StringElement, d.maxHeight)

	for h := d.maxHeight - 1; h >= 0; h-- {

		if h == d.maxHeight-1 || left[h+1] == nil {
			iter = d.front[h]
		} else {
			left[h] = left[h+1]
			iter = left[h].nexts[h]
		}

		for {
			if iter == nil || iter == el || el.key < iter.key {
				break
			} else {
				left[h] = iter
				iter = iter.nexts[h]
			}
		}
	}

	return
}

func (d *Uint64String) del(left []*Uint64StringElement, el *Uint64StringElement) {
	for i := 0; i < len(el.nexts); i++ {
		d.reassignLeftAtIndex(i, left, el.nexts[i])
	}
}

func (d *Uint64String) Insert(key uint64, val string) *Uint64StringElement {

	el := newUint64StringElement(key, val, d.maxHeight)

	if d.front[0] == nil {

		d.insert(d.front, el, nil)

	} else {

		d.searchAndInsert(el)
	}
	return el
}

func (d *Uint64String) searchAndInsert(el *Uint64StringElement) {
	left, iter := d.search(el.key)
	d.insert(left, el, iter)
}

func (d *Uint64String) search(key uint64) (left []*Uint64StringElement, iter *Uint64StringElement) {
	left = make([]*Uint64StringElement, d.maxHeight)

	for h := d.maxHeight - 1; h >= 0; h-- {

		if h == d.maxHeight-1 || left[h+1] == nil {
			iter = d.front[h]
		} else {
			left[h] = left[h+1]
			iter = left[h].nexts[h]
		}

		for {
			if iter == nil || key <= iter.key {
				break
			} else {
				left[h] = iter
				iter = iter.nexts[h]
			}
		}
	}

	return
}

func (d *Uint64String) insert(left []*Uint64StringElement, el, right *Uint64StringElement) {
	for i := 0; i < len(el.nexts); i++ {
		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			d.takeNextsFromLeftAtIndex(i, left, el)
		}

		d.reassignLeftAtIndex(i, left, el)
	}
}

func (d *Uint64String) takeNextsFromLeftAtIndex(i int, left []*Uint64StringElement, el *Uint64StringElement) {
	if left[i] != nil {
		el.nexts[i] = left[i].nexts[i]
	} else {
		el.nexts[i] = d.front[i]
	}
}

func (d *Uint64String) reassignLeftAtIndex(i int, left []*Uint64StringElement, el *Uint64StringElement) {
	if left[i] == nil {
		d.front[i] = el
	} else {
		left[i].nexts[i] = el
	}
}

func (d *Uint64String) DelFirst() {

	for i := 0; i < d.maxHeight; i++ {
		if d.front[i] == nil || d.front[i] != d.front[0] {
			continue
		}
		d.front[i] = d.front[i].nexts[i]
	}
}
