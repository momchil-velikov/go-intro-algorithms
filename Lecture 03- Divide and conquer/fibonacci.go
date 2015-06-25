package main

import "fmt"

func fibonacci_rec(n uint) uint {
    if n <= 1 {
        return 1
    }
    return fibonacci_rec(n - 2) + fibonacci_rec(n - 1)
}

func fibonacci_rec1(n uint) (uint, uint) {
    if n <= 1 {
        return 1, 1
    }
    a, b := fibonacci_rec1(n - 1)
    return b, a + b
}

func fibonacci_iter(n uint) uint {
    var a, b uint = 1, 1
    for i := uint(1); i < n; i++ {
        a, b = b, a + b
    }
    return b
}

type mat2 struct {
    a, b, c, d uint
}

func (m mat2) String() string {
    return fmt.Sprintf("[%d %d %d %d]", m.a, m.b, m.c, m.d)
}

func mul(a, b mat2) (r mat2) {
    r.a = a.a*b.a + a.b*b.c
    r.b = a.a*b.b + a.b*b.d
    r.c = a.c*b.a + a.d*b.c
    r.d = a.c*b.b + a.d*b.d
    return
}

func power(a mat2, n uint) mat2 {
    switch {
    case n == 1:
        return a
    case (n & 1) == 0:
        y := power(a, n / 2)
        return mul(y, y)
    default:
        y := power(a, (n - 1) / 2)
        return mul(y, mul(y, a))
    }
}

func fibonacci_mat(n uint) uint {
    if n < 2 {
        return 1
    }
    a := power(mat2{1, 1, 1, 0}, n)
    return a.a
}

func main() {
    n := []uint{0, 1, 2, 3, 4, 5, 6, 8, 20, 30, 38 }
    for _, v := range n {
        fmt.Println(v, fibonacci_rec(v))
    }
    for _, v := range n {
        _, k := fibonacci_rec1(v)
        fmt.Println(v, k)
    }
    for _, v := range n {
        fmt.Println(v, fibonacci_iter(v))
    }
    for _, v := range n {
        fmt.Println(v, fibonacci_mat(v))
    }
    fmt.Println(fibonacci_iter(300000000))
    fmt.Println(fibonacci_mat (300000000))
}
