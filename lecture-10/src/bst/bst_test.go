package bst

import (
    "math/rand"
    "testing"
)

func randomShuffle(a []int64) []int64 {
    for i := len(a); i > 1; i-- {
        j := rand.Intn(i)
        a[i - 1], a[j] = a[j], a[i - 1]
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

func TestInsert(t *testing.T) {
    a := makeValues(_N)
    tr := New()
    if !tr.verifyBSTProperty() {
        t.Error("empty tree fails BST check")
    }
    for _, v := range a {
        tr.Insert(KeyType(v))
        if !tr.verifyBSTProperty() {
            t.Errorf("BST check fail after inserting key %d\n", v)
        }
    }
    tr.Delete(KeyType(_N))
    randomShuffle(a)
    for _, v := range a {
        tr.Delete(KeyType(v))
        if !tr.verifyBSTProperty() {
            t.Errorf("BST check fail after deleting key %d\n", v)
        }
    }
}
