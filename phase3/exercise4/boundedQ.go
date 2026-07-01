package main

import "context"

type Q struct {
	jobs chan int
}

func NewQ(size int) *Q {
	return &Q{
		jobs: make(chan int, size),
	}
}

func (q *Q) push(ctx context.Context, item int) {
	select {
	case <-ctx.Done():
		return
	case q.jobs <- item:
	}
}

func (q *Q) shutdown() {
	close(q.jobs)
}
