package datastructure

type Stack[T any] interface {
	Top() T
	Push(k T)
	Pop() T
	IsEmpty() bool
}

type stack[T any] struct {
	arr []T
	i   int
}

func NewStack[T any]() Stack[T] {
	return &stack[T]{
		arr: make([]T, 0),
		i:   -1,
	}
}

// IsEmpty implements Stack.
func (s *stack[T]) IsEmpty() bool {
	return s.i == -1
}

// Pop implements Stack. please check IsEmpty before calling this function
func (s *stack[T]) Pop() T {
	if s.i == -1 {
		panic("no objects in queue")
	}
	obj := s.arr[s.i]
	s.i--
	return obj
}

// Push implements Stack.
func (s *stack[T]) Push(k T) {
	s.i++
	if s.i == cap(s.arr) {
		s.arr = append(s.arr, k)
	} else {
		s.arr[s.i] = k
	}
}

// Top implements Stack. Check IsEmpty before calling this function
func (s *stack[T]) Top() T {
	if s.i == -1 {
		panic("no objects in queue")
	}
	return s.arr[s.i]
}
