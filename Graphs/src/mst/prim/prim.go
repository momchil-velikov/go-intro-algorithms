package main

import (
	"bufio"
	"fmt"
	gio "graph/io"
	"io"
	"os"
)

type edge struct {
	src    int
	dst    int
	tree   bool
	weight float64
}

type node struct {
	id        int
	heapIndex int
	tree      bool
	treeEdge  *edge
	label     string
	weight    float64
	in        []*edge
	out       []*edge
}

type graph struct {
	nodes []node
	edges []*edge
}

func (g *graph) CreateNode(name string) int {
	id := len(g.nodes)
	g.nodes = append(g.nodes, node{id: id, label: name, heapIndex: -1})
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

func findIn(np *node, src int) int {
	for i := range np.in {
		if np.in[i].src == src {
			return i
		}
	}
	return -1
}

func findOut(np *node, dst int) int {
	for i := range np.out {
		if np.out[i].dst == dst {
			return i
		}
	}
	return -1
}

func (g *graph) CreateEdge(src int, dst int, attrs []string) {
	sp, dp := &g.nodes[src], &g.nodes[dst]
	if findOut(sp, dst) == -1 && findIn(dp, src) == -1 {
		e := &edge{src: src, dst: dst}
		// e.weight = w
		for i := 0; i < len(attrs); i += 2 {
			if attrs[i] == "weight" {
				fmt.Sscanf(attrs[i+1], "%f", &e.weight)
				break
			}
		}
		g.edges = append(g.edges, e)
		sp.out = append(sp.out, e)
		dp.in = append(dp.in, e)
	}
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
	return len(g.nodes[id].out)
}

func (g *graph) EdgeAttrs(id int, i int) (dst int, attrs []string) {
	e := g.nodes[id].out[i]
	s := fmt.Sprintf("%.5f", e.weight)
	var color string
	if e.tree {
		color = "black"
	} else {
		color = "#f0f0f0"
	}
	return e.dst, []string{"color", color, "label", s, "weight", s}
}

func heapPush(h []*node, np *node) []*node {
	h = append(h, np)
	np.heapIndex = len(h) - 1
	heapUp(h, len(h)-1)
	return h
}

func heapPop(h []*node) ([]*node, *node) {
	if len(h) == 0 {
		return h, nil
	}
	np := h[0]
	n := len(h) - 1
	h[0] = h[n]
	h[0].heapIndex = 0
	h = h[:n]
	heapDown(h, 0)
	return h, np
}

func heapUp(h []*node, i int) {
	w := h[i].weight
	for i > 0 {
		j := (i - 1) / 2
		if w >= h[j].weight {
			break
		}
		h[i], h[j] = h[j], h[i]
		h[i].heapIndex = i
		h[j].heapIndex = j
		i = j
	}
}

func heapDown(h []*node, i int) {
	n := len(h)
	for i < n {
		x := i
		w := h[x].weight
		j := 2*i + 1
		if j < n && w > h[j].weight {
			x = j
			w = h[x].weight
		}
		if j+1 < n && w > h[j+1].weight {
			x = j + 1
			w = h[x].weight
		}
		if i == x {
			break
		}
		h[i], h[x] = h[x], h[i]
		h[i].heapIndex = i
		h[x].heapIndex = x
		i = x
	}
}

func mstPrim(g *graph) float64 {
	np := &g.nodes[0]
	np.heapIndex = 0
	np.weight = 0.0
	pq := []*node{np}

	var relax = func(ep *edge, np *node) {
		if np.tree {
			return
		}
		if np.heapIndex == -1 {
			np.weight = ep.weight
			np.treeEdge = ep
			pq = heapPush(pq, np)
		} else {
			if ep.weight < np.weight {
				np.weight = ep.weight
				np.treeEdge = ep
				heapUp(pq, np.heapIndex)
			}
		}
	}

	weight := 0.0
	comp := 0.0
	for len(pq) != 0 {
		pq, np = heapPop(pq)
		np.tree = true

		y := np.weight - comp
		t := weight + y
		comp = (t - weight) - y
		weight = t

		for _, e := range np.out {
			dp := &g.nodes[e.dst]
			relax(e, dp)
		}
		for _, e := range np.in {
			sp := &g.nodes[e.src]
			relax(e, sp)
		}
	}

	for i := range g.nodes {
		e := g.nodes[i].treeEdge
		if e != nil {
			e.tree = true
		}
	}
	return weight
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
	}
	weight := mstPrim(g)
	w := bufio.NewWriter(os.Stdout)
	gio.Write(w, g)
	w.Flush()
	fmt.Fprintf(os.Stderr, "weight = %.5f\n", weight)
}
