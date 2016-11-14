/*
Package sortb provides a data structure to store sorted generic values
using a balanced binary search tree.

Stored objects must implement the Less(Value) bool comparison and the
Equal(Value) bool identity check. For tree inserts and traversal, only
the Less comparison is used. This allows use cases where the comparison
and the equality is based on different attributes of the stored objects.
*/
package sortb

type node struct {
	value   Value
	less    *node
	greater *node
	depth   int
}

// Tree objects store sorted values.
type Tree struct {
	node *node
}

// Iterator is used to iterate over values in the tree
// from the smallest to the greatest or in reverse order.
//
// The iterator does not operate over a snapshot, so
// the effects of calling Insert() and Delete() while
// iterating may result in invalid iteration.
type Iterator struct {
	node    *node
	child   *Iterator
	from    Value
	reverse bool
}

// Value must be implemented by the objects to be stored by the tree.
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

func insert(to *node, n *node) (*node, bool) {
	if to == nil {
		return n, true
	}

	if n.value.Equal(to.value) {
		return to, false
	}

	var changed bool
	if n.value.Less(to.value) {
		to.less, changed = insert(to.less, n)
	} else {
		to.greater, changed = insert(to.greater, n)
	}

	to.updateDepth()
	to = balance(to)
	return to, changed
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

func prev(n *node, v Value) *node {
	if n == nil || v == nil {
		return nil
	}

	if n.value.Less(v) {
		if n.greater == nil {
			return n
		}

		ng := prev(n.greater, v)
		if ng == nil {
			return n
		}

		return ng
	}

	return prev(n.less, v)
}

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

// Insert a value in the tree. If the value is already
// a member of the tree, the tree stays unchanged. It
// returns true if the tree was changed.
func (t *Tree) Insert(v Value) bool {
	if v == nil {
		return false
	}

	var changed bool
	t.node, changed = insert(t.node, &node{value: v, depth: 1})
	return changed
}

// Find returns true if a value is a member of the tree.
func (t *Tree) Find(v Value) bool {
	return find(t.node, v) != nil
}

// Next returns the next value in order and true, or it returns
// false as the second argument if such value was not found. The
// value of the v argument does not need to be a member of the
// tree.
func (t *Tree) Next(v Value) (Value, bool) {
	n := next(t.node, v)
	if n == nil {
		return nil, false
	}

	return n.value, true
}

// Prev returns the previous value in order and true, or it
// returns false as the second argument if such value was not
// found. The value of the v argument does not need to be a
// member of the tree.
func (t *Tree) Prev(v Value) (Value, bool) {
	n := prev(t.node, v)
	if n == nil {
		return nil, false
	}

	return n.value, true
}

// Delete a value from the tree. It returns true if the tree was
// changed.
func (t *Tree) Delete(v Value) bool {
	var found bool
	t.node, found = del(t.node, v)
	return found
}

// Iterate returns a new iterator to iterate over the sorted values
// stored by the tree.
//
// The from argument tells the iterator where to start the iteration.
// If it is nil, it starts at the beginning. If it is not nil, it
// starts with the first value than which from is less. The value of
// from does not need to be a member of the tree.
func (t *Tree) Iterate(from Value) *Iterator {
	return newIterator(t.node, from, false)
}

// Reverse is like Iterate but with reverse order.
func (t *Tree) Reverse(from Value) *Iterator {
	return newIterator(t.node, from, true)
}

func newIterator(n *node, from Value, reverse bool) *Iterator {
	i := &Iterator{node: n, from: from, reverse: reverse}
	if i.node == nil {
		return i
	}

	if i.reverse {
		i.childrenBeforeReverse()
	} else {
		i.childrenBefore()
	}

	return i
}

func (i *Iterator) childrenBefore() {
	if i.node.less == nil {
		return
	}

	if i.from != nil && !i.from.Less(i.node.value) {
		return
	}

	i.child = newIterator(i.node.less, i.from, false)
}

func (i *Iterator) childrenBeforeReverse() {
	if i.node.greater == nil {
		return
	}

	if i.from != nil && !i.node.value.Less(i.from) {
		return
	}

	i.child = newIterator(i.node.greater, i.from, true)
}

func (i *Iterator) childrenAfter() {
	if i.node.greater != nil {
		i.child = newIterator(i.node.greater, i.from, false)
	}
}

func (i *Iterator) childrenAfterReverse() {
	if i.node.less != nil {
		i.child = newIterator(i.node.less, i.from, true)
	}
}

func (i *Iterator) afterFrom(v Value) bool {
	if i.from == nil {
		return true
	}

	if i.reverse && v.Less(i.from) {
		return true
	}

	if !i.reverse && i.from.Less(v) {
		return true
	}

	return false
}

// Next returns the next value in order (or the previous in case of
// the reverse iterator). If the iterator has reached the last stored
// value, Next() returns false as the second argument, otherwise true.
func (i *Iterator) Next() (Value, bool) {
	if i.child != nil {
		if v, ok := i.child.Next(); ok && i.afterFrom(v) {
			return v, true
		}
	}

	if i.node == nil {
		return nil, false
	}

	v := i.node.value
	if i.reverse {
		i.childrenAfterReverse()
	} else {
		i.childrenAfter()
	}

	i.node = nil

	if i.afterFrom(v) {
		return v, true
	}

	return i.Next()
}
