package main

import "fmt"

func merge(a, b, d []int) {
    i, j, k := 0, 0, 0
    na, nb := len(a), len(b)
    for i < na && j < nb {
        if a[i] < b[j] {
            d[k] = a[i]
            i++
        } else {
            d[k] = b[j]
            j++;
        }
        k++;
    }
    for i < na {
        d[k] = a[i]
        k++
        i++
    }
    for j < nb {
        d[k] = b[j]
        k++
        j++
    }
}

func merge_sort1(a, b []int) {
    n := len(a)
    if n == 1 {
        b[0] = a[0]
        return
    }
    merge_sort1(a[:n/2], b[:n/2])
    merge_sort1(a[n/2:], b[n/2:])
    merge(b[:n/2], b[n/2:], a)
    copy(b,a)
}

func merge_sort(a []int) {
    n := len(a)
    if n <= 1 {
        return
    }
    b := make([]int, n)
    merge_sort1(a[:n/2], b[:n/2])
    merge_sort1(a[n/2:], b[n/2:])
    merge(b[:n/2], b[n/2:], a)
}

func main() {
    a := []int{3,4,178,642,85,3,45,78,5,976,34,51,65,761,414,1,87,612,53,14,550}
    merge_sort(a)
    fmt.Println(a)
}

