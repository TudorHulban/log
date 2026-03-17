package arena_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/arena"
)

func TestHowToUse(t *testing.T) {
	rawLogger := arena.NewRawLogger(arena.Size100K, os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())

	chIngestionEnd := rawLogger.StartIngestion(ctx)

	payload := "xxx"

	bytesWritten, errWrite := rawLogger.Write([]byte(payload))
	require.NoError(t, errWrite)
	require.EqualValues(t,
		len(payload),
		bytesWritten,
	)

	cancel()

	// Wait for consumer shutdown flush.
	<-chIngestionEnd
}
