package main

import (
	"fmt"
	"sync"
)

func producer(id int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			out <- i * id
		}
		close(out)
	}()

	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	result := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go forward(ch, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func forward(ch <-chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		result <- v
	}
}

func main() {
	out := merge(
		producer(1),
		producer(10),
		producer(100),
	)

	for v := range out {
		fmt.Println(v)
	}
}
