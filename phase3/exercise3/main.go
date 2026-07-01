package main

import (
	"context"
	"time"
)

type Limiter struct {
	tokens chan struct{}
}

func NewLimiter(rateLimit int) *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, rateLimit),
	}

	go func() {
		tickTime := time.Second / time.Duration(rateLimit)
		timer := time.NewTicker(tickTime)
		defer timer.Stop()

		for range timer.C {
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket is full
			}
		}
	}()

	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
		return nil
	}
}
