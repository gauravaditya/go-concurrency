package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type event struct {
	workerId  int
	workId    int
	name      string
	at        int64
	timeTaken int32
}

func (e event) String() string {
	return fmt.Sprintf("{workerId: %d, workId: %d,name: %s, at: %d, timeTaken: %d}", e.workerId, e.workId, e.name, e.at, e.timeTaken)
}

func pool(
	ctx context.Context,
	minWorkers int32,
	maxWorkers int32,
	jobs <-chan int,
) <-chan int {
	results := make(chan int)
	events := make(chan event)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	idCounter := atomic.Int32{}
	workerCount := atomic.Int32{}
	desiredWorkerCount := atomic.Int32{}
	desiredWorkerCount.Store(minWorkers)

	jobsReceived := atomic.Int32{}
	jobsCompleted := atomic.Int32{}
	receivedElapsedTime := atomic.Int32{}
	completionElapsedTime := atomic.Int32{}

	for range 1 {
		workerCount.Add(1)
		id := idCounter.Add(1)
		worker(ctx, int(id), jobs, results, events, &wg, &workerCount, &desiredWorkerCount)
	}

	go func(wg *sync.WaitGroup, wc *atomic.Int32, dwc *atomic.Int32) {
		ticker := time.NewTicker(time.Millisecond * 500)
		arrivalRate := 0.0
		completionRate := 0.0

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if jobsReceived.Load() > 0 {
					arrivalRate = float64(jobsReceived.Load()) * 1000 / float64(receivedElapsedTime.Load())
				}

				if jobsCompleted.Load() > 0 {
					completionRate = float64(jobsCompleted.Load()) * 1000 / float64(completionElapsedTime.Load())
				}

				// fmt.Printf("arrivalRate: %.8f, %.8f :completionRate, %v\n", arrivalRate, completionRate, arrivalRate > completionRate)
			case <-doneCh:
				fmt.Println("Done.")
				return
			default:
				switch true {
				case workerCount.Load() < desiredWorkerCount.Load():
					fmt.Printf("need more workers: %d, %d\n", dwc.Load(), maxWorkers)
					workerCount.Add(1)
					id := idCounter.Add(1)
					worker(ctx, int(id), jobs, results, events, wg, wc, dwc)
					fmt.Printf("added worker %d\n", id)
				case workerCount.Load() > desiredWorkerCount.Load():
					fmt.Printf("need less workers: %d, %d", wc.Load(), dwc.Load())
					desiredWorkerCount.Add(-1)
				default:
					if completionRate > arrivalRate && desiredWorkerCount.Load() > minWorkers {
						desiredWorkerCount.Add(-1)
					}
					if arrivalRate < completionRate && desiredWorkerCount.Load() < maxWorkers {
						desiredWorkerCount.Add(1)
					}
				}

			}
		}
	}(&wg, &workerCount, &desiredWorkerCount)

	go func(wg *sync.WaitGroup) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context cancelled...")
				return
			case e, ok := <-events:
				if !ok {
					return
				}
				switch e.name {
				case "received":
					jobsReceived.Add(1)
					receivedElapsedTime.Add(e.timeTaken)
				case "completed":
					jobsCompleted.Add(1)
					completionElapsedTime.Add(e.timeTaken)
				}
			}
		}
	}(&wg)

	go func() {
		wg.Wait()
		doneCh <- struct{}{}
		close(doneCh)
		close(results)
		close(events)
	}()

	return results
}

func worker(
	ctx context.Context,
	id int,
	jobs <-chan int,
	results chan<- int,
	events chan<- event,
	wg *sync.WaitGroup,
	currentWorkerCount *atomic.Int32,
	desiredWorkerCount *atomic.Int32,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastReceivedAt := time.Now().UnixNano()

		for {
			if currentWorkerCount.Load() > desiredWorkerCount.Load() {
				fmt.Printf("worker %d: currentWorkerCount > desiredWorkerCount\n", id)
				currentWorkerCount.Add(-1)
				return
			}

			select {
			case <-ctx.Done():
				return
			case data, ok := <-jobs:
				if !ok {
					fmt.Printf("worker %d: returning..\n", id)
					return
				}

				select {
				case <-ctx.Done():
					return
				default:
					receivedAt := time.Now().UnixNano()
					events <- event{
						workerId:  id,
						workId:    data,
						name:      "received",
						at:        receivedAt,
						timeTaken: int32(receivedAt + 1 - lastReceivedAt),
					}
					lastReceivedAt = receivedAt

					// time.Sleep(500 * time.Millisecond) //work delay
					results <- data

					completedAt := time.Now().UnixNano()
					events <- event{
						workerId:  id,
						workId:    data,
						name:      "completed",
						at:        completedAt,
						timeTaken: int32(completedAt - receivedAt),
					}
				}
			}
		}
	}()
}
