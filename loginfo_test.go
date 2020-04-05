package loginfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1Logger(t *testing.T) {
	a := assert.New(t)

	l, err := New(-99)
	a.Nil(err)

	l.Debug(1)
	l.Info(2)
}
