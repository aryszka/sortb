package sortb

type node struct {
	value   Value
	less    *node
	greater *node
	depth   int
}

// Tree objects store values sorted as a binary tree.
type Tree struct {
	node *node
}

// Iterator is used to iterate over values in the tree
// from the smallest to the greatest.
//
// The iterator does not operate over a snapshot, so
// the effects of calling Insert() and Delete() while
// iterating may result in invalid iteration.
type Iterator struct {
	node  *node
	child *Iterator
}

// Value must be implemented by the objects to be stored.
type Value interface {

	// Less is used to identify the sorting position of
	// a stored value.
	Less(i Value) bool

	// Equal is used to identify values. Identical values are
	// stored only once. It is OK to return always false, if
	// the Delete() and Find() functions are not used.
	Equal(i Value) bool
}

func max(i, j int) int {
	if i < j {
		return j
	}

	return i
}

func (n *node) getDepth() int {
	if n == nil {
		return 0
	}

	return n.depth
}

func (n *node) updateDepth() {
	n.depth = 1 + max(n.less.getDepth(), n.greater.getDepth())
}

func swapLeft(n *node) *node {
	l := n.less
	n.less, l.greater = l.greater, n
	n.updateDepth()
	l.updateDepth()
	return l
}

func swapRight(n *node) *node {
	l := n.greater
	n.greater, l.less = l.less, n
	n.updateDepth()
	l.updateDepth()
	return l
}

func balance(n *node) *node {
	ld := n.less.getDepth()
	gd := n.greater.getDepth()
	dd := ld - gd

	if dd > 1 {
		if n.less.less.getDepth() < n.less.greater.getDepth() {
			n.less = swapRight(n.less)
		}

		n = swapLeft(n)
	} else if dd < -1 {
		if n.greater.greater.getDepth() < n.greater.less.getDepth() {
			n.greater = swapLeft(n.greater)
		}

		n = swapRight(n)
	}

	return n
}

// TODO: eliminate recursion
func insert(to *node, n *node) *node {
	if to == nil {
		return n
	}

	if n.value.Equal(to.value) {
		return to
	}

	if n.value.Less(to.value) {
		to.less = insert(to.less, n)
	} else {
		to.greater = insert(to.greater, n)
	}

	to.updateDepth()
	to = balance(to)
	return to
}

func find(n *node, v Value) *node {
	if n == nil || v == nil {
		return nil
	}

	if n.value.Equal(v) {
		return n
	}

	if n.value.Less(v) {
		return find(n.greater, v)
	}

	return find(n.less, v)
}

func next(n *node, v Value) *node {
	if n == nil || v == nil {
		return nil
	}

	if v.Less(n.value) {
		if n.less == nil {
			return n
		}

		nl := next(n.less, v)
		if nl == nil {
			return n
		}

		return nl
	}

	return next(n.greater, v)
}

// TODO:
// - optimize
// - eliminate recursion
func del(n *node, v Value) (*node, bool) {
	if n == nil {
		return nil, false
	}

	if n.value.Equal(v) {
		switch {
		case n.less == nil && n.greater == nil:
			n = nil
		case n.less == nil:
			n = n.greater
		case n.greater == nil:
			n = n.less
		default:
			nn := next(n, n.value)
			n.value = nn.value
			n.greater, _ = del(n.greater, n.value)
			n.updateDepth()
			n = balance(n)
		}

		return n, true
	}

	var deleted bool
	if n.value.Less(v) {
		n.greater, deleted = del(n.greater, v)
	} else {
		n.less, deleted = del(n.less, v)
	}

	if deleted {
		n.updateDepth()
		n = balance(n)
	}

	return n, deleted
}

// Insert a value in the tree.
func (t *Tree) Insert(v Value) {
	if v != nil {
		t.node = insert(t.node, &node{value: v, depth: 1})
	}
}

// Find returns true if a value is a member of the tree.
func (t *Tree) Find(v Value) bool {
	return find(t.node, v) != nil
}

func (t *Tree) Next(v Value) Value {
	n := next(t.node, v)
	if n == nil {
		return nil
	}

	return n.value
}

// Delete a value from the tree.
func (t *Tree) Delete(v Value) bool {
	var found bool
	t.node, found = del(t.node, v)
	return found
}

// Iterate returns a new iterator.
func (t *Tree) Iterate() *Iterator {
	return newIterator(t.node)
}

func newIterator(n *node) *Iterator {
	i := &Iterator{node: n}
	if i.node != nil && i.node.less != nil {
		i.child = newIterator(i.node.less)
	}

	return i
}

// Next returns the next value.
func (i *Iterator) Next() (Value, bool) {
	if i.child != nil {
		if n, ok := i.child.Next(); ok {
			return n, true
		}
	}

	if i.node == nil {
		return nil, false
	}

	n := i.node.value
	if i.node.greater != nil {
		i.child = newIterator(i.node.greater)
	}

	i.node = nil
	return n, true
}
