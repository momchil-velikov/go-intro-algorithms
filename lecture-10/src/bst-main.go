package main

import (
    "bst"
    "math/rand"
    "os"
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

const _N = 1000

func main() {
    a := makeValues(_N)
    tr := bst.New()
    for _, v := range a {
        tr.Insert(bst.KeyType(v))
    }
    f, err := os.Create("bst-1.dot")
    if err == nil {
        f.WriteString(tr.WriteDot("PhaseOne"))
        f.Close()
    }
    randomShuffle(a)
    for i := 0; i < 3 * len(a) / 4; i++ {
        tr.Delete(bst.KeyType(a[i]))
    }
    f, err = os.Create("bst-2.dot")
    if err == nil {
        f.WriteString(tr.WriteDot("PhaseTwo"))
        f.Close()
    }
}
