package bytearena_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
)

func TestHowToUse(t *testing.T) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	bytesWritten, errWrite := ingestor.Write([]byte(payload))
	require.NoError(t, errWrite)
	require.Equal(t, len(payload), bytesWritten)

	cancel()
	<-chIngestionEnd

	require.Contains(t, sink.String(), payload)
}
