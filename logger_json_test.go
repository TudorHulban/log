package log

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	p1 := paramsJSONWCaller{
		timestamp: "12345",
		file:      "/home/some_file",
		level:     "INFO",
		message:   "abcd is xyz",

		line: 70,
	}

	fmt.Println(string(jsonWCaller(&p1)))

	p2 := paramsJSONWCaller{
		timestamp: "12345",
		file:      "/home/some_file",
		level:     "INFO",
		message:   "abcd is xyz",

		line: 90,
	}

	fmt.Println(string(jsonWCaller(&p2)))
}
