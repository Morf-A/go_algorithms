package array

func Search(n int, f func(int) bool) int {
	p, r := 0, n
	for p < r {
		q := (p + r) / 2
		if f(q) {
			r = q
		} else {
			p = q + 1
		}
	}
	return r
}

func InsertionSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func SelectionSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		smallest := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[smallest] {
				smallest = j
			}
		}
		arr[i], arr[smallest] = arr[smallest], arr[i]
	}
}

func BinarySearch(arr []int, x int) int {
	p, r := 0, len(arr)-1
	for p < r {
		q := (p + r) / 2
		if arr[q] >= x {
			r = q
		} else {
			p = q + 1
		}
	}
	if arr[r] == x {
		return r
	} else {
		return -1
	}
}

func BinarySearchRec(arr []int, x, p, r int) int {
	if p > r {
		return -1
	}
	q := (p + r) / 2

	if arr[q] == x {
		return q
	}

	if arr[q] < x {
		p = q + 1
	} else { // arr[q] > x
		r = q - 1
	}

	return BinarySearchRec(arr, x, p, r)
}
