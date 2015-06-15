package main

import (
	"bytes"
	"fmt"
)

type Matrix struct {
	nr, nc uint
	data   []string
}

func NewMatrix(n, m uint) *Matrix {
	return &Matrix{nr: n, nc: m, data: make([]string, n*m)}
}

func (m *Matrix) Indexer() func(uint, uint) *string {
	return func(i, j uint) *string {
		return &m.data[i*m.nc+j]
	}
}

func (mat *Matrix) String() string {
	m := mat.Indexer()
	buf := bytes.Buffer{}
	for i := uint(0); i < mat.nr; i++ {
		buf.WriteString("[")
		for j := uint(0); j+1 < mat.nc; j++ {
			buf.WriteString(fmt.Sprintf("%4d, ", *m(i, j)))
		}
		buf.WriteString(fmt.Sprintf("%4d]\n", *m(i, mat.nc-1)))
	}
	return buf.String()
}

func lcs(a, b string) string {
	n, m := uint(len(a)), uint(len(b))
	mat := NewMatrix(n+1, m+1)
	C := mat.Indexer()
	for i := uint(1); i < n+1; i++ {
		for j := uint(1); j < m+1; j++ {
			if a[i-1] == b[j-1] {
				*C(i, j) = *C(i-1, j-1) + string(a[i-1])
			} else {
				c, d := *C(i, j-1), *C(i-1, j)
				if len(c) < len(d) {
					*C(i, j) = d
				} else {
					*C(i, j) = c
				}
			}
		}
	}
	return *C(n, m)
}

func main() {
	fmt.Println(lcs("abcbdab", "bdcaba"))
	fmt.Println(lcs("accggtcgagtgcgcggaagccggccgaa", "gtcgttcggaatgccgttgctctgtaaa"))
}
