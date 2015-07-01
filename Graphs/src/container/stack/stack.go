package stack

type Stack struct {
	stk []interface{}
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Size() int {
	return len(s.stk)
}

func (s *Stack) Push(item interface{}) {
	s.stk = append(s.stk, item)
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		panic("empty stack")
	}
	n := len(s.stk) - 1
	item := s.stk[n]
	if 4*n < cap(s.stk) {
		stk := make([]interface{}, n)
		copy(stk, s.stk)
		s.stk = stk
	} else {
		s.stk = s.stk[:n]
	}
	return item
}

func (s *Stack) Top() interface{} {
	if s.IsEmpty() {
		panic("empty stack")
	}
	return s.stk[len(s.stk)-1]
}
