package main

import "sync"

type BlockingQueue[T any] struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	buf      []T
	capacity int
}

func NewBlockingQueue[T any](cap int) *BlockingQueue[T] {
	q := &BlockingQueue[T]{
		capacity: cap,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *BlockingQueue[T]) Put(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// TODO: wait while full
	for len(q.buf) == q.capacity {
		q.notFull.Wait()
	}
	// TODO: add item
	q.buf = append(q.buf, v)
	// TODO: signal notEmpty
	q.notEmpty.Signal()
}

func (q *BlockingQueue[T]) Get() T {
	q.mu.Lock()
	defer q.mu.Unlock()

	// TODO: wait while empty
	for len(q.buf) == 0 {
		q.notEmpty.Wait()
	}
	// TODO: remove item
	item := q.buf[0]
	q.buf = q.buf[1:]
	// TODO: signal notFull
	q.notFull.Signal()

	return item
}
