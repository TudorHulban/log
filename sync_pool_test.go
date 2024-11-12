package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	p := NewPool[bytes.Buffer]()

	buf1 := p.Get()

	buf1.WriteString("x")

	_, _ = buf1.WriteTo(os.Stdout)

	require.NoError(t,
		p.Put(buf1),
	)

	buf2 := p.Get()

	_, _ = buf2.WriteTo(os.Stdout)

	require.NoError(t,
		p.Put(buf2),
	)

	require.Error(t,
		p.Put(buf2),
	)
}
