package deque

import (
	"math/rand"
	"testing"
)

func shuffle(a []int) []int {
	for i := len(a); i > 1; i-- {
		j := rand.Intn(i)
		a[i-1], a[j] = a[j], a[i-1]
	}
	return a
}

func makeValues(n uint) []int {
	a := make([]int, n)
	for i := uint(0); i < n; i++ {
		a[i] = int(i) + 1
	}
	return shuffle(a)
}

const _N = 1000000

func TestDeque(t *testing.T) {
	a := makeValues(_N)
	b := []int{}
	q := Deque{}
	i := 0
	for i < _N {
		nin := rand.Intn(1000)
		nout := rand.Intn(1000)
		for i < _N && nin > 0 {
			q.Push(a[i])
			i++
			nin--
		}
		for nout > 0 && !q.IsEmpty() {
			b = append(b, q.Pop().(int))
			nout--
		}
	}
	for !q.IsEmpty() {
		b = append(b, q.Pop().(int))
	}
	if len(a) != len(b) {
		t.Fatal("input and output length differ")
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("input and output sequences not identical\n%v\n%v\n", a, b)
		}
	}
}
