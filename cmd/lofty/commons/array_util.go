package commons

// LastIndexOf Get the last index of the first element in the array that matches the predicate
func LastIndexOf[E any](array []E, predicate func(E) bool) int {
	for i := len(array) - 1; i >= 0; i-- {
		if predicate(array[i]) {
			return i
		}
	}
	return -1
}
