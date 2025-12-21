package exercise5

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker -> does an assigned work from a given source queue/channel
// pool -> controls workers
// autoscaler -> decides when to scale up/down

type state uint

const (
	idle state = iota
	busy
	down
	kill
)

type event struct {
	workerId int
	state    state
	message  string
	at       time.Time
	ack      chan struct{}
}

func newEvent(id int, s state, msg string) event {
	return event{
		workerId: id,
		state:    s,
		message:  msg,
	}
}

type worker struct {
	id      int
	jobs    <-chan int
	results chan<- int
	events  chan event
	wg      *sync.WaitGroup
}

func newWorker(id int, jobs <-chan int, results chan<- int, events chan event, wg *sync.WaitGroup) *worker {
	return &worker{
		id:      id,
		jobs:    jobs,
		results: results,
		events:  events,
		wg:      wg,
	}
}

func (w *worker) run(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case e := <-w.events:
			if e.state == kill {
				e.ack <- struct{}{}
				return
			}
		case <-ctx.Done():
			// blocks until sent before exit
			w.events <- newEvent(w.id, down, fmt.Sprintf("worker %d: exiting, context cancelled", w.id))
			return
		case data, ok := <-w.jobs:
			if !ok {
				// blocks until sent before exit
				w.events <- newEvent(w.id, down, fmt.Sprintf("worker %d: exiting, jobs finished", w.id))
				return
			}

			select {
			case <-ctx.Done():
				// blocks until sent before exit
				w.events <- newEvent(w.id, down, fmt.Sprintf("worker %d: exiting, context cancelled", w.id))
				return
			case w.results <- data:
			}
		}
	}
}

type pool struct {
	minWorkers int
	maxWorkers int
	jobs       <-chan int
	workers    map[int]*worker
}

func newWorkerPool(minWorkers, maxWorkers int, jobs <-chan int) *pool {
	return &pool{
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		jobs:       jobs,
		workers:    make(map[int]*worker),
	}
}

func (wp *pool) start(ctx context.Context) <-chan int {
	results := make(chan int)

	for range wp.minWorkers {
		w := newWorker(rand.Intn(99),jobs, results,)
	}

	return results
}
