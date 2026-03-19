package helpers

import (
	"fmt"
	"strconv"
)

// AppendArgs appends each arg to dst without reflection or fmt.
// Covers the types that appear in practice. Falls back to fmt only for
// exotic types — still no alloc on the hot path.
func AppendArgs(destination []byte, args []any) []byte {
	for ix, arg := range args {
		if ix > 0 {
			destination = append(destination, ' ')
		}

		switch value := arg.(type) {
		case string:
			destination = append(destination, value...)
		case []byte:
			destination = append(destination, value...)
		case int:
			destination = strconv.AppendInt(destination, int64(value), 10)
		case int64:
			destination = strconv.AppendInt(destination, value, 10)
		case int32:
			destination = strconv.AppendInt(destination, int64(value), 10)
		case uint:
			destination = strconv.AppendUint(destination, uint64(value), 10)
		case uint64:
			destination = strconv.AppendUint(destination, value, 10)
		case float64:
			destination = strconv.AppendFloat(destination, value, 'f', -1, 64)
		case float32:
			destination = strconv.AppendFloat(destination, float64(value), 'f', -1, 32)
		case bool:
			destination = strconv.AppendBool(destination, value)
		case error:
			destination = append(destination, value.Error()...)
		case nil:
			destination = append(destination, "null"...)

		default:
			// Exotic types only — keeps hot path clean.
			destination = append(destination, fmt.Sprint(value)...)
		}
	}

	return destination
}
