package skip

import (
	"fmt"
	"math/rand"
)

type KeyT int64

type node struct {
	key  KeyT
	next []*node
}

type List struct {
	head *node
	rng  *rand.Rand
}

func New() *List {
	var lst List
	lst.Init()
	return &lst
}

func (l *List) Init() {
	l.head = &node{key: 0, next: make([]*node, 1)}
	l.rng = rand.New(rand.NewSource(1))
}

func (l *List) newNode(k KeyT) *node {
	h := int(1)
	for {
		r := l.rng.Intn(100)
		if r < 50 {
			break
		}
		h++
	}
	return &node{key: k, next: make([]*node, h)}
}

func (l *List) Insert(k KeyT) {
	n := len(l.head.next)
	z := l.insert(n-1, l.head, k)
	m := len(z.next)
	if m > n {
		next := make([]*node, m)
		copy(next, l.head.next)
		for i := n; i < m; i++ {
			next[i] = z
		}
		l.head.next = next
	}
}

func (l *List) insert(n int, x *node, k KeyT) *node {
	y := x.next[n]
	for y != nil && k >= y.key {
		x = y
		y = x.next[n]
	}
	if n == 0 {
		y = l.newNode(k)
		y.next[0] = x.next[0]
		x.next[0] = y
	} else {
		y = l.insert(n-1, x, k)
		if n < len(y.next) {
			y.next[n] = x.next[n]
			x.next[n] = y
		}
	}
	return y
}

func (l *List) Delete(k KeyT) {
	n := len(l.head.next)
	z := l.delete(n-1, l.head, k)
	if z != nil {
		for i := 0; i < len(z.next); i++ {
			if l.head.next[i] == z {
				l.head.next[i] = z.next[i]
			}
		}
	}
}

func (l *List) delete(n int, x *node, k KeyT) *node {
	y := x.next[n]
	for y != nil && y.key < k {
		x = y
		y = x.next[n]
	}
	if n == 0 {
		if y != nil && k == y.key {
			x.next[0] = y.next[0]
			return y
		} else {
			return nil
		}
	} else {
		y = l.delete(n-1, x, k)
		if y != nil && n < len(y.next) {
			x.next[n] = y.next[n]
		}
		return y
	}
}

func (l *List) VerifyOrder() bool {
	x := l.head.next[0]
	for x != nil {
		for i := 0; i < len(x.next); i++ {
			y := x.next[i]
			if y != nil && x.key > y.key {
				return false
			}
		}
		x = x.next[0]
	}
	return true
}

func (l *List) WriteDot() string {
	s := "digraph Skip {\n"
	x := l.head
	for x != nil {
		s += fmt.Sprintf("n%d[label=\"%d[%d]\"]\n", x.key, x.key, len(x.next))
		for i := 0; i < len(x.next); i++ {
			if x.next[i] != nil {
				s += fmt.Sprintf("n%d -> n%d\n", x.key, x.next[i].key)
			}
		}
		x = x.next[0]
	}
	s += "}\n"
	return s
}
