package helpers

func AppendInt(destination []byte, value int) []byte {
	var buf [20]byte // Enough space for int64 decimal representation

	i := len(buf)

	neg := value < 0
	if neg {
		value = -value
	}

	for {
		i--
		buf[i] = byte('0' + value%10)

		value = value / 10
		if value == 0 {
			break
		}
	}

	if neg {
		i--
		buf[i] = '-'
	}

	return append(destination, buf[i:]...)
}

func AppendUint64(destination []byte, value uint64) []byte {
	var buf [20]byte

	i := len(buf)

	for {
		i--
		buf[i] = byte('0' + value%10)

		value = value / 10
		if value == 0 {
			break
		}
	}

	return append(destination, buf[i:]...)
}

func AppendFloat(destination []byte, value float64, prec int) []byte {
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
	destination = AppendUint64(destination, intPart)

	// Fractional part
	if prec > 0 {
		destination = append(destination, '.')
		fractional := value - float64(intPart)

		for range prec {
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
