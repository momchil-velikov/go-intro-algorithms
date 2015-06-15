package main

import (
	"bytes"
	"fmt"
)

type Sol struct {
	cost, pos uint
}

type Matrix struct {
	nr, nc uint
	data   []Sol
}

func NewMatrix(n, m uint) *Matrix {
	return &Matrix{nr: n, nc: m, data: make([]Sol, n*m)}
}

type IndexerT func(uint, uint) *Sol

func (m *Matrix) Indexer() func(uint, uint) *Sol {
	return func(i, j uint) *Sol {
		return &m.data[i*m.nc+j]
	}
}

func (mat *Matrix) String() string {
	m := mat.Indexer()
	buf := bytes.Buffer{}
	for i := uint(0); i < mat.nr; i++ {
		buf.WriteString("[")
		for j := uint(0); j+1 < mat.nc; j++ {
			buf.WriteString(fmt.Sprintf("%v, ", *m(i, j)))
		}
		buf.WriteString(fmt.Sprintf("%v]\n", *m(i, mat.nc-1)))
	}
	return buf.String()
}

func getBreaks(i, j uint, c IndexerT, br []uint) []uint {
	b := c(i, j).pos
	if b != 0 {
		br = append(br, b)
		br = getBreaks(i, b, c, br)
		br = getBreaks(b+1, j, c, br)
	}
	return br
}

func stringBreak(n uint, L []uint) (uint, []uint) {
	mat := NewMatrix(n, n)
	c := mat.Indexer()
	for l := uint(2); l <= n; l++ {
		for i := uint(0); i < n-l+1; i++ {
			j := i + l - 1
			m := uint(100000000)
			b := -1
			for _, k := range L {
				if i <= k && k < j {
					v := l + c(i, k).cost + c(k+1, j).cost
					if v < m {
						m = v
						b = int(k)
					}
				}
			}
			if b >= 0 {
				c(i, j).cost = m
				c(i, j).pos = uint(b)
			}
		}
	}

	return c(0, n-1).cost, getBreaks(0, n-1, c, nil)
}

func main() {
	L := []uint{1, 7, 9}
	c, b := stringBreak(20, L)
	fmt.Println("cost =", c, "breaks =", b)
}
