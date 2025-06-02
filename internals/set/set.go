// Wrapper around  'map' to make it work like a Set
package set

type Set[K comparable] struct {
	internal map[K]bool
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{
		internal: make(map[K]bool),
	}
}

func (s *Set[K]) Add(key K) {
	_, ok := s.internal[key]
	if !ok {
		s.internal[key] = true
	}
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

func (s *Set[K]) ToSlice() []K {
	slice := make([]K, 0, len(s.internal))
	for key := range s.internal {
		slice = append(slice, key)
	}
	return slice
}

func (s *Set[K]) Size() int {
	return len(s.internal)
}
