package main

import "fmt"

func max_subarray(a []int) (low, high int, sm int) {
    n := len(a)
    if n == 1 {
        return 0, 1, a[0]
    }
    m := len(a) / 2
    left_low, left_high, left_sum := max_subarray(a[:m])
    right_low, right_high, right_sum := max_subarray(a[m:])
    cross_low, cross_high, cross_sum := m, m + 1, a[m]
    s := cross_sum
    for i := m; i > 0; i-- {
        s += a[i - 1]
        if s > cross_sum {
            cross_sum = s
            cross_low = i - 1
        }
    }
    for i := m + 1; i < n; i++ {
        s += a[i]
        if s > cross_sum {
            cross_sum = s
            cross_high = i + 1
        }
    }
    if left_sum > right_sum && left_sum > cross_sum {
        return left_low, left_high, left_sum
    } else if right_sum > left_sum && right_sum > cross_sum {
        return m + right_low, m + right_high, right_sum
    } else {
        return cross_low, cross_high, cross_sum
    }
}

type ending_at struct {
    index int
    sum int
}

func max_subarray_linear(a []int) (lo, hi, sum int) {
    if len(a) == 0 {
        return 0, 0, 0
    } else if len(a) == 1 {
        return 0, 1, a[0]
    } else {
        n := len(a)
        ends := make([]ending_at, n)
        ends[0].index = 0
        ends[0].sum = a[0]
        for i := 1; i < n; i++ {
            if ends[i-1].sum > 0 {
                ends[i].index = ends[i-1].index
                ends[i].sum = a[i] + ends[i-1].sum
            } else {
                ends[i].index = i
                ends[i].sum = a[i]
            }
        }
        max := 0
        for i := 1; i < n; i++ {
            if ends[i].sum > ends[max].sum {
                max = i
            }
        }
        return ends[max].index, max + 1, ends[max].sum
    }
}

func main() {
    a := []int{12, -3, -25, 20, -3, -16, -23, 18, 20, -7, 12, -5 -22, 15, -4, 7}
    fmt.Println(max_subarray(a))
    fmt.Println(max_subarray_linear(a))
}
