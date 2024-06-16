package utils

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Float64InSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func EqualSlices(a []string, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true

}
