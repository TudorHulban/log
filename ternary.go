package log

func ternary[T any](condition bool, value1, value2 T) T {
	if condition {
		return value1
	}

	return value2
}
