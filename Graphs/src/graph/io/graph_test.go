package io

type node struct {
	label string
	succ  []int
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
	return g.nodes[id].succ[i], nil
}
