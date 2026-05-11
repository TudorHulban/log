package query

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	ansiRegex         = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	nonPrintableRegex = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F-\x9F]`)
)

// valuesMatch provides a "fuzzy" equality check for JSON-parsed data
func valuesMatch(actual, expected any) bool {
	// 1. Try standard deep equal first (handles strings, bools, and complex objects)
	if reflect.DeepEqual(actual, expected) {
		return true
	}

	// 2. Handle the "JSON Number" problem (float64 vs int/int64/etc)
	// We convert both to strings via %v to see if they represent the same value
	// This makes 12345 (int) match 12345 (float64)
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)

	return actualStr == expectedStr
}
