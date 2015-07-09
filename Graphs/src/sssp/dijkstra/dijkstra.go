package main

import (
	"bufio"
	"fmt"
	gio "graph/io"
	"io"
	"math"
	"os"
)

type edge struct {
	dst    int
	weight float64
}

type node struct {
	id    int
	label string
	dist  float64
	pred  int
	succ  []edge
}

type graph struct {
	nodes []node
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
	np := &g.nodes[src]
	for i := range np.succ {
		if np.succ[i].dst == dst {
			return
		}
	}
	weight := 0.0
	for i := 0; i < len(attrs); i += 2 {
		if attrs[i] == "weight" {
			fmt.Sscanf(attrs[i+1], "%f", &weight)
			break
		}
	}
	np.succ = append(np.succ, edge{dst: dst, weight: weight})
}

func (g *graph) NodeCount() int {
	return len(g.nodes)
}

func (g *graph) NodeAttrs(id int) []string {
	np := &g.nodes[id]
	if len(np.label) == 0 {
		return nil
	}

	dist := ""
	if np.dist == math.MaxFloat64 {
		dist = "\u221e"
	} else {
		dist = fmt.Sprintf("%.2f", np.dist)
	}
	label := fmt.Sprintf("%s(%s)", np.label, dist)
	return []string{"label", label}
}

func (g *graph) EdgeCount(id int) int {
	return len(g.nodes[id].succ)
}

func (g *graph) EdgeAttrs(id int, i int) (dst int, attrs []string) {
	e := g.nodes[id].succ[i]
	s := fmt.Sprintf("%.2f", e.weight)
	var color string
	if g.nodes[e.dst].pred == id {
		color = "#808080"
	} else {
		color = "#e0e0e0"
	}
	return e.dst, []string{"label", s, "weight", s, "color", color}
}

type nodeRep struct {
	np     *node
	index  int
	weight float64
	pred   int
}

func heapPush(h []*nodeRep, np *node) []*nodeRep {
	h = append(h, &nodeRep{np: np, index: len(h) - 1, weight: math.MaxFloat64})
	heapUp(h, len(h)-1)
	return h
}

func heapPop(h []*nodeRep) ([]*nodeRep, *nodeRep) {
	if len(h) == 0 {
		panic("pop from an empty heap")
	}
	nd := h[0]
	n := len(h) - 1
	h[0] = h[n]
	h[0].index = 0
	h = h[:n]
	heapDown(h, 0)
	return h, nd
}

func heapUp(h []*nodeRep, i int) {
	w := h[i].weight
	for i > 0 {
		j := (i - 1) / 2
		if w >= h[j].weight {
			break
		}
		h[i], h[j] = h[j], h[i]
		h[i].index = i
		h[j].index = j
		i = j
	}
}

func heapDown(h []*nodeRep, i int) {
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
		h[i].index = i
		h[x].index = x
		i = x
	}
}

func ssspDijkstra(g *graph) {
	reps := make([]nodeRep, len(g.nodes))
	pq := make([]*nodeRep, len(g.nodes))
	relax := func(np *node, weight float64, rep *nodeRep) {
		if np.dist+weight < rep.weight {
			rep.weight = np.dist + weight
			rep.pred = np.id
			heapUp(pq, rep.index)
		}
	}
	for i := range g.nodes {
		pq[i] = &reps[i]
		pq[i].np = &g.nodes[i]
		pq[i].index = i
		pq[i].weight = math.MaxFloat64
	}
	pq[0].weight = 0.0

	for len(pq) != 0 {
		var rep *nodeRep
		pq, rep = heapPop(pq)
		np := rep.np
		np.dist = rep.weight
		np.pred = rep.pred
		for _, e := range np.succ {
			dp := &g.nodes[e.dst]
			rep = &reps[dp.id]
			relax(np, e.weight, rep)
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
	}
	ssspDijkstra(g)
	w := bufio.NewWriter(os.Stdout)
	gio.Write(w, g)
	w.Flush()
}
