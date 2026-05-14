package timestamp

import (
	"context"
	"time"
)

// use as:
// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	go timestamp.Start(ctx)

// 	// ...
// }

// in tests
// go Start(t.Context())

// Start begins background cache updates and blocks until ctx is cancelled.
// Call once from main, typically as: go timestamp.Start(ctx)
func Start(ctx context.Context) {
	// Pre-warm so readers never see a nil pointer.
	now := time.Now()

	buildRFC3339Cache(now.UTC())
	buildStandardCache(now)
	buildYYYYMonthCache(now)

	ticker := time.NewTicker(500 * time.Microsecond)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			utc := t.UTC()
			buildRFC3339Cache(utc)
			buildStandardCache(t)
			buildYYYYMonthCache(t)

		case <-ctx.Done():
			return
		}
	}
}
