package treap

import "math/rand"

type KeyT int64

const (
	maxKey = KeyT(0x7fffffffffffffff)
	minKey = -maxKey - 1
)

type node struct {
	key         KeyT
	left, right *node
	prio        uint
}

type Treap struct {
	rng  *rand.Rand
	root *node
}

func New() *Treap {
	var t Treap
	t.Init()
	return &t
}

func (t *Treap) Init() {
	t.root = nil
	t.rng = rand.New(rand.NewSource(1))
}

func (t *Treap) newNode(k KeyT) *node {
	return &node{key: k, left: nil, right: nil, prio: uint(t.rng.Int())}
}

func (t *Treap) Insert(k KeyT) {
	t.root = t.insert(t.root, k)
}

func (t *Treap) insert(x *node, k KeyT) *node {
	if x == nil {
		return t.newNode(k)
	}
	if k <= x.key {
		x.left = t.insert(x.left, k)
		if x.left.prio < x.prio {
			x = rotateRight(x)
		}
	} else {
		x.right = t.insert(x.right, k)
		if x.right.prio < x.prio {
			x = rotateLeft(x)
		}
	}
	return x
}

func (t *Treap) Delete(k KeyT) {
	var p *node = nil
	x := t.root
	for x != nil && x.key != k {
		p = x
		if k < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	switch {
	case x == nil:
		return
	case p == nil:
		t.root = t.down(x)
	case x == p.left:
		p.left = t.down(x)
	default: // x == p.right
		p.right = t.down(x)
	}
}

func rotateLeft(x *node) *node {
	y := x.right
	x.right = y.left
	y.left = x
	return y
}

func rotateRight(y *node) *node {
	x := y.left
	y.left = x.right
	x.right = y
	return x
}

func (t *Treap) down(x *node) *node {
	var y *node
	switch {
	case x.left == nil:
		return x.right
	case x.right == nil:
		return x.left
	case x.left.prio < x.right.prio:
		y = rotateRight(x)
		y.right = t.down(x)
	default:
		y = rotateLeft(x)
		y.left = t.down(x)
	}
	return y
}

func (t *Treap) verifyBSTProperty() bool {
	return verifyBSTProperty(t.root, minKey, maxKey)
}

func verifyBSTProperty(p *node, low, high KeyT) bool {
	if p == nil {
		return true
	}
	if p.key < low || p.key > high {
		return false
	}
	return verifyBSTProperty(p.left, low, p.key) && verifyBSTProperty(p.right, p.key+1, high)
}

func (t *Treap) verifyHeapProperty() bool {
	return verifyHeapProperty(t.root)
}

func verifyHeapProperty(x *node) bool {
	if x == nil {
		return true
	}
	if !verifyHeapProperty(x.left) || !verifyHeapProperty(x.right) {
		return false
	}
	if x.left != nil && x.left.prio < x.prio {
		return false
	}
	if x.right != nil && x.right.prio < x.prio {
		return false
	}
	return true
}

func (t *Treap) height() uint {
	return height(t.root)
}

func max(a, b uint) uint {
	if a > b {
		return a
	} else {
		return b
	}
}

func height(x *node) uint {
	if x == nil {
		return 0
	} else {
		return 1 + max(height(x.left), height(x.right))
	}
}
