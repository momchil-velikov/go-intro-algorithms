package ost

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

const _N = 10000

func TestInsertDelete(t *testing.T) {
	a := makeValues(_N)
	tr := New()
	for i, v := range a {
		tr.Insert(KeyType(v))
		if !tr.verifyTreeSize() {
			t.Errorf("Invalid tree size, after insert [%d] = %d", i, v)
		}
	}
	tr.Delete(KeyType(_N))
	randomShuffle(a)
	for i, v := range a {
		tr.Delete(KeyType(v))
		if !tr.verifyTreeSize() {
			t.Errorf("Invalid tree size, after delete [%d] = %d", i, v)
		}
	}
}

func TestSelect(t *testing.T) {
	a := makeValues(_N)
	tr := New()
	for _, v := range a {
		tr.Insert(KeyType(v))
	}
	for i := uint(1); i <= _N; i++ {
		v, ok := tr.Select(i)
		if !ok {
			t.Errorf("Select fail at rank %d", i)
		}
		if v != KeyType(i-1) {
			t.Errorf("Incorect rank %d, value %d", i, v)
		}
	}
}
