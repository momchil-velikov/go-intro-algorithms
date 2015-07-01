package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type node struct {
	id   int
	next [2]*node
}

const maxDepth = 13

var nextId int = 0
var blocks []*node

func genBasicBlock() *node {
	n := &node{id: nextId}
	nextId++
	blocks = append(blocks, n)
	return n
}

func genRandom(d uint) (*node, *node) {
	if d == 0 {
		n := genBasicBlock()
		return n, n
	}
	r := rand.Intn(100)
	switch {
	case r < 16:
		n := genBasicBlock()
		return n, n
	case r < 32:
		return genSeq(d)
	case r < 48:
		return genIrreducibleLoop(d)
	case r < 56:
		return genLoop(d)
	case r < 72:
		return genIf(d)
	default:
		return genIfElse(d)
	}
}

func genSeq(d uint) (*node, *node) {
	n1, n2 := genRandom(d - 1)
	n3, n4 := genRandom(d - 1)
	n2.next[0] = n3
	return n1, n4
}

func genIfElse(d uint) (*node, *node) {
	ent0, ent1 := genRandom(d - 1)
	yes0, yes1 := genRandom(d - 1)
	no0, no1 := genRandom(d - 1)
	ex0, ex1 := genRandom(d - 1)
	ent1.next[0] = yes0
	ent1.next[1] = no0
	yes1.next[0] = ex0
	no1.next[0] = ex0
	return ent0, ex1
}

func genIf(d uint) (*node, *node) {
	ent0, ent1 := genRandom(d - 1)
	yes0, yes1 := genRandom(d - 1)
	ex0, ex1 := genRandom(d - 1)
	ent1.next[0] = yes0
	ent1.next[1] = ex0
	yes1.next[0] = ex0
	return ent0, ex1
}

func genLoop(d uint) (*node, *node) {
	body0, body1 := genRandom(d - 1)
	cond0, cond1 := genRandom(d - 1)
	exit := genBasicBlock()
	cond1.next[0] = body0
	cond1.next[1] = exit
	body1.next[0] = cond0
	return cond0, exit
}

func genIrreducibleLoop(d uint) (*node, *node) {
	entry := genBasicBlock()
	left0, left1 := genRandom(d)
	right0, right1 := genRandom(d)
	exit := genBasicBlock()
	entry.next[0] = left0
	entry.next[1] = right0
	left1.next[0] = right0
	left1.next[1] = exit
	right1.next[0] = left0
	right1.next[1] = exit
	return entry, exit
}

func main() {
	depth := flag.Uint("d", 3, "max subgraph nesting level")
	seed := flag.Int64("s", 0, "random seed")
	flag.Parse()
	if *seed == 0 {
		rand.Seed(time.Now().Unix())
	} else {
		rand.Seed(int64(*seed))
	}

	genRandom(*depth)

	fmt.Println("digraph CFG {")
	fmt.Println("n0[shape=diamond]")
	fmt.Printf("n%d[shape=diamond]\n", len(blocks)-1)
	for _, s := range blocks {
		if e := s.next[0]; e != nil {
			fmt.Printf("n%d -> n%d\n", s.id, e.id)
		}
		if e := s.next[1]; e != nil {
			fmt.Printf("n%d -> n%d\n", s.id, e.id)
		}
	}
	fmt.Println("}\n")
}
