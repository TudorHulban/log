package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type testEntry struct {
	timestamp string
	keyValues map[string]any
}

func (e testEntry) hasTimestamp() bool {
	return len(e.timestamp) != 0
}

func (e testEntry) hasKey(key string) (bool, any) {
	val, exists := e.keyValues[key]

	return exists, val
}

type testEntries []testEntry

func (e testEntries) haveTimestamp() bool {
	for _, item := range e {
		if !item.hasTimestamp() {
			return false
		}
	}

	return true
}

// hasKey checks if a key exists a specific number of times across all entries.
func (e testEntries) hasKey(name string, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, _ := item.hasKey(name); exists {
			count++
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q to appear %d times, but found it %d times",
			name,
			noTimes,
			count,
		)
	}

	return nil
}

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

// hasKeyWithValue checks if a key with a specific value exists a specific number of times.
func (e testEntries) hasKeyWithValue(name string, value any, noTimes uint) error {
	var count uint
	for _, item := range e {
		if exists, val := item.hasKey(name); exists {
			if valuesMatch(val, value) {
				count++
			}
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected key %q with value %v (%T) to appear %d times, but found %d",
			name,
			value,
			value,
			noTimes,
			count,
		)
	}

	return nil
}

func (e testEntries) hasKeysWithValues(noTimes uint, kv ...any) error {
	if len(kv)%2 != 0 {
		return fmt.Errorf("hasKeysWithValues requires an even number of kv arguments")
	}

	var count uint
	for _, item := range e {
		matchAll := true
		for i := 0; i < len(kv); i += 2 {
			key, ok := kv[i].(string)
			if !ok {
				return fmt.Errorf(
					"key at index %d must be a string",
					i,
				)
			}

			if exists, actual := item.hasKey(key); !exists || !valuesMatch(actual, kv[i+1]) {
				matchAll = false
				break
			}
		}

		if matchAll {
			count++
		}
	}

	if count != noTimes {
		return fmt.Errorf(
			"expected %d entries matching %v, but found %d",
			noTimes,
			kv,
			count,
		)
	}

	return nil
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func newTestEntries(from string) (testEntries, error) {
	lines := strings.Split(strings.TrimSpace(from), "\n")
	entries := make(testEntries, 0, len(lines))

	for _, line := range lines {
		cleanLine := strings.TrimSpace(line)
		if len(cleanLine) == 0 {
			continue
		}

		cleanLine = ansiRegex.ReplaceAllString(line, "")

		idx := strings.IndexByte(cleanLine, '{')
		if idx == -1 {
			continue
		}

		cleanLine = cleanLine[idx:]

		// Parse JSON into a temporary map
		var rawMap map[string]any
		if errUnmarshal := json.Unmarshal([]byte(cleanLine), &rawMap); errUnmarshal != nil {
			return nil,
				fmt.Errorf(
					"unmarshaling error: %w for line: %s",
					errUnmarshal,
					cleanLine,
				)
		}

		entry := testEntry{
			keyValues: make(map[string]any),
		}

		// Separate the timestamp from the rest of the data
		for k, v := range rawMap {
			if k == "ts" {
				if tsStr, ok := v.(string); ok {
					entry.timestamp = tsStr
				}
			} else {
				entry.keyValues[k] = v
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
