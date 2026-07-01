Perfect.

Let's move to **Rate Limiting — Token Bucket**.

I'm going to keep the theory to the minimum and let the implementation drive the learning.

---

## Goal

Build a rate limiter that allows:

```text
5 requests / second
```

Example:

```go
limiter := NewLimiter(5)
```

Then:

```go
limiter.Acquire(ctx)
```

must block until a token is available.

---

## Design Hint

Think of a token bucket as:

```text
bucket
```

containing tokens.

A background goroutine periodically adds tokens.

Consumers:

```go
Acquire(ctx)
```

remove a token.

No token available?

```text
wait
```

until one arrives.

---

## Constraints

For the first version:

### Constructor

```go
func NewLimiter(rate int) *Limiter
```

where:

```text
rate = tokens per second
```

---

### Acquire

```go
func (l *Limiter) Acquire(ctx context.Context) error
```

Requirements:

* If a token exists, consume it immediately.
* Otherwise wait.
* Respect context cancellation.

---

### Don't Worry About Yet

Ignore:

* burst capacity
* dynamic rate changes
* shutdown
* fairness
* distributed rate limiting

We'll get to those later.

---

## One Design Question

Before you code:

Would you represent the bucket as:

### Option A

```go
tokens int
mutex
cond
```

### Option B

```go 
tokens chan struct{}
```

Given everything you've built so far, which one are you instinctively drawn toward, and why?
