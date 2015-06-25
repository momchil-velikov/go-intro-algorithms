package bst

import "fmt"

type KeyType int64
const (
    maxKey = KeyType(0x7fffffffffffffff)
    minKey = -maxKey - 1
)

type Node struct {
        Key KeyType
        Left, Right *Node
}

type Tree struct {
    Root *Node
}

func makeNode(k KeyType) *Node {
    return &Node{Key: k, Left: nil, Right: nil}
}

func New() *Tree {
    return &Tree{Root: nil}
}

func (t *Tree) Insert(k KeyType) {
    var p *Node = nil
    x := t.Root
    for x != nil {
        p = x
        if k <= x.Key {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    z := makeNode(k)
    if p == nil {
        t.Root = z
    } else if k < p.Key {
        p.Left = z
    } else {
        p.Right = z
    }
}

func (t *Tree) Delete(k KeyType) {
    // Find node with the given key
    var p *Node = nil
    x := t.Root
    for x != nil && x.Key != k {
        p = x
        if k < x.Key {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    switch {
    // Key not found.
    case x == nil: /* nothing */;
    // Node to delete has no left child.
    case x.Left == nil:
        t.transplant(p, x, x.Right)
    // Node to delete has left, but no right child.
    case x.Right == nil:
        t.transplant(p, x, x.Left)
    default:
        // Node to delete has two children.
        z := detachMin(x, x.Right)
        z.Left = x.Left
        z.Right = x.Right
        t.transplant(p, x, z)
    }
}

func (t *Tree) transplant(p, u, v *Node) {
    if p == nil {
        t.Root = v
    } else if u == p.Left {
        p.Left = v
    } else { // u == p.Right
        p.Right = v
    }
}

func detachMin(p, x *Node) *Node {
    for x.Left != nil {
        p = x
        x = x.Left
    }
    if x == p.Right {
        p.Right = x.Right
    } else {
        p.Left = x.Right
    }
    return x
}

func (t *Tree) WriteDot(name string) string {
    var sn, se string
    _, sn = writeDotNodes(1, t.Root)
    _, se = writeDotEdges(1, t.Root)
    return "digraph " + name + "{\n" + sn + se + "}\n"
}

func writeDotNodes(n uint, nd *Node) (uint, string) {
    if nd == nil {
        return n + 1, fmt.Sprintf("n%d[shape=point]\n", n)
    }
    var s, sl, sr string
    s = fmt.Sprintf("n%d[label=\"%d\"]\n", n, int64(nd.Key))
    n, sl = writeDotNodes(n + 1, nd.Left)
    n, sr = writeDotNodes(n, nd.Right)
    return n, s + sl + sr
}

func writeDotEdges(n uint, nd *Node) (uint, string) {
    if nd == nil {
        return n + 1, ""
    }
    nl, sl := writeDotEdges(n + 1, nd.Left)
    nr, sr := writeDotEdges(nl, nd.Right)
    s := fmt.Sprintf("n%d -> n%d[arrowhead=none]\n", n, n + 1)
    s += fmt.Sprintf("n%d -> n%d[arrowhead=none]\n", n, nl)
    return nr, s + sl + sr
}

func (t *Tree) verifyBSTProperty() bool {
    return verifyBSTProperty(t.Root, minKey, maxKey)
}

func verifyBSTProperty(p *Node, low, high KeyType) bool {
    if p == nil {
        return true
    }
    if p.Key < low || p.Key > high {
        return false
    }
    return verifyBSTProperty(p.Left, low, p.Key) && verifyBSTProperty(p.Right, p.Key + 1, high)
}
