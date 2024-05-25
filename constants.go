package log

const delim = ": "

const (
	LevelNONE  = 0
	LevelINFO  = 1
	LevelWARN  = 2
	LevelDEBUG = 3
)

var logLevels = [4]string{"NONE", "INFO", "WARN", "DEBUG"}
