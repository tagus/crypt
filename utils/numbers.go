package utils

// Min determines the minimum among the provided values
func Min(vals ...int) int {
	min := vals[0]
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min
}
