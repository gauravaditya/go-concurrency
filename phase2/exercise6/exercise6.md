Let's move to **Fan-Out / Fan-In**.

You have actually already implemented pieces of it before, but now we'll treat it as its own pattern.

---

## Exercise 1: Fan-In

You are given multiple input channels:

```go
func producer(id int) <-chan int
```

Example:

```go
ch1 := producer(1)
ch2 := producer(10)
ch3 := producer(100)
```

Each producer emits 5 values and then closes.

Example output:

```text
1
2
3
4
5
```

```text
10
20
30
40
50
```

```text
100
200
300
400
500
```

---

### Task

Implement:

```go
func merge(channels ...<-chan int) <-chan int
```

Requirements:

* Read from all input channels concurrently.
* Emit values onto a single output channel.
* Close the output channel when all input channels are exhausted.
* No goroutine leaks.
* No cancellation support yet.

---

### Important

I know you've implemented something similar before.

I don't want you to copy the old code from memory.

Write it again from scratch.