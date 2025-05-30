// Wrapper around  'map' to make it work like a Set
package set

type Set[K comparable] struct {
	internal map[K]struct{}
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{
		internal: make(map[K]struct{}),
	}
}

func (s *Set[K]) Add(key K) {
	s.internal[key] = struct{}{}
}

func (s *Set[K]) Delete(key K) {
	delete(s.internal, key)
}

func (s *Set[K]) Contains(key K) bool {
	_, found := s.internal[key]
	return found
}

func (s *Set[K]) ForEach(f func(key K)) {
	for key := range s.internal {
		f(key)
	}
}

func (s *Set[K]) Len() int {
	return len(s.internal)
}
