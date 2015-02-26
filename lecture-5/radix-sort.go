package main

import "fmt"

func counting_sort(m uint, s uint, a []uint, t []uint) []uint {
    k := len(t)
    for i := 0; i < k; i++ {
        t[i] = 0
    }
    n := len(a)
    for i := 0; i < n; i++ {
        v := (a[i] >> s) & m
        t[v]++
    }
    for i := 1; i < k; i++ {
        t[i] += t[i-1]
    }
    b := make([]uint, n)
    for i := n; i > 0; i-- {
        v := (a[i - 1] >> s) & m
        b[t[v] - 1] = a[i - 1]
        t[v]--
    }
    return b
}

func radix_sort(a []uint) []uint {
    t := make([]uint, 256)
    a = counting_sort(255, 0, a, t)
    a = counting_sort(255, 8, a, t)
    return a
}

func main() {
    a := []uint{934,858,739,48905,12,3,7503,4603,84,714,65,16253,200,456,81,365,812,
        56125,81,101}
    b := radix_sort(a)
    fmt.Println(b)
}
