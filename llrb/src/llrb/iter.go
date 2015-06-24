package llrb

type Iter struct {
	nd *node
}

func (t *Tree) Min() Iter {
	y := (*node)(nil)
	x := t.root
	for x != nil {
		y = x
		x = x.left
	}
	return Iter{y}
}

func (t *Tree) Max() Iter {
	y := (*node)(nil)
	x := t.root
	for x != nil {
		y = x
		x = x.right
	}
	return Iter{y}
}

func (t *Tree) Find(k KeyT) bool {
	p := lowerBound(t.root, k)
	return p != nil && k == p.key
}

func (t *Tree) Count(k KeyT) uint {
	lo, hi := t.LowerBound(k), t.UpperBound(k)
	c := uint(0)
	for lo != hi {
		c++
		lo = lo.Succ()
	}
	return c
}

func (t *Tree) LowerBound(k KeyT) Iter {
	return Iter{lowerBound(t.root, k)}
}

func lowerBound(x *node, k KeyT) *node {
	y := (*node)(nil)
	for x != nil {
		if k <= x.key {
			y = x
			x = x.left
		} else {
			x = x.right
		}
	}
	return y
}

func (t *Tree) UpperBound(k KeyT) Iter {
	return Iter{upperBound(t.root, k)}
}

func upperBound(x *node, k KeyT) *node {
	y := (*node)(nil)
	for x != nil {
		if k < x.key {
			y = x
			x = x.left
		} else {
			x = x.right
		}
	}
	return y
}

func (i Iter) Succ() Iter {
	return Iter{succ(i.nd)}
}

func succ(x *node) *node {
	if x == nil {
		return nil
	}
	if x.right == nil {
		p := x.parent
		for p != nil && x == p.right {
			x = p
			p = p.parent
		}
		x = p
	} else {
		x = x.right
		for x.left != nil {
			x = x.left
		}
	}
	return x
}

func (i Iter) Pred() Iter {
	return Iter{pred(i.nd)}
}

func pred(x *node) *node {
	if x == nil {
		return nil
	}
	if x.left == nil {
		p := x.parent
		for p != nil && x == p.left {
			x = p
			p = p.parent
		}
		x = p
	} else {
		x = x.left
		for x.right != nil {
			x = x.right
		}
	}
	return x
}

func (i Iter) Key() KeyT {
	return i.nd.key
}

func (i Iter) IsNil() bool {
	return i.nd == nil
}
