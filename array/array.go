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
