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

func lps(a string) string {
	n := uint(len(a))
	mat := NewMatrix(n, n)
	s := mat.Indexer()
	for l := uint(1); l <= n; l++ {
		for i := uint(0); i < n-l+1; i++ {
			j := i + l - 1
			if i == j {
				*s(i, j) = string(a[i])
			} else if a[i] == a[j] {
				if i+1 == j {
					*s(i, j) = string(a[i]) + string(a[i])
				} else {
					*s(i, j) = string(a[i]) + *s(i+1, j-1) + string(a[i])
				}
			} else {
				s1 := *s(i+1, j)
				s2 := *s(i, j-1)
				if len(s1) > len(s2) {
					*s(i, j) = s1
				} else {
					*s(i, j) = s2
				}
			}
		}
	}
	return *s(0, n-1)
}

func main() {
	fmt.Println(lps("character"))
}
