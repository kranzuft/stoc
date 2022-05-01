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
func StartsWith(first []rune, index int, second []rune) bool {
	if (len(first) - index) < (len(second)) {
		return false
	}

	// determine the argument with the smallest length
	var minLength = (len(first) - index)
	if len(second) < minLength {
		minLength = len(second)
	}

	for i := 0; i < minLength; i++ {
		if first[index+i] != second[i] {
			return false
		}
	}

	return true
}
