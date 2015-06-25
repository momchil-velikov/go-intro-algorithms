package rbt

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

const _N = 1000000

func TestInsert(t *testing.T) {
    a := makeValues(_N)
    tr := New()
    if !tr.verifyTreeBSTProperty() {
        t.Error("empty tree fails BST check")
    }
    for _, v := range a {
        tr.Insert(KeyType(v))
    }
    if !tr.verifyTreeBSTProperty() {
        t.Errorf("Invalid tree: not BST")
    }
    if !tr.validateTreeConnectivity() {
        t.Error("Invalid tree: connectivity")
    }
    if !tr.validateTreeBlackHeight() {
        t.Error("Invalid tree: balance")
    }
    if !tr.validateTreeRedChildren() {
        t.Error("Invalid tree: adjacent red nodes")
    }
    tr.Delete(KeyType(_N))
    randomShuffle(a)
    for i, v := range a {
        tr.Delete(KeyType(v))
        if (i % 5000) == 0 {
            if !tr.verifyTreeBSTProperty() {
                t.Errorf("Invalid tree: not BST")
            }
            if !tr.validateTreeConnectivity() {
                t.Error("Invalid tree: connectivity")
            }
            if !tr.validateTreeBlackHeight() {
                t.Error("Invalid tree: balance")
            }
            if !tr.validateTreeRedChildren() {
                t.Error("Invalid tree: adjacent red nodes")
            }
        }
    }
}
