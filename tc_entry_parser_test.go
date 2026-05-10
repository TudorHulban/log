package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type LogEntry struct {
	timestamp string
	keyValues map[string]any
	raw       string // The original, untouched line
}

func (e LogEntry) String() string {
	return e.raw
}

func (e LogEntry) IsRAW() bool {
	return len(e.raw) != 0 && (len(e.timestamp) == 0 || len(e.keyValues) == 0)
}

func (e LogEntry) HasTimestamp() bool {
	return len(e.keyValues) == 0
}

func (e LogEntry) HasKey(key string) (bool, any) {
	val, exists := e.keyValues[key]

	return exists, val
}

type LogEntries []LogEntry

func (e LogEntries) String() string {
	entries := make([]string, len(e))

	for ix, entry := range e {
		entries[ix] = entry.String()
	}

	return strings.Join(entries, "\n")
}

func (e LogEntries) haveTimestamp() bool {
	for _, item := range e {
		if !item.HasTimestamp() {
			return false
		}
	}

	return true
}

// HasKey checks if a key exists a specific number of times across all entries.
func (e LogEntries) HasKey(name string, noTimes uint) error {
	var count uint

	for _, item := range e {
		if exists, _ := item.HasKey(name); exists {
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

// HasKeyWithValue checks if a key with a specific value exists a specific number of times.
func (e LogEntries) HasKeyWithValue(name string, value any, noTimes uint) error {
	var count uint
	for _, item := range e {
		if exists, val := item.HasKey(name); exists {
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

func (e LogEntries) HasKeysWithValues(noTimes uint, kv ...any) error {
	if len(kv)%2 != 0 {
		return fmt.Errorf(
			"hasKeysWithValues requires an even number of kv arguments",
		)
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

			if exists, actual := item.HasKey(key); !exists || !valuesMatch(actual, kv[i+1]) {
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

// FilterBy returns a new subset of entries where the key matches the expected value.
func (e LogEntries) FilterBy(key string, value any) LogEntries {
	var filtered LogEntries

	for _, item := range e {
		if exists, val := item.HasKey(key); exists {
			if valuesMatch(val, value) {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

func (e LogEntries) First() LogEntry {
	if len(e) == 0 {
		return LogEntry{}
	}

	return e[0]
}

func (e LogEntries) Last() LogEntry {
	if len(e) == 0 {
		return LogEntry{}
	}

	return e[len(e)-1]
}

// SortByTimestamp reorders the entries based on the ts field.
// If desc is true, it sorts newest to oldest.
func (e LogEntries) SortByTimestamp(desc bool) LogEntries {
	// Create a copy to avoid mutating the original slice during a test
	sorted := make(LogEntries, len(e))
	copy(sorted, e)

	sort.SliceStable(
		sorted,
		func(i, j int) bool {
			if desc {
				return sorted[i].timestamp > sorted[j].timestamp
			}
			return sorted[i].timestamp < sorted[j].timestamp
		},
	)

	return sorted
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func NewTestEntries(from string) (LogEntries, error) {
	lines := strings.Split(strings.TrimSpace(from), "\n")
	entries := make(LogEntries, 0, len(lines))

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		entry := LogEntry{
			raw:       line, // Store the original line immediately
			keyValues: make(map[string]any),
		}

		// Strip ANSI and find JSON
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		idx := strings.IndexByte(cleanLine, '{')

		// If it is not JSON, we still keep it as a 'raw' entry
		// but it will not have keyValues or a timestamp.
		if idx != -1 {
			jsonPart := cleanLine[idx:]
			var rawMap map[string]any

			if err := json.Unmarshal([]byte(jsonPart), &rawMap); err == nil {
				for k, v := range rawMap {
					if k == "ts" {
						if tsStr, ok := v.(string); ok {
							entry.timestamp = tsStr
						}
					} else {
						entry.keyValues[k] = v
					}
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
