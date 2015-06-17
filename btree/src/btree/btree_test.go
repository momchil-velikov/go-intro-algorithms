package btree

import (
	"math/rand"
	"testing"
)

func assert(t *testing.T, b bool, args ...interface{}) {
	if !b {
		t.Error(args...)
	}
}

func assertf(t *testing.T, b bool, f string, args ...interface{}) {
	if !b {
		t.Errorf(f, args...)
	}
}

func randomShuffle(a []int64) []int64 {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
	}
	return a
}

func makeValues(n uint) []int64 {
	a := make([]int64, n)
	for i := uint(0); i < n; i++ {
		a[i] = int64(i)
	}
	return randomShuffle(a)
}

const _N = 100000

var (
	_A = makeValues(_N)
	_B = randomShuffle(_A)
)

func TestInsert(t *testing.T) {
	tr := New()
	for _, v := range _A {
		tr.Insert(KeyT(v))
	}
	for _, v := range _B {
		assertf(t, tr.Search(KeyT(v)), "key %d not present", v)
	}
}

func TestInOrder(t *testing.T) {
	tr := New()
	for _, v := range _A {
		tr.Insert(KeyT(v))
	}
	b := []KeyT(nil)
	tr.InOrder(func(k KeyT) {
		b = append(b, k)
	})

	assert(t, len(b) == _N, "Invalid length")
	for i := 1; i < _N; i++ {
		assert(t, b[i-1] <= b[i], "Inorder sequence not sorted")
	}
}

func TestSplit(t *testing.T) {
	y := &node{}
	y.n = 2*_T - 1
	for i := uint(0); i < y.n; i++ {
		y.key[i] = KeyT(i)
	}
	x := &node{}
	x.n = 1
	x.key[0] = 100
	x.child[0] = y
	split(x, 0)

	z := x.child[1]
	assert(t, x.n == 2 && y.n == _T-1 && z.n == _T-1, "Invalid key count in nodes")
	assert(t, x.key[0] == _T-1, "Invalid key pulled up")
	for i := uint(0); i < y.n; i++ {
		assert(t, y.key[i] == KeyT(i), "Invalid key in left child")
	}
	for i := uint(0); i < z.n; i++ {
		assert(t, z.key[i] == KeyT(i+_T), "Invalid key in right child")
	}
}

func TestMerge(t *testing.T) {
	y := &node{}
	y.n = _T - 1
	for i := uint(0); i < y.n; i++ {
		y.key[i] = KeyT(i)
	}
	z := &node{}
	z.n = _T - 1
	for i := uint(0); i < z.n; i++ {
		z.key[i] = KeyT(_T + i)
	}
	x := &node{}
	x.n = 1
	x.key[0] = _T - 1
	x.child[0] = y
	x.child[1] = z
	merge(x, 0)

	assert(t, x.n == 0, "Parent node not empty")
	assert(t, x.child[0] == y, "Lost pointer to left child")
	assert(t, x.child[1] == nil, "Stale pointer in parent")
	for i := uint(0); i <= y.n; i++ {
		assert(t, y.child[i] == nil, "Leaf contains non-null pointer at", i)
	}
	assert(t, y.n == 2*_T-1, "Merged node not full")
	for i := uint(0); i < y.n; i++ {
		assertf(t, y.key[i] == KeyT(i), "Wrong key in merged node: %d instead of %d", y.key[i], i)
	}
}

func TestRotateLeft(t *testing.T) {
	x := &node{n: 2, key: [2*_T - 1]KeyT{2, 4}}
	y := &node{n: 1, key: [2*_T - 1]KeyT{1}}
	z := &node{n: 1, key: [2*_T - 1]KeyT{3}}
	x.child[0] = y
	x.child[1] = z
	u := &node{}
	y.child[1] = u
	v := &node{}
	z.child[0] = v
	rotateLeft(x, 0)

	assert(t, x.n == 2, "Invalid key count in root: ", x.n)
	assert(t, x.key[0] == 3, "Invalid key at x[0]: ", x.key[0])
	assert(t, x.key[1] == 4, "Invalid key at x[1]: ", x.key[1])
	assert(t, x.child[0] == y, "Invalid left ptr")
	assert(t, x.child[1] == z, "Invalid right ptr")
	assert(t, y.n == 2, "Invalid key count in left: ", y.n)
	assert(t, y.key[0] == 1, "Invalid key at x[0]: ", y.key[0])
	assert(t, y.key[1] == 2, "Invalid key at x[0]: ", y.key[0])
	assert(t, y.child[1] == u, "Corrupted old last pointer on left")
	assert(t, y.child[2] == v, "Invalid last pointer on left")
	assert(t, z.n == 0, "Invalid key count in right: ", z.n)
}

func TestRotateRight(t *testing.T) {
	x := &node{n: 2, key: [2*_T - 1]KeyT{2, 4}}
	y := &node{n: 1, key: [2*_T - 1]KeyT{1}}
	z := &node{n: 1, key: [2*_T - 1]KeyT{3}}
	x.child[0] = y
	x.child[1] = z
	u := &node{}
	y.child[0] = &node{}
	y.child[1] = u
	v := &node{}
	z.child[0] = v
	z.child[1] = &node{}
	rotateRight(x, 0)

	assert(t, x.n == 2, "Invalid key count in root: ", x.n)
	assert(t, x.key[0] == 1, "Invalid key at x[0]: ", x.key[0])
	assert(t, x.key[1] == 4, "Invalid key at x[1]: ", x.key[1])
	assert(t, x.child[0] == y, "Invalid left ptr")
	assert(t, x.child[1] == z, "Invalid right ptr")
	assert(t, y.n == 0, "Invalid key count on left: ", y.n)
	assert(t, z.n == 2, "Invalid key count on right: ", z.n)
	assert(t, z.key[0] == 2, "Invalid key at z[0]: ", z.key[0])
	assert(t, z.key[1] == 3, "Invalid key at z[1]: ", z.key[1])
	assert(t, z.child[0] == u, "Invalid first pointer on left")
	assert(t, z.child[1] == v, "Corrupted second pointer on left")
}

func TestDeleteMinMax(t *testing.T) {
	const n = 50
	tr := New()
	for i := 0; i < n; i++ {
		tr.Insert(KeyT(i))
	}
	x := tr.root
	if x.n < 4 || x.child[3] == nil || x.child[4] == nil {
		t.Skip("Test skipped, tweak the value of _T")
	} else {
		k1 := deleteMax(x.child[3])
		assert(t, k1+1 == x.key[3], "Failed to find predeccessor of", x.key[3])
		k2 := deleteMin(x.child[4])
		assert(t, k2 == x.key[3]+1, "Failed to find successor of", x.key[3])
		assert(t, !tr.Search(k1), "Failed to remove", k1)
		assert(t, !tr.Search(k2), "Failed to remove", k2)
		b := []KeyT(nil)
		tr.InOrder(func(k KeyT) {
			b = append(b, k)
		})
		assert(t, len(b)+2 == n, "Invalid number of elements in the tree", len(b))
	}
}

func TestDelete(t *testing.T) {
	tr := New()
	for _, v := range _A {
		tr.Insert(KeyT(v))
	}
	for _, v := range _B {
		assertf(t, tr.Search(KeyT(v)), "key %d not present", v)
		assert(t, tr.Delete(KeyT(v)), "Delete failure")
		assertf(t, !tr.Search(KeyT(v)), "key %d present", v)
		assert(t, !tr.Delete(KeyT(v)), "Unexpected delete success")
	}
	assert(t, tr.root.n == 0, "Tree not empty")
}

func TestWriteDot(t *testing.T) {
	const n = 300
	tr := New()
	for i := 0; i < n; i++ {
		tr.Insert(KeyT(i))
	}
	s := tr.WriteDot()
	e := `digraph BTree { node[shape=record]
n[label="<p0> 63 |<p1> 127 |<p2> 191 |<p3>"]
n:p0 -> n00
n00[label="<p0> 7 |<p1> 15 |<p2> 23 |<p3> 31 |<p4> 39 |<p5> 47 |<p6> 55 |<p7>"]
n00:p0 -> n0000
n0000[label="<p0> 0 |<p1> 1 |<p2> 2 |<p3> 3 |<p4> 4 |<p5> 5 |<p6> 6 |<p7>"]
n00:p1 -> n0001
n0001[label="<p0> 8 |<p1> 9 |<p2> 10 |<p3> 11 |<p4> 12 |<p5> 13 |<p6> 14 |<p7>"]
n00:p2 -> n0002
n0002[label="<p0> 16 |<p1> 17 |<p2> 18 |<p3> 19 |<p4> 20 |<p5> 21 |<p6> 22 |<p7>"]
n00:p3 -> n0003
n0003[label="<p0> 24 |<p1> 25 |<p2> 26 |<p3> 27 |<p4> 28 |<p5> 29 |<p6> 30 |<p7>"]
n00:p4 -> n0004
n0004[label="<p0> 32 |<p1> 33 |<p2> 34 |<p3> 35 |<p4> 36 |<p5> 37 |<p6> 38 |<p7>"]
n00:p5 -> n0005
n0005[label="<p0> 40 |<p1> 41 |<p2> 42 |<p3> 43 |<p4> 44 |<p5> 45 |<p6> 46 |<p7>"]
n00:p6 -> n0006
n0006[label="<p0> 48 |<p1> 49 |<p2> 50 |<p3> 51 |<p4> 52 |<p5> 53 |<p6> 54 |<p7>"]
n00:p7 -> n0007
n0007[label="<p0> 56 |<p1> 57 |<p2> 58 |<p3> 59 |<p4> 60 |<p5> 61 |<p6> 62 |<p7>"]
n:p1 -> n01
n01[label="<p0> 71 |<p1> 79 |<p2> 87 |<p3> 95 |<p4> 103 |<p5> 111 |<p6> 119 |<p7>"]
n01:p0 -> n0100
n0100[label="<p0> 64 |<p1> 65 |<p2> 66 |<p3> 67 |<p4> 68 |<p5> 69 |<p6> 70 |<p7>"]
n01:p1 -> n0101
n0101[label="<p0> 72 |<p1> 73 |<p2> 74 |<p3> 75 |<p4> 76 |<p5> 77 |<p6> 78 |<p7>"]
n01:p2 -> n0102
n0102[label="<p0> 80 |<p1> 81 |<p2> 82 |<p3> 83 |<p4> 84 |<p5> 85 |<p6> 86 |<p7>"]
n01:p3 -> n0103
n0103[label="<p0> 88 |<p1> 89 |<p2> 90 |<p3> 91 |<p4> 92 |<p5> 93 |<p6> 94 |<p7>"]
n01:p4 -> n0104
n0104[label="<p0> 96 |<p1> 97 |<p2> 98 |<p3> 99 |<p4> 100 |<p5> 101 |<p6> 102 |<p7>"]
n01:p5 -> n0105
n0105[label="<p0> 104 |<p1> 105 |<p2> 106 |<p3> 107 |<p4> 108 |<p5> 109 |<p6> 110 |<p7>"]
n01:p6 -> n0106
n0106[label="<p0> 112 |<p1> 113 |<p2> 114 |<p3> 115 |<p4> 116 |<p5> 117 |<p6> 118 |<p7>"]
n01:p7 -> n0107
n0107[label="<p0> 120 |<p1> 121 |<p2> 122 |<p3> 123 |<p4> 124 |<p5> 125 |<p6> 126 |<p7>"]
n:p2 -> n02
n02[label="<p0> 135 |<p1> 143 |<p2> 151 |<p3> 159 |<p4> 167 |<p5> 175 |<p6> 183 |<p7>"]
n02:p0 -> n0200
n0200[label="<p0> 128 |<p1> 129 |<p2> 130 |<p3> 131 |<p4> 132 |<p5> 133 |<p6> 134 |<p7>"]
n02:p1 -> n0201
n0201[label="<p0> 136 |<p1> 137 |<p2> 138 |<p3> 139 |<p4> 140 |<p5> 141 |<p6> 142 |<p7>"]
n02:p2 -> n0202
n0202[label="<p0> 144 |<p1> 145 |<p2> 146 |<p3> 147 |<p4> 148 |<p5> 149 |<p6> 150 |<p7>"]
n02:p3 -> n0203
n0203[label="<p0> 152 |<p1> 153 |<p2> 154 |<p3> 155 |<p4> 156 |<p5> 157 |<p6> 158 |<p7>"]
n02:p4 -> n0204
n0204[label="<p0> 160 |<p1> 161 |<p2> 162 |<p3> 163 |<p4> 164 |<p5> 165 |<p6> 166 |<p7>"]
n02:p5 -> n0205
n0205[label="<p0> 168 |<p1> 169 |<p2> 170 |<p3> 171 |<p4> 172 |<p5> 173 |<p6> 174 |<p7>"]
n02:p6 -> n0206
n0206[label="<p0> 176 |<p1> 177 |<p2> 178 |<p3> 179 |<p4> 180 |<p5> 181 |<p6> 182 |<p7>"]
n02:p7 -> n0207
n0207[label="<p0> 184 |<p1> 185 |<p2> 186 |<p3> 187 |<p4> 188 |<p5> 189 |<p6> 190 |<p7>"]
n:p3 -> n03
n03[label="<p0> 199 |<p1> 207 |<p2> 215 |<p3> 223 |<p4> 231 |<p5> 239 |<p6> 247 |<p7> 255 |<p8> 263 |<p9> 271 |<p10> 279 |<p11> 287 |<p12>"]
n03:p0 -> n0300
n0300[label="<p0> 192 |<p1> 193 |<p2> 194 |<p3> 195 |<p4> 196 |<p5> 197 |<p6> 198 |<p7>"]
n03:p1 -> n0301
n0301[label="<p0> 200 |<p1> 201 |<p2> 202 |<p3> 203 |<p4> 204 |<p5> 205 |<p6> 206 |<p7>"]
n03:p2 -> n0302
n0302[label="<p0> 208 |<p1> 209 |<p2> 210 |<p3> 211 |<p4> 212 |<p5> 213 |<p6> 214 |<p7>"]
n03:p3 -> n0303
n0303[label="<p0> 216 |<p1> 217 |<p2> 218 |<p3> 219 |<p4> 220 |<p5> 221 |<p6> 222 |<p7>"]
n03:p4 -> n0304
n0304[label="<p0> 224 |<p1> 225 |<p2> 226 |<p3> 227 |<p4> 228 |<p5> 229 |<p6> 230 |<p7>"]
n03:p5 -> n0305
n0305[label="<p0> 232 |<p1> 233 |<p2> 234 |<p3> 235 |<p4> 236 |<p5> 237 |<p6> 238 |<p7>"]
n03:p6 -> n0306
n0306[label="<p0> 240 |<p1> 241 |<p2> 242 |<p3> 243 |<p4> 244 |<p5> 245 |<p6> 246 |<p7>"]
n03:p7 -> n0307
n0307[label="<p0> 248 |<p1> 249 |<p2> 250 |<p3> 251 |<p4> 252 |<p5> 253 |<p6> 254 |<p7>"]
n03:p8 -> n0308
n0308[label="<p0> 256 |<p1> 257 |<p2> 258 |<p3> 259 |<p4> 260 |<p5> 261 |<p6> 262 |<p7>"]
n03:p9 -> n0309
n0309[label="<p0> 264 |<p1> 265 |<p2> 266 |<p3> 267 |<p4> 268 |<p5> 269 |<p6> 270 |<p7>"]
n03:p10 -> n0310
n0310[label="<p0> 272 |<p1> 273 |<p2> 274 |<p3> 275 |<p4> 276 |<p5> 277 |<p6> 278 |<p7>"]
n03:p11 -> n0311
n0311[label="<p0> 280 |<p1> 281 |<p2> 282 |<p3> 283 |<p4> 284 |<p5> 285 |<p6> 286 |<p7>"]
n03:p12 -> n0312
n0312[label="<p0> 288 |<p1> 289 |<p2> 290 |<p3> 291 |<p4> 292 |<p5> 293 |<p6> 294 |<p7> 295 |<p8> 296 |<p9> 297 |<p10> 298 |<p11> 299 |<p12>"]
}
`
	assert(t, s == e, "Invalid Graphviz output")
}

func BenchmarkBTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tr := New()
		for _, v := range _A {
			tr.Insert(KeyT(v))
		}
		for i := 0; i < _N/2; i++ {
			tr.Delete(KeyT(_B[i]))
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int64]bool)
		for _, v := range _A {
			m[v] = true
		}
		for i := 0; i < _N/2; i++ {
			delete(m, _B[i])
		}
	}
}
