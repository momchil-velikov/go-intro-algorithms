package main

import "fmt"
import "math"

func power(x float64, n int) float64 {
    switch {
    case n == 0:
        return 1.0
    case n == 1:
        return x
    case (n & 1) == 0:
        x = power(x, n / 2)
        return x * x
    default:
        y := power(x, (n - 1) / 2)
        return y * y * x
    }
}

func main() {
    xs := []float64{1.1001, 3.1415, 2.7165}
    ns := []int{0, 1, 2, 3, 4, 5}

    for _, x := range xs {
        for _, n := range ns {
            fmt.Println(x, "**", n, power(x, n), math.Pow(x, float64(n)))
        }
    }
}
