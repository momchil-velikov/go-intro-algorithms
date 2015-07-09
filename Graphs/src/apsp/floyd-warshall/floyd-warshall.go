package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

type IncidenceMatrixElt struct {
	dist float64
	prev int
}

type IncidenceMatrix struct {
	N int
	d []IncidenceMatrixElt
}

func (m *IncidenceMatrix) Init(n int) *IncidenceMatrix {
	m.N = n
	m.d = make([]IncidenceMatrixElt, n*n)
	for i := 0; i < n*n; i++ {
		m.d[i].dist = math.Inf(0)
		m.d[i].prev = -1
	}
	for i := 0; i < n; i++ {
		m.d[i*n+i].dist = 0.0
		m.d[i*n+i].prev = i
	}
	return m
}

func (m *IncidenceMatrix) Index() func(int, int) *IncidenceMatrixElt {
	return func(i, j int) *IncidenceMatrixElt {
		return &m.d[i*m.N+j]
	}
}

func (m *IncidenceMatrix) Copy() *IncidenceMatrix {
	r := new(IncidenceMatrix)
	r.N = m.N
	r.d = make([]IncidenceMatrixElt, len(m.d))
	copy(r.d, m.d)
	return r
}

func (m *IncidenceMatrix) Move() *IncidenceMatrix {
	r := new(IncidenceMatrix)
	r.N = m.N
	r.d = m.d
	m.N = 0
	m.d = nil
	return r
}

func (a *IncidenceMatrix) Swap(b *IncidenceMatrix) {
	a.N, b.N = b.N, a.N
	a.d, b.d = b.d, a.d
}

func readGraph(in io.Reader) (*IncidenceMatrix, error) {
	nodeCount, edgeCount := 0, 0
	if _, err := fmt.Fscan(in, &nodeCount); err != nil {
		return nil, err
	}
	if _, err := fmt.Fscan(in, &edgeCount); err != nil {
		return nil, err
	}
	a := new(IncidenceMatrix).Init(nodeCount)
	m := a.Index()
	s, d, w := 0, 0, float64(0.0)
	_, err := fmt.Fscan(in, &s, &d, &w)
	for err == nil {
		e := m(s, d)
		e.dist = w
		e.prev = s
		_, err = fmt.Fscan(in, &s, &d, &w)
	}
	if err == nil || err == io.EOF {
		return a, nil
	} else {
		return nil, err
	}
}

func apspFloydWarshall(m *IncidenceMatrix) *IncidenceMatrix {
	a := m
	b := a.Copy()
	n := a.N
	for k := 0; k < n; k++ {
		ma, mb := a.Index(), b.Index()
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				x, y := ma(i, k), ma(k, j)
				p, q := ma(i, j), mb(i, j)
				d := x.dist + y.dist
				if d < p.dist {
					q.dist = d
					q.prev = y.prev
				} else {
					q.dist = p.dist
					q.prev = p.prev
				}
			}
		}
		a, b = b, a
	}
	return a
}

func main() {
	var in io.Reader
	var filename string
	if len(os.Args) < 2 {
		filename = "stdin"
		in = os.Stdin
	} else {
		filename = os.Args[1]
		if f, err := os.Open(filename); err != nil {
			fmt.Println(filename, ":", err)
			return
		} else {
			defer f.Close()
			in = bufio.NewReader(f)
		}
	}
	a, err := readGraph(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename, err)
	}
	a = apspFloydWarshall(a)
	m := a.Index()
	for i := 0; i < a.N; i++ {
		for j := 0; j < a.N; j++ {
			d := m(i, j)
			if math.IsInf(d.dist, 0) {
				fmt.Printf("%d to %d(\u221e)\n", i, j)
			} else {
				fmt.Printf("%d to %d(%5.2f): ", i, j, d.dist)
				q := j
				for q != i {
					x := m(i, q)
					fmt.Printf("%d->%d(%5.2f) ", q, x.prev, x.dist)
					q = x.prev
				}
				fmt.Println()
			}
		}
	}
}
