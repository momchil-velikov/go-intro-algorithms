package main

import "fmt"

func counting_sort(k uint, a []uint) []uint {
    c := make([]uint, k)
    n := uint(len(a))
    for i := uint(0); i < n; i++ {
        v := a[i]
        if v < k {
            c[v]++
        }
    }
    for i := uint(1); i < k; i++ {
        c[i] += c[i - 1]
    }
    b := make([]uint, n)
    for i := n; i > 0; i-- {
        v := a[i-1]
        if v < k {
            b[c[v]-1] = v
            c[v]--
        }
    }
    return b
}

func main() {
    a := []uint{9,3,4,87,9,58,7,3,200,45,6,81,36,58,1,25,6,12,5,81,101}
    b := counting_sort(100, a)
    fmt.Println(b)
}
