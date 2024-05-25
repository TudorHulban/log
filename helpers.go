package log

func convertLevel(level Level) int8 {
	switch {
	case level < 1:
		return 0
	case level == 1:
		return 1
	case level == 2:
		return 2
	default:
		return 3
	}
}
