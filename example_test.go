package sortb_test

import (
	"fmt"

	"github.com/aryszka/sortb"
)

type intt int

func (i intt) Less(j sortb.Value) bool  { return i < j.(intt) }
func (i intt) Equal(j sortb.Value) bool { return i == j.(intt) }

func Example() {
	t := new(sortb.Tree)
	for _, i := range []intt{-2, 5, 3, 0, 3, -1} {
		t.Insert(i)
	}

	iter := t.Iterate(nil)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		fmt.Println(v)
	}

	// Output:
	// -2
	// -1
	// 0
	// 3
	// 5
}
