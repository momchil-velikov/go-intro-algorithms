package llrb

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

func makeSeqValues(n uint) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = int64(i)
	}
	return a
}

func makeValues(n uint, low, high int) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = int64(low + rand.Intn(high-low))
	}
	return randomShuffle(a)
}

const _R = 10
const _N = 10000
const _NN = 1000000

func TestInsert(t *testing.T) {
	a := makeValues(_N, -_N, _N)
	tr := &Tree{}
	m := make(map[KeyT]bool)
	for i, k := range a {
		m[KeyT(k)] = true
		tr.Insert(KeyT(k))
		if !tr.isBST() {
			t.Errorf("tree is not BST after insert #%d, %d", i, k)
			t.Log(tr.WriteDot())
		}
	}
	b := makeValues(2*_N, -_N, _N)
	for _, k := range b {
		if _, ok := m[KeyT(k)]; ok != tr.Find(KeyT(k)) {
			t.Error("key %d is erroneously found or not found", k)
		}
	}
}

func TestInsertStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping insert stress test")
	}
	for r := 0; r < _R; r++ {
		a := makeValues(_NN, -_NN, _NN)
		tr := &Tree{}
		m := make(map[KeyT]bool)
		for _, k := range a {
			m[KeyT(k)] = true
			tr.Insert(KeyT(k))
		}
		b := makeValues(2*_NN, -_NN, _NN)
		for _, k := range b {
			if _, ok := m[KeyT(k)]; ok != tr.Find(KeyT(k)) {
				t.Error("key %d is erroneously found or not found", k)
			}
		}
	}
}

func TestIter(t *testing.T) {
	a := makeValues(_N, -_N, _N)
	tr := &Tree{}
	m := make(map[KeyT]uint)
	for _, k := range a {
		tr.Insert(KeyT(k))
		m[KeyT(k)]++
	}
	var b []KeyT = nil
	for i := tr.Min(); !i.IsNil(); i = i.Succ() {
		b = append(b, i.Key())
	}
	for i := 1; i < len(b); i++ {
		if b[i-1] > b[i] {
			t.Error("invalid order in tree")
		}
	}
	b = nil
	for i := tr.Max(); !i.IsNil(); i = i.Pred() {
		b = append(b, i.Key())
	}
	for i := 1; i < len(b); i++ {
		if b[i-1] < b[i] {
			t.Error("invalid order in tree")
		}
	}
	c := makeValues(2*_NN, -_NN, _NN)
	for _, k := range c {
		if m[KeyT(k)] != tr.Count(KeyT(k)) {
			t.Error("Wrong count of key", k)
		}
	}
}

func TestDeleteMin(t *testing.T) {
	a := makeSeqValues(_N)
	tr := Tree{}
	for _, k := range a {
		tr.Insert(KeyT(k))
	}
	// t.Log(tr.WriteDot())
	// tr.DeleteMin()

	for _, k := range a {
		tr.DeleteMin()
		if tr.Find(KeyT(k)) {
			t.Errorf("failed to remove %d from tree", k)
		}
		if !tr.isBST() {
			t.Errorf("tree is not BST after deleteMin %d", k)
		}
	}
}
