package loginfo

import (
	"fmt"
	"testing"
)

func TestDebugColor(t *testing.T) {
	fmt.Println(ColourizeDebug("Hello"))
}
