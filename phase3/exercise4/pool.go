package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type workerPool struct {
	Q       *Q
	limiter *Limiter
	results *Q
}

type poolOption func(*workerPool)

func withQSize(size int) poolOption {
	return func(wp *workerPool) {
		wp.Q = NewQ(size)
	}
}

func withResultQSize(size int) poolOption {
	return func(wp *workerPool) {
		wp.results = NewQ(size)
	}
}

func withRateLimit(limit int) poolOption {
	return func(wp *workerPool) {
		wp.limiter = NewLimiter(limit)
	}
}

func NewWorkerPool(ctx context.Context, workerCount int, options ...poolOption) *workerPool {
	var wg sync.WaitGroup
	p := &workerPool{
		Q:       NewQ(5),
		results: NewQ(5),
		limiter: NewLimiter(2),
	}

	for _, optn := range options {
		optn(p)
	}

	for id := range workerCount {
		wg.Add(1)
		go p.worker(ctx, id, &wg)
	}

	go func() {
		wg.Wait()
		p.results.shutdown()
	}()

	return p
}

func (p *workerPool) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-p.Q.jobs:
			if !ok {
				return
			}

			if err := p.limiter.Acquire(ctx); err != nil {
				fmt.Println(err.Error())
				return
			}
			process(id)
			p.results.push(ctx, v)
		}
	}
}

func process(job int) error {
	time.Sleep(time.Second)

	if job%3 == 0 {
		return fmt.Errorf("temporary failure")
	}

	return nil
}
