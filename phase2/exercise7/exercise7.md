We have:

```text
producer -> producer channel -> forwarder -> merged channel -> consumer
```

Let's assume the consumer reads exactly one value:

```go
fmt.Println(<-out)
return
```

and then imagine the process does **not** exit.

---

### Step 1

A producer sends a value:

```go
out <- value
```

to its own channel.

A forwarder receives it:

```go
for v := range ch
```

and then tries:

```go
result <- v
```

One of those sends succeeds because the consumer reads one value.

So far, so good.

---

### Step 2

Consumer stops reading.

Now every forwarder eventually reaches:

```go
result <- v
```

and blocks.

This is the first place things get stuck.

---

### Step 3

What are the producers doing?

The forwarders are no longer reading from:

```go
ch
```

because they're stuck trying to write to `result`.

Therefore the producers eventually reach:

```go
out <- value
```

and block too.

---

So the final state is:

```text
producer blocked sending
        ↓
forwarder blocked sending
        ↓
consumer gone
```

Everything backs up through the pipeline.

---

One subtle distinction:

I wouldn't call this a **deadlock**.

A deadlock usually implies:

> Goroutine A waits for B, B waits for C, C waits for A.

A cycle.

Here we have something slightly different:

```text
consumer stopped
    ↓
forwarders blocked
    ↓
producers blocked
```

There's no cycle.

It's more accurately described as:

* blocked goroutines
* goroutine leak
* stalled pipeline
* backpressure propagation

Depending on context.

---

And this is the first major lesson of Fan-In:

> Fan-In is only safe if somebody continues draining the merged output.

If the consumer disappears, the entire upstream graph can become stuck.

---

Now for the next exercise, I want you to solve exactly this problem.

Modify the design conceptually (don't code yet) so that:

```go
ctx, cancel := context.WithCancel(...)
```

can be used to stop:

* producers
* forwarders
* merge goroutine

even if the consumer abandons the output channel.

How would you thread cancellation through the system?
