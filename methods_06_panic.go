package log

import (
	"fmt"
	"strconv"
	"strings"
)

// panic methods are not using ingestor due to the nature of panic.

func (l *Logger) Panic(args ...any) {
	if l.withJSON {
		var b strings.Builder

		b.WriteString(`{"level":`)
		b.WriteString(strconv.Quote(logLevels[LevelPanic]))
		b.WriteString(`,"msg":`)
		b.WriteString(strconv.Quote(fmt.Sprint(args...)))
		b.WriteByte('}')

		panic(b.String())
	}

	panic(fmt.Sprint(args...))
}

func (l *Logger) Panicf(format string, args ...any) {
	if l.withJSON {
		var b strings.Builder

		b.WriteString(`{"level":`)
		b.WriteString(strconv.Quote(logLevels[LevelPanic]))
		b.WriteString(`,"msg":`)
		b.WriteString(strconv.Quote(fmt.Sprintf(format, args...)))
		b.WriteByte('}')

		panic(b.String())
	}

	panic(fmt.Sprintf(format, args...))
}

func (l *Logger) Panicw(msg string, keysAndValues ...any) {
	if l.withJSON {
		if (len(keysAndValues) & 1) != 0 {
			panic(`{"level":"panic","error":"odd number of key-value arguments"}`)
		}

		var b strings.Builder

		b.WriteString(`{"level":`)
		b.WriteString(strconv.Quote(logLevels[LevelPanic]))
		b.WriteString(`,"msg":`)
		b.WriteString(strconv.Quote(msg))

		for i := 0; i < len(keysAndValues); i += 2 {
			key := keysAndValues[i]
			val := keysAndValues[i+1]

			ks, ok := key.(string)
			if !ok {
				panic(`{"level":"panic","error":"panicw: key must be string"}`)
			}

			b.WriteByte(',')
			b.WriteString(strconv.Quote(ks))
			b.WriteByte(':')

			switch v := val.(type) {
			case string:
				b.WriteString(strconv.Quote(v))
			case int:
				b.WriteString(strconv.Itoa(v))
			case bool:
				if v {
					b.WriteString("true")
				} else {
					b.WriteString("false")
				}
			default:
				b.WriteString(strconv.Quote(fmt.Sprint(v)))
			}
		}

		b.WriteByte('}')
		panic(b.String())
	}

	// non-JSON path
	if (len(keysAndValues) & 1) != 0 {
		panic("panicw: odd number of key-value arguments")
	}

	var b strings.Builder
	b.Grow(len(msg) + len(keysAndValues)*8)

	b.WriteString(msg)

	for i := 0; i < len(keysAndValues); i += 2 {
		b.WriteByte(' ')
		fmt.Fprint(&b, keysAndValues[i])
		b.WriteByte('=')
		fmt.Fprint(&b, keysAndValues[i+1])
	}

	panic(b.String())
}
