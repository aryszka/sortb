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
// the effects of calling Insert() and Remove() while
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

	// Equal is used to identify identical values. Identical
	// values are stored only once. To be able to remove a
	// value from the tree, it must be identified by Equal().
	// It is OK to return always false, if the Remove()
	// function is not used.
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

func balance(n *node) *node {
	ld := n.less.getDepth()
	gd := n.greater.getDepth()
	if ld > 1 && ld > 2*gd {
		to := n.less
		n.less = nil
		n.depth = gd
		n = insert(to, n)
	} else if gd > 1 && gd > 2*ld {
		to := n.greater
		n.greater = nil
		n.depth = ld
		n = insert(to, n)
	}

	n.depth = max(n.less.getDepth(), n.greater.getDepth())
	return n
}

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

	to = balance(to)
	return to
}

func (n *node) remove(value Value) (*node, bool) {
	if n == nil {
		return nil, false
	}

	if n.value.Equal(value) {
		switch {
		case n.less != nil:
			n = insert(n.greater, n.less)
		case n.greater != nil:
			n = insert(n.less, n.greater)
		default:
			n = nil
		}

		return n, true
	}

	var removed bool
	if n.value.Less(value) {
		n.greater, removed = n.greater.remove(value)
	} else {
		n.less, removed = n.less.remove(value)
	}

	n = balance(n)
	return n, removed
}

// Insert a value in the tree.
func (t *Tree) Insert(value Value) {
	t.node = insert(t.node, &node{value: value, depth: 1})
}

// Remove a value from the tree.
func (t *Tree) Remove(value Value) bool {
	var found bool
	t.node, found = t.node.remove(value)
	return found
}

// Get a new iterator.
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

// Get the next value.
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
