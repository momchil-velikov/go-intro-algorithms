package main

import (
	"fmt"
)

var p []uint = []uint{0, 1, 5, 8, 9, 10, 17, 17, 20, 24, 26}

func cuts(n uint) (uint, []uint) {
	type sol struct {
		cost uint
		cuts []uint
	}
	c := make([]sol, n+1)
	for i := uint(1); i <= n; i++ {
		m := p[i]
		cuts := []uint(nil)
		for j := uint(1); j < i; j++ {
			if m < p[j]+c[i-j].cost {
				m = p[j] + c[i-j].cost
				cuts = append(c[i-j].cuts, j)
			}
		}
		c[i].cost = m
		s := uint(0)
		for j := 0; j < len(cuts); j++ {
			s += cuts[j]
		}
		if i > s {
			c[i].cuts = append(cuts, i-s)
		} else {
			c[i].cuts = cuts
		}
	}
	return c[n].cost, c[n].cuts
}

func main() {
	for k := uint(1); k <= 10; k++ {
		c, v := cuts(k)
		fmt.Println(k, "->", c, ":", v)
	}
}
