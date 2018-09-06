package duplist

// Uint64String is repetition....cuz generics.
// Uint64String is required by access stats implementation to pool keys which
// are no longer relevant.
type Uint64String struct {
	front     []*Uint64StringElement
	rh        *randomHeight
	maxHeight int
}

func NewUint64String(maxHeight int) *Uint64String {
	us := &Uint64String{}
	us.Init(maxHeight)
	return us
}

func (us *Uint64String) Init(maxHeight int) {
	us.maxHeight = maxHeight
	us.front = make([]*Uint64StringElement, maxHeight)
	if maxHeight < 2 || maxHeight >= 64 {
		panic("maxHeight must be between 2 and 64")
	}
	us.rh = newRandomHeight(maxHeight, nil)
}

func (us *Uint64String) First() *Uint64StringElement {
	return us.front[0]
}

func (us *Uint64String) DelElement(el *Uint64StringElement) {
	if el == nil {
		return
	}
	left, it := us.iterSearch(el)
	if it == el {
		us.del(left, it)
	}
}

// for deletes
func (us *Uint64String) iterSearch(el *Uint64StringElement) (
	left []*Uint64StringElement,
	iter *Uint64StringElement,
) {

	left = make([]*Uint64StringElement, us.maxHeight)

	for h := us.maxHeight - 1; h >= 0; h-- {

		if h == us.maxHeight-1 || left[h+1] == nil {
			iter = us.front[h]
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

func (us *Uint64String) del(left []*Uint64StringElement, el *Uint64StringElement) {
	for i := 0; i < len(el.nexts); i++ {
		us.reassignLeftAtIndex(i, left, el.nexts[i])
	}
}

func (us *Uint64String) Insert(key uint64, val string) *Uint64StringElement {

	el := newUint64StringElement(key, val, us.maxHeight)

	if us.front[0] == nil {

		us.insert(us.front, el, nil)

	} else {

		us.searchAndInsert(el)
	}
	return el
}

func (us *Uint64String) searchAndInsert(el *Uint64StringElement) {
	left, iter := us.search(el.key)
	us.insert(left, el, iter)
}

// for inserts
func (us *Uint64String) search(key uint64) (
	left []*Uint64StringElement,
	iter *Uint64StringElement,
) {

	left = make([]*Uint64StringElement, us.maxHeight)

	for h := us.maxHeight - 1; h >= 0; h-- {

		if h == us.maxHeight-1 || left[h+1] == nil {
			iter = us.front[h]
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

func (us *Uint64String) insert(left []*Uint64StringElement, el, right *Uint64StringElement) {
	for i := 0; i < len(el.nexts); i++ {
		if right != nil && i < len(right.nexts) {

			el.nexts[i] = right

		} else {

			us.takeNextsFromLeftAtIndex(i, left, el)
		}

		us.reassignLeftAtIndex(i, left, el)
	}
}

func (us *Uint64String) takeNextsFromLeftAtIndex(i int, left []*Uint64StringElement, el *Uint64StringElement) {
	if left[i] != nil {
		el.nexts[i] = left[i].nexts[i]
	} else {
		el.nexts[i] = us.front[i]
	}
}

func (us *Uint64String) reassignLeftAtIndex(i int, left []*Uint64StringElement, el *Uint64StringElement) {
	if left[i] == nil {
		us.front[i] = el
	} else {
		left[i].nexts[i] = el
	}
}

func (us *Uint64String) DelFirst() {

	for i := 0; i < us.maxHeight; i++ {
		if us.front[i] == nil || us.front[i] != us.front[0] {
			continue
		}
		us.front[i] = us.front[i].nexts[i]
	}
}
