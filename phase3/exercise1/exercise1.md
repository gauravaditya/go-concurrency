The next natural step is **Pipelines**.

This is where multiple stages are chained together:

```text id="0qpnww"
source
  ↓
parse
  ↓
transform
  ↓
persist
```

and each stage can have multiple workers.

---

### Pipeline Exercise 1

Let's keep it simple.

Build:

```text id="b50w04"
numbers
  ↓
square
  ↓
double
  ↓
print
```

Requirements:

```go id="5h2yrg"
generate(ctx, n)
```

Produces:

```text id="1vfj5x"
1 2 3 ... n
```

---

```go id="m4pn4m"
square(ctx, in)
```

Transforms:

```text id="j79x9y"
1 -> 1
2 -> 4
3 -> 9
```

---

```go id="48g76p"
double(ctx, in)
```

Transforms:

```text id="j20xhd"
1 -> 2
4 -> 8
9 -> 18
```

---

Each stage should:

* run in its own goroutine
* close its output when finished
* respect context cancellation

Don't over-engineer it.

The goal is simply to build your first proper pipeline and experience how channel ownership and closure work across stages.
