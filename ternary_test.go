package log

import (
	"fmt"
	"testing"
)

func TestTernary(t *testing.T) {
	conditionTrue, conditionFalse := "true", "false"

	fmt.Println(
		ternary(
			true,

			conditionTrue,
			conditionFalse,
		),
	)

	fmt.Println(
		ternary(
			false,

			conditionTrue,
			conditionFalse,
		),
	)
}
