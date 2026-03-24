package log

import "fmt"

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
