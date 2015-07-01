package deque

type Deque struct {
	front []interface{}
	back  []interface{}
}

func (q *Deque) IsEmpty() bool {
	return q.Size() == 0
}

func (q *Deque) Size() int {
	return len(q.front) + len(q.back)
}

func (q *Deque) Push(item interface{}) {
	q.back = append(q.back, item)
}

func moveToFront(q *Deque) {
	n := len(q.back)
	for i := n; i > 0; i-- {
		q.front = append(q.front, q.back[i-1])
	}
	if 4*n < cap(q.back) {
		q.back = []interface{}{}
	} else {
		q.back = q.back[:0]
	}
}

func (q *Deque) Pop() interface{} {
	if q.IsEmpty() {
		panic("empty deque")
	}
	if len(q.front) == 0 {
		moveToFront(q)
	}
	n := len(q.front) - 1
	item := q.front[n]
	if 4*n < cap(q.front) {
		f := make([]interface{}, n)
		copy(f, q.front)
		q.front = f
	} else {
		q.front = q.front[:n]
	}
	return item
}

func (q *Deque) Top() interface{} {
	if q.IsEmpty() {
		panic("empty deque")
	}
	if len(q.front) == 0 {
		moveToFront(q)
	}
	return q.front[len(q.front)-1]
}
