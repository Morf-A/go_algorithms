package array

import (
	"fmt"
	"sort"
)

func ExampleRadixSort() {
	input := []string{"XI7FS6", "PL4ZQ2", "JI8FR9", "XL8FQ6", "PY2ZR5", "KV7WS9", "JL2ZV3", "KI4WR2"}
	fmt.Println(input)
	if err := RadixSort(input, 6); err != nil {
		panic(err)
	}
	fmt.Println(input)
}

func ExampleCountingSort() {
	input := []int{3, 6, 0, 2, 3, 0, 5, 6, 0, 3, 3, 2, 4}
	fmt.Println(input)
	if err := CountingSort(input); err != nil {
		panic(err)
	}
	fmt.Println(input)
}

func ExampleQuickSort() {
	toPartition := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8, 3}
	fmt.Println(toPartition)
	q := Partition(toPartition, 0, len(toPartition)-1)
	fmt.Println(q, toPartition)
	fmt.Println("-------------")
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	QuickSort(input)
	fmt.Println(input)
}

func ExampleHeapSort() {
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	HeapSort(input)
	fmt.Println(input)
}

func ExampleMergeSort() {
	toMerge := []int{-1, 1, 3, 5, 9, 10, 2, 4, 6, -1}
	fmt.Println(toMerge)
	Merge(toMerge, 1, 5, 8)
	fmt.Println(toMerge)
	fmt.Println("-------------")

	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	MergeSort(input)
	fmt.Println(input)
}

func ExampleInsertionSort() {
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	InsertionSort(input)
	fmt.Println(input)
}

func ExampleSelectionSort() {
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	SelectionSort(input)
	fmt.Println(input)
}

func ExampleBinarySearch() {
	input := []int{1, 4, 6, 8, 9, 13, 16, 19, 30, 35, 48, 50, 51}

	i := sort.Search(len(input), func(i int) bool { return input[i] >= 35 })
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = BinarySearch(input, 35)
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = Search(len(input), func(i int) bool { return input[i] >= 35 })
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = BinarySearch(input, 8)
	fmt.Printf("Number %d found at index %d\n", 8, i)

	i = BinarySearch(input, 5)
	fmt.Printf("Number %d found at index %d\n", 5, i)

	i = BinarySearchRec(input, 35, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = BinarySearchRec(input, 8, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 8, i)

	i = BinarySearchRec(input, 5, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 5, i)

}
