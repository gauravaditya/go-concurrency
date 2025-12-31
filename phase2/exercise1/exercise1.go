package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Pool[T any, R any] struct {
	jobs    chan job[T]
	results chan R

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	maxJobTimeout time.Duration
	closed        atomic.Bool
}

type job[T any] struct {
	ctx context.Context
	val T
}

func NewPool[T any, R any](
	workers int,
	queueSize int,
	maxJobTimeout time.Duration,
	work func(context.Context, T) (R, error),
) *Pool[T, R] {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool[T, R]{
		jobs:          make(chan job[T], queueSize),
		results:       make(chan R),
		ctx:           ctx,
		cancel:        cancel,
		maxJobTimeout: maxJobTimeout,
	}

	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(work)
	}

	go func() {
		p.wg.Wait()
		close(p.results)
	}()

	return p
}

func (p *Pool[T, R]) worker(
	work func(context.Context, T) (R, error),
) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return

		case j, ok := <-p.jobs:
			if !ok {
				return
			}

			jobCtx := j.ctx
			var cancel context.CancelFunc

			if p.maxJobTimeout > 0 {
				jobCtx, cancel = context.WithTimeout(j.ctx, p.maxJobTimeout)
			}

			r, err := work(jobCtx, j.val)

			if cancel != nil {
				cancel() // ✅ per-job cleanup (no defer)
			}

			if err == nil {
				select {
				case p.results <- r:
				case <-p.ctx.Done():
				case <-jobCtx.Done():
				}
			}
		}
	}
}

func (p *Pool[T, R]) Submit(ctx context.Context, v T) error {
	if p.closed.Load() {
		return errors.New("pool is shut down")
	}

	j := job[T]{ctx: ctx, val: v}

	select {
	case p.jobs <- j:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return errors.New("pool is shut down")
	}
}

func (p *Pool[T, R]) Results() <-chan R {
	return p.results
}

func (p *Pool[T, R]) Shutdown() {
	if p.closed.CompareAndSwap(false, true) {
		p.cancel()
		close(p.jobs)
	}
}
