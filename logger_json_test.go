package log

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	p1 := paramsJSONWCaller{
		timestamp: "12345",
		file:      "/home/some_file1",
		level:     "ERROR",
		message:   "abcd is above xyz",

		line: 70,
	}

	fmt.Println(string(jsonWCaller(&p1)))

	p2 := paramsJSONWCaller{
		timestamp: "67890",
		file:      "/home/some_file2",
		level:     "INFO",
		message:   "abcd is below xyz",

		line: 90,
	}

	fmt.Println(
		string(jsonWCaller(&p2)),
	)
}
