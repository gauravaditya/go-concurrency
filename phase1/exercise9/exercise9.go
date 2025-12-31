package main

import (
	"context"
	"time"
)

func StartPeriodicWorker(
	ctx context.Context,
	interval time.Duration,
	work func(),
) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			work()
		}
	}
}
