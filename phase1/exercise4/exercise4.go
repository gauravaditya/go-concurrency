package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type worker[T, R any] struct {
	cfg workerCfg[T, R]
}

type workerCfg[T, R any] struct {
	id      int
	in      <-chan T
	out     chan<- R
	sigkill chan struct{}
	lag     chan<- int
	wg      *sync.WaitGroup
}

func newWorkerCfg[T, R any](id int, in <-chan T, out chan<- R, lag chan<- int, wg *sync.WaitGroup) workerCfg[T, R] {
	return workerCfg[T, R]{
		id:      id,
		in:      in,
		out:     out,
		sigkill: make(chan struct{}),
		wg:      wg,
		lag:     lag,
	}
}

func newWorker[T, R any](cfg workerCfg[T, R]) *worker[T, R] {
	return &worker[T, R]{cfg: cfg}
}

func (w *worker[T, R]) run(ctx context.Context, apply func(T) R) {
	defer w.cfg.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.cfg.sigkill:
			fmt.Printf("worker %d: received sigkill!\n", w.cfg.id)
			return
		case <-time.After(time.Millisecond * 200):
			select {
			case <-ctx.Done():
				return
			case w.cfg.lag <- w.cfg.id * -1:
			}
		case data, ok := <-w.cfg.in:
			if !ok {
				fmt.Printf("worker %d: input channel drained...exiting\n", w.cfg.id)
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-w.cfg.sigkill:
				fmt.Printf("worker %d: received sigkill!\n", w.cfg.id)
				return
			case w.cfg.out <- apply(data):
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func autoScalerPool[T, R any](
	ctx context.Context,
	minWorkers int,
	maxWorkers int,
	jobs <-chan T,
) (<-chan R, chan<- int) {
	wg := sync.WaitGroup{}
	lag := make(chan int)
	results := make(chan R)
	workers := make(map[int]workerCfg[T, R])
	workerCount := 0
	workerId := 1
	for range minWorkers {
		cfg := newWorkerCfg(workerId, jobs, results, lag, &wg)
		workers[workerId] = cfg
		workerId++
		workerCount++

		wg.Add(1)
		go newWorker(cfg).run(ctx, work[T, R])
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case scale, ok := <-lag:
				if !ok {
					fmt.Printf("lag worker: lag channel closed..\n")
					return
				}

				if scale > 0 && workerCount < maxWorkers {
					cfg := newWorkerCfg(workerId, jobs, results, lag, &wg)
					workers[workerId] = cfg
					workerId++
					workerCount++

					wg.Add(1)
					go newWorker(cfg).run(ctx, work[T, R])
				}

				if scale < 0 && workerCount > minWorkers {
					id := scale * -1 // derive id of the worker
					workers[id].sigkill <- struct{}{}
					delete(workers, id)
					workerCount--
				}
			}
		}

	}(ctx)

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, lag
}

func work[T, R any](t T) (r R) {
	switch v := any(t).(type) {
	case R:
		r = v
		return
	default:
		fmt.Printf("worker %d: unknown datatype received\n")
		return
	}
}
