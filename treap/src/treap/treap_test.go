package treap

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
		a[i] = int64(i)
	}
	return a
}

const _N = 1000000

func TestInsertDelete(t *testing.T) {
	a := makeValues(_N)
	tr := New()
	for _, v := range a {
		tr.Insert(KeyT(v))
	}
	if !tr.verifyBSTProperty() || !tr.verifyHeapProperty() {
		t.Error("treap check fail after insert")
	}
	t.Log(_N, tr.height())
	randomShuffle(a)
	for i := 0; i < _N/2; i++ {
		tr.Delete(KeyT(a[i]))
	}
	if !tr.verifyBSTProperty() || !tr.verifyHeapProperty() {
		t.Error("treap check fail after delete")
	}
	t.Log(_N/2, tr.height())
}
