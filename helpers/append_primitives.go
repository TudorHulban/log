package helpers

import "strconv"

func appendFloat(destination []byte, value float64, precision int) []byte {
	// Handle NaN and Inf explicitly
	if value != value {
		return append(destination, 'n', 'a', 'n')
	}

	if value > 1e308 {
		return append(destination, 'i', 'n', 'f')
	}

	if value < -1e308 {
		return append(destination, '-', 'i', 'n', 'f')
	}

	// Sign
	if value < 0 {
		destination = append(destination, '-')
		value = -value
	}

	// Integer part
	intPart := uint64(value)
	destination = strconv.AppendUint(destination, intPart, 10)

	// Fractional part
	if precision > 0 {
		destination = append(destination, '.')
		fractional := value - float64(intPart)

		for range precision {
			fractional *= 10
			digit := uint64(fractional)
			destination = append(destination, byte('0'+digit))
			fractional -= float64(digit)
		}
	}

	return destination
}

func AppendBool(destination []byte, value bool) []byte {
	if value {
		return append(destination, 't', 'r', 'u', 'e')
	}

	return append(destination, 'f', 'a', 'l', 's', 'e')
}

func AppendError(destination []byte, err error) []byte {
	if err == nil {
		return append(destination, 'n', 'i', 'l')
	}

	return append(destination, err.Error()...)
}
