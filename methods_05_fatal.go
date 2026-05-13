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

	var buffer strings.Builder

	if l.withJSON {
		// 1. Start the JSON object
		buffer.WriteByte('{')

		// 2. Handle timestamp inside the object
		if l.fnTimestamp != nil {
			// Use a small temporary buffer for the timestamp
			tsBuf := l.fnTimestamp(nil)

			buffer.WriteString(`"ts":"`)
			// Assuming fnTimestamp appends the time string to the builder
			// or returns a string you write to it.
			buffer.Write(tsBuf)
			buffer.WriteString(`",`)
		}

		// 3. Add level and msg
		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[LevelFatal]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(fmt.Sprint(args...)))

		// 4. Close the object
		buffer.WriteByte('}')
	} else {
		// Non-JSON / Raw section
		if l.fnTimestamp != nil {
			tsBuf := l.fnTimestamp(nil)
			buffer.Write(tsBuf)
			buffer.WriteByte(' ')
		}

		// log level
		buffer.WriteByte('[')
		buffer.WriteString(logLevels[LevelFatal])
		buffer.WriteString("] ")

		// message
		buffer.WriteString(fmt.Sprint(args...))
	}

	msg = buffer.String()

	// Bypass ingestion: write directly to fatal writer
	_, _ = l.fatalWriter.Write([]byte(msg))

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	var msg string

	var buffer strings.Builder

	if l.withJSON {
		// 1. Start the JSON object
		buffer.WriteByte('{')

		// 2. Handle timestamp inside the object
		if l.fnTimestamp != nil {
			// Use a small temporary buffer for the timestamp
			tsBuf := l.fnTimestamp(nil)

			buffer.WriteString(`"ts":"`)
			// Assuming fnTimestamp appends the time string to the builder
			// or returns a string you write to it.
			buffer.Write(tsBuf)
			buffer.WriteString(`",`)
		}

		// 3. Add level and msg
		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[LevelFatal]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(fmt.Sprintf(format, args...)))

		// 4. Close the object
		buffer.WriteByte('}')
	} else {
		// Non-JSON / Raw section
		if l.fnTimestamp != nil {
			tsBuf := l.fnTimestamp(nil)
			buffer.Write(tsBuf)
			buffer.WriteByte(' ')
		}

		// log level
		buffer.WriteByte('[')
		buffer.WriteString(logLevels[LevelFatal])
		buffer.WriteString("] ")

		// message
		buffer.WriteString(fmt.Sprintf(format, args...))
	}

	msg = buffer.String()

	// Bypass ingestion: write directly to fatal writer
	_, _ = l.fatalWriter.Write([]byte(msg))

	os.Exit(1)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	var out string

	if (len(keysAndValues) & 1) != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")

		msg = "LOG_ERR(odd_args): " + msg
	}

	var buffer strings.Builder

	if l.withJSON {
		// 1. Start the JSON object
		buffer.WriteByte('{')

		// 2. Handle timestamp inside the object
		if l.fnTimestamp != nil {
			// Use a small temporary buffer for the timestamp
			tsBuf := l.fnTimestamp(nil)

			buffer.WriteString(`"ts":"`)
			// Assuming fnTimestamp appends the time string to the builder
			// or returns a string you write to it.
			buffer.Write(tsBuf)
			buffer.WriteString(`",`)
		}

		// 3. Add level and msg
		buffer.WriteString(`"level":`)
		buffer.WriteString(strconv.Quote(logLevels[LevelFatal]))
		buffer.WriteString(`,"msg":`)
		buffer.WriteString(strconv.Quote(msg))

		for i := 0; i < len(keysAndValues); i += 2 {
			key := keysAndValues[i]
			val := keysAndValues[i+1]

			ks, couldCast := key.(string)
			if !couldCast {
				out = `{"level":"fatal","error":"fatalw: key must be string"}`

				goto write_and_exit
			}

			buffer.WriteByte(',')
			buffer.WriteString(strconv.Quote(ks))
			buffer.WriteByte(':')

			switch value := val.(type) {
			case string:
				buffer.WriteString(strconv.Quote(value))
			case int:
				buffer.WriteString(strconv.Itoa(value))
			case bool:
				if value {
					buffer.WriteString("true")
				} else {
					buffer.WriteString("false")
				}
			default:
				buffer.WriteString(strconv.Quote(fmt.Sprint(value)))
			}
		}

		// 4. Close the object
		buffer.WriteByte('}')
	} else {
		// Non-JSON fatal path
		if l.fnTimestamp != nil {
			tsBuf := l.fnTimestamp(nil)
			buffer.Write(tsBuf)
			buffer.WriteByte(' ')
		}

		// log level
		buffer.WriteByte('[')
		buffer.WriteString(logLevels[LevelFatal])
		buffer.WriteString("] ")

		buffer.Grow(len(msg) + len(keysAndValues)*8)

		buffer.WriteString(msg)

		for i := 0; i < len(keysAndValues); i = i + 2 {
			buffer.WriteByte(' ')
			fmt.Fprint(&buffer, keysAndValues[i])
			buffer.WriteByte('=')
			fmt.Fprint(&buffer, keysAndValues[i+1])
		}
	}

	out = buffer.String()

write_and_exit:
	_, _ = l.fatalWriter.Write([]byte(out))

	os.Exit(1)
}
