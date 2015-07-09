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

func ssspBellmanFord(g *graph) bool {
	relax := func(sp *node, weight float64, dp *node) {
		if sp.dist+weight < dp.dist {
			dp.dist = sp.dist + weight
			dp.pred = sp.id
		}
	}
	for i := range g.nodes {
		g.nodes[i].dist = math.MaxFloat64
		g.nodes[i].pred = -1
	}
	g.nodes[0].dist = 0
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
	if !ssspBellmanFord(g) {
		fmt.Fprintln(os.Stderr, "negative cycle")
	}
	w := bufio.NewWriter(os.Stdout)
	gio.Write(w, g)
	w.Flush()
}
