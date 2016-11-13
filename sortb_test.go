package sortb

import (
	"math/rand"
	"testing"
)

type intt int

func (i intt) Less(j Value) bool  { return i < j.(intt) }
func (i intt) Equal(j Value) bool { return i == j.(intt) }

func testBalance(t *testing.T, n int, f func(int) intt) {
	all := make([]intt, 0, n)
	tree := new(Tree)

	for i := 0; i < n; i++ {
		ii := f(i)
		all = append(all, ii)

		tree.Insert(ii)

		dd := tree.node.less.getDepth() - tree.node.greater.getDepth()
		if dd < -1 || dd > 1 {
			t.Error("failed to balance tree, insert", n)
		}
	}

	for _, ii := range all {
		tree.Delete(ii)
		if tree.node.getDepth() == 0 {
			break
		}

		dd := tree.node.less.getDepth() - tree.node.greater.getDepth()
		if dd < -1 || dd > 1 {
			t.Error("failed to balance tree, delete", dd, n)
		}
	}
}

func TestBalanceLinear(t *testing.T) {
	for _, ti := range []int{
		1,
		2,
		5,
		10,
		20,
		50,
		100,
		200,
		500,
		1000,
	} {
		testBalance(t, ti, func(i int) intt { return intt(i) })
	}
}

func TestBalanceRandom(t *testing.T) {
	sizes := []int{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	for i := 0; i < 10; i++ {
		rnd := rand.New(rand.NewSource(int64(i)))
		s := sizes
		if i > 2 || testing.Short() {
			sizes = sizes[:9]
		}

		for _, j := range s {
			testBalance(t, j, func(int) intt { return intt(rnd.Intn(j)) })
		}
	}
}

func TestInsertIter(t *testing.T) {
	for _, ti := range []struct {
		insert []intt
		check  []intt
	}{{
		nil,
		nil,
	}, {
		[]intt{42},
		[]intt{42},
	}, {
		[]intt{42, 42},
		[]intt{42},
	}, {
		[]intt{-5, 42, -42, 42, 3, -18},
		[]intt{-42, -18, -5, 3, 42},
	}} {
		tree := new(Tree)
		for _, i := range ti.insert {
			tree.Insert(i)
		}

		iter := tree.Iterate(nil)
		for _, i := range ti.check {
			if v, ok := iter.Next(); !ok || v.(intt) != i {
				t.Error("failed to retrieve", i)
			}
		}

		if v, ok := iter.Next(); ok {
			t.Error("unexpected value", v.(intt))
		}
	}
}

func TestFind(t *testing.T) {
	for _, ti := range []struct {
		init   []intt
		find   Value
		expect bool
	}{{
		nil,
		nil,
		false,
	}, {
		nil,
		intt(42),
		false,
	}, {
		[]intt{42},
		nil,
		false,
	}, {
		[]intt{-18},
		intt(42),
		false,
	}, {
		[]intt{-18, 42},
		intt(42),
		true,
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		if tree.Find(ti.find) != ti.expect {
			t.Error("invalid find result", !ti.expect)
		}
	}
}

func TestNext(t *testing.T) {
	for _, ti := range []struct {
		init        []intt
		value, next Value
	}{{
		nil,
		nil,
		nil,
	}, {
		nil,
		intt(42),
		nil,
	}, {
		[]intt{42},
		nil,
		nil,
	}, {
		[]intt{42},
		intt(81),
		nil,
	}, {
		[]intt{42},
		intt(42),
		nil,
	}, {
		[]intt{42},
		intt(18),
		intt(42),
	}, {
		[]intt{-18, -5, 3, 42},
		nil,
		nil,
	}, {
		[]intt{-18, -5, 3, 42},
		intt(-42),
		intt(-18),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(-18),
		intt(-5),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(1),
		intt(3),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(3),
		intt(42),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(42),
		nil,
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		if n := tree.Next(ti.value); n != ti.next {
			t.Error("failed to find next value", n, ti.next)
		}
	}
}

func TestPrev(t *testing.T) {
	for _, ti := range []struct {
		init        []intt
		value, next Value
	}{{
		nil,
		nil,
		nil,
	}, {
		nil,
		intt(42),
		nil,
	}, {
		[]intt{42},
		nil,
		nil,
	}, {
		[]intt{42},
		intt(18),
		nil,
	}, {
		[]intt{42},
		intt(42),
		nil,
	}, {
		[]intt{42},
		intt(81),
		intt(42),
	}, {
		[]intt{-18, -5, 3, 42},
		nil,
		nil,
	}, {
		[]intt{-18, -5, 3, 42},
		intt(81),
		intt(42),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(42),
		intt(3),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(1),
		intt(-5),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(-5),
		intt(-18),
	}, {
		[]intt{-18, -5, 3, 42},
		intt(-18),
		nil,
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		if n := tree.Prev(ti.value); n != ti.next {
			t.Error("failed to find previous value", n, ti.next)
		}
	}
}

func TestDelete(t *testing.T) {
	for _, ti := range []struct {
		init    []intt
		delete  []intt
		deleted []bool
		check   []intt
	}{{
		nil,
		nil,
		nil,
		nil,
	}, {
		nil,
		[]intt{42},
		[]bool{false},
		nil,
	}, {
		[]intt{42},
		[]intt{42},
		[]bool{true},
		nil,
	}, {
		[]intt{42},
		[]intt{-42, -42},
		[]bool{false, false},
		[]intt{42},
	}, {
		[]intt{42},
		[]intt{42, 42},
		[]bool{true, false},
		nil,
	}, {
		[]intt{-5, 42, -42, 42, 3, -18},
		[]intt{-42, 18, -5, 3},
		[]bool{true, false, true, true},
		[]intt{-18, 42},
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		for i, ii := range ti.delete {
			if deleted := tree.Delete(ii); deleted != ti.deleted[i] {
				t.Error("invalid deleted result", ti.delete, ii, deleted)
			}
		}

		iter := tree.Iterate(nil)
		for _, i := range ti.check {
			if v, ok := iter.Next(); !ok || v.(intt) != i {
				t.Error("failed to retrieve", i)
			}
		}

		if v, ok := iter.Next(); ok {
			t.Error("unexpected value", v.(intt))
		}
	}
}

func TestIterate(t *testing.T) {
	for _, ti := range []struct {
		init   []intt
		from   Value
		expect []intt
	}{{
		nil,
		nil,
		nil,
	}, {
		nil,
		intt(42),
		nil,
	}, {
		[]intt{42},
		nil,
		[]intt{42},
	}, {
		[]intt{42},
		intt(42),
		nil,
	}, {
		[]intt{42, -5},
		nil,
		[]intt{-5, 42},
	}, {
		[]intt{42, -5},
		intt(42),
		nil,
	}, {
		[]intt{42, -5},
		intt(-5),
		[]intt{42},
	}, {
		[]intt{42, -5},
		intt(3),
		[]intt{42},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		nil,
		[]intt{-18, -5, 3, 42, 81},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(-42),
		[]intt{-18, -5, 3, 42, 81},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(-18),
		[]intt{-5, 3, 42, 81},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(3),
		[]intt{42, 81},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(5),
		[]intt{42, 81},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(81),
		nil,
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(128),
		nil,
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		iter := tree.Iterate(ti.from)
		i := 0
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}

			if len(ti.expect) <= i || v != ti.expect[i] {
				var e Value
				if len(ti.expect) > i {
					e = ti.expect[i]
				}

				t.Error("failed to return the right value", v, e)
			}

			i++
		}
	}
}

func TestReverse(t *testing.T) {
	for _, ti := range []struct {
		init   []intt
		from   Value
		expect []intt
	}{{
		nil,
		nil,
		nil,
	}, {
		nil,
		intt(42),
		nil,
	}, {
		[]intt{42},
		nil,
		[]intt{42},
	}, {
		[]intt{42},
		intt(42),
		nil,
	}, {
		[]intt{42, -5},
		nil,
		[]intt{42, -5},
	}, {
		[]intt{42, -5},
		intt(-5),
		nil,
	}, {
		[]intt{42, -5},
		intt(-5),
		nil,
	}, {
		[]intt{42, -5},
		intt(3),
		[]intt{-5},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		nil,
		[]intt{81, 42, 3, -5, -18},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(-42),
		nil,
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(-18),
		nil,
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(3),
		[]intt{-5, -18},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(5),
		[]intt{3, -5, -18},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(81),
		[]intt{42, 3, -5, -18},
	}, {
		[]intt{-18, 42, -5, 3, 81},
		intt(128),
		[]intt{81, 42, 3, -5, -18},
	}} {
		tree := new(Tree)
		for _, i := range ti.init {
			tree.Insert(i)
		}

		iter := tree.Reverse(ti.from)
		i := 0
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}

			if len(ti.expect) <= i || v != ti.expect[i] {
				var e Value
				if len(ti.expect) > i {
					e = ti.expect[i]
				}

				t.Error("failed to return the right value", v, e)
			}

			i++
		}
	}
}
