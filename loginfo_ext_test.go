package loginfo_test

/*
File details how to use logger.
*/

import (
	"os"
	"testing"

	"github.com/TudorHulban/loginfo"
)

func Test1ELogger(t *testing.T) {
	logger := loginfo.New(3, os.Stderr)
	logger.Print("0")
}
