package log

import (
	"fmt"
	"strconv"
)

type fieldKind uint8

const (
	kindString fieldKind = iota
	kindInt
	kindBool
	kindFloat
)

type field struct {
	key         string
	valueString string
	valueInt    int64
	valueFloat  float64

	kind fieldKind

	valueBool bool
}

func makeField(key string, value any) field {
	switch v := value.(type) {
	case string:
		return field{
			key:         key,
			kind:        kindString,
			valueString: v,
		}

	case int64:
		return field{
			key:      key,
			kind:     kindInt,
			valueInt: v,
		}

	case bool:
		return field{
			key:       key,
			kind:      kindBool,
			valueBool: v,
		}

	case float64:
		return field{
			key:        key,
			kind:       kindFloat,
			valueFloat: v,
		}

	// TODO: add more typed cases here?
	// case uint:
	// case error:
	// etc.

	default:
		// Fallback: convert to string once
		return field{
			key:         key,
			kind:        kindString,
			valueString: fmt.Sprint(v),
		}
	}
}

func makeFieldPtr(key string, value any) *field {
	fld := makeField(key, value)

	return &fld
}

func appendField(buf []byte, fld *field) []byte {
	buf = append(buf, fld.key...)
	buf = append(buf, '=')

	switch fld.kind {
	case kindString:
		buf = append(buf, fld.valueString...)

	case kindInt:
		buf = strconv.AppendInt(buf, fld.valueInt, 10)

	case kindBool:
		buf = strconv.AppendBool(buf, fld.valueBool)
	}

	return append(buf, ' ')
}
