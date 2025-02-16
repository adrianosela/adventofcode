package set

type Set[T comparable] map[T]struct{}

func New[T comparable](elems ...T) Set[T] {
	s := make(Set[T])
	for _, elem := range elems {
		s.Put(elem)
	}
	return s
}

func (s Set[T]) Copy() Set[T] {
	s2 := New[T]()
	for k := range s {
		s2[k] = struct{}{}
	}
	return s2
}

func (s Set[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

func (s Set[T]) Put(e T) {
	s[e] = struct{}{}
}

func (s Set[T]) Remove(e T) {
	delete(s, e)
}

func (s Set[T]) Size() int {
	return len(s)
}
