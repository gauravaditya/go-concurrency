# Mutex vs RWMutex (measured, not guessed)

You’ll write two implementations + benchmarks and see the trade-offs.

## Goal

### Empirically answer

When does RWMutex help?
When is it worse than Mutex?

## What to build

A shared map with two variants

Variant 1: sync.Mutex

Variant 2: sync.RWMutex

Both expose:

```go
Get(k int) (int, bool)
Set(k int, v int)
```
