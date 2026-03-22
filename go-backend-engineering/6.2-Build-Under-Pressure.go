# Challenge 6.2 — Build Under Pressure

> **Rules:**
> - Each problem has a **45-minute hard limit**. Set a timer.
> - No hints, no searching, no AI.
> - Write code that **compiles and runs correctly**.
> - After each problem — stop. Review before moving to the next.
> - Say what you're doing out loud as you type — practice explaining your reasoning clearly.

---

## Problem 1 — Interfaces
### `🟡 Beginner → Intermediate` | ⏱ 45 minutes

### Scenario
You are building a notification system for an e-commerce platform. The platform needs to send notifications through different channels — but the business logic that triggers notifications should not care which channel is used.

### Task
Build a notification dispatcher with the following requirements:

**1. Define a `Notifier` interface with:**
```
Send(to string, subject string, body string) error
```

**2. Implement 3 notifiers:**
- `EmailNotifier` — prints: `[EMAIL] to: <to> | subject: <subject> | body: <body>`
- `SMSNotifier` — prints: `[SMS] to: <to> | <body>` (SMS has no subject)
- `SlackNotifier` — prints: `[SLACK] #<to> | *<subject>* | <body>` (to = channel name)

**3. `SMSNotifier` must reject messages longer than 160 characters — return an error.**

**4. Write a function:**
```go
func Dispatch(notifiers []Notifier, to, subject, body string) {
    // send to ALL notifiers
    // if one fails, print the error and continue — don't stop
}
```

**5. In `main()`:**
- Create one of each notifier
- Dispatch a normal message to all three
- Dispatch a message with a 200-character body — show that SMS fails gracefully

### Expected Output
```
[EMAIL] to: user@example.com | subject: Order Shipped | body: Your order #1234 is on the way!
[SMS]   to: +49123456789 | Your order #1234 is on the way!
[SLACK] #ops-alerts | *Order Shipped* | Your order #1234 is on the way!

[EMAIL] to: user@example.com | subject: Newsletter | body: Lorem ipsum....(200 chars)
[SMS]   ERROR: message too long (200 chars, max 160)
[SLACK] #ops-alerts | *Newsletter* | Lorem ipsum...(200 chars)
```

---

## Problem 2 — Concurrency
### `🟠 Intermediate` | ⏱ 45 minutes

### Scenario
You are building a deployment pipeline tool. When a deployment is triggered, it must run 4 checks simultaneously — all must pass before deployment proceeds. If any check fails, the entire deployment must be cancelled immediately, and the reason reported.

### Task
Build a concurrent pipeline checker with these requirements:

**1. Define a `Check` struct:**
```go
type Check struct {
    Name string
    Run  func() error  // the check logic
}
```

**2. Write:**
```go
func RunChecks(checks []Check) error
```
- Runs ALL checks **concurrently**
- If **any** check returns an error → immediately cancel remaining checks and return that error
- If **all** pass → return nil
- Must use `context` for cancellation
- Must NOT leak goroutines — all goroutines must exit

**3. In `main()`:**

Scenario A — all pass:
```go
checks := []Check{
    {Name: "unit tests",     Run: func() error { time.Sleep(100ms); return nil }},
    {Name: "lint",           Run: func() error { time.Sleep(50ms);  return nil }},
    {Name: "security scan",  Run: func() error { time.Sleep(200ms); return nil }},
    {Name: "docker build",   Run: func() error { time.Sleep(150ms); return nil }},
}
```

Scenario B — one fails:
```go
// security scan returns an error after 80ms
```

### Expected Output
```
// Scenario A:
[check] running: unit tests
[check] running: lint
[check] running: security scan
[check] running: docker build
[check] passed: lint (51ms)
[check] passed: unit tests (101ms)
[check] passed: docker build (151ms)
[check] passed: security scan (201ms)
✅ all checks passed — deploying!

// Scenario B:
[check] running: unit tests
[check] running: lint
[check] running: security scan
[check] running: docker build
[check] passed: lint (51ms)
❌ deployment cancelled: security scan failed: vulnerability found in dependency
```

---

## Problem 3 — Design + Performance
### `🔴 Intermediate → Advanced` | ⏱ 45 minutes

### Scenario
You are joining a team that runs a high-traffic API. One endpoint — `GET /stats` — is called 10,000 times per second. It reads from an in-memory store, computes some statistics, and returns JSON. The team says it's causing GC pressure. You need to redesign it.

### Task
You are given this working but slow implementation:

```go
type Store struct {
    mu    sync.RWMutex
    items map[string]int
}

func (s *Store) Add(key string, value int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items[key] = value
}

func (s *Store) GetStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()

    total := 0
    for _, v := range s.items {
        total += v
    }

    return map[string]interface{}{
        "count": len(s.items),
        "total": total,
        "avg":   float64(total) / float64(len(s.items)),
    }
}

func statsHandler(store *Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        stats := store.GetStats()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    }
}
```

**Requirements:**

1. Identify **all performance problems** in the code above (there are at least 3)

2. Rewrite `GetStats()` and `statsHandler()` to:
   - Use a **typed struct** instead of `map[string]interface{}` for the response
   - Pre-allocate the JSON encoder using a `sync.Pool`
   - Cache the stats result for 1 second — don't recompute on every call

3. Write a **benchmark** comparing old vs new `GetStats()` under concurrent load:
```go
func BenchmarkGetStatsOld(b *testing.B) { ... }
func BenchmarkGetStatsNew(b *testing.B) { ... }
```

4. In comments, explain: **why does `map[string]interface{}` cause GC pressure?**

### Expected benchmark improvement
```
BenchmarkGetStatsOld-8    100000    15234 ns/op    640 B/op    8 allocs/op
BenchmarkGetStatsNew-8   1000000     1823 ns/op     64 B/op    1 allocs/op
```

---

## After Each Problem — Self Review Checklist

Before looking at the solution, ask yourself:

```
□ Does it compile?
□ Does it handle errors — all of them?
□ Is there any goroutine that could leak?
□ Would `go run -race .` pass?
□ Did I use the simplest solution, or did I over-engineer?
□ Can I explain every line I wrote?
```

---

## Self-Assessment Guide

| Result | What it means |
|---|---|
| Solved all 3 in time, clean code | Advanced — strong grasp of Go under pressure |
| Solved 2 fully, 1 partially | Solid — intermediate-to-advanced understanding |
| Solved 2 fully | Good — comfortable with core Go patterns |
| Solved 1 fully, explained approach for others | Progressing — concepts are there, speed will come |
| Could not finish but explained clearly | Keep going — understanding matters more than speed |

> Note: **explaining your thinking** counts as much as the code.
> A partial solution with clear reasoning beats silent perfect code.