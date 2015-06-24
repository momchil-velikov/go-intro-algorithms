package llrb

import "fmt"

type KeyT int64

const (
	maxkey = KeyT(0x7fffffffffffffff)
	minkey = -maxkey - 1
)

type node struct {
	key                 KeyT
	parent, left, right *node
	black               bool
}

type Tree struct {
	root *node
}

func (t *Tree) Insert(k KeyT) {
	t.root = insert(t.root, k)
	t.root.parent = nil
	t.root.black = true
}

func insert(x *node, k KeyT) *node {
	if x == nil {
		return &node{key: k}
	}

	var p *node
	for x != nil {
		p = x
		if k <= p.key {
			x = p.left
			if x == nil {
				p.left = &node{key: k}
				p.left.parent = p
			}
		} else {
			x = p.right
			if x == nil {
				p.right = &node{key: k}
				p.right.parent = p
			}
		}
	}

	for p != nil {
		if isRed(p.right) {
			p = rotateLeft(p)
		}

		if isRed(p.left) && isRed(p.left.left) {
			p = rotateRight(p)
		}

		if isRed(p.left) && isRed(p.right) {
			colorFlip(p)
		}
		x = p
		p = p.parent
	}
	return x
}

func rotateLeft(x *node) *node {
	y := x.right
	y.parent = x.parent
	if x.parent != nil {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	y.left = x
	x.parent = y
	x.black, y.black = y.black, x.black
	return y
}

func rotateRight(y *node) *node {
	x := y.left
	x.parent = y.parent
	if y.parent != nil {
		if y == y.parent.left {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
	}
	y.left = x.right
	if y.left != nil {
		y.left.parent = y
	}
	x.right = y
	y.parent = x
	x.black, y.black = y.black, x.black
	return x
}

func isRed(x *node) bool {
	if x == nil {
		return false
	} else {
		return !x.black
	}
}

func colorFlip(x *node) {
	x.black = !x.black
	x.left.black = !x.left.black
	x.right.black = !x.right.black
}

func fixUp(x *node) *node {
	if isRed(x.right) {
		x = rotateLeft(x)
	}

	if isRed(x.left) && isRed(x.left.left) {
		x = rotateRight(x)
	}

	if isRed(x.left) && isRed(x.right) {
		colorFlip(x)
	}
	return x
}

func moveRedRight(x *node) *node {
	// invariant: either |x| or |x.right| is red
	colorFlip(x)
	if isRed(x.left.left) {
		x = rotateRight(x)
		colorFlip(x)
	}
	return x
}

func moveRedLeft(x *node) *node {
	//invariant: either |x| or |x.right| is red
	colorFlip(x)
	if isRed(x.right.left) {
		x.right = rotateRight(x.right)
		x = rotateLeft(x)
		colorFlip(x)
	}
	return x
}

func (t *Tree) DeleteMin() {
	r, _ := deleteMin(t.root)
	if r != nil {
		r.black = true
	}
	t.root = r
}

func deleteMin(x *node) (*node, *node) {
	if x == nil {
		return nil, nil
	}
	for x.left != nil {
		if !isRed(x.left) && !isRed(x.left.left) {
			x = moveRedLeft(x)
		}
		x = x.left
	}
	p := x.parent
	r := x.right
	if p != nil {
		if x == p.left {
			p.left = r
		} else {
			p.right = r
		}
	}
	p = nil
	for r != nil {
		p = fixUp(r)
		r = p.parent
	}
	return p, x
}

func (t *Tree) Delete(k KeyT) bool {
	r, d := delete(t.root, k)
	if r != nil {
		r.black = true
	}
	t.root = r
	return d != nil
}

// NOT WORKING
func delete(x *node, k KeyT) (*node, *node) {
	if x == nil {
		return nil, nil
	}
	var del *node = nil
	if k < x.key {
		if !isRed(x.left) && !isRed(x.left.left) {
			x = moveRedLeft(x)
		}
		x.left, del = delete(x.left, k)
	} else {
		if isRed(x.left) {
			x = rotateRight(x)
		}
		if k == x.key && x.right == nil {
			return nil, x
		}
		if !isRed(x.right) && !isRed(x.right.left) {
			x = moveRedRight(x)
		}
		if k == x.key {
			r, d := deleteMin(x.right)
			d.left = x.left
			d.right = r
			d.parent = x.parent
			if d.left != nil {
				d.left.parent = d
			}
			if d.right != nil {
				d.right.parent = d
			}
			d.black = x.black
			del = x
			x = d
		} else {
			x.right, del = delete(x.right, k)
		}
	}
	return fixUp(x), del
}

func (t *Tree) isBST() bool {
	return isBST(t.root, minkey, maxkey)
}

func isBST(x *node, min, max KeyT) bool {
	switch {
	case x == nil:
		return true
	case x.key < min || x.key > max:
		return false
	default:
		return isBST(x.left, min, x.key) && isBST(x.right, x.key, max)
	}
}

func (t *Tree) WriteDot() string {
	return "digraph LLRB {\n" + writeDot(t.root) + "}\n"
}

func writeDot(x *node) string {
	if x == nil {
		return fmt.Sprintf("n%p[shape=point]", x)
	}
	c := "red"
	if x.black {
		c = "black"
	}
	s := fmt.Sprintf("n%p[label=%d, color=%s]\n", x, x.key, c)
	if x.left != nil {
		s += fmt.Sprintf("n%p -> n%p\n", x, x.left)
		s += writeDot(x.left)
	} else {
		s += fmt.Sprintf("n%p -> n%pleft\n", x, x)
		s += fmt.Sprintf("n%pleft[shape=point]\n", x)
	}
	if x.right != nil {
		s += fmt.Sprintf("n%p -> n%p\n", x, x.right)
		s += writeDot(x.right)
	} else {
		s += fmt.Sprintf("n%p -> n%pright\n", x, x)
		s += fmt.Sprintf("n%pright[shape=point]\n", x)
	}

	return s
}
