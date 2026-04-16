package log

import "fmt"

/*
LOG LEVEL TRUTH TABLE
---------------------

Numeric ordering (lowest → highest severity):
    NONE(0) < DEBUG(1) < INFO(2) < WARN(3) < ERROR(4)

Filtering rule:
    A log entry is emitted if entry.level >= logger.threshold

PRINT is outside the severity hierarchy and always emitted.

+----------------+---------+---------+---------+---------+---------+
| Threshold ↓    | DEBUG   | INFO    | WARN    | ERROR   | PRINT   |
+----------------+---------+---------+---------+---------+---------+
| NONE (0)       |   NO    |   NO    |   NO    |   NO    |  YES    |
| DEBUG (1)      |  YES    |  YES    |  YES    |  YES    |  YES    |
| INFO (2)       |   NO    |  YES    |  YES    |  YES    |  YES    |
| WARN (3)       |   NO    |   NO    |  YES    |  YES    |  YES    |
| ERROR (4)      |   NO    |   NO    |   NO    |  YES    |  YES    |
+----------------+---------+---------+---------+---------+---------+

Legend:
    YES → log entry is emitted
    NO  → log entry is suppressed
*/

type Level uint8

const (
	LevelNONE  Level = 0 // no logs except PRINT
	LevelDEBUG Level = 1 // everything logs
	LevelINFO  Level = 2 // suppress DEBUG
	LevelWARN  Level = 3 // suppress DEBUG, INFO
	LevelERROR Level = 4 // suppress DEBUG, INFO, WARN
)

func (l Level) String() string {
	switch l {
	case LevelNONE:
		return "NONE"
	case LevelDEBUG:
		return "DEBUG"
	case LevelINFO:
		return "INFO"
	case LevelWARN:
		return "WARN"
	case LevelERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

var logLevels = [5]string{
	"NONE",  // 0
	"DEBUG", // 1
	"INFO",  // 2
	"WARN",  // 3
	"ERROR", // 4
}

func convertLevel(level Level) Level {
	if level < LevelNONE {
		return LevelNONE
	}

	if level > LevelERROR {
		return LevelERROR
	}

	return level
}
