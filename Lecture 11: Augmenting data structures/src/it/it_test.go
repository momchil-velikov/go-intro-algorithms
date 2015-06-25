package it

import (
	"math/rand"
	"testing"
)

func randomShuffle(a []int64) {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
	}
}

func randomShuffle2(a, b []int64) {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
		b[i-1], b[j] = b[j], b[i-1]
	}
}

func makeLow(n uint) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = int64(i) * 10
	}
	randomShuffle(a)
	return a
}

func makeSize(n uint) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = rand.Int63n(20)
	}
	return a
}

const _N = 1000

func TestInsertDelete(t *testing.T) {
	low := makeLow(_N)
	sz := makeSize(_N)
	tr := New()
	for i := range low {
		tr.Insert(KeyType(low[i]), KeyType(low[i]+sz[i]))
		if !tr.verifyTreeUpper() {
			t.Errorf("Invalid tree upper, after insert #%d:[%d, %d]", i, low[i], low[i]+sz[i])
		}
	}
	randomShuffle2(low, sz)
	for i := range low {
		tr.Delete(KeyType(low[i]), KeyType(low[i]+sz[i]))
		if !tr.verifyTreeUpper() {
			t.Errorf("Invalid tree upper, after delete #%d:[%d, %d]", i, low[i], low[i]+sz[i])
		}
	}
}

type interval struct {
	low, high KeyType
}

func randomShuffleIntervals(a []interval) {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
	}
}

func TestSearch(t *testing.T) {
	ints := []interval{{0, 3}, {5, 8}, {6, 10}, {8, 9}, {15, 23}, {16, 21}, {17, 19}, {25, 30}, {26, 26}}
	randomShuffleIntervals(ints)
	tr := New()
	for i := range ints {
		tr.Insert(ints[i].low, ints[i].high)
	}
	args := []interval{{10, 10}, {11, 12}, {11, 16}}
	res := []interval{{6, 10}, {0, 0}, {15, 23}}
	for i := range args {
		s, e, ok := tr.Search(args[i].low, args[i].high)
		if res[i].low == 0 && res[i].high == 0 {
			if ok {
				t.Errorf("Found [%d, %d] with arg [%d, %d]", s, e, args[i].low, args[i].high)
			}
		} else {
			if !ok {
				t.Errorf("Not found: [%d, %d]", args[i].low, args[i].high)
			} else if s != res[i].low || e != res[i].high {
				t.Errorf(
					"Found [%d, %d] instead of [%d, %d] with args [%d, %d]",
					s, e, res[i].low, res[i].high, args[i].low, args[i].high,
				)
			}
		}
	}
}
