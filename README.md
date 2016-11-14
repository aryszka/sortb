[![GoDoc](https://godoc.org/github.com/aryszka/sortb?status.svg)](https://godoc.org/github.com/aryszka/sortb)
[![Go Report Card](https://goreportcard.com/badge/github.com/aryszka/sortb)](https://goreportcard.com/report/github.com/aryszka/sortb)
[![Coverage](http://gocover.io/_badge/github.com/aryszka/sortb)](http://gocover.io/github.com/aryszka/sortb)

# sortb

Sortb, a Go package, provides a data structure to store sorted generic values
using a balanced binary search tree.

Documentation: https://godoc.org/github.com/aryszka/sortb

Example:

```
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
```
