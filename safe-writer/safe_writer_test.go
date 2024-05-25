package safewriter

import (
	"os"
	"strconv"
	"testing"
)

func TestSafeWriter(t *testing.T) {
	w := NewSafeWriterInfo(os.Stdout)

	go w.Writer.Listen()

	numberWorkers := 5

	worker := func(work <-chan int) {
		w.Writer.write(
			[]byte(
				strconv.Itoa(
					<-work,
				),
			),
		)
	}

	chPayload := make(chan int)
	defer close(chPayload)

	for range numberWorkers {
		go worker(
			chPayload,
		)
	}

	for ix := range numberWorkers {
		chPayload <- ix
	}
}
