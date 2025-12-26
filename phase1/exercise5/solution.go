package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type event struct {
	workerId  int
	workId    int
	name      string
	at        int64
	timeTaken int64
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
		worker(ctx, int(id), jobs, results, events, &workerCount, &desiredWorkerCount)
	}

	go func(wc *atomic.Int32, dwc *atomic.Int32) {
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

				fmt.Println("workerCount is ->", workerCount.Load())
				fmt.Println("arrivalRate, completionRate ->", arrivalRate, completionRate)
				switch true {
				case desiredWorkerCount.Load() <= 0:
					return
				case workerCount.Load() < desiredWorkerCount.Load():
					workerCount.Add(1)
					id := idCounter.Add(1)
					worker(ctx, int(id), jobs, results, events, wc, dwc)
					fmt.Printf("added worker %d\n", id)
				default:
					if completionRate > arrivalRate && desiredWorkerCount.Load() > minWorkers {
						desiredWorkerCount.Add(-1)
					}
					if arrivalRate > completionRate && desiredWorkerCount.Load() < maxWorkers {
						desiredWorkerCount.Add(1)
					}
				}

			}
		}
	}(&workerCount, &desiredWorkerCount)

	go func() {
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
					receivedElapsedTime.Add(int32(e.timeTaken))
				case "completed":
					jobsCompleted.Add(1)
					completionElapsedTime.Add(int32(e.timeTaken))
				}
			}
		}
	}()

	return results
}

func worker(
	ctx context.Context,
	id int,
	jobs <-chan int,
	results chan<- int,
	events chan<- event,
	currentWorkerCount *atomic.Int32,
	desiredWorkerCount *atomic.Int32,
) {
	go func() {
		defer func() {
			if currentWorkerCount.Add(-1) <= 0 {
				close(results)
				close(events)
			}
		}()

		for {
			if currentWorkerCount.Load() > desiredWorkerCount.Load() {
				fmt.Printf("worker %d: currentWorkerCount > desiredWorkerCount\n", id)
				return
			}

			select {
			case <-ctx.Done():
				return
			case data, ok := <-jobs:
				if !ok {
					fmt.Printf("worker %d: returning..\n", id)
					desiredWorkerCount.Swap(0)
					return
				}

				select {
				case <-ctx.Done():
					return
				default:
					receivedAt := time.Now().UnixNano()
					select {
					case <-ctx.Done():
						return
					case events <- newEvent(id, data, "received", receivedAt):
					}

					// time.Sleep(500 * time.Millisecond) //work delay
					results <- data

					completedAt := time.Now().UnixNano()
					select {
					case <-ctx.Done():
						return
					case events <- newEvent(id, data, "completed", completedAt):
					}
				}
			}
		}
	}()
}

func newEvent(id, data int, eventName string, receivedAt int64) event {
	now := time.Now().UnixNano()
	timeTaken := int64(0)
	switch eventName {
	case "received":
		timeTaken = now + 1 - receivedAt
	case "completed":
		timeTaken = receivedAt - now
	}

	return event{
		workerId:  id,
		workId:    data,
		name:      eventName,
		at:        receivedAt,
		timeTaken: timeTaken,
	}
}
