package main

import (
    "fmt"
    "sort"
    "math/rand"
)

type Elt struct {
    x float32
    w float32
}

type ByX []Elt

func (a ByX) Len() int           { return len(a) }
func (a ByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByX) Less(i, j int) bool { return a[i].x < a[j].x }

func partition(a[] Elt) int {
    k := rand.Intn(len(a))
    a[0], a[k] = a[k], a[0]
    n := len(a)
    i := 0
    x := a[0].x
    for j := 1; j < n; j++ {
        if a[j].x <= x {
            i++
            a[i], a[j] = a[j], a[i]
        }
    }
    a[0], a[i] = a[i], a[0]
    return i
}

func weighted_median(a[] Elt) (float32, float32) {
    p := 0
    q := len(a)
    for {
        k := p + partition(a[p:q])
        w := float32(0.0)
        for i := 0; i < k; i++ {
            w += a[i].w
        }
        if w >= 0.5 {
            q = k;
        } else if w + a[k].w < 0.5 {
            p = k + 1;
        } else {
            return a[k].x, a[k].w
        }
    }
}

func makeRandomWeights(n uint) []float32 {
    w := make([]float32, n)
    v := float32(0.0)
    for i := uint(1); i < n; i++ {
        w[i] = rand.Float32() * (1.0 - v)
        v += w[i]
    }
    w[0] = 1.0 - v;
    return w
}

func makeRandomValues(n uint) []float32 {
    a := make([]float32, n)
    for i := uint(0); i < n; i++ {
        a[i] = rand.Float32()
    }
    return a
}

func main() {
    n := uint(20)
    w := makeRandomWeights(n)
    v := makeRandomValues(n)

    a := make([]Elt, n)
    for i := range v {
        a[i].x = v[i]
        a[i].w = w[i]
    }

    b := make([]Elt, len(a))
    copy(b, a)
    sort.Sort(ByX(b))
    u := float32(0)
    for i := range(b) {
        fmt.Printf(".x = %8.6f .w = %8.6f u = %8.6f\n", b[i].x, b[i].w, u)
        u += b[i].w
    }
    fmt.Println("u =", u)

    x, _ := weighted_median(a)
    fmt.Println("weighted median ", x)
}
