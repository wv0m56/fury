package duplist

// Uint64String is repetition....cuz generics.
// Uint64String is required by access stats implementation to pool keys which
// are no longer relevant.
type Uint64String struct {
	front     []*Uint64StringElement
	rhg       *randomHeightGenerator
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
	us.rhg = newRandomHeightGenerator(maxHeight, nil)
}

func (us *Uint64String) First() *Uint64StringElement {
	return us.front[0]
}

func (us *Uint64String) DelElement(el *Uint64StringElement) {
	if el == nil {
		return
	}

	for i := 0; i < len(el.nexts); i++ {

		if el.prevs[i] == nil {
			us.front[i] = el.nexts[i]
		} else {
			el.prevs[i].nexts[i] = el.nexts[i]
		}

		if el.nexts[i] != nil {
			el.nexts[i].prevs[i] = el.prevs[i]
		}
	}
}

func (us *Uint64String) Insert(key uint64, val string) *Uint64StringElement {

	el := newUint64StringElement(key, val, us.rhg)

	if us.front[0] == nil {
		us.insert(make([]*Uint64StringElement, us.maxHeight), el, nil)
	} else {
		us.searchAndInsert(el)
	}
	return el
}

func (us *Uint64String) searchAndInsert(el *Uint64StringElement) {
	leftAll, right := us.search(el.key)
	us.insert(leftAll, el, right)
}

func (us *Uint64String) search(key uint64) (
	leftAll []*Uint64StringElement,
	right *Uint64StringElement, // iterator and result
) {

	leftAll = make([]*Uint64StringElement, us.maxHeight)

	for h := us.maxHeight - 1; h >= 0; h-- {

		if h == us.maxHeight-1 || leftAll[h+1] == nil {
			right = us.front[h]
		} else {
			leftAll[h] = leftAll[h+1]
			right = leftAll[h].nexts[h]
		}

		for {
			if right == nil || key <= right.key {
				break
			} else {
				leftAll[h] = right
				right = right.nexts[h]
			}
		}
	}

	return
}

func (us *Uint64String) insert(leftAll []*Uint64StringElement, el, right *Uint64StringElement) {

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
					if us.front[i] != nil {
						us.front[i].prevs[i] = el
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
				el.nexts[i] = us.front[i]
			}
		}

		// reassign what leftAll[i].nexts[i] points to
		if leftAll[i] == nil {
			us.front[i] = el
		} else {
			leftAll[i].nexts[i] = el
		}
	}
}
