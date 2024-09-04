package log

func convertLevel(level Level) int8 {
	switch {
	case level < 1:
		return LevelNONE
	case level == 1:
		return LevelINFO
	case level == 2:
		return LevelWARN
	case level == 3:
		return LevelDEBUG

	default:
		return LevelERROR
	}
}
