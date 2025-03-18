package datastructure

import (
	"container/heap"
	"errors"
)

// HeapType represents the type of heap (Min or Max).
type HeapType int

const (
	MinHeap HeapType = iota // Min-Heap (smallest value at the top)
	MaxHeap                 // Max-Heap (largest value at the top)
)

// Heap interface exposing public methods.
type Heap[T any] interface {
	Insert(value T)       // Insert or update an element in the heap
	Extract() (T, error)  // Extract the top element
	Peek() (T, error)     // Get the top element without removing it
	Remove(value T) error // Remove an element from the heap
	Len() int             // Get the size of the heap
}

// heapImpl struct implementing the Heap interface.
type heapImpl[T any] struct {
	data    []T
	hType   HeapType
	compare func(a, b T) int // Comparison function: returns -1, 0, or 1
}

// NewHeap creates a new heap (MinHeap/MaxHeap).
func NewHeap[T any](hType HeapType, compare func(a, b T) int) Heap[T] {
	h := &heapImpl[T]{
		data:    []T{},
		hType:   hType,
		compare: compare,
	}
	heap.Init(h) // Initialize the heap
	return h
}

// Len returns the number of elements in the heap.
func (h *heapImpl[T]) Len() int {
	return len(h.data)
}

// Less determines heap order based on heap type (private method).
func (h *heapImpl[T]) Less(i, j int) bool {
	cmp := h.compare(h.data[i], h.data[j])
	if h.hType == MinHeap {
		return cmp < 0 // MinHeap: smaller value at top
	}
	return cmp > 0 // MaxHeap: larger value at top
}

// Swap swaps elements at index i and j (private method).
func (h *heapImpl[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Push adds an element to the heap (private method).
func (h *heapImpl[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

// Pop removes and returns the root element (private method).
func (h *heapImpl[T]) Pop() any {
	if len(h.data) == 0 {
		return nil
	}
	n := len(h.data) - 1
	item := h.data[n]
	h.data = h.data[:n]
	return item
}

// Insert adds a new element or updates it if already present.
func (h *heapImpl[T]) Insert(value T) {
	for i, v := range h.data {
		if h.compare(v, value) == 0 { // Element already exists, update it
			h.data[i] = value
			heap.Fix(h, i) // Reorder the heap
			return
		}
	}
	heap.Push(h, value)
}

// Extract removes and returns the root element.
func (h *heapImpl[T]) Extract() (T, error) {
	if len(h.data) == 0 {
		var zeroValue T
		return zeroValue, errors.New("heap is empty")
	}
	return heap.Pop(h).(T), nil
}

// Peek returns the root element without removing it.
func (h *heapImpl[T]) Peek() (T, error) {
	if len(h.data) == 0 {
		var zeroValue T
		return zeroValue, errors.New("heap is empty")
	}
	return h.data[0], nil
}

// Remove finds and removes an element from the heap.
func (h *heapImpl[T]) Remove(value T) error {
	for i, v := range h.data {
		if h.compare(v, value) == 0 {
			heap.Remove(h, i) // Remove element while maintaining heap order
			return nil
		}
	}
	return errors.New("element not found in heap")
}
