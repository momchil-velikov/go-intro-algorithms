package main

import (
	"bytes"
	"fmt"
)

type Matrix struct {
	nr, nc uint
	data   []uint
}

func NewMatrix(n, m uint) *Matrix {
	return &Matrix{nr: n, nc: m, data: make([]uint, n*m)}
}

type IndexerType func(uint, uint) *uint

func (m *Matrix) Indexer() func(uint, uint) *uint {
	return func(i, j uint) *uint {
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

// Matrices A_0, A_1, ..., A_n-1
// Matrix A_i has dimensions p[i]xp[i+1]
var p []uint = []uint{30, 35, 15, 5, 10, 20, 25}

func parens(i, j uint, s IndexerType) string {
	if i == j {
		return fmt.Sprintf("A%d", i)
	}
	k := *s(i, j)
	return fmt.Sprintf("(%s%s)", parens(i, k, s), parens(k+1, j, s))
}

func matrixChainOpt(p []uint) uint /*[]uint */ {
	n := uint(len(p)) - 1
	mat := NewMatrix(n, n)
	idx := NewMatrix(n, n)
	c := mat.Indexer()
	s := idx.Indexer()
	for l := uint(2); l <= n; l++ {
		for i := uint(0); i < n-l+1; i++ {
			j := i + l - 1
			m := uint(10000000000)
			mk := uint(0)
			for k := i; k < j; k++ {
				if m > *c(i, k)+*c(k+1, j)+p[i]*p[k+1]*p[j+1] {
					m = *c(i, k) + *c(k+1, j) + p[i]*p[k+1]*p[j+1]
					mk = k
				}
			}
			*c(i, j) = m
			*s(i, j) = mk
		}
	}
	fmt.Println(mat, idx)
	fmt.Println(parens(0, n-1, s))
	return *c(0, n-1)
}

func main() {
	fmt.Println(matrixChainOpt(p))
}
