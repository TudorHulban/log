package log

import (
	"fmt"

	"github.com/tudorhulban/log/helpers"
)

type fieldKind uint8

const (
	fieldKindString fieldKind = iota
	fieldKindNumeric
	fieldKindBool
)

type field struct {
	key          string
	valueString  string
	valueNumeric int

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

	case int:
		return field{
			key:          key,
			kind:         kindInt,
			valueNumeric: v,
		}

	case bool:
		return field{
			key:       key,
			kind:      kindBool,
			valueBool: v,
		}

	// You can add more typed cases here:
	// case uint:
	// case float64:
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
		buf = helpers.AppendInt(buf, fld.valueNumeric)

	case kindBool:
		buf = helpers.AppendBool(buf, fld.valueBool)
	}

	return append(buf, ' ')
}
