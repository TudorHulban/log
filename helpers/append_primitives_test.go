package helpers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendFloat(t *testing.T) {
	tests := []struct { //nolint:govet
		description string
		value       float64
		prec        int
		expected    string
	}{
		// 1. Error-like cases
		{"nan", math.NaN(), 4, "nan"},
		{"pos_inf", math.Inf(+1), 4, "inf"},
		{"neg_inf", math.Inf(-1), 4, "-inf"},

		// 2. Normal cases (truncate, do not round)
		{"small_int", 42, 0, "42"},
		{"small_frac", 3.14159, 4, "3.1415"},
		{"large_int", 123456789012345.0, 0, "123456789012345"},
		{"large_frac", 987654321.123456, 6, "987654321.123456"},
		{"tiny", 0.00000012345, 10, "0.0000001234"},
		{"negative", -123.456, 3, "-123.456"},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				var dst []byte

				out := appendFloat(dst, tc.value, tc.prec)

				// Normalize to lowercase for safety
				actual := strings.ToLower(string(out))
				expected := strings.ToLower(tc.expected)

				// NaN requires special handling
				if math.IsNaN(tc.value) {
					require.Equal(t, expected, actual)
					return
				}

				require.Equal(t, expected, actual)
			},
		)
	}
}

func BenchmarkAppendFloat(b *testing.B) {
	tests := []struct {
		description string
		value       float64
		prec        int
	}{
		// 1. Error-like cases
		{"nan", math.NaN(), 4},
		{"pos_inf", math.Inf(+1), 4},
		{"neg_inf", math.Inf(-1), 4},

		// 2. Normal cases
		{"small_int", 42, 0},
		{"small_frac", 3.14159, 4},
		{"large_int", 123456789012345.0, 0},
		{"large_frac", 987654321.123456, 6},
		{"tiny", 0.00000012345, 10},
		{"negative", -123.456, 3},
	}

	for _, tc := range tests {
		fmt.Println("")

		b.Run(
			"1_custom_"+tc.description,
			func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				var dst []byte
				for i := 0; i < b.N; i++ {
					dst = dst[:0]
					out := appendFloat(dst, tc.value, tc.prec)
					require.NotNil(b, out)
				}
			},
		)

		b.Run(
			"2_stdlib_"+tc.description,
			func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				var dst []byte
				for i := 0; i < b.N; i++ {
					dst = dst[:0]
					out := strconv.AppendFloat(dst, tc.value, 'f', tc.prec, 64)
					require.NotNil(b, out)
				}
			},
		)
	}
}
