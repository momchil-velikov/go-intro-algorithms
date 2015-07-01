package main

import (
	"bufio"
	"bytes"
	"container/stack"
	"flag"
	"fmt"
	gio "graph/io"
	"io"
	"os"
)

type edge struct {
	dst   int
	color string
}

type colorT int

const (
	WHITE colorT = iota
	GREY
	BLACK
)

type node struct {
	id           int
	label        string
	color        colorT
	dtime, ftime int
	scc          int
	lowlink      int
	succ         []edge
	pred         []edge
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
	for i := 0; i < len(np.succ); i++ {
		if np.succ[i].dst == dst {
			return
		}
	}
	np.succ = append(np.succ, edge{dst: dst, color: "#c0c0c0"})
}

func (g *graph) CreateRevEdge(src int, dst int, attrs []string) {
	np := &g.nodes[src]
	for i := 0; i < len(np.pred); i++ {
		if np.pred[i].dst == dst {
			return
		}
	}
	np.pred = append(np.pred, edge{dst: dst, color: "#c0c0c0"})
}

func (g *graph) NodeCount() int {
	return len(g.nodes)
}

func (g *graph) NodeAttrs(id int) []string {
	np := &g.nodes[id]
	if len(np.label) == 0 {
		return nil
	} else {
		return []string{
			"label", fmt.Sprintf("%s(%d)", np.label, np.scc),
			"color", fmt.Sprintf("/rdbu9/%d", 1+np.scc%9),
			"style", "filled",
		}
	}
}

func (g *graph) EdgeCount(id int) int {
	return len(g.nodes[id].succ)
}

func (g *graph) EdgeAttrs(id int, i int) (dst int, attrs []string) {
	e := &g.nodes[id].succ[i]
	if len(e.color) == 0 {
		e.color = "#e0e0e0"
	}
	return e.dst, []string{"color", e.color}
}

func dfsIterInner(g *graph, r int, tn int) {
	type state struct {
		id int // node id
		i  int // edge index
	}
	q := stack.Stack{}
	q.Push(&state{id: r, i: 0})
	g.nodes[r].color = GREY
	for !q.IsEmpty() {
		s := q.Top().(*state)
		a := &g.nodes[s.id]
		if s.i == len(a.succ) {
			q.Pop()
			a.color = BLACK
			a.scc = tn
		} else {
			e := &a.succ[s.i]
			s.i++
			b := &g.nodes[e.dst]
			if b.color == WHITE {
				e.color = "green"
				b.color = GREY
				q.Push(&state{id: e.dst, i: 0})
			} else if b.color == GREY {
				e.color = "red"
			} else {
				e.color = "yellow"
			}
		}
	}
}

func dfsIter(g *graph) {
	tn := 1
	for i := range g.nodes {
		if g.nodes[i].color == WHITE {
			dfsIterInner(g, i, tn)
			tn++
		}
	}
}

func dfsInner(g *graph, r int, tn int, tm int) int {
	a := &g.nodes[r]
	a.color = GREY
	tm++
	a.dtime = tm
	a.scc = tn
	for i := range a.succ {
		e := &a.succ[i]
		b := &g.nodes[e.dst]
		switch b.color {
		case WHITE:
			e.color = "green"
			tm = dfsInner(g, e.dst, tn, tm)
		case GREY:
			e.color = "red"
		default:
			if a.dtime < b.dtime {
				e.color = "lightblue"
			} else {
				e.color = "yellow"
			}
		}
	}
	a.color = BLACK
	tm++
	a.ftime = tm
	return tm
}

func dfs(g *graph) {
	tn := 1
	tm := 0
	for i := range g.nodes {
		if g.nodes[i].color == WHITE {
			tm = dfsInner(g, i, tn, tm)
			tn++
		}
	}
}

func topSortInner(g *graph, r int, a []int) []int {
	h := &g.nodes[r]
	h.color = BLACK
	for i := range h.succ {
		e := &h.succ[i]
		t := &g.nodes[e.dst]
		if t.color == WHITE {
			a = topSortInner(g, e.dst, a)
		}
	}
	return append(a, r)
}

func topSort(g *graph) (a []int) {
	for i := range g.nodes {
		if g.nodes[i].color == WHITE {
			a = topSortInner(g, i, a)
		}
	}
	return
}

func simplePathRec(g *graph, s, t int) {
	if s == t {
		g.nodes[t].dtime++
		return
	}
	a := &g.nodes[s]
	a.color = BLACK
	for i := range a.succ {
		e := &a.succ[i]
		b := &g.nodes[e.dst]
		if b.color == WHITE {
			simplePathRec(g, e.dst, t)
		}
	}
	a.color = WHITE
}

func simplePath(g *graph, s, t int) int {
	simplePathRec(g, s, t)
	return g.nodes[t].dtime
}

func sccInner(g *graph, stk *stack.Stack, r int, tm int, scc int) (int, int) {
	a := &g.nodes[r]
	tm++
	a.dtime = tm
	a.lowlink = tm
	stk.Push(a)
	a.color = GREY
	for i := range a.succ {
		e := &a.succ[i]
		b := &g.nodes[e.dst]
		if b.color == WHITE {
			tm, scc = sccInner(g, stk, e.dst, tm, scc)
			if a.lowlink > b.lowlink {
				a.lowlink = b.lowlink
			}
		} else if b.color == GREY {
			if a.lowlink > b.lowlink {
				a.lowlink = b.lowlink
			}
		}
	}
	if a.lowlink == a.dtime {
		scc++
	loop:
		w := stk.Pop().(*node)
		w.color = BLACK
		w.scc = scc
		if w != a {
			goto loop
		}
	}
	return tm, scc
}

func scc(g *graph) {
	tm, scc := 0, 0
	stk := stack.Stack{}
	for i := range g.nodes {
		if g.nodes[i].color == WHITE {
			tm, scc = sccInner(g, &stk, i, tm, scc)
		}
	}
}

func addReverseEdges(g *graph) {
	for i := range g.nodes {
		np := &g.nodes[i]
		for j := range np.succ {
			e := &np.succ[j]
			g.CreateRevEdge(e.dst, i, nil)
		}
	}
}

func dfsKosarajuInner(g *graph, r int, tm int, ns []*node) (int, []*node) {
	a := &g.nodes[r]
	a.color = GREY
	for i := range a.succ {
		e := &a.succ[i]
		b := &g.nodes[e.dst]
		if b.color == WHITE {
			tm, ns = dfsKosarajuInner(g, e.dst, tm, ns)
		}
	}
	tm++
	a.ftime = tm
	return tm, append(ns, a)
}

func dfsKosaraju(g *graph) (a []*node) {
	tm := 0
	for i := range g.nodes {
		if g.nodes[i].color == WHITE {
			tm, a = dfsKosarajuInner(g, i, tm, a)
		}
	}
	return
}

func sccKosarajuInner(g *graph, r int, scc int) {
	a := &g.nodes[r]
	a.color = BLACK
	a.scc = scc
	for i := range a.pred {
		e := &a.pred[i]
		b := &g.nodes[e.dst]
		if b.color == GREY {
			sccKosarajuInner(g, e.dst, scc)
		}
	}
}

func sccKosaraju(g *graph) {
	a := dfsKosaraju(g)
	addReverseEdges(g)
	scc := 0
	for i := len(a) - 1; i >= 0; i-- {
		np := a[i]
		if np.color == GREY {
			sccKosarajuInner(g, np.id, scc)
			scc++
		}
	}
}

func main() {
	flagDFS := flag.Bool("dfs", false, "run DFS travarsal")
	flagTarjan := flag.Bool("scc-tarjan", false, "run Tarjan's SCC algorithm")
	flagKosaraju := flag.Bool("scc-kosaraju", false, "run Kosaraju's SCC algorithm")
	flag.Parse()

	var in io.Reader
	var filename string
	if len(flag.Args()) == 0 {
		filename = "stdin"
		in = os.Stdin
	} else {
		filename = flag.Args()[0]
		if f, err := os.Open(filename); err != nil {
			fmt.Println(filename, ":", err)
			return
		} else {
			defer f.Close()
			in = f
		}
	}

	g := new(graph)
	if err := gio.Read(filename, bufio.NewReader(in), g); err == nil {
		if *flagDFS {
			dfs(g)
		} else if *flagTarjan {
			scc(g)
		} else if *flagKosaraju {
			sccKosaraju(g)
		}
		b := new(bytes.Buffer)
		if err := gio.Write(b, g); err == nil {
			fmt.Println(b.String())
		}
	} else {
		fmt.Println(err)
	}
}
