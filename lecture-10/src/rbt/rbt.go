package rbt

import "fmt"

type KeyType int64
const (
    maxKey = KeyType(0x7fffffffffffffff)
    minKey = -maxKey - 1
)

const (
    Red = 0
    Black = 1
)

type Node struct {
    Key KeyType
    Parent, Left, Right *Node
    Color uint8
}

type Tree struct {
    Root, Nil *Node
    nilN Node
}


func New() *Tree {
    var t Tree
    t.nilN.Color = Black
    t.Nil = &t.nilN
    t.Root = t.Nil
    return &t
}

func (t *Tree) makeNode(k KeyType) *Node {
    return &Node{Key: k, Parent: t.Nil, Left: t.Nil, Right: t.Nil, Color: Red}
}

func (t *Tree) Insert(k KeyType) {
    y := t.Nil
    x := t.Root
    for x != t.Nil {
        y = x
        if k <= x.Key {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    z := t.makeNode(k)
    z.Parent = y
    if y == t.Nil {
        t.Root = z
    } else if k <= y.Key {
        y.Left = z
    } else {
        y.Right = z
    }
    t.insertFixup(z)
}

func (t *Tree) insertFixup(z *Node) {
    for z.Parent.Color == Red {
        y := z.Parent.Parent
        if z.Parent == y.Left {
            if y.Right.Color == Red {
                // Case 1: recolor
                y.Color = Red
                y.Left.Color = Black
                y.Right.Color = Black
                z = y
            } else { // y.Right.Color == Black
                if z == z.Parent.Right {
                // Case 2: straighten
                    z = z.Parent
                    t.leftRotate(z)
                }
                // Case 3: rotate/balance
                z.Parent.Color = Black
                y.Color = Red
                t.rightRotate(y)
            }
        } else { // z.Parent == y.Right
            if y.Left.Color == Red {
                // Case 1: recolor
                y.Color = Red
                y.Left.Color = Black
                y.Right.Color = Black
                z = y
            } else { // y.Left.Color == Black
                if z == z.Parent.Left {
                    // Case 2: straighten
                    z = z.Parent
                    t.rightRotate(z)
                }
                // Case 3: rotate/balance
                z.Parent.Color = Black
                y.Color = Red
                t.leftRotate(y)
            }
        }
    }
    t.Root.Color = Black
}

func (t *Tree) Delete(k KeyType) {
    z := t.Root
    for z != t.Nil && z.Key != k {
        if k < z.Key {
            z = z.Left
        } else {
            z = z.Right
        }
    }
    var x *Node
    color := z.Color
    switch {
    case z == t.Nil:
        return // Key not found.
    case z.Left == t.Nil:
        // Node to delete has at most a right child.
        x = z.Right
        t.transplant(z, x)
    case z.Right == t.Nil:
        // Node to delete has only a left child
        x = z.Left
        t.transplant(z, x)
    default:
        // Node to delete has two children
        y := t.min(z.Right)
        color = y.Color
        x = y.Right
        if y.Parent == z {
            x.Parent = y
        } else {
            t.transplant(y, y.Right)
            y.Right = z.Right
            y.Right.Parent = y
        }
        t.transplant(z, y)
        y.Left = z.Left
        y.Left.Parent = y
        y.Color = z.Color
    }
    if color == Black {
        t.deleteFixup(x)
    }
}

func (t *Tree) deleteFixup(x *Node) {
    for x != t.Root && x.Color == Black {
        if x == x.Parent.Left {
            w := x.Parent.Right
            if w.Color == Red {
                // Case 1: make the sibling black
                t.leftRotate(x.Parent)
                w.Color = Black
                x.Parent.Color = Red
                w = x.Parent.Right
            }
            if w.Left.Color == Black && w.Right.Color == Black {
                // Case 2: push the double blackness up
                w.Color = Red
                x = x.Parent
            } else {
                // Case 3:
                if w.Left.Color == Red && w.Right.Color == Black {
                    w.Color = Red
                    w.Left.Color = Black
                    t.rightRotate(w)
                    w = x.Parent.Right
                }
                // Case 4: remove the double blackness; end
                w.Color = x.Parent.Color
                w.Right.Color = Black
                x.Parent.Color = Black
                t.leftRotate(x.Parent)
                x = t.Root
            }
        } else { // x == x.Parent.Right
            w := x.Parent.Left
            if w.Color == Red {
                // Case 1: make the sibling black
                t.rightRotate(x.Parent)
                w.Color = Black
                x.Parent.Color = Red
                w = x.Parent.Left
            }
            if w.Left.Color == Black && w.Right.Color == Black {
                // Case 2: push the double blackness up
                w.Color = Red
                x = x.Parent
            } else {
                // Case 3:
                if w.Left.Color == Black && w.Right.Color == Red {
                    w.Color = Red
                    w.Right.Color = Black
                    t.leftRotate(w)
                    w = x.Parent.Left
                }
                // Case 4: remove the double blackness; end
                w.Color = x.Parent.Color
                w.Left.Color = Black
                x.Parent.Color = Black
                t.rightRotate(x.Parent)
                x = t.Root
            }
        }
    }
    x.Color = Black
}

func (t *Tree) leftRotate(x *Node) {
    y := x.Right
    y.Parent = x.Parent
    if x == t.Root {
        t.Root = y
    } else {
        if x == x.Parent.Left {
            x.Parent.Left = y
        } else {
            x.Parent.Right = y
        }
    }
    x.Right = y.Left
    if y.Left != t.Nil {
        y.Left.Parent = x
    }
    y.Left = x
    x.Parent = y
}

func (t *Tree) rightRotate(y *Node) {
    x := y.Left
    x.Parent = y.Parent
    if y == t.Root {
        t.Root = x
    } else {
        if y == y.Parent.Left {
            y.Parent.Left = x
        } else {
            y.Parent.Right = x
        }
    }
    y.Left = x.Right
    if x.Right != t.Nil {
        x.Right.Parent = y
    }
    x.Right = y
    y.Parent = x
}

func (t *Tree) transplant(u, v *Node) {
    if u.Parent == t.Nil {
        t.Root = v
    } else if u == u.Parent.Left {
        u.Parent.Left = v
    } else {
        u.Parent.Right = v
    }
    v.Parent = u.Parent
}

func (t *Tree) min(u *Node) *Node {
    for u.Left != t.Nil {
        u = u.Left
    }
    return u
}

func (t *Tree) WriteDot() string {
    var sn, se string
    _, sn = t.writeDotNodes(1, t.Root)
    _, se = t.writeDotEdges(1, t.Root)
    return "digraph RBTree {\n" + sn + se + "}\n"
}

const leafString = "[shape=point]"
func (t *Tree) writeDotNodes(n uint, nd *Node) (uint, string) {
    if nd == t.Nil {
        return n + 1, fmt.Sprintf("n%d%s\n", n, leafString)
    }
    var color string
    if nd.Color == Red {
        color = "red"
    } else {
        color = "black"
    }
    var s, sl, sr string
    s = fmt.Sprintf("n%d[label=\"%d\", color=%s]\n", n, int64(nd.Key), color)
    n, sl = t.writeDotNodes(n + 1, nd.Left)
    n, sr = t.writeDotNodes(n, nd.Right)
    return n, s + sl + sr
}

func (t *Tree) writeDotEdges(n uint, nd *Node) (uint, string) {
    if nd == t.Nil {
        return n + 1, ""
    }
    nl, sl := t.writeDotEdges(n + 1, nd.Left)
    nr, sr := t.writeDotEdges(nl, nd.Right)
    s := fmt.Sprintf("n%d -> n%d[arrowhead=none]\n", n, n + 1)
    s += fmt.Sprintf("n%d -> n%d[arrowhead=none]\n", n, nl)
    return nr, s + sl + sr
}

func (t *Tree) validateTreeConnectivity() bool {
    return t.validateConnectivity(t.Root)
}

func (t *Tree) validateConnectivity(nd *Node) bool {
    if nd == t.Nil {
        return true
    }
    if nd == nil {
        return false
    }
    if !t.validateConnectivity(nd.Left) || !t.validateConnectivity(nd.Right) {
        return false
    }
    if nd.Left != t.Nil && nd.Left.Parent != nd {
        return false
    }
    if nd.Right != t.Nil && nd.Right.Parent != nd {
        return false
    }
    return true
}

func (t *Tree) validateTreeBlackHeight() bool {
    _, e := t.validateBlackHeight(t.Root)
    return e
}

func (t *Tree) validateBlackHeight(x *Node) (uint, bool) {
    if x == t.Nil {
        return 0, true
    }
    h1, e1 := t.validateBlackHeight(x.Left)
    h2, e2 := t.validateBlackHeight(x.Left)
    if !e1 || !e2 {
        return 0, false
    }
    if h1 != h2 {
        return 0, false
    }
    return h1 + uint(x.Color), true
}

func (t *Tree) validateTreeRedChildren() bool {
    return t.validateRedChildren(t.Root)
}

func (t *Tree) validateRedChildren(x *Node) bool {
    if x == t.Nil {
        return true
    }
    if x.Color == Red && (x.Left.Color != Black || x.Right.Color != Black) {
        return false
    }
    return t.validateRedChildren(x.Left) && t.validateRedChildren(x.Right)
}

func (t *Tree) verifyTreeBSTProperty() bool {
    return t.verifyBSTProperty(t.Root, minKey, maxKey)
}

func (t *Tree) verifyBSTProperty(p *Node, low, high KeyType) bool {
    if p == t.Nil {
        return true
    }
    if p.Key < low || p.Key > high {
        return false
    }
    return t.verifyBSTProperty(p.Left, low, p.Key) &&
           t.verifyBSTProperty(p.Right, p.Key + 1, high)
}
