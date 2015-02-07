package main

import "fmt"

func insertion_sort(a []int) {
    n := len(a)
    if n <= 1 {
        return;
    }
    for i := 0; i < n - 1; i++ {
        j := i + 1
        key := a[j]
        for ;j > 0 && key < a[j-1]; j-- {
            a[j] = a[j - 1]
        }
        a[j] = key
    }
}

func main() {
    a := []int{3,4,178,642,85,3,45,78,5,976,34,51,65,761,414,1,87,612,53,14,550}
    insertion_sort(a)
    fmt.Println(a)
}
