package heap

type IntHeap struct {
	Size int
	Data []int64
}

func (h *IntHeap) Push(v int64) {
	h.Data = append(h.Data, v)
	h.Size++
	h.up(h.Size - 1)
}

func (h *IntHeap) Pop() int64 {
	v := h.Data[0]
	h.Data[0] = h.Data[h.Size-1]
	h.Size--
	h.down(0)
	return v
}

func (h *IntHeap) Len() int {
	return h.Size
}

func (h *IntHeap) Init(data []int64) {
	h.Size = len(data)
	h.Data = data
	h.makeHeap()
}

func (h *IntHeap) makeHeap() {
	n := h.Size / 2
	for i := n; i >= 0; i-- {
		h.down(i)
	}
}

func Sort(data []int64) {
	var h IntHeap
	h.Init(data)
	for i := h.Size; i > 0; i-- {
		h.Data[0], h.Data[i-1] = h.Data[i-1], h.Data[0]
		h.Size--
		h.down(0)
	}
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}

func parent(i int) int {
	return (i - 1) / 2
}

func (h *IntHeap) down(i int) {
	for {
		x := i
		l := left(i)
		r := right(i)
		if l < h.Size && h.Data[x] > h.Data[l] {
			x = l
		}
		if r < h.Size && h.Data[x] > h.Data[r] {
			x = r
		}
		if x == i {
			break
		}
		h.Data[x], h.Data[i] = h.Data[i], h.Data[x]
		i = x
	}
}

func (h *IntHeap) up(i int) {
	for i != 0 {
		x := parent(i)
		if h.Data[x] < h.Data[i] {
			break
		}
		h.Data[i], h.Data[x] = h.Data[x], h.Data[i]
		i = x
	}
}

func (h *IntHeap) Check() bool {
	n := h.Size / 2
	for i := 0; i < n; i++ {
		l, r := left(i), right(i)
		if l < h.Size && h.Data[i] > h.Data[l] {
			return false
		}
		if r < h.Size && h.Data[i] > h.Data[r] {
			return false
		}
	}
	return true
}
