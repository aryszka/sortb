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
	tree := &Tree{}

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
	sizes := []int{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000}
	if testing.Short() {
		sizes = sizes[:9]
	}

	for i := 0; i < 10; i++ {
		rnd := rand.New(rand.NewSource(int64(i)))
		for _, j := range sizes {
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
		tree := &Tree{}
		for _, i := range ti.insert {
			tree.Insert(i)
		}

		iter := tree.Iterate()
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

func TestDelete(t *testing.T) {
	for _, ti := range []struct {
		init    []intt
		delete  []intt
		deleted []bool
		check   []intt
	}{{
		// 	nil,
		// 	nil,
		// 	nil,
		// 	nil,
		// }, {
		// 	nil,
		// 	[]intt{42},
		// 	[]bool{false},
		// 	nil,
		// }, {
		// 	[]intt{42},
		// 	[]intt{42},
		// 	[]bool{true},
		// 	nil,
		// }, {
		// 	[]intt{42},
		// 	[]intt{-42, -42},
		// 	[]bool{false, false},
		// 	[]intt{42},
		// }, {
		// 	[]intt{42},
		// 	[]intt{42, 42},
		// 	[]bool{true, false},
		// 	nil,
		// }, {
		[]intt{-5, 42, -42, 42, 3, -18},
		[]intt{-42, 18, -5, 3},
		[]bool{true, false, true, true},
		[]intt{-18, 42},
	}} {
		tree := &Tree{}
		for _, i := range ti.init {
			tree.Insert(i)
		}

		for i, ii := range ti.delete {
			if deleted := tree.Delete(ii); deleted != ti.deleted[i] {
				t.Error("invalid deleted result", ti.delete, ii, deleted)
			}
		}

		iter := tree.Iterate()
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
