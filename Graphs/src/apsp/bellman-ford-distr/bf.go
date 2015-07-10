package main

import (
	"bufio"
	"fmt"
	gio "graph/io"
	"io"
	"math"
	"os"
	"runtime"
)

type msg struct {
	src  int
	dist float64
}

type edge struct {
	dst    int
	weight float64
}

type node struct {
	id    int
	label string
	dist  float64
	prev  int
	succ  []edge
	pred  []edge
	port  chan msg
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
	if src == dst {
		return
	}
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
	np = &g.nodes[dst]
	np.pred = append(np.pred, edge{dst: src, weight: weight})
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
	s := fmt.Sprintf("%.2f", e.weight)
	var color string
	if g.nodes[e.dst].prev == id {
		color = "#808080"
	} else {
		color = "#e0e0e0"
	}
	return e.dst, []string{"label", s, "weight", s, "color", color}
}

func bellmanFordLeader(nn int, ng int, sync chan bool, done chan struct{}) {
	for k := 1; k < nn; k++ {
		for i := 0; i < ng; i++ {
			sync <- true
		}
		for i := 0; i < ng; i++ {
			<-done
		}
	}
	for i := 0; i < ng; i++ {
		sync <- false
	}
	for i := 0; i < ng; i++ {
		<-done
	}
}

func bellmanFordNode(g *graph, start int, end int, sync chan bool, done chan struct{}) {
	for {
		if b := <-sync; !b {
			done <- struct{}{}
			return
		}

		for id := start; id < end; id++ {
			self := &g.nodes[id]
			for i := range self.succ {
				e := &self.succ[i]
				g.nodes[e.dst].port <- msg{self.id, self.dist + self.succ[i].weight}
			}
		}

		for id := start; id < end; id++ {
			self := &g.nodes[id]
			n := len(self.pred)
			for i := 0; i < n; i++ {
				m := <-self.port
				if m.dist < self.dist {
					self.dist = m.dist
					self.prev = m.src
				}
			}
		}
		done <- struct{}{}
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

	sync := make(chan bool)
	done := make(chan struct{})
	for i := range g.nodes {
		np := &g.nodes[i]
		np.dist = math.Inf(0)
		np.prev = -1
		np.port = make(chan msg, len(np.pred))
	}
	g.nodes[0].dist = 0.0

	N := runtime.NumCPU()
	runtime.GOMAXPROCS(N)
	n := len(g.nodes) / N
	for i := 0; i < N; i++ {
		go bellmanFordNode(g, i*n, (i+1)*n, sync, done)
	}
	if n*N < len(g.nodes) {
		go bellmanFordNode(g, n*N, len(g.nodes), sync, done)
		N++
	}
	bellmanFordLeader(len(g.nodes), N, sync, done)

	w := bufio.NewWriter(os.Stdout)
	gio.Write(w, g)
	w.Flush()
}
