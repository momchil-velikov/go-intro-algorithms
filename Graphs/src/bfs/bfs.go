package main

import (
	"bufio"
	"bytes"
	"container/deque"
	"fmt"
	gio "graph/io"
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
	label string
	color colorT
	tn    int
	succ  []edge
}

type graph struct {
	nodes []node
}

func (g *graph) CreateNode(name string) int {
	id := len(g.nodes)
	g.nodes = append(g.nodes, node{})
	g.nodes[id].label = name
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
	np.succ = append(np.succ, edge{dst: dst})
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
			"label", fmt.Sprintf("%s(%d)", np.label, np.tn),
			"color", fmt.Sprintf("/rdbu9/%d", 1+np.tn%9),
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

func bfsInner(g *graph, r int, tn int) int {
	q := deque.Deque{}
	q.Push(r)
	for !q.IsEmpty() {
		a := &g.nodes[q.Pop().(int)]
		a.tn = tn
		for i := range a.succ {
			e := &a.succ[i]
			b := &g.nodes[e.dst]
			if b.color == WHITE {
				b.color = GREY
				e.color = "#a0a0a0"
				q.Push(e.dst)
			} else {
				e.color = "#f0f0f0"
			}
		}
		a.color = BLACK
	}
	return tn
}

func bfs(g *graph) {
	tn := 0
	for i := range g.nodes {
		np := &g.nodes[i]
		if np.color == WHITE {
			tn = bfsInner(g, i, tn)
			tn++
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: bfs <filename>.dot")
		return
	}
	filename := os.Args[1]
	if f, err := os.Open(filename); err == nil {
		g := new(graph)
		if err = gio.Read(filename, bufio.NewReader(f), g); err == nil {
			bfs(g)
			b := new(bytes.Buffer)
			if err := gio.Write(b, g); err == nil {
				fmt.Println(b.String())
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(filename, ":", err)
	}
}
