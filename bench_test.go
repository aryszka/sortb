package sortb

import (
	"math/rand"
	"testing"
)

func benchmarkInsert(b *testing.B, n int) {
	tree := new(Tree)
	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < n; i++ {
		tree.Insert(intt(2 * rnd.Intn(n)))
	}

	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		tree.Insert(intt(2*rnd.Intn(n) + 1))
	}
}

func benchmarkFind(b *testing.B, n int) {
	tree := new(Tree)
	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < n; i++ {
		tree.Insert(intt(2 * rnd.Intn(n)))
	}

	rnd = rand.New(rand.NewSource(0))
	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		tree.Find(intt(2 * rnd.Intn(n)))
	}
}

func benchmarkDelete(b *testing.B, n int) {
	tree := new(Tree)
	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < n; i++ {
		tree.Insert(intt(2 * rnd.Intn(n)))
	}

	rnd = rand.New(rand.NewSource(0))
	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		tree.Delete(intt(2 * rnd.Intn(n)))
	}
}

func BenchmarkInsert1(b *testing.B)       { benchmarkInsert(b, 1) }
func BenchmarkInsert100(b *testing.B)     { benchmarkInsert(b, 100) }
func BenchmarkInsert10000(b *testing.B)   { benchmarkInsert(b, 10000) }
func BenchmarkInsert1000000(b *testing.B) { benchmarkInsert(b, 1000000) }

func BenchmarkFind1(b *testing.B)       { benchmarkFind(b, 1) }
func BenchmarkFind100(b *testing.B)     { benchmarkFind(b, 100) }
func BenchmarkFind10000(b *testing.B)   { benchmarkFind(b, 10000) }
func BenchmarkFind1000000(b *testing.B) { benchmarkFind(b, 1000000) }

func BenchmarkDelete1(b *testing.B)       { benchmarkDelete(b, 1) }
func BenchmarkDelete100(b *testing.B)     { benchmarkDelete(b, 100) }
func BenchmarkDelete10000(b *testing.B)   { benchmarkDelete(b, 10000) }
func BenchmarkDelete1000000(b *testing.B) { benchmarkDelete(b, 1000000) }
