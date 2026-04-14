package log

const (
	LevelNONE  uint8 = 0
	LevelINFO  uint8 = 1
	LevelWARN  uint8 = 2
	LevelDEBUG uint8 = 3
	LevelERROR uint8 = 4
)

var logLevels = [5]string{
	"NONE",
	"INFO",
	"WARN",
	"DEBUG",
	"ERROR",
}

func convertLevel(level Level) uint8 {
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
