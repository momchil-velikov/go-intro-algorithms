package main

import "fmt"

func binary_search(k int, a []int) int {
    i, j := 0, len(a)
    for i < j {
        m := (i + j) / 2
        switch {
        case k == a[m]:
            return m
        case k < a[m]:
            j = m
        default:
            i = m + 1
        }
    }
    return -1
}

// Pre:  i = lower_bound(k, a)
// Post: a[i:] >= k
func lower_bound(k int, a []int) int {
    i, j := 0, len(a)
    for i < j {
        m := (i + j) / 2
        if k <= a[m] {
            j = m
        } else {
            i = m + 1
        }
    }
    return i
}

// Pre:  i = upper_bound(k, a)
// Post: a[:i] < k
func upper_bound(k int, a []int) int {
    i, j := 0, len(a)
    for i < j {
        m := (i + j) / 2
        if k < a[m] {
            j = m
        } else {
            i = m + 1
        }
    }
    return i
}

func main() {
    f := upper_bound
    a := []int{
        1,    1,   1,   1,   3 ,  3,   3,   4,   5,  14,
        34,  45,  45,  45,  51,  53,  65,  78,  78,  85,
        87, 178, 414, 550, 612, 612, 612, 642, 761, 976}
    for _, v := range a {
        fmt.Println(v, f(v, a))
    }

    b := []int{0, 2, 46, 980}
    for _, v := range b {
        fmt.Println(v, f(v, a))
    }
}
