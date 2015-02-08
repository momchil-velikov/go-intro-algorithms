package main

import "fmt"

type Matrix struct {
    n uint
    v []float64
}

func (m *Matrix) Init(n uint) {
    m.n = n
    m.v = make([]float64, n * n)
}

func (m *Matrix) At(i, j uint) float64 {
    return m.v[i*m.n + j]
}

func (m *Matrix) Set(i, j uint, v float64) {
    m.v[i*m.n + j] = v
}

func (a *Matrix) MulBasic(b *Matrix) (c Matrix) {
    n := a.n
    c.Init(n)
    for i := uint(0); i < n; i++ {
        for j := uint(0); j < n; j++ {
            v := 0.0
            for k := uint(0); k < n; k++ {
                v += a.At(i, k) * b.At(k, j)
            }
            c.Set(i, j, v)
        }
    }
    return
}


// cba to implement Strassen's algorithm :(

func main() {
    var a, b Matrix
    a.Init(2)
    a.Set(0, 0, 1)
    a.Set(0, 1, 1)
    a.Set(1, 0, 1)
    a.Set(1, 1, 0)

    b.Init(2)
    b.Set(0, 0, 1)
    b.Set(0, 1, 1)
    b.Set(1, 0, 1)
    b.Set(1, 1, 0)

    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)
    a = a.MulBasic(&b)

    fmt.Println(a)
}
