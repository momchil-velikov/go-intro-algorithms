package main

import (
	"bufio"
	"container/deque"
	"fmt"
	gio "graph/io"
	"io"
	"math"
	"os"
)

type edge struct {
	src       int
	dst       int
	flow, cap float64
}

type node struct {
	id      int
	label   string
	visited bool
	link    *node
	flow    float64
	succ    []*edge
	pred    []*edge
}

type graph struct {
	nodes []node
	edges []*edge
}

func (g *graph) CreateNode(name string) int {
	id := len(g.nodes)
	g.nodes = append(g.nodes, node{id: id, label: name})
	return id
}

func (g *graph) SetNodeAttrs(id int, attrs []string) {
	np := &g.nodes[id]
	for i := 0; i < len(attrs); i += 2 {
		if attrs[i] == "label" {
			np.label = attrs[i+1]
			break
		}
	}
}

func (g *graph) CreateEdge(src int, dst int, attrs []string) {
	sp, dp := &g.nodes[src], &g.nodes[dst]
	e := &edge{src: src, dst: dst}
	for i := 0; i < len(attrs); i += 2 {
		if attrs[i] == "weight" {
			fmt.Sscan(attrs[i+1], &e.cap)
			break
		}
	}
	g.edges = append(g.edges, e)
	sp.succ = append(sp.succ, e)
	dp.pred = append(dp.pred, e)
}

func (g *graph) NodeCount() int {
	return len(g.nodes)
}

func (g *graph) NodeAttrs(id int) []string {
	np := &g.nodes[id]
	if len(np.label) == 0 {
		return nil
	} else {
		return []string{"label", np.label}
	}
}

func (g *graph) EdgeCount(id int) int {
	return len(g.nodes[id].succ)
}

func (g *graph) EdgeAttrs(id int, i int) (dst int, attrs []string) {
	e := g.nodes[id].succ[i]
	return e.dst, []string{
		"label", fmt.Sprintf("%.2f/%.2f", e.flow, e.cap),
		"weight", fmt.Sprint(e.cap),
	}
}

func augPath(g *graph, s, t int) bool {
	q := deque.Deque{}
	g.nodes[s].visited = true
	q.Push(&g.nodes[s])
	for !q.IsEmpty() {
		sp := q.Pop().(*node)
		if sp.id == t {
			return true
		}
		for i := range sp.succ {
			e := sp.succ[i]
			if e.flow < e.cap {
				dp := &g.nodes[e.dst]
				if !dp.visited {
					dp.link = sp
					dp.flow = e.cap - e.flow
					dp.visited = true
					q.Push(dp)
				}
			}
		}
		for i := range sp.pred {
			e := sp.pred[i]
			if e.flow > 0 {
				dp := &g.nodes[e.src]
				if !dp.visited {
					dp.link = sp
					dp.flow = e.flow
					dp.visited = true
					q.Push(dp)
				}
			}
		}
	}
	return false
}

func findEdge(s, d *node) *edge {
	for i := range s.succ {
		if s.succ[i].dst == d.id {
			return s.succ[i]
		}
	}
	return nil
}

func edmondsKarp(g *graph, s, t int) {
	for augPath(g, s, t) {
		v := &g.nodes[t]
		u := v.link
		v.link = nil
		f := math.Inf(0)
		for u != nil {
			if v.flow < f {
				f = v.flow
			}
			u, v, u.link = u.link, u, v
		}
		u = &g.nodes[s]
		v = u.link
		for v != nil {
			if e := findEdge(u, v); e == nil {
				e = findEdge(v, u)
				e.flow -= f
			} else {
				e.flow += f
			}
			u, v = v, v.link
		}
		for i := range g.nodes {
			g.nodes[i].visited = false
			g.nodes[i].link = nil
		}
	}
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
			in = f
		}
	}

	g := new(graph)
	if err := gio.Read(filename, bufio.NewReader(in), g); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		edmondsKarp(g, 0, 1)
		w := bufio.NewWriter(os.Stdout)
		gio.Write(w, g)
		w.Flush()
	}
}
