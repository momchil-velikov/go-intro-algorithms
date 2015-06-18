package btree

import (
	"fmt"
)

type KeyT int64

// Degree of the B-Tree; maximum number of keys per node is 2*_T-1, minimum is _T-1
const _T = 8

type node struct {
	n     uint           // number of keys in the node
	key   [2*_T - 1]KeyT // keys
	child [2 * _T]*node  // child nodes
}

type Tree struct {
	root *node // root node of the tree; always present
}

// Creates a B-Tree
func New() *Tree {
	t := &Tree{}
	t.Init()
	return t
}

// Initializes a B-Tree
func (t *Tree) Init() {
	t.root = &node{}
}

// Searches for a key in the tree. Returns |true| if the key is present.
func (t *Tree) Search(k KeyT) bool {
	x := t.root
	for x != nil {
		i := find(x, k)
		if i < x.n && k == x.key[i] {
			return true
		}
		x = x.child[i]
	}
	return false
}

// Inserts a key into the tree; allows duplicates.
func (t *Tree) Insert(k KeyT) {
	// Check if the root is full and split it. This is where the tree height increases.
	if t.root.n == 2*_T-1 {
		z := &node{}
		z.child[0] = t.root
		split(z, 0)
		t.root = z
	}
	// On each loop iteration we maintain the invariant that the node |x| is not full.
	x := t.root
	for {
		// Look for the key in the current node.
		i := find(x, k)
		y := x.child[i]
		if y == nil {
			// We cannot descend further down the tree; insert the key here.
			copy(x.key[i+1:x.n+1], x.key[i:x.n])
			x.n++
			x.key[i] = k
			return
		}
		// We are about to descend into the subtree rooted at |y|. Make sure
		// the next node down the path is not full.
		if y.n == 2*_T-1 {
			split(x, i)
			if k <= x.key[i] {
				y = x.child[i]
			} else {
				y = x.child[i+1]
			}
		}
		x = y
	}
}

// Deletes a key from the tree. Returns true if the key was present.
func (t *Tree) Delete(k KeyT) bool {
	// Call the function that does the actual work.
	b := del(t.root, k)
	// If the root has no keys but the tree is not empty, delete the root.
	// This is where the tree height decreases.
	if t.root.n == 0 && t.root.child[0] != nil {
		t.root = t.root.child[0]
	}
	return b
}

// Deletes a key from the tree.
func del(x *node, k KeyT) bool {
	// On each loop iteration we maintain the invariant that |x| holds
	// at least _T keys.
	for x != nil {
		// Search in the current node.
		i := find(x, k)
		y := x.child[i]
		// Check if we have found the key.
		if i < x.n && k == x.key[i] {
			if y == nil {
				// Key found in a leaf node: just remove the key.
				copy(x.key[i:x.n-1], x.key[i+1:x.n])
				copy(x.child[i:x.n-1], x.child[i+1:x.n])
				x.n--
				return true
			}
			// Key found in an internal node: try to replace it with its immediate
			// predeccessor or successor.
			if y.n >= _T {
				// can delete on the left
				x.key[i] = deleteMax(y)
				return true
			}
			z := x.child[i+1]
			if z.n >= _T {
				// can delete on the right
				x.key[i] = deleteMin(z)
				return true
			}
			// Nodes on both sides don't have enough keys: merge the nodes
			// together with the key and continue deleting from the merged node.
			merge(x, i)
		}
		// The search continues at |y|; make sure it has enough keys
		if y != nil && y.n < _T {
			if i > 0 {
				if x.child[i-1].n >= _T {
					// refill from the left sibling
					rotateRight(x, i-1)
				} else {
					// merge with the left sibling
					merge(x, i-1)
					y = x.child[i-1]
				}
			} else if i < x.n {
				if x.child[i+1].n >= _T {
					// refill from the right sibling
					rotateLeft(x, i)
				} else {
					// merge with the right sibling
					merge(x, i)
				}
			}
		}
		x = y
	}
	return false
}

// Performs an in-order traversal of the tree, calling the given function
// with each key as an argument.
func (t *Tree) InOrder(fn func(KeyT)) {
	inorder(t.root, fn)
}

func inorder(x *node, fn func(KeyT)) {
	if x == nil {
		return
	}
	for i := uint(0); i < x.n; i++ {
		inorder(x.child[i], fn)
		fn(x.key[i])
	}
	inorder(x.child[x.n], fn)
}

// Deletes the maximum element in the subtree and returns it.
func deleteMax(x *node) KeyT {
	for {
		// On each loop iteration we maintain the invariant that |x| holds
		// at least _T keys.
		z := x.child[x.n]
		if z == nil {
			// The node has is a leaf; the maximum key is the last one.
			x.n--
			return x.key[x.n]
		}
		// Descend into the rightmost child. Make sure it holds at least _T keys.
		if z.n < _T {
			y := x.child[x.n-1]
			if y.n < _T {
				// merge with the left sibling
				merge(x, x.n-1)
				z = y
			} else {
				// refill from the left sibling
				rotateRight(x, x.n-1)
			}
		}
		x = z
	}
}

// Deletes the minimum element in the subtree and returns it.
func deleteMin(x *node) KeyT {
	for {
		// On each loop iteration we maintain the invariant that |x| holds
		// at least _T keys.
		y := x.child[0]
		if y == nil {
			// The node is a leaf; the minimum key is the first one.
			k := x.key[0]
			copy(x.key[:], x.key[1:x.n])
			copy(x.child[:], x.child[1:x.n+1])
			x.child[x.n] = nil
			x.n--
			return k
		}
		// Descend into the leftmost child. Make sure it holds at least _T keys.
		if y.n < _T {
			z := x.child[1]
			if z.n < _T {
				// merge with the right sibling
				merge(x, 0)
			} else {
				// refill from the right sibling
				rotateLeft(x, 0)
			}
		}
		x = y
	}
}

// Finds a key in the current node. Retuns the first position, such that all keys
// on the left of it are strictly less than the given one.
func find(x *node, k KeyT) uint {
	i := uint(0)
	for i < x.n && k > x.key[i] {
		i++
	}
	return i
}

// Splits the i'th child of |x| into two nodes, containing _T-1 keys
// and lifts the median into |x|. The split node must be full and the parent
// node |x| must not be full.
func split(x *node, i uint) {
	y := x.child[i]
	z := &node{}
	z.n = _T - 1
	copy(z.key[:], y.key[_T:])
	copy(z.child[:], y.child[_T:])
	for j := uint(_T); j < _T+_T; j++ {
		y.child[j] = nil
	}
	y.n = _T - 1
	copy(x.key[i+1:x.n+1], x.key[i:x.n])
	copy(x.child[i+2:x.n+2], x.child[i+1:x.n+1])
	x.n++
	x.key[i] = y.key[_T-1]
	x.child[i+1] = z
}

// Merges the i'th and the i+1'th children of |x| rogether with
// the i'th key of |x|. The merged nodes must hold exactly _T-1 keys and
// the parent node must hold at least _T keys.
func merge(x *node, i uint) {
	y := x.child[i]
	z := x.child[i+1]
	y.key[y.n] = x.key[i]
	copy(y.key[y.n+1:], z.key[:z.n])
	copy(y.child[y.n+1:], z.child[:z.n+1])
	y.n += 1 + z.n
	copy(x.key[i:], x.key[i+1:x.n])
	copy(x.child[i+1:], x.child[i+2:x.n+1])
	x.child[x.n] = nil
	x.n--
}

// Refills the i'th child of |x| by moving the i'th key from the parent
// to the i'th child and the first key from the i+1'th child to the parent.
func rotateLeft(x *node, i uint) {
	y := x.child[i]
	z := x.child[i+1]
	y.n++
	y.key[y.n-1] = x.key[i]
	y.child[y.n] = z.child[0]
	x.key[i] = z.key[0]
	copy(z.key[:], z.key[1:z.n])
	copy(z.child[:], z.child[1:z.n+1])
	z.child[z.n] = nil
	z.n--
}

// Refills the i+1'th child of |x| by moving the i'th key from the parent
// to the i+1'th child and the last key from the i'th child to the parent.
func rotateRight(x *node, i uint) {
	y := x.child[i]
	z := x.child[i+1]
	z.n++
	copy(z.key[1:], z.key[:z.n])
	copy(z.child[1:], z.child[:z.n+1])
	z.key[0] = x.key[i]
	z.child[0] = y.child[y.n]
	x.key[i] = y.key[y.n-1]
	y.child[y.n] = nil
	y.n--
}

// Outputs a Graphviz representantion of the tree.
func (t *Tree) WriteDot() string {
	return "digraph BTree { node[shape=record]\n" + writeDot("n", t.root) + "}\n"
}

func writeDot(name string, x *node) string {
	var s string = fmt.Sprintf("%s[label=\"", name)
	for i := uint(0); i < x.n; i++ {
		s += fmt.Sprintf("<p%d> %d |", i, x.key[i])
	}
	s += fmt.Sprintf("<p%d>\"]\n", x.n)
	for i := uint(0); i <= x.n; i++ {
		if x.child[i] == nil {
			return s
		}
		cn := fmt.Sprintf("%s%02d", name, i)
		s += fmt.Sprintf("%s:p%d -> %s\n", name, i, cn)
		s += writeDot(cn, x.child[i])
	}
	return s
}
