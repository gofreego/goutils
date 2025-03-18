package main

import (
	"fmt"

	"github.com/gofreego/goutils/datastructure"
)

func main() {
	f := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	// Min-Heap for integers

	minHeap := datastructure.NewHeap[int](datastructure.MinHeap, f)

	minHeap.Insert(10)
	minHeap.Insert(5)
	minHeap.Insert(15)
	minHeap.Insert(2)

	value, _ := minHeap.Peek()
	fmt.Println("Min:", value) // Output: Min: 2
	val, _ := minHeap.Extract()
	fmt.Println("Extracted:", val) // Output: Extracted: 2

	// Max-Heap for integers
	maxHeap := datastructure.NewHeap[int](datastructure.MaxHeap, f)

	maxHeap.Insert(10)
	maxHeap.Insert(5)
	maxHeap.Insert(15)
	maxHeap.Insert(2)

	value, _ = maxHeap.Peek()
	fmt.Println("Max:", value) // Output: Max: 15
}
