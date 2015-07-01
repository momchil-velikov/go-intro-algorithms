package main

import (
	"bufio"
	djs "container/disjoint_set"
	"fmt"
	gio "graph/io"
	"io"
	"os"
)

type node struct {
	label string
	id    int
	set   djs.SetId
	grp   djs.SetId
	succ  []int
}

type graph struct {
	nodes []*node
}

func (g *graph) CreateNode(name string) int {
	id := len(g.nodes)
	g.nodes = append(g.nodes, &node{set: -1})
	g.nodes[id].label = name
	g.nodes[id].id = id
	return id
}

func (g *graph) SetNodeAttrs(id int, attrs []string) {
	np := g.nodes[id]
	for i := 0; i < len(attrs); i += 2 {
		if attrs[i] == "label" {
			np.label = attrs[i+1]
			break
		}
	}
}

func (g *graph) CreateEdge(src int, dst int, attrs []string) {
	np := g.nodes[src]
	for i := 0; i < len(np.succ); i++ {
		if np.succ[i] == dst {
			return
		}
	}
	np.succ = append(np.succ, dst)
}

func (g *graph) NodeCount() int {
	return len(g.nodes)
}

func (g *graph) NodeAttrs(id int) []string {
	np := g.nodes[id]
	if len(np.label) == 0 {
		return nil
	} else {
		s := int(np.grp)
		return []string{
			"label", fmt.Sprintf("%s(%d)", np.label, s),
			"color", fmt.Sprintf("/rdbu9/%d", 1+s%9),
			"style", "filled",
		}
	}
}

func (g *graph) EdgeCount(id int) int {
	return len(g.nodes[id].succ)
}

func (g *graph) EdgeAttrs(id int, i int) (dst int, attrs []string) {
	return g.nodes[id].succ[i], []string{"color", "#c0c0c0"}
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
	var eltCount int
	if _, err := fmt.Fscan(in, &eltCount); err == nil {
		g := new(graph)
		sets := djs.Sets{}
		for i := 0; i < eltCount; i++ {
			id := g.CreateNode(fmt.Sprintf("v%d", i))
			np := g.nodes[id]
			np.set = sets.Create(np)
		}
		var s, e djs.SetId
		_, err := fmt.Fscan(in, &s, &e)
		for err == nil {
			sp, ep := g.nodes[s], g.nodes[e]
			g.CreateEdge(sp.id, ep.id, nil)
			sets.Union(sp.set, ep.set)
			_, err = fmt.Fscan(in, &s, &e)
		}
		if err != nil && err != io.EOF {
			fmt.Println(filename, ":error reading connections:", err)
			return
		}
		for _, np := range g.nodes {
			np.grp = sets.Find(np.set)
		}
		w := bufio.NewWriter(os.Stdout)
		gio.Write(w, g)
		w.Flush()
		fmt.Fprintln(os.Stderr, sets.Count(), " components")
	} else {
		fmt.Println(filename, ":", err)
	}
}
