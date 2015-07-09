package main

import (
	"bufio"
	djs "container/disjoint_set"
	"fmt"
	gio "graph/io"
	"io"
	"os"
	"sort"
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
	set       djs.SetId
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

type byWeight []*edge

func (s byWeight) Len() int {
	return len(s)
}

func (s byWeight) Less(i, j int) bool {
	return s[i].weight < s[j].weight
}

func (s byWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func mstKruskal(g *graph) float64 {
	sets := djs.Sets{}
	for i := range g.nodes {
		np := &g.nodes[i]
		np.set = sets.Create(np.id)
	}
	sort.Sort(byWeight(g.edges))
	weight := 0.0
	comp := 0.0
	for _, e := range g.edges {
		sp, dp := &g.nodes[e.src], &g.nodes[e.dst]
		if sets.Find(sp.set) != sets.Find(dp.set) {
			y := e.weight - comp
			t := weight + y
			comp = (t - weight) - y
			weight = t
			// weight += e.weight
			e.tree = true
			sets.Union(sp.set, dp.set)
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
			in = bufio.NewReader(f)
		}
	}

	g := new(graph)
	if err := gio.Read(filename, bufio.NewReader(in), g); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	weight := mstKruskal(g)
	w := bufio.NewWriter(os.Stdout)
	gio.Write(w, g)
	w.Flush()
	fmt.Fprintf(os.Stderr, "weight = %.5f\n", weight)
}
