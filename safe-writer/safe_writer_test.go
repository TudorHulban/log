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

		w.Write(
			[]byte(
				strconv.Itoa(
					<-chPayload,
				),
			),
		)
	}

	for ix := 0; ix < numberWorkers; ix++ {
		wg.Add(1)

		go worker()
	}

	for ix := 0; ix < numberWorkers; ix++ {
		chPayload <- ix
	}

	wg.Wait()
}
