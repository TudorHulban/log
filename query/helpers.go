package query

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

	// Matches:
	// 0x00-0x08 (Control chars: Null, Bell, Backspace, etc.)
	// 0x0B-0x0C (Vertical tab, Form feed)
	// 0x0E-0x1F (Remaining C0 controls: Shift Out, Unit Separator, etc.)
	// 0x7F      (DEL - Delete character)
	// 0x80-0x9F (C1 Control chars: Padding, Break, Start of string, etc.)
	// 0xA0-0xFF (Extended ASCII: Non-breaking space, symbols, and Latin letters)
	nonPrintableRegex = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F-\xFF]`)
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
