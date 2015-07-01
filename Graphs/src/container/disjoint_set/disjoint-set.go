package disjoint_set

type SetId int

type rep struct {
	rank   uint
	parent SetId
}

type Sets struct {
	count uint
	reps  []rep
	data  []interface{}
}

func (s *Sets) Count() uint {
	return s.count
}

func (s *Sets) Create(data interface{}) SetId {
	id := SetId(len(s.reps))
	s.reps = append(s.reps, rep{rank: 1, parent: id})
	s.data = append(s.data, data)
	s.count++
	return id
}

func (s *Sets) Get(id SetId) interface{} {
	return s.data[id]
}

func (s *Sets) Union(a SetId, b SetId) {
	a, b = s.Find(a), s.Find(b)
	if a == b {
		return
	}
	if s.reps[a].rank > s.reps[b].rank {
		s.reps[b].parent = a
	} else {
		s.reps[a].parent = b
		if s.reps[a].rank == s.reps[b].rank {
			s.reps[b].rank++
		}
	}
	s.count--
}

func (s *Sets) Find(id SetId) SetId {
	x := int(id)
	for id != s.reps[id].parent {
		id = s.reps[id].parent
	}
	for id != s.reps[x].parent {
		x, s.reps[x].parent = int(s.reps[x].parent), id
	}
	return id
}
