package log

const delim = ": "

const (
	LevelNONE  = 0
	LevelINFO  = 1
	LevelWARN  = 2
	LevelDEBUG = 3
	LevelERROR = 4
)

var logLevels = [5]string{
	"NONE",
	"INFO",
	"WARN",
	"DEBUG",
	"ERROR",
}
