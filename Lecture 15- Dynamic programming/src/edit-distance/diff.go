package main

import (
	"bytes"
	"fmt"
)

const (
	Nop = iota
	Copy
	Insert
	Replace
	Delete
)

type Edit struct {
	op       uint
	src, dst uint
	ch       uint8
}

type Sol struct {
	cost uint
	seq  []Edit
}

type Matrix struct {
	nr, nc uint
	data   []Sol
}

func NewMatrix(n, m uint) *Matrix {
	return &Matrix{nr: n, nc: m, data: make([]Sol, n*m)}
}

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

func diff(x, y string) []Edit {
	n, m := uint(len(x)), uint(len(y))
	mat := NewMatrix(n+1, m+1)
	e := mat.Indexer()
	for j := uint(1); j <= m; j++ {
		e(0, j).cost = e(0, j-1).cost + 1
		e(0, j).seq = append(e(0, j-1).seq, Edit{op: Insert, src: 0, dst: j - 1, ch: y[j-1]})
	}
	for i := uint(1); i <= n; i++ {
		e(i, 0).cost = e(i-1, 0).cost + 1
		e(i, 0).seq = append(e(i-1, 0).seq, Edit{op: Delete, src: i - 1, dst: 0, ch: x[i-1]})
	}
	for i := uint(1); i <= n; i++ {
		for j := uint(1); j <= m; j++ {
			if x[i-1] == y[j-1] {
				e(i, j).cost = e(i-1, j-1).cost
				e(i, j).seq = append(e(i-1, j-1).seq, Edit{op: Copy, src: i - 1, dst: j - 1, ch: x[i-1]})
			} else {
				cost := e(i, j-1).cost + 1
				op := append(e(i, j-1).seq, Edit{op: Insert, src: 0, dst: i - 1, ch: y[j-1]})
				if e(i-1, j-1).cost+1 < cost {
					cost = e(i-1, j-1).cost + 1
					op = append(e(i-1, j-1).seq, Edit{op: Replace, src: i - 1, dst: j - 1, ch: y[j-1]})
				}
				if e(i-1, j).cost+1 < cost {
					cost = e(i-1, j).cost + 1
					op = append(e(i-1, j).seq, Edit{op: Delete, src: i - 1, dst: 0, ch: x[i-1]})
				}
				e(i, j).cost = cost
				e(i, j).seq = op
			}
		}
	}
	return e(n, m).seq
}

func main() {
	e := diff("algorithm", "altruistic")
	//e := diff("a", "ab")
	//e := diff("abdc", "adbc")
	fmt.Println(len(e))
	for _, op := range e {
		switch op.op {
		case Nop:
			fmt.Println("nop")
		case Copy:
			fmt.Printf("copy: src[%d] -> dst[%d], %c\n", op.src, op.dst, op.ch)
		case Insert:
			fmt.Printf("insert: dst[%d], %c\n", op.dst, op.ch)
		case Replace:
			fmt.Printf("replace: src[%d] -> dst[%d], %c\n", op.src, op.dst, op.ch)
		case Delete:
			fmt.Printf("delete: src[%d], %c\n", op.src, op.ch)
		}
	}
}
