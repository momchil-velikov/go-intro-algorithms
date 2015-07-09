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
	id      int
	label   string
	dist, h float64
	pred    int
	index   int
	succ    []edge
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
	if math.IsInf(np.dist, 0) {
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
	s := fmt.Sprintf("%.2f/%.2f", e.weight, e.weight+g.nodes[id].h-g.nodes[e.dst].h)
	return e.dst, []string{"label", s, "weight", s}
}

func heapPop(h []*node) ([]*node, *node) {
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

func heapUp(h []*node, i int) {
	w := h[i].dist
	for i > 0 {
		j := (i - 1) / 2
		if w >= h[j].dist {
			break
		}
		h[i], h[j] = h[j], h[i]
		h[i].index = i
		h[j].index = j
		i = j
	}
}

func heapDown(h []*node, i int) {
	n := len(h)
	for i < n {
		x := i
		w := h[x].dist
		j := 2*i + 1
		if j < n && w > h[j].dist {
			x = j
			w = h[x].dist
		}
		if j+1 < n && w > h[j+1].dist {
			x = j + 1
			w = h[x].dist
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

func dijkstra(s int, g *graph) {
	pq := make([]*node, len(g.nodes))
	relax := func(sp *node, weight float64, dp *node) {
		if weight < 0.0 {
			panic("negative weight")
		}
		if sp.dist+weight < dp.dist {
			dp.dist = sp.dist + weight
			dp.pred = sp.id
			heapUp(pq, dp.index)
		}
	}
	for i := range g.nodes {
		pq[i] = &g.nodes[i]
		pq[i].index = i
		pq[i].dist = math.Inf(0)
		pq[i].pred = -1
	}
	pq[0], pq[s] = pq[s], pq[0]
	pq[0].dist = 0.0
	for len(pq) != 0 {
		var sp *node
		pq, sp = heapPop(pq)
		for _, e := range sp.succ {
			dp := &g.nodes[e.dst]
			relax(sp, e.weight+sp.h-dp.h, dp)
		}
	}
}

func bellmanFord(s int, g *graph) bool {
	relax := func(sp *node, weight float64, dp *node) {
		if sp.dist+weight < dp.dist {
			dp.dist = sp.dist + weight
		}
	}
	for i := range g.nodes {
		g.nodes[i].dist = math.Inf(0)
	}
	g.nodes[s].dist = 0
	for step := 0; step < len(g.nodes)-1; step++ {
		for i := range g.nodes {
			sp := &g.nodes[i]
			for _, e := range sp.succ {
				dp := &g.nodes[e.dst]
				relax(sp, e.weight, dp)
			}
		}
	}
	for i := range g.nodes {
		sp := &g.nodes[i]
		for _, e := range sp.succ {
			dp := &g.nodes[e.dst]
			if dp.dist > sp.dist+e.weight {
				return false
			}
		}
	}
	return true
}

func findEdge(np *node, dst int) *edge {
	for i := range np.succ {
		e := &np.succ[i]
		if e.dst == dst {
			return e
		}
	}
	return nil
}

func johnson(g *graph) bool {
	n := len(g.nodes)
	s := g.CreateNode("r3wE1gH7")
	for i := 0; i < n; i++ {
		g.CreateEdge(s, i, nil)
	}
	if !bellmanFord(s, g) {
		return false
	}
	for i := range g.nodes {
		g.nodes[i].h = g.nodes[i].dist
	}
	for i := 0; i < len(g.nodes)-1; i++ {
		sp := &g.nodes[i]
		dijkstra(i, g)
		for j := range g.nodes {
			dp := &g.nodes[j]
			if math.IsInf(dp.dist, 0) {
				fmt.Printf("%s to %s(\u221e)\n", sp.label, dp.label)
			} else {
				fmt.Printf("%s to %s(%5.2f): ", sp.label, dp.label, dp.dist-sp.h+dp.h)
				for dp.pred >= 0 {
					sp := &g.nodes[dp.pred]
					e := findEdge(sp, dp.id)
					fmt.Printf("%s<-%s(%5.2f) ", dp.label, sp.label, e.weight)
					dp = sp
				}
				fmt.Println()
			}
		}
	}
	return true
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
	johnson(g)
	// w := bufio.NewWriter(os.Stdout)
	// gio.Write(w, g)
	// w.Flush()
}
