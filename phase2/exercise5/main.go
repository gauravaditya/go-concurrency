package main

import (
	"context"
	"fmt"
	"go-concurrency/phase2/types"
	"time"
)

type pool struct {
	jobs chan types.Job
}

func NewPool(ctx context.Context, workers int, queueSize int) *pool {
	p := &pool{
		jobs: make(chan types.Job, queueSize),
	}

	for range workers {
		go p.worker(ctx)
	}

	return p
}

func (p *pool) submit(ctx context.Context, job types.Job) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled..")
	case p.jobs <- job:
		return nil
	}
}

func (p *pool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				fmt.Printf("worker processed %v\n", job)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewPool(ctx, 2, 5)

	time.AfterFunc(5*time.Second, cancel)
	var err error
	for i := 1; i <= 20; i++ {
		fmt.Println("submitting", i)
		err = pool.submit(ctx, types.Job{ID: i})
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("submitted", i)
	}

	// time.Sleep(5 * time.Second)
}
