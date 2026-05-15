package timestamp

import (
	"context"
	"time"
)

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

func StartRFC3339Cache(ctx context.Context) <-chan struct{} {
	chReady := make(chan struct{})

	go func() {
		now := time.Now() // Pre-warm so readers never see a nil pointer.

		buildRFC3339Cache(now.UTC())

		// Signal readiness to the caller
		close(chReady)

		ticker := time.NewTicker(500 * time.Microsecond)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				buildRFC3339Cache(
					t.UTC(),
				)

			case <-ctx.Done():
				return
			}
		}
	}()

	return chReady
}

func StartRFC3339BucharestCache(ctx context.Context) <-chan struct{} {
	chReady := make(chan struct{})

	go func() {
		loc, errLoad := time.LoadLocation("Europe/Bucharest")
		if errLoad != nil {
			loc = time.UTC // Fallback safety
		}

		buildRFC3339Cache(time.Now().In(loc))

		// Signal readiness to the caller
		close(chReady)

		ticker := time.NewTicker(500 * time.Microsecond)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				buildRFC3339Cache(
					t.In(loc),
				)

			case <-ctx.Done():
				return
			}
		}
	}()

	return chReady
}

func StartRFC3339CustomLocationCache(ctx context.Context, location string) <-chan struct{} {
	chReady := make(chan struct{})

	go func() {
		loc, errLoad := time.LoadLocation(location)
		if errLoad != nil {
			loc = time.UTC // Fallback safety
		}

		buildRFC3339Cache(time.Now().In(loc))

		// Signal readiness to the caller
		close(chReady)

		ticker := time.NewTicker(500 * time.Microsecond)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				buildRFC3339Cache(
					t.In(loc),
				)

			case <-ctx.Done():
				return
			}
		}
	}()

	return chReady
}
