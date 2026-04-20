package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// terminal event, should not have dependencies like the ingestor.

func (l *Logger) Fatal(args ...any) {
	var msg string

	if l.withJSON {
		var b strings.Builder

		b.WriteString(`{"level":`)
		b.WriteString(strconv.Quote(logLevels[LevelFatal]))
		b.WriteString(`,"msg":`)
		b.WriteString(strconv.Quote(fmt.Sprint(args...)))
		b.WriteByte('}')

		msg = b.String()
	} else {
		msg = fmt.Sprint(args...)
	}

	// Bypass ingestion: write directly to fatal writer
	_, _ = l.fatalWriter.Write([]byte(msg))

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	var msg string

	if l.withJSON {
		var b strings.Builder

		b.WriteString(`{"level":`)
		b.WriteString(strconv.Quote(logLevels[LevelFatal]))
		b.WriteString(`,"msg":`)
		b.WriteString(strconv.Quote(fmt.Sprintf(format, args...)))
		b.WriteByte('}')

		msg = b.String()
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	// Direct fatal bypass — never ingest
	_, _ = l.fatalWriter.Write([]byte(msg))

	os.Exit(1)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	var out string

	if l.withJSON {
		if (len(keysAndValues) & 1) != 0 {
			out = `{"level":"fatal","error":"odd number of key-value arguments"}`
		} else {
			var b strings.Builder

			b.WriteString(`{"level":`)
			b.WriteString(strconv.Quote(logLevels[LevelFatal]))
			b.WriteString(`,"msg":`)
			b.WriteString(strconv.Quote(msg))

			for i := 0; i < len(keysAndValues); i += 2 {
				key := keysAndValues[i]
				val := keysAndValues[i+1]

				ks, couldCast := key.(string)
				if !couldCast {
					out = `{"level":"fatal","error":"fatalw: key must be string"}`

					goto write_and_exit
				}

				b.WriteByte(',')
				b.WriteString(strconv.Quote(ks))
				b.WriteByte(':')

				switch value := val.(type) {
				case string:
					b.WriteString(strconv.Quote(value))
				case int:
					b.WriteString(strconv.Itoa(value))
				case bool:
					if value {
						b.WriteString("true")
					} else {
						b.WriteString("false")
					}
				default:
					b.WriteString(strconv.Quote(fmt.Sprint(value)))
				}
			}

			b.WriteByte('}')
			out = b.String()
		}
	} else {
		// Non-JSON fatal path
		if (len(keysAndValues) & 1) != 0 {
			out = "fatalw: odd number of key-value arguments"
		} else {
			var b strings.Builder

			b.Grow(len(msg) + len(keysAndValues)*8)

			b.WriteString(msg)

			for i := 0; i < len(keysAndValues); i += 2 {
				b.WriteByte(' ')
				fmt.Fprint(&b, keysAndValues[i])
				b.WriteByte('=')
				fmt.Fprint(&b, keysAndValues[i+1])
			}

			out = b.String()
		}
	}

write_and_exit:
	_, _ = l.fatalWriter.Write([]byte(out))

	os.Exit(1)
}
