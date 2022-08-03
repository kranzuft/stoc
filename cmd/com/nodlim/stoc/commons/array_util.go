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

func StartsWith(first []rune, index int, seconds ...[]rune) bool {
	sLength := len(seconds)

	var skips []int
	// continue until we skip all indexes in seconds array
	for i := 0; len(skips) != len(seconds); i++ {
		for a := 0; a < sLength; a++ {
			if !Contains(skips, a) && charsAreNotEqualAt(first, index+i, seconds[a], i) {
				if len(seconds[a]) <= i {
					// a == i in length implying 'first' starts with the ath list in 'seconds'
					// since we are only looking for one match this is fine
					// len(a) < i is impossible since we return true as soon as it equals i but just in case we
					// capture it here since if len(a) < i it should not go to else condition
					return true
				} else {
					// if len(a) > i then it implies either
					// 		- List a did not match first
					//		- List a is larger than first
					skips = append(skips, a)
				}
			}
		}
	}

	return false
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func charsAreNotEqualAt(first []rune, fi int, second []rune, si int) bool {
	return len(second) <= si || len(first) <= fi || first[fi] != second[si]
}
