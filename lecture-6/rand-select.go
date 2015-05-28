package main

import (
    "fmt"
    "math/rand"
)

func partition(a[] int) int {
    k := rand.Intn(len(a))
    a[0], a[k] = a[k], a[0]
    n := len(a)
    i := 0
    x := a[0]
    for j := 1; j < n; j++ {
        if a[j] <= x {
            i++
            a[i], a[j] = a[j], a[i]
        }
    }
    a[0], a[i] = a[i], a[0]
    return i
}

func rand_select(a []int, k int) int {
    n := len(a)
    if n == 0 {
        return 0
    }
    p := partition(a)
    if p == k {
        return a[k];
    } else if p < k {
        return rand_select(a[p + 1:], k - p - 1)
    } else {
        return rand_select(a[:p], k);
    }
}

func test(k int) {
    a := []int{9,48,7,571,26,5,61,561,25,625,1,6,928,61, 56,3,90,86}
    v := rand_select(a, k)
    fmt.Println("k =", k, "v =", v, "a =", a);
}

func main() {
    for i := 0; i < 18; i++ {
        test(i)
    }
}
