# Go Mastery — Real-World Challenge Handbook

> A hands-on curriculum for developers who want to go from Go beginner to production-ready engineer.
> Built for people who learn by doing — not by reading slides.

---

## Who is this for?

This handbook is for developers who:
- Already know at least one programming language
- Want to learn Go properly — not just syntax, but how Go is actually used in production
- Want to understand Go backend and systems engineering deeply

## What you'll be able to do after this

By the end of these challenges you will be able to:
- Design systems using Go's interface model (the way the standard library does it)
- Write concurrent programs without race conditions
- Profile and optimize Go code for production workloads
- Build HTTP services that handle real traffic patterns
- Read and review Go code like a advanced engineer
- Understand diverse Go topics such as concurrency model, memory model, and type system

## Skill progression

| After Phase | What you can do |
|---|---|
| Phase 1 ✅ | Write basic Go — structs, slices, methods, file I/O |
| Phase 2 | Design with interfaces, read standard library code |
| Phase 3 | Build concurrent systems without race conditions |
| Phase 4 | Profile and optimize Go code for performance |
| Phase 5 | Ship a production-grade HTTP service end-to-end |
| Phase 6 | Read, review, and debug Go code under pressure |

---

# PHASE 1 — Foundations ✅
Check

---

# PHASE 2 — Interfaces & Type System

> **Why this phase matters**
> Interfaces are the backbone of Go's entire standard library. `http.Handler`, `io.Reader`, `io.Writer`, `error` — all interfaces. If you don't understand interfaces deeply, you'll struggle to read Go code written by others, and you'll write brittle code yourself. This phase teaches you to think the Go way: *program to behavior, not to concrete types.*

---

## Challenge 2.1 — Build a Multi-Format Logger
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 8–10 hours**

### 1. Context
Every production system logs events. But *where* those logs go changes depending on the environment: console during local development, structured files in staging, JSON for log aggregation tools like Datadog, Grafana Loki, or AWS CloudWatch in production.

A well-designed logging system should let you swap the destination without changing the code that *uses* the logger. This is exactly how Go's `io.Writer`, `log/slog`, and popular libraries like `uber-go/zap` work internally.

### 2. Goal
Build a logging system that can write to multiple destinations (console, file, JSON) using a common interface. The logger must be swappable — the rest of the code should not care where logs go.

### 3. Scope
- Define a `Logger` interface with at least one method: `Log(level, message string)`
- Implement 3 concrete loggers, all satisfying the `Logger` interface:
  - `ConsoleLogger` — prints to terminal with timestamp
  - `FileLogger` — writes to a `.log` file
  - `JSONLogger` — writes structured JSON lines to a file
- Write a function `RunApp(l Logger)` that takes any logger and logs 3 events (startup, a warning, a shutdown)
- In `main()`, call `RunApp` three times — once with each logger type
- No `if/else` based on logger type anywhere in `RunApp` — it must work purely through the interface

### 4. Expected Output
Console:
```
[2026-03-22 10:00:01] INFO  app started
[2026-03-22 10:00:01] WARN  high memory usage
[2026-03-22 10:00:01] INFO  app shutdown
```
File (`app.log`):
```
2026-03-22 10:00:01 | INFO  | app started
2026-03-22 10:00:01 | WARN  | high memory usage
2026-03-22 10:00:01 | INFO  | app shutdown
```
JSON (`app.json`):
```json
{"time":"2026-03-22T10:00:01","level":"INFO","message":"app started"}
{"time":"2026-03-22T10:00:01","level":"WARN","message":"high memory usage"}
{"time":"2026-03-22T10:00:01","level":"INFO","message":"app shutdown"}
```

### 5. Hints & Knowledge
- In Go, interfaces are implemented **implicitly** — no `implements` keyword. If your struct has the right methods, it satisfies the interface automatically.
- `io.Writer` is Go's most important interface: `Write(p []byte) (n int, err error)`. `os.Stdout` and `os.File` both implement it — that's why you can write to both the same way.
- `time.Now().Format("2006-01-02 15:04:05")` — Go uses a reference time for formatting (Jan 2, 2006 = Go's birthday).
- `encoding/json` — use `json.Marshal(struct)` to convert a struct to JSON bytes.
- `os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)` — open a file for appending.

### 6. Sources
- Go interfaces explained: https://go.dev/tour/methods/9
- `io.Writer`: https://pkg.go.dev/io#Writer
- `encoding/json`: https://pkg.go.dev/encoding/json
- `time.Format`: https://pkg.go.dev/time#Time.Format
- `os.OpenFile`: https://pkg.go.dev/os#OpenFile

### 7. Knowledge Gained
- ✅ How Go interfaces work (implicit implementation)
- ✅ Writing to `io.Writer` — the foundation of all Go I/O
- ✅ `encoding/json` for structured data
- ✅ Dependency injection via interfaces (pass behavior, not implementation)
- ✅ The design pattern used by `net/http`, `os`, `bufio`, and most Go packages

---

## Challenge 2.2 — Fix the Shape Calculator
### `🟢 Beginner`
**🕐 Expected duration: 3–4 hours**

### 1. Context
A junior developer tried to build a geometry calculator that computes the area of different shapes using Go interfaces. The code compiles in some places and panics in others. Your job is to fix it and make it robust.

### 2. Goal
Fix all bugs in the provided broken code, understand why each bug exists, and add one defensive improvement using a **type switch**.

### 3. Scope
Here is the broken code:

```go
package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Describe() string
}

type Circle struct {
    Radius float64
}

type Rectangle struct {
    Width, Height float64
}

type Triangle struct {
    Base, Height float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Describe() string {
    return fmt.Sprintf("Rectangle %.1f x %.1f", r.Width, r.Height)
}

func (c Circle) Describe() string {
    return fmt.Sprintf("Circle r=%.1f", c.Radius)
}

func printArea(s Shape) {
    fmt.Printf("%s → area: %.2f\n", s.Describe, s.Area)
}

func totalArea(shapes []Shape) float64 {
    total := 0
    for _, s := range shapes {
        total += s.Area()
    }
    return total
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 3, Height: 4},
        Triangle{Base: 6, Height: 8},
    }
    for _, s := range shapes {
        printArea(s)
    }
    fmt.Printf("Total area: %.2f\n", totalArea(shapes))
}
```

Find ALL bugs (there are 5), fix them, then add:
- A `type switch` inside `printArea` that prints `"(is a circle)"` if the shape is a `Circle`

### 4. Expected Output
```
Circle r=5.0 (is a circle) → area: 78.54
Rectangle 3.0 x 4.0 → area: 12.00
Triangle b=6.0 h=8.0 → area: 24.00
Total area: 114.54
```

### 5. Hints & Knowledge
- Missing method on a type = does NOT satisfy the interface → compile error
- `s.Describe` vs `s.Describe()` — calling a method needs `()`
- `total := 0` makes `total` an `int` — can't add `float64` to it
- Type switch syntax: `switch v := s.(type) { case Circle: ... }`
- A `Triangle` must implement ALL methods of `Shape` to be used as one

### 6. Sources
- Type assertions: https://go.dev/tour/methods/15
- Type switches: https://go.dev/tour/methods/16

### 7. Knowledge Gained
- ✅ Interface satisfaction rules — ALL methods must be implemented
- ✅ Type assertion and type switch
- ✅ Common interface bugs and how to spot them
- ✅ Zero values and type mismatches

---

# PHASE 3 — Goroutines & Channels

> **Why this phase matters**
> Concurrency is Go's killer feature. Goroutines are why Go is chosen popular. A goroutine costs ~2KB of memory vs ~8MB for a thread — you can run hundreds of thousands of them. Channels replace shared memory with message passing, eliminating entire classes of bugs. If you can write correct concurrent Go, you are got what is nice of Go.

---

## Challenge 3.1 — Build a Concurrent Log Scanner
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 15–20 hours**

### 1. Context
Imagine you're on a platform team. An incident just happened and you need to scan 50 server log files for `ERROR` lines — fast. Doing it sequentially takes 50× longer than necessary. The fix: scan all files at the same time, each in its own goroutine, and collect results through a channel.

This is a fan-out/fan-in pattern — one of the two most important concurrency patterns in Go. You'll see it in log processors, data pipelines, and CI/CD systems.

For myself and 42er: This is the Go equivalent of running multiple shell pipelines at the same time in your minishell — except Go manages the "processes" (goroutines) for you, and channels replace your pipes.

### 2. Goal
Build a concurrent log file scanner that processes multiple files simultaneously, collects all ERROR lines via a channel, and prints a final report.

### 3. Scope
- Generate 5 fake `.log` files at startup (100 lines each, random `INFO`/`WARN`/`ERROR`)
- Launch one **goroutine per file** to scan it — fan-out
- Each goroutine sends its result into a shared results channel — fan-in
- Use `sync.WaitGroup` to know when all workers are done
- Close the channel after all goroutines finish
- Print a sorted final report: errors per file + total errors found
- Must pass `go run -race .` with zero race conditions

### 4. Expected Output
```
Scanning 5 files concurrently...

[worker] log_1.log → 7 errors found
[worker] log_3.log → 3 errors found
[worker] log_2.log → 9 errors found
[worker] log_5.log → 4 errors found
[worker] log_4.log → 6 errors found

=== ERROR REPORT ===
log_1.log : 7 errors
log_2.log : 9 errors
log_3.log : 3 errors
log_4.log : 6 errors
log_5.log : 4 errors
Total     : 29 errors
```


### 5. Why This Matters in Production
Fan-out/fan-in is the core of Go's concurrency model. It's used in:
- **CI/CD systems** — run tests for N packages in parallel, collect results
- **Data pipelines** — process N files/records simultaneously
- **Web scrapers** — fetch N URLs concurrently
- **Kubernetes controllers** — reconcile N resources at the same time
Understanding this pattern is essential for writing real Go backend services.

### 6. Common Mistakes to Avoid
- Closing the channel from inside a goroutine instead of the coordinator — causes panic if another goroutine writes after close
- Calling `wg.Wait()` in the main goroutine before launching the close-channel goroutine — causes deadlock
- Using an unbuffered channel when workers produce faster than collector consumes — causes goroutine leak
- Not passing `&wg` (pointer) to the worker — `wg.Done()` on a copy does nothing

### 7. What a Senior Would Do Differently
- Add a `context.Context` with cancel — so if one worker hits a fatal error, all others stop
- Use `errgroup` from `golang.org/x/sync/errgroup` — handles WaitGroup + error collection in one
- Cap the number of goroutines with a semaphore channel instead of one-per-file, for when N is large
- Use `bufio.Scanner` with a custom buffer size for large log files

### 8. Hints & Knowledge
- `go func()` launches a goroutine — 2KB stack, multiplexed onto OS threads by the Go scheduler
- `wg.Add(1)` before `go`, `wg.Done()` inside goroutine, `wg.Wait()` to block until all done
- `make(chan Result, N)` — buffered channel, workers don't block waiting for collector
- `for result := range ch` — reads until `close(ch)` is called
- `go run -race .` — the race detector, essential tool for concurrent Go

### 9. Sources
- Goroutines: https://go.dev/tour/concurrency/1
- Channels: https://go.dev/tour/concurrency/2
- `sync.WaitGroup`: https://pkg.go.dev/sync#WaitGroup
- Race detector: https://go.dev/doc/articles/race_detector
- `bufio.Scanner`: https://pkg.go.dev/bufio#Scanner

### 10. Knowledge Gained
- ✅ Goroutines — launching and managing concurrent work
- ✅ Channels — type-safe message passing between goroutines
- ✅ `sync.WaitGroup` — coordinating goroutine completion
- ✅ Fan-out / fan-in — the most common Go concurrency pattern
- ✅ Race condition detection

---

## Challenge 3.2 — Build a Worker Pool URL Checker
### `🟠 Intermediate`
**🕐 Expected duration: 15–20 hours**

### 1. Context
In production, you never spawn unlimited goroutines. If 10,000 requests come in and you launch 10,000 goroutines to handle them simultaneously, your server runs out of memory. The solution: a **worker pool** — a fixed number of goroutines that pick jobs from a queue, process them, and stay alive for the next job.

Worker pools are how Go HTTP servers, job queues (like Faktory, Asynq), and data processors work under the hood. This is the second of the two most important Go concurrency patterns.

### 2. Goal
Build a URL health checker using a fixed worker pool of 5 goroutines that processes 20 URLs concurrently with timeout control.

### 3. Scope
- Define exactly **5 worker goroutines** — no more, no less
- Feed 20 URLs into a jobs channel (mix of valid/invalid/timeout URLs)
- Each worker performs an HTTP GET with a **3-second timeout** using `context`
- Results sent to a results channel, printed as they arrive
- Graceful handling: timeouts, unreachable hosts, invalid URLs
- Print final summary: total success vs failed
- Workers must stop cleanly when there are no more jobs

### 4. Expected Output
```
[worker 2] ✅ https://google.com       → 200 OK    (121ms)
[worker 1] ❌ https://notarealsite.xyz → timeout   (3001ms)
...

=== SUMMARY ===
✅ Success : 14
❌ Failed  : 6
```

### 5. Why This Matters in Production
Worker pools are used everywhere:
- **Payment processors** — exactly N workers process transactions to avoid overloading downstream APIs
- **Web crawlers** — N goroutines fetch pages, respecting rate limits
- **Background job systems** — N workers drain a Redis/SQS queue
- **Health checkers** — this exact challenge, in production (Prometheus blackbox exporter does this)

Knowing when to use a worker pool vs fan-out is a key advanced Go topic.

### 6. Common Mistakes to Avoid
- Not closing the jobs channel — workers loop forever waiting for jobs that never come (goroutine leak)
- Closing the results channel from a worker goroutine — if multiple workers do this, panic
- Not using `context` for timeouts — HTTP requests can hang forever without it
- Using `time.Sleep` for timeouts instead of `context.WithTimeout` — never do this
- Sharing the HTTP client between goroutines without knowing it's safe (it is — `http.Client` is safe for concurrent use)

### 7. What a Senior Would Do Differently
- Parameterize pool size and pass it via config or environment variable
- Use `errgroup` with a semaphore instead of explicit channels
- Add retry logic with exponential backoff for transient failures
- Export metrics (success rate, p95 latency) to Prometheus
- Use `http.Client` with a custom `Transport` for connection pooling tuning

### 8. Hints & Knowledge
- `context.WithTimeout(context.Background(), 3*time.Second)` — cancels after 3s
- `http.NewRequestWithContext(ctx, "GET", url, nil)` — attaches context to request
- `close(jobs)` — workers reading `for job := range jobs` stop automatically
- Channel direction: `jobs <-chan Job` (receive-only), `results chan<- Result` (send-only)
- `time.Since(start)` — measure elapsed time

### 9. Sources
- Worker pools: https://gobyexample.com/worker-pools
- `context` package: https://pkg.go.dev/context
- `net/http` client: https://pkg.go.dev/net/http
- `select` statement: https://go.dev/tour/concurrency/5

### 10. Knowledge Gained
- ✅ Worker pool — fixed concurrency pattern
- ✅ `context` — timeout and cancellation (essential for all network code)
- ✅ `net/http` client — making HTTP requests in Go
- ✅ Channel directionality — enforcing send/receive contracts
- ✅ Graceful goroutine shutdown

---

# PHASE 4 — Memory & Performance

> **Why this phase matters**
> Go is used for high-performance systems because it gives you control over memory — without the danger of C. Companies running Go at scale (Cloudflare processes 50M+ req/s with Go) care deeply about allocations per request, GC pauses, and heap pressure. Knowing how to measure, profile, and reduce allocations is what separates surface-level Go from deep Go.

---

## Challenge 4.1 — The Benchmark Battle
### `🟠 Intermediate → Advanced`
**🕐 Expected duration: 15–20 hours**

### 1. Context
A data pipeline at your company processes millions of log lines per day. Each line is parsed into a key-value map. The current implementation is correct but slow — and it's causing GC pressure because it allocates a new map on every single call. Your job: measure it, understand why it's slow, and fix it.

This is a real scenario. Datadog, Cloudflare, and similar companies do this kind of optimization routinely on their log ingestion pipelines.

### 2. Goal
Benchmark two existing implementations, analyze their memory behavior using Go's built-in tooling, then write a faster third version that wins on both time and allocations.

### 3. Scope
- Write proper benchmark tests for Version A and Version B (provided)
- Run `go test -bench=. -benchmem` and record: `ns/op`, `B/op`, `allocs/op`
- Run `go build -gcflags="-m"` to see escape analysis output — what goes to heap?
- Write Version C using `sync.Pool` to reduce heap allocations
- Version C must have fewer `allocs/op` than both A and B
- Write a short explanation (as comments) of what each optimization does and why

### 4. Expected Output
```
BenchmarkParseA-8    500000    2341 ns/op    512 B/op    8 allocs/op
BenchmarkParseB-8    800000    1823 ns/op    384 B/op    6 allocs/op
BenchmarkParseC-8   2000000     601 ns/op     64 B/op    1 allocs/op
```

### 5. Why This Matters in Production
At 1M requests/day, the difference between 8 allocs/op and 1 alloc/op is 7 million fewer heap allocations. Each allocation the GC doesn't have to track = less GC pause = lower tail latency. This is why high-performance Go services obsess over allocations per request.

### 6. Common Mistakes to Avoid
- Optimizing without measuring first — "premature optimization is the root of all evil"
- Not calling `mapPool.Put(result)` after use — defeats the purpose of `sync.Pool`
- Using `sync.Pool` for objects that have identity or state that shouldn't be shared
- Confusing `b.N` — Go determines the right N automatically, never hardcode it
- Not running benchmarks multiple times — use `-count=5` for stable results

### 7. What a Senior Would Do Differently
- Use `pprof` for CPU and heap profiling on a running service, not just microbenchmarks
- Consider `[]byte` instead of `string` for the entire pipeline — avoids string→[]byte conversion
- Use `go test -benchmem -cpuprofile cpu.out` then `go tool pprof cpu.out` to see flamegraphs
- Know when `sync.Pool` is *not* appropriate — for objects with finalizers or long-lived state

### 8. Hints & Knowledge
- `func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { ... } }` — standard shape
- `go test -bench=. -benchmem` — run all benchmarks with memory stats
- `B/op` = bytes allocated per operation, `allocs/op` = number of heap allocations
- Escape analysis: if a local variable is used after the function returns, it "escapes" to heap
- `sync.Pool`: Get() → use → Put() — reuse without allocation

### 9. Sources
- Go benchmarks: https://pkg.go.dev/testing#hdr-Benchmarks
- `sync.Pool`: https://pkg.go.dev/sync#Pool
- Escape analysis deep dive: https://go.dev/doc/faq#stack_or_heap
- pprof tutorial: https://go.dev/blog/pprof

### 10. Knowledge Gained
- ✅ Writing and interpreting Go benchmark tests
- ✅ `benchmem` — reading allocation output
- ✅ Escape analysis — understanding stack vs heap
- ✅ `sync.Pool` — object reuse pattern
- ✅ How to approach performance work: measure → profile → optimize → re-measure

---

# PHASE 5 — Standard Library & Systems Integration

> **Why this phase matters**
> Building real Go services means combining everything: goroutines for concurrency, interfaces for flexibility, HTTP for APIs, JSON for data, and mutexes for safe state. This phase simulates the architecture of a real microservice — the kind you'd find at any company running Go in production.

---

## Challenge 5.1 — Build a Mini DevOps Dashboard
### `🔴 Intermediate → Advanced`
**🕐 Expected duration: 25–30 hours**

### 1. Context
You're joining a platform team at a mid-size company. They have dozens of services writing log files to a shared directory. The ops team needs a lightweight internal tool that automatically picks up new log files, analyzes them for error rates, and exposes the results via a simple HTTP API — without restarting the service.

This is a simplified version of what tools like Fluentd, Logstash, and Vector do. You're building the Go-native version from scratch.

### 2. Goal
Build a self-contained HTTP service that watches a directory for log files, processes them concurrently, and exposes results through a REST JSON API.

### 3. Scope
The service has 3 components working together:

**Component A — File Watcher** (goroutine)
- Every 5 seconds, scan `./logs/` for new `.log` files
- Send new filenames to a jobs channel
- Track already-seen files — don't reprocess

**Component B — Worker Pool** (3 goroutines)
- Read filenames from jobs channel
- Count: total lines, ERROR, WARN, INFO per file
- Store results in a thread-safe store (`sync.Mutex`)

**Component C — HTTP API**
- `GET /status` → JSON of all processed files + counts
- `GET /errors` → JSON of only files with at least 1 error
- `POST /scan` → trigger an immediate re-scan without waiting for ticker

### 4. Expected Output
```bash
$ curl http://localhost:8080/status
{
  "total_files": 3,
  "files": [
    {"name": "log_1.log", "total": 100, "errors": 7, "warns": 12, "infos": 81}
  ]
}

$ curl http://localhost:8080/errors
{
  "files_with_errors": [
    {"name": "log_1.log", "errors": 7}
  ]
}
```

### 5. Why This Matters in Production
This challenge's architecture is the architecture of most Go microservices:
- A **background goroutine** doing periodic work (cron jobs, health checks, cache refresh)
- A **worker pool** processing a queue (job processors, event consumers)
- A **mutex-protected store** as the source of truth
- An **HTTP API** as the external interface

Understanding how these three components communicate safely is the difference between Go code that works and Go code that works *under load*.

### 6. Common Mistakes to Avoid
- Protecting reads *and* writes with the mutex — a common mistake is only locking writes
- Launching a new goroutine on every HTTP request instead of using the existing worker pool
- Not generating test log files — the service has nothing to do without them
- Using `map[string]FileStats` without a mutex — concurrent map writes cause runtime panic
- Blocking the HTTP handler while waiting for a scan to complete — use `go` and return `202 Accepted`

### 7. What a Senior Would Do Differently
- Replace `sync.Mutex` + `map` with a proper repository interface — easier to swap for Redis/Postgres later
- Use `chi` or `net/http`'s `ServeMux` patterns for cleaner routing
- Add structured logging with `log/slog` — every request logged with duration and status
- Add `/healthz` and `/readyz` endpoints — standard in any Kubernetes-deployed service
- Use `context` propagation from HTTP request through to file processing

### 8. Hints & Knowledge
- `http.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {})` — register a handler
- `json.NewEncoder(w).Encode(data)` — write JSON directly to response writer
- `w.Header().Set("Content-Type", "application/json")` — always set before writing body
- `time.NewTicker(5 * time.Second)` — fires every 5s, like a cron job
- `os.ReadDir("./logs/")` — returns `[]os.DirEntry`
- `sync.Mutex`: `mu.Lock()` before read/write, `defer mu.Unlock()` immediately after

### 9. Sources
- `net/http`: https://pkg.go.dev/net/http
- `encoding/json`: https://pkg.go.dev/encoding/json
- `sync.Mutex`: https://pkg.go.dev/sync#Mutex
- `time.Ticker`: https://pkg.go.dev/time#Ticker
- `os.ReadDir`: https://pkg.go.dev/os#ReadDir
- Go HTTP patterns: https://gobyexample.com/http-servers

### 10. Knowledge Gained
- ✅ `net/http` — building production HTTP servers
- ✅ JSON encoding for REST APIs
- ✅ `sync.Mutex` — protecting shared state under concurrent access
- ✅ `time.Ticker` — background periodic tasks
- ✅ Wiring goroutines + channels + HTTP into a cohesive service architecture
- ✅ The standard Go microservice architecture pattern

---

# PHASE 6 — Stress test

> **Why this phase matters**
> Technical skill and performance under pressure are different skills. Phase 6 trains the second one: reading code under pressure, explaining your decisions out loud, and writing correct Go quickly. After Phase 5 you have the knowledge — Phase 6 makes sure you can demonstrate it.

---

## Challenge 6.1 — Code Review Gauntlet
### `🔴 Advanced`
**🕐 Expected duration: 10 hours**

### 1. Context
Every developer will face hard time while developing. You're shown real-looking code with subtle bugs — goroutine leaks, race conditions, nil panics, interface misuse — and asked to spot and explain them. No running the code.

### 2. Goal
Review 5 broken Go programs. For each: identify all issues, explain why each is a problem, and write the fix.

### 3. Scope
*(Programs provided when you reach this challenge)*
Covers: goroutine leak, race condition, nil interface panic, bad interface design, performance anti-pattern.

### 4. Why This Matters in Production
Every Go team does code review. The ability to read a PR and say *"this goroutine leaks if the context is cancelled"* or *"this map access needs a mutex"* is what separates a beginner from an advanced engineer.

### 5. Common Mistakes to Avoid
- Assuming code is correct because it compiles — Go's concurrency bugs are runtime bugs
- Missing nil interface subtlety: an interface holding a nil pointer is NOT nil itself
- Not checking if channels are ever closed — the most common goroutine leak

### 6. What a Senior Would Do Differently
- Use `go vet` and `staticcheck` as automated first passes before manual review
- Reference https://100go.co — the canonical Go mistakes resource

### 7. Knowledge Gained
- ✅ Critical code reading skills
- ✅ Goroutine leak patterns
- ✅ Race condition identification
- ✅ Nil interface gotchas

---

## Challenge 6.2 — Build Under Pressure
### `🔴 Advanced`
**🕐 Expected duration: 8–10 hours**

### 1. Context
The final challenge. 3 timed problems, 45 minutes each, no hints. Designed to simulate real pressure — the kind you face during incidents, tight deadlines, or live debugging.

### 2. Goal
Solve 3 problems under time pressure. After each attempt: review together — what you got right, what you missed.

### 3. Scope
*(Problems provided when you reach this challenge)*
One problem each on: interfaces, concurrency, and performance.

### 4. Why This Matters
Performance under pressure is a skill. Writing correct Go under a 45-minute clock while explaining your thinking out loud is different from writing it comfortably at home. This challenge trains that specific skill.

### 5. Tips for Performing Well
- Read the problem twice — misunderstanding costs more time than reading slowly
- Define your data structures before writing logic
- Write the happy path first, then add error handling
- Name things clearly — `workerCount` not `wc`
- Say what you're doing out loud as you type - so that others can also understand

### 6. Knowledge Gained
- ✅ Performing under time pressure
- ✅ Structuring solutions quickly
- ✅ Clear Go style and technical communication

---

# Full Roadmap Summary

| Phase | Challenge | Hours | Level | Topics |
|---|---|---|---|---|
| 2 | Multi-Format Logger | 8–10h | 🟡 Beginner→Intermediate | interfaces, io.Writer, json |
| 2 | Shape Calculator Fix | 3–4h | 🟢 Beginner | type switch, interface bugs |
| 3 | Concurrent Log Scanner | 15–20h | 🟡 Beginner→Intermediate | goroutines, channels, WaitGroup |
| 3 | Worker Pool URL Checker | 15–20h | 🟠 Intermediate | worker pool, context, http client |
| 4 | Benchmark Battle | 15–20h | 🟠 Intermediate→Advanced | benchmarks, sync.Pool, pprof |
| 5 | Mini DevOps Dashboard | 25–30h | 🔴 Intermediate→Advanced | http, json API, mutex, ticker |
| 6 | Code Review Gauntlet | 10h | 🔴 Advanced | spot bugs, leaks, races |
| 6 | Build Under Pressure | 8–10h | 🔴 Advanced | pressure simulation |
| | **Total** | **~110–130h** | | |

---

*Start with Challenge 2.1. Come back when you're done — or when you're stuck.*
*This handbook is designed to be worked through in order. Don't skip phases.*