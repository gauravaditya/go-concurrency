package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type eventName uint

const (
	job_started eventName = iota
	job_finished
)

type event struct {
	name eventName
	at   time.Time
}

type workerCfg struct {
	id      int
	jobs    <-chan int
	results chan<- int
	updates chan<- event
	token   chan time.Duration
	wg      *sync.WaitGroup
}

func worker(
	ctx context.Context,
	cfg *workerCfg,
) {
	defer cfg.wg.Done()
	d := <-cfg.token // blocks till token received
	ticker := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				return
			case cfg.token <- d:
				fmt.Printf("worker %d: returned token\n", cfg.id)
				return
			}
		case data, ok := <-cfg.jobs:
			if !ok {
				select {
				case <-ctx.Done():
					return
				case cfg.token <- d:
					fmt.Printf("worker %d: returned token\n", cfg.id)
					return
				}
			}
			fmt.Printf("worker %d: started job\n", cfg.id)
			cfg.updates <- event{name: job_started, at: time.Now()}
			select {
			case <-ctx.Done():
				select {
				case <-ctx.Done():
					return
				case cfg.token <- d:
					fmt.Printf("worker %d: returned token\n", cfg.id)
					return
				}
			case cfg.results <- data:
				fmt.Printf("worker %d: finished job\n", cfg.id)
				cfg.updates <- event{name: job_finished, at: time.Now()}
			}
		}
	}
}

func autoScalerPool(
	ctx context.Context,
	minWorkers int,
	maxWorkers int,
	jobs <-chan int,
) <-chan int {
	results := make(chan int)
	updates := make(chan event)
	tokens := make(chan time.Duration, maxWorkers)
	workers := sync.Map{}
	var counter atomic.Int32
	var wg sync.WaitGroup

	for range minWorkers {
		cfg := &workerCfg{
			int(counter.Add(1)),
			jobs,
			results,
			updates,
			tokens,
			&wg,
		}
		workers.Store(cfg.id, cfg)
		wg.Add(1)
		go worker(ctx, cfg)
	}

	for range minWorkers {
		tokens <- time.Second * 3
	}

	wg.Add(1)
	go func() {
		startedAt := time.Now()
		started, finished := 0, 0
		for {
			select {
			case <-ctx.Done():
				return
			case u, ok := <-updates:
				if !ok {
					return
				}
				elapsed := time.Now().UnixMilli() - startedAt.UnixMilli()
				switch u.name {
				case job_started:
					started++
				case job_finished:
					finished++
				}

				arrivalRate := math.Mod(float64(started), float64(elapsed))
				completionRate := math.Mod(float64(finished), float64(elapsed))
				if arrivalRate > completionRate {
					tokens <- time.Second * 3

					cfg := &workerCfg{
						int(counter.Add(1)),
						jobs,
						results,
						updates,
						tokens,
						&wg,
					}
					workers.Store(cfg.id, cfg)
					wg.Add(1)
					go worker(ctx, cfg)
				} else {
					<-tokens // reduce token
				}
			case d, ok := <-tokens:
				if !ok {
					return
				}
				tokens <- d

			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(updates)
		for range <-tokens { // drain tokens
			fmt.Println("draining tokens...")
		}
		close(tokens)
	}()

	return results
}
