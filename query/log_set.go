package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type LogSet []LogRecord

func NewLogset(from string) (LogSet, error) {
	lines := strings.Split(strings.TrimSpace(from), "\n")
	entries := make(LogSet, 0, len(lines))

	for _, line := range lines {
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		cleanLine = nonPrintableRegex.ReplaceAllString(cleanLine, "")

		if len(strings.TrimSpace(cleanLine)) == 0 {
			continue
		}

		entry := LogRecord{
			raw:       line, // Store the original line immediately
			keyValues: make(map[string]any),
		}

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

func (e LogSet) String() string {
	entries := make([]string, len(e))

	for ix, entry := range e {
		entries[ix] = entry.String()
	}

	return strings.Join(entries, "\n")
}

func (e LogSet) WithTimestamp() LogSet {
	var filtered LogSet

	for _, item := range e {
		if item.HasTimestamp() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (e LogSet) WithNoTimestamp() LogSet {
	var filtered LogSet

	for _, item := range e {
		if !item.HasTimestamp() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// HasKey checks if a key exists a specific number of times across all entries.
func (e LogSet) HasKey(name string, noTimes uint) error {
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

// HasKeyWithValue checks if a key with a specific value exists a specific number of times.
func (e LogSet) HasKeyWithValue(name string, value any, noTimes uint) error {
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

func (e LogSet) HasKeysWithValues(noTimes uint, kv ...any) error {
	if len(kv)%2 != 0 {
		return errors.New(
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
func (e LogSet) FilterBy(key string, value any) LogSet {
	var filtered LogSet

	for _, item := range e {
		if exists, val := item.HasKey(key); exists {
			if valuesMatch(val, value) {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

func (e LogSet) First() LogRecord {
	if len(e) == 0 {
		return LogRecord{}
	}

	return e[0]
}

func (e LogSet) Last() LogRecord {
	if len(e) == 0 {
		return LogRecord{}
	}

	return e[len(e)-1]
}

// SortByTimestamp reorders the entries based on the ts field.
// If desc is true, it sorts newest to oldest.
func (e LogSet) SortByTimestamp(desc bool) LogSet {
	// Create a copy to avoid mutating the original slice during a test
	sorted := make(LogSet, len(e))
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
