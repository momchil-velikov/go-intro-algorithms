package heap

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
	return randomShuffle(a)
}

const _N = 20000

func TestHeap(t *testing.T) {
	a := makeValues(_N)
	b := makeValues(_N)
	var h IntHeap
	if !h.Check() {
		t.Error("violation of the heap property")
	}
	for i := range a {
		h.Push(a[i] + b[i])
		if !h.Check() {
			t.Error("violation of the heap property")
		}
	}
	i := 0
	for h.Len() > 0 {
		a[i] = h.Pop()
		i++
		if !h.Check() {
			t.Error("violation of the heap property")
		}
	}
	for i = 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			t.Errorf("output array not sorted")
		}
	}
}

const _NN = 1000000

func TestInit(t *testing.T) {
	a := makeValues(_NN)
	var h IntHeap
	h.Init(a)
	if !h.Check() {
		t.Error("violation of the heap property: ", h.Data)
	}
}

func TestSort(t *testing.T) {
	a := makeValues(_NN)
	b := makeValues(_NN)
	for i := range a {
		a[i] = (a[i] - b[i])

	}
	Sort(a)
	a[0] = -a[0]
	for i := 1; i < len(a); i++ {
		a[i] = -a[i]
		if a[i-1] > a[i] {
			t.Fatalf("output array not sorted")
		}
	}
}
