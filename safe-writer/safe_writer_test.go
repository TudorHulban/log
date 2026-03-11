package safewriter

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestSafeWriter(t *testing.T) {
	w := NewSafeWriter(os.Stdout)
	defer w.Stop()

	numberWorkers := 5

	chPayload := make(chan int)
	defer close(chPayload)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		_, _ = w.Write(
			[]byte(
				strconv.Itoa(
					<-chPayload,
				),
			),
		)
	}

	for range numberWorkers {
		wg.Add(1)

		go worker()
	}

	for ix := range numberWorkers {
		chPayload <- ix
	}

	wg.Wait()
}
