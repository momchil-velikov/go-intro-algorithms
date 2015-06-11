package skip

import (
	"math/rand"
	"testing"
)

func randomShuffle(a []int64) []int64 {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
	}
	return a
}

func makeValues(n uint) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = int64(i) + 1
	}
	return randomShuffle(a)
}

const _N = 1000000

func TestInsert(t *testing.T) {
	a := makeValues(_N)
	lst := New()
	for _, v := range a {
		lst.Insert(KeyT(v))
	}
	if !lst.VerifyOrder() {
		t.Error("order verification error after insert")
	}
	lst.Delete(KeyT(_N))
	randomShuffle(a)
	for i := 0; i < _N-50; i++ {
		lst.Delete(KeyT(a[i]))
	}
	if !lst.VerifyOrder() {
		t.Error("order verification error after delete")
	}
	// t.Log(lst.WriteDot())
}
