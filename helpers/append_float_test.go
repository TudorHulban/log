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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAppendFloat/1_custom_nan-16         	 8092822	       149.7 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_nan-16         	 7834503	       151.4 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_pos_inf-16     	 8180378	       145.6 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_pos_inf-16     	 7891600	       151.6 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_neg_inf-16     	 8386286	       144.4 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_neg_inf-16     	 7895463	       150.8 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_small_int-16   	 8076417	       148.8 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_small_int-16   	 7080556	       169.2 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_small_frac-16  	 7612082	       157.4 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_small_frac-16  	 6586550	       180.5 ns/op	       8 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_large_int-16   	 7126051	       167.4 ns/op	      16 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_large_int-16   	 6664712	       179.5 ns/op	      16 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_large_frac-16  	 6707352	       178.5 ns/op	      16 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_large_frac-16  	 6325773	       189.5 ns/op	      16 B/op	       1 allocs/op

// BenchmarkAppendFloat/1_custom_tiny-16        	 6114620	       195.9 ns/op	      24 B/op	       2 allocs/op
// BenchmarkAppendFloat/2_stdlib_tiny-16        	 5817366	       205.4 ns/op	      24 B/op	       2 allocs/op

// BenchmarkAppendFloat/1_custom_negative-16    	 7627508	       156.9 ns/op	       8 B/op	       1 allocs/op
// BenchmarkAppendFloat/2_stdlib_negative-16    	 6720366	       177.7 ns/op	       8 B/op	       1 allocs/op

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

				for b.Loop() {
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

				for b.Loop() {
					dst = dst[:0]

					out := strconv.AppendFloat(dst, tc.value, 'f', tc.prec, 64)
					require.NotNil(b, out)
				}
			},
		)
	}
}
