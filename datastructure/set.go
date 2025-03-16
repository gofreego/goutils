package datastructure

type Set[T comparable] interface {
	Add(value T) bool
	Remove(value T) bool
	Contains(value T) bool
	Values() []T
}

type set[T comparable] struct {
	m map[T]struct{}
}

// Add implements Set. Returns true if the value was added to the set, false if it was already present.
func (s *set[T]) Add(value T) bool {
	if _, ok := s.m[value]; ok {
		return false
	}
	s.m[value] = struct{}{}
	return true
}

// Contains implements Set. Returns true if the value is in the set, false otherwise.
func (s *set[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}

// Remove implements Set. Returns true if the value was removed from the set, false if it was not present.
func (s *set[T]) Remove(value T) bool {
	if _, ok := s.m[value]; ok {
		delete(s.m, value)
		return true
	}
	return false
}

// Values implements Set. Returns all the values in the set.
func (s *set[T]) Values() []T {
	values := make([]T, 0, len(s.m))
	for value := range s.m {
		values = append(values, value)
	}
	return values
}

// NewSet returns a new Set.
func NewSet[T comparable](keys ...T) Set[T] {
	s := &set[T]{
		m: make(map[T]struct{}),
	}
	for _, key := range keys {
		s.Add(key)
	}
	return s
}
