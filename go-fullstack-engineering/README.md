# Go Full-Stack Engineering Curriculum

This curriculum builds a production-grade, cloud-native full-stack engineering skillset from the ground up using Go and Vue.js. You will design and ship real backend services with a focus on Go's concurrency model, interface-driven architecture, and memory performance — then extend them into complete products with a typed Vue.js frontend, containerized delivery via Docker, automated CI/CD pipelines, and Kubernetes-based deployment. Every challenge is grounded in real production scenarios used by engineering teams. By the end, you will have designed, tested, containerized, and deployed a full-stack distributed service end-to-end — and be able to own it in production under a YBIYRI culture.

---

## What this curriculum covers

| Skill | Source |
|---|---|
| Go backend engineering | Phases 1–6 (existing) |
| Interface-driven design | Phase 2 |
| Concurrency model | Phase 3 |
| Memory & performance | Phase 4 |
| HTTP APIs & standard library | Phase 5 |
| TypeScript fundamentals | Phase 7 |
| Vue.js frontend | Phase 8 |
| Unit testing (Go + Vue) | Phase 9 |
| Docker & containerization | Phase 10 |
| CI/CD with GitHub Actions | Phase 11 |
| Kubernetes (conceptual + hands-on) | Phase 12 |
| Full-stack capstone project | Phase 13 |

## Skill progression

Each phase builds on the last. The table below shows what you can do at each milestone.

| Milestone | What you can do |
|---|---|
| After Phase 1 | Write idiomatic Go — structs, slices, error handling, file I/O |
| After Phase 3 | Design and implement concurrent systems without race conditions |
| After Phase 5 | Ship a tested, production-grade HTTP service using Go's standard library |
| After Phase 8 | Build and connect a typed Vue.js frontend to a Go API |
| After Phase 11 | Containerize, test, and deliver software through an automated CI/CD pipeline |
| After Phase 13 | Design, build, test, containerize, and operate a full-stack distributed service end-to-end |

## YBIYRI
> *"You Build It, You Run It"* — a principle from Werner Vogels (AWS CTO).
> The team that builds a feature owns it in production. No hand-off to ops.
> This curriculum is designed around YBIYRI: you will build AND deploy AND monitor everything yourself.

---

# PART 1 — Go Engineering
### Phases 1–6 | ~130 hours | ~4 weeks at 5h/day

This track builds a deep, production-ready understanding of Go as a systems and backend language. You will work through Go's type system and interface-driven design, master its concurrency model using goroutines and channels, profile and optimize memory allocation under realistic workloads, and ship a fully tested HTTP service using Go's standard library. Every challenge is modeled on patterns used in real Go codebases. By the end of this track, you can design, implement, test, and reason about concurrent Go services at an advanced level.

---

# PHASE 7 — TypeScript Fundamentals

> **Why this phase matters**
> TypeScript is JavaScript with a type system — and you already think in types from Go and C. The concepts are familiar; the syntax is new. TypeScript is the language of the Vue.js ecosystem. Every modern frontend codebase uses TypeScript.

---

## Challenge 7.1 — Port Your Go Types to TypeScript
### `🟢 Beginner`
**🕐 Expected duration: 8–10 hours**

### 1. Context
You've spent weeks thinking in Go's type system — structs, interfaces, methods, error handling. TypeScript has equivalents for all of these. The fastest way to learn TypeScript is to take something you've already built in Go and rebuild the type layer in TypeScript. This challenge does exactly that.

### 2. Goal
Recreate the type system from your Go DevOps Dashboard (Phase 5) in TypeScript — types, interfaces, and a small data transformation layer.

### 3. Scope
- Install TypeScript: `npm install -g typescript`
- Create a `types.ts` file that mirrors your Go structs:
  ```typescript
  // Go:  type FileStats struct { Name string; Errors int }
  // TS:  interface FileStats { name: string; errors: number }
  ```
- Create these TypeScript types: `FileStats`, `StatusResponse`, `ErrorsResponse`
- Write a function `filterErrors(files: FileStats[]): FileStats[]` that returns only files with errors
- Write a function `totalErrors(files: FileStats[]): number`
- Write a function `formatReport(response: StatusResponse): string` that returns a human-readable summary
- Add a `type LogLevel = "INFO" | "WARN" | "ERROR"` union type (Go equivalent: `const` iota)
- All functions must have explicit return types — no implicit `any`
- Run with `ts-node index.ts`

### 4. Expected Output
```
Total files: 3
Files with errors: 2
Total errors: 10

log_1.log — 7 errors (INFO: 81, WARN: 12, ERROR: 7)
log_3.log — 3 errors (INFO: 89, WARN: 8, ERROR: 3)
```

### 5. Why This Matters in Production
TypeScript's interface system is Go's interface system — minus the implicit satisfaction. When your Go backend returns JSON and your Vue frontend consumes it, TypeScript interfaces are the contract between them. Mismatched types between backend and frontend cause runtime bugs that are nearly impossible to debug without TypeScript.

### 6. Go → TypeScript Comparison

| Go | TypeScript |
|---|---|
| `struct` | `interface` or `type` |
| `interface` | `interface` (but explicit) |
| `string`, `int`, `bool` | `string`, `number`, `boolean` |
| `[]string` | `string[]` or `Array<string>` |
| `map[string]int` | `Record<string, number>` |
| `error` return | `throw` or `Result<T, E>` pattern |
| `:=` inference | `const`/`let` inference |
| `iota` const | `union type` or `enum` |

### 7. Common Mistakes to Avoid
- Using `any` everywhere — defeats the purpose of TypeScript
- Confusing `interface` and `type` — use `interface` for objects, `type` for unions/primitives
- Forgetting that TypeScript types are erased at runtime — they only exist at compile time
- Using `==` instead of `===` — TypeScript inherits JavaScript's loose equality trap

### 8. What a Senior Would Do Differently
- Generate TypeScript types automatically from Go structs using `tygo` or `openapi-typescript`
- Use `zod` for runtime validation — TypeScript types disappear at runtime, zod catches bad API responses
- Use `strict: true` in `tsconfig.json` — enables all strict type checks

### 9. Hints & Knowledge
- `tsc --init` — creates a `tsconfig.json`
- `ts-node file.ts` — run TypeScript directly without compiling
- `interface` vs `type`: both work for objects; `type` is needed for unions: `type Status = "ok" | "error"`
- Optional fields: `name?: string` (same concept as Go pointer fields)
- `Array.filter()`, `Array.reduce()` — functional equivalents of Go's range loops

### 10. Sources
- TypeScript handbook: https://www.typescriptlang.org/docs/handbook/intro.html
- TS for Go devs: https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html
- `ts-node`: https://typestrong.org/ts-node/

### 11. Knowledge Gained
- ✅ TypeScript type system fundamentals
- ✅ Interfaces and type aliases
- ✅ Union types and type narrowing
- ✅ Go↔TypeScript mental model mapping
- ✅ Working with typed JSON data

---

# PHASE 8 — Vue.js Frontend

> **Why this phase matters**
> Vue.js is the frontend framework used in this project. It's component-based — you build small, reusable UI pieces (components) that react to data changes automatically. If you've done Flutter or React, you already understand reactive UI: state changes → UI updates. Vue is that idea on the web. Combined with TypeScript, Vue gives you a fully typed, component-driven frontend.

---

## Challenge 8.1 — Build a Dashboard UI for Your Go API
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 20–25 hours**

### 1. Context
Your Go DevOps Dashboard from Phase 5 has a JSON API. Right now, to see the data you have to `curl` it from the terminal. That's not usable by a real team. This challenge adds a Vue.js frontend that visualizes the same data — making it a real internal tool that any team member could use.

### 2. Goal
Build a Vue.js + TypeScript frontend that connects to your Phase 5 Go API and displays the log file analysis results in a usable dashboard.

### 3. Scope
**Setup:**
- `npm create vue@latest` — choose TypeScript + Vue Router
- Connect to your running Go API (enable CORS in Go first)

**Pages:**
- `/` — Dashboard: total files processed, total errors, a summary table
- `/errors` — Error view: only files with errors, sorted by error count
- `/files/:name` — Detail view: stats for one specific file

**Components to build:**
- `StatCard.vue` — reusable card showing a number + label (used for total files, total errors, etc.)
- `FileTable.vue` — table of files with columns: name, total lines, errors, warns, infos
- `ErrorBadge.vue` — colored badge: red if errors > 0, green if clean
- `RefreshButton.vue` — button that calls `POST /scan` and refreshes data

**Requirements:**
- All API calls typed with TypeScript interfaces from Challenge 7.1
- Loading state shown while fetching
- Error state shown if API is unreachable
- Auto-refresh every 10 seconds using `setInterval`

### 4. Expected Output
A running web UI at `http://localhost:5173` showing:
```
┌─────────────────────────────────────────┐
│  DevOps Log Dashboard                   │
│                                         │
│  [3 Files]  [10 Errors]  [25 Warns]    │
│                                         │
│  File         Errors  Warns  Infos     │
│  log_1.log    🔴 7     12     81       │
│  log_2.log    ✅ 0      5     95       │
│  log_3.log    🔴 3      8     89       │
│                                         │
│  [🔄 Trigger Scan]                      │
└─────────────────────────────────────────┘
```

### 5. Why This Matters in Production
Every internal tool needs a UI. The Go backend provides data; Vue provides visibility. Under YBIYRI, you own the whole stack — you can't say "that's the frontend team's problem." Being able to quickly build a Vue dashboard for your Go service is a core full-stack skill.

### 6. Common Mistakes to Avoid
- Calling the API directly from components — use a `composable` (`useFiles.ts`) to encapsulate fetch logic
- Not handling loading and error states — always show the user what's happening
- Forgetting CORS — Go must allow requests from `localhost:5173`
- Putting all code in `App.vue` — split into components from the start

### 7. What a Senior Would Do Differently
- Use `pinia` for state management instead of prop-drilling
- Use `axios` or `ofetch` instead of raw `fetch` for better error handling
- Add `vue-query` (TanStack Query) for caching and background refetch
- Use `vitest` for component testing (covered in Phase 9)

### 8. Hints & Knowledge
- CORS in Go: `w.Header().Set("Access-Control-Allow-Origin", "*")`
- `fetch` in Vue: use `onMounted` + `ref` for reactive data
- `v-for` = Go's `range`, `v-if` = Go's `if`, `v-bind` = dynamic attributes
- `defineProps<{ errors: number }>()` — typed props, like Go function parameters
- `ref<FileStats[]>([])` — reactive array with TypeScript type

### 9. Sources
- Vue 3 official tutorial: https://vuejs.org/tutorial/
- Vue + TypeScript: https://vuejs.org/guide/typescript/overview
- Vue Router: https://router.vuejs.org/
- Pinia (state management): https://pinia.vuejs.org/

### 10. Knowledge Gained
- ✅ Vue 3 component model (Composition API)
- ✅ TypeScript in Vue components
- ✅ Reactive state with `ref` and `computed`
- ✅ HTTP calls from a frontend to a Go backend
- ✅ Component decomposition
- ✅ Vue Router for multi-page apps

---

# PHASE 9 — Testing

> **Why this phase matters**
> Testing is not optional in a YBIYRI team. When you own your service in production, untested code = 3am incident. Go has first-class testing built into the language. Vue has Vitest. A proper engineer is expected to write tests as naturally as they write feature code.

---

## Challenge 9.1 — Test Your Go API
### `🟠 Intermediate`
**🕐 Expected duration: 12–15 hours**

### 1. Context
Your Phase 5 Go API works — but how do you *know* it works after every change? Manual `curl` testing doesn't scale. This challenge adds a proper test suite to your Go API so you can refactor with confidence.

### 2. Goal
Write a full test suite for the Phase 5 DevOps Dashboard API covering unit tests, integration tests, and table-driven tests.

### 3. Scope
**Unit tests:**
- Test `processFile()` — mock a file with known content, assert correct counts
- Test the store's `Set()` and `GetAll()` — concurrent writes, no race conditions
- Test `filterErrors()` logic in isolation

**HTTP handler tests** (integration):
- Test `GET /status` → returns correct JSON structure
- Test `GET /errors` → returns only files with errors
- Test `POST /scan` → returns 202, not 405 for GET
- Use `httptest.NewRecorder()` — no real server needed

**Table-driven tests:**
```go
tests := []struct {
    name     string
    input    string
    expected FileStats
}{
    {"all info",  "INFO: a\nINFO: b",        FileStats{Total: 2, Infos: 2}},
    {"has error", "ERROR: x\nINFO: y",       FileStats{Total: 2, Errors: 1, Infos: 1}},
    {"empty",     "",                         FileStats{Total: 0}},
}
```

**Run:**
- `go test ./...` — all tests pass
- `go test -race ./...` — no race conditions
- `go test -cover ./...` — achieve >80% coverage

### 4. Why This Matters in Production
Table-driven tests are idiomatic Go — you'll see them in every serious Go codebase. `httptest` means you can test your full HTTP layer without spinning up a real server. In a CI/CD pipeline (Phase 11), these tests run on every push — catching regressions before they hit production.

### 5. Common Mistakes to Avoid
- Testing implementation details instead of behavior — test what the function *does*, not *how* it does it
- Writing one test case per test function — use table-driven tests for multiple cases
- Not testing error paths — test what happens when the file doesn't exist, the store is empty, etc.
- Using `t.Fatal` vs `t.Error` wrong — `Fatal` stops the test immediately, `Error` continues

### 6. What a Senior Would Do Differently
- Use `testify/assert` for cleaner assertions: `assert.Equal(t, expected, actual)`
- Use interfaces + mocks for external dependencies — mock the file system, not real files
- Generate mocks with `mockgen` from `google/mock`
- Measure test coverage per package: `go test -coverprofile=coverage.out && go tool cover -html=coverage.out`

### 7. Hints & Knowledge
- `testing.T` is the test context — pass it everywhere
- `httptest.NewRecorder()` captures HTTP responses without a real server
- `httptest.NewRequest("GET", "/status", nil)` — fake HTTP request
- `t.Run("name", func(t *testing.T) {...})` — subtests, great for table-driven

### 8. Sources
- Go testing package: https://pkg.go.dev/testing
- `httptest`: https://pkg.go.dev/net/http/httptest
- Table-driven tests: https://go.dev/wiki/TableDrivenTests
- testify: https://github.com/stretchr/testify

### 9. Knowledge Gained
- ✅ Go testing fundamentals — `testing.T`, `t.Run`, `t.Fatal`
- ✅ Table-driven tests — idiomatic Go test pattern
- ✅ HTTP handler testing with `httptest`
- ✅ Test coverage measurement
- ✅ Testing concurrent code with `-race`

---

## Challenge 9.2 — Test Your Vue Frontend
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 8–10 hours**

### 1. Context
Frontend bugs cost user trust. A button that doesn't work, a number that shows `NaN`, a loading spinner that never stops — these are caught by component tests. Vitest + Vue Test Utils is the standard Vue testing stack at many companies currently.

### 2. Goal
Add a unit test suite to your Vue dashboard from Phase 8.

### 3. Scope
- Install: `npm install -D vitest @vue/test-utils jsdom`
- Configure Vitest in `vite.config.ts`

**Tests to write:**
- `StatCard.vue` — renders the correct number and label from props
- `ErrorBadge.vue` — renders red when errors > 0, green when 0
- `FileTable.vue` — renders correct number of rows from a list of files
- `useFiles.ts` composable — mock `fetch`, assert correct state transitions:
  - loading=true while fetching
  - data populated on success
  - error state set on failure
- `filterErrors()` utility — same table-driven approach as Go tests

**Run:** `npx vitest run --coverage`
**Target:** >75% component coverage

### 4. Why This Matters in Production
Component tests catch: props that don't render, computed values that break, conditional rendering bugs. In a YBIYRI team shipping multiple times per day, untested Vue components = production incidents.

### 5. Common Mistakes to Avoid
- Testing Vue internals (refs, computed) instead of rendered output
- Not mocking `fetch` — tests should never make real HTTP calls
- Snapshot testing everything — snapshots break on any UI change, even intentional ones

### 6. Hints & Knowledge
- `mount(Component, { props: {...} })` — mount a component with props
- `wrapper.text()` — get rendered text
- `wrapper.find('.class')` — find DOM element
- `vi.fn()` — mock a function; `vi.stubGlobal('fetch', mockFetch)` — mock fetch

### 7. Sources
- Vitest: https://vitest.dev/
- Vue Test Utils: https://test-utils.vuejs.org/
- Testing composables: https://vuejs.org/guide/scaling-up/testing

### 8. Knowledge Gained
- ✅ Vitest setup and configuration
- ✅ Component testing with Vue Test Utils
- ✅ Mocking HTTP calls in tests
- ✅ Testing composables (Vue's equivalent of Go's testable functions)
- ✅ Frontend test coverage

---

# PHASE 10 — Docker & Containerization

> **Why this phase matters**
> Docker means your code runs the same everywhere — your laptop, your colleague's laptop, staging, production. It eliminates "works on my machine." Kubernetes (Phase 12) runs your containers at scale. You can't do Kubernetes without Docker.

---

## Challenge 10.1 — Containerize the Full Stack
### `🟠 Intermediate`
**🕐 Expected duration: 15–18 hours**

### 1. Context
Your Go API and Vue frontend run locally. But right now, anyone wanting to use them needs to install Go, Node.js, configure ports, run both commands manually. Docker packages everything into containers — isolated, reproducible, deployable anywhere. This challenge puts your entire application into Docker.

### 2. Goal
Containerize your Go API and Vue frontend using Docker, wire them together with Docker Compose, and run the complete stack with one command.

### 3. Scope
**Go API container:**
- Write a `Dockerfile` for the Go API
- Use multi-stage build: build stage (Go compiler) + runtime stage (minimal alpine)
- The final image must be under 20MB
- Expose port 8080

**Vue frontend container:**
- Write a `Dockerfile` for the Vue app
- Build stage: `node` image compiles Vue
- Runtime stage: `nginx` serves the static files
- Expose port 80

**Docker Compose:**
- `docker-compose.yml` wires both containers together
- Go API at `http://api:8080` (internal), `http://localhost:8080` (external)
- Vue frontend at `http://localhost:3000`
- A named volume for the `./logs` directory so log files persist
- Health check on the Go API

**Run everything:**
```bash
docker compose up --build
# → Go API running
# → Vue dashboard at localhost:3000 showing live data
```

### 4. Expected Dockerfile (Go — multi-stage)
```dockerfile
# Stage 1: build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# Stage 2: runtime (no Go compiler — tiny image)
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### 5. Why This Matters in Production
Multi-stage builds are how every serious Go shop builds containers. Without them, your image includes the Go compiler (~300MB). With them, the final image is just your binary + alpine (~15MB). Smaller images = faster deployments, less attack surface, lower cloud storage costs. In a YBIYRI team, you own the Dockerfile — ops doesn't write it for you.

### 6. Common Mistakes to Avoid
- Single-stage Go builds — includes the compiler, huge image
- Not using `.dockerignore` — accidentally copying `node_modules` or Go build cache into the image
- Hardcoding ports — use environment variables: `PORT=8080`
- Running containers as root — always add `USER nonroot` in production
- Not handling graceful shutdown — Go API should handle `SIGTERM` before Docker kills it

### 7. What a Senior Would Do Differently
- Use `docker buildx` for multi-platform builds (arm64 for M-series Macs, amd64 for servers)
- Pin exact image versions: `alpine:3.19.0` not `alpine:latest`
- Use `distroless` instead of alpine for even smaller, more secure images
- Add `HEALTHCHECK` instruction to the Dockerfile, not just compose
- Use build cache mounts: `RUN --mount=type=cache,target=/go/pkg/mod go build`

### 8. Hints & Knowledge
- `docker build -t myapp .` — build an image
- `docker run -p 8080:8080 myapp` — run with port mapping
- `docker compose up --build` — build and start all services
- `docker compose logs -f api` — stream logs from a service
- `CGO_ENABLED=0` — static binary, no C dependencies, works in alpine

### 9. Sources
- Docker multi-stage builds: https://docs.docker.com/build/building/multi-stage/
- Docker Compose: https://docs.docker.com/compose/
- Go Docker best practices: https://docs.docker.com/language/golang/
- `.dockerignore`: https://docs.docker.com/engine/reference/builder/#dockerignore-file

### 10. Knowledge Gained
- ✅ Docker fundamentals — images, containers, layers
- ✅ Multi-stage builds for minimal Go images
- ✅ Docker Compose for multi-service orchestration
- ✅ Volumes for persistent data
- ✅ Health checks and container lifecycle

---

# PHASE 11 — CI/CD with GitHub Actions

> **Why this phase matters**
> CI/CD means every push to the repository automatically runs tests, builds the Docker image, and (optionally) deploys. No manual steps. CI/CD is the safety net that makes fast delivery safe.

---

## Challenge 11.1 — Build a Full CI/CD Pipeline
### `🟠 Intermediate → Advanced`
**🕐 Expected duration: 15–18 hours**

### 1. Context
Right now, to run tests you type `go test ./...` manually. To build Docker images you type `docker build` manually. Every time you forget, broken code might reach production. GitHub Actions automates all of this — triggered on every push, pull request, or release.

### 2. Goal
Build a complete CI/CD pipeline using GitHub Actions that automatically tests, builds, and packages your application on every push.

### 3. Scope
**Pipeline structure (`.github/workflows/ci.yml`):**

```
On: push to main, pull_request

Jobs:
  test-go:
    - Checkout code
    - Set up Go 1.24
    - Run go test -race ./...
    - Run go vet ./...
    - Upload coverage report

  test-vue:
    - Checkout code
    - Set up Node 20
    - npm install
    - npx vitest run --coverage
    - Upload coverage report

  build-docker:
    needs: [test-go, test-vue]  ← only if tests pass
    - Build Go API Docker image
    - Build Vue Docker image
    - Push to GitHub Container Registry (ghcr.io)

  security-scan:
    needs: build-docker
    - Run trivy scanner on Docker image
    - Fail if HIGH/CRITICAL vulnerabilities found
```

**Branch protection rule (configure in GitHub):**
- `main` branch requires: `test-go` + `test-vue` to pass before merge

**Add a badge to your README:**
```markdown
![CI](https://github.com/yourname/repo/actions/workflows/ci.yml/badge.svg)
```

### 4. Expected Pipeline Output
```
✅ test-go        (1m 23s)  — 47 tests passed, 82% coverage
✅ test-vue       (0m 45s)  — 23 tests passed, 76% coverage
✅ build-docker   (2m 12s)  — images pushed to ghcr.io
✅ security-scan  (0m 58s)  — 0 critical vulnerabilities
```

### 5. Why This Matters in Production
Every Pull Request gets automatic test results. No PR merges without green tests. Docker images are built automatically — no manual `docker push`. CI/CD is the mechanism that makes "you run it" safe — you don't ship without passing tests.

### 6. Common Mistakes to Avoid
- Running tests without `-race` flag in CI — race conditions only appear in concurrent tests
- Not caching `node_modules` and Go module cache — CI is slow without caching
- Pushing Docker images on every commit including feature branches — only push on `main`
- Not pinning Action versions: use `actions/checkout@v4` not `actions/checkout@latest`
- Storing secrets in code — use GitHub Secrets for tokens and passwords

### 7. What a Senior Would Do Differently
- Add `golangci-lint` — the standard Go linter bundle used at most companies
- Add `dependabot` for automatic dependency updates
- Use matrix builds to test on Go 1.22, 1.23, 1.24 simultaneously
- Add a `deploy` job that pushes to a real server or Kubernetes cluster
- Use semantic versioning: tag releases, use `docker/metadata-action` for image tags

### 8. Hints & Knowledge
- GitHub Actions YAML: `on`, `jobs`, `steps`, `uses`, `run`
- `needs: [job1, job2]` — job dependency (fan-in)
- `secrets.GITHUB_TOKEN` — automatically available, no setup needed for ghcr.io
- `actions/cache@v4` — cache Go modules and Node modules between runs
- `docker/build-push-action@v5` — builds and pushes Docker image in one step

### 9. Sources
- GitHub Actions: https://docs.github.com/en/actions
- Go CI workflow: https://github.com/actions/setup-go
- Docker + GitHub Actions: https://docs.docker.com/build/ci/github-actions/
- Trivy security scanner: https://github.com/aquasecurity/trivy-action

### 10. Knowledge Gained
- ✅ GitHub Actions — triggers, jobs, steps, secrets
- ✅ Automated testing in CI
- ✅ Docker image building and publishing in CI
- ✅ Security scanning with Trivy
- ✅ Branch protection and PR gating
- ✅ CI/CD pipeline design

---

# PHASE 12 — Kubernetes (Conceptual + Hands-On Basics)

> **Why this phase matters**
> Kubernetes runs your Docker containers at scale — multiple copies, automatic restarts, load balancing, rolling deployments. Kubernetes is what cloud-native distributed systems run on. You don't need to be a Kubernetes admin — but understanding each component and writing basic manifests for your own services is part of owning the full stack.

---

## Challenge 12.1 — Understand and Deploy to Kubernetes
### `🟠 Intermediate`
**🕐 Expected duration: 15–18 hours**

### 1. Context
Your Go API and Vue frontend are containerized. Now imagine your company runs 50 such services. Someone needs to manage them — restart crashed containers, route traffic, manage secrets, scale under load. That's Kubernetes. In a YBIYRI culture, you write the Kubernetes manifests for your own service.

### 2. Goal
Understand the Kubernetes architecture conceptually, then write the manifests to deploy your full-stack application to a local Kubernetes cluster.

### 3. Scope

**Part A — Conceptual (no code, draw diagrams, explain components):**
Be able to explain (in your own words, possibly with a diagram):
- What is a **Pod** — smallest deployable unit (your container)
- What is a **Deployment** — manages multiple copies of a Pod
- What is a **Service** — stable network address for a set of Pods
- What is a **Ingress** — routes external HTTP traffic to Services
- What is a **ConfigMap** — non-secret configuration (like env vars)
- What is a **Secret** — sensitive config (passwords, tokens)
- What is a **Namespace** — logical isolation (like separate folders)
- What is a **Node** — a machine (VM or physical) in the cluster
- The control plane: **API Server**, **Scheduler**, **etcd**, **Controller Manager**

**Part B — Hands-on (write manifests, deploy locally):**
- Install `minikube` for local Kubernetes
- Write a `deployment.yaml` for your Go API:
  - 2 replicas
  - Resources: `requests: cpu 100m, memory 64Mi` / `limits: cpu 200m, memory 128Mi`
  - Liveness probe: `GET /healthz`
  - Environment variable from ConfigMap: `LOG_DIR=/logs`
- Write a `service.yaml` — ClusterIP for the API
- Write a `deployment.yaml` for Vue (nginx serving static files)
- Write a `service.yaml` — LoadBalancer for Vue (exposed externally)
- Deploy: `kubectl apply -f k8s/`
- Verify: `kubectl get pods`, `kubectl logs`, `kubectl describe`
- Simulate a crash: `kubectl delete pod <name>` — watch Kubernetes restart it automatically

### 4. Expected Architecture
```
                    Internet
                        │
                   [ Ingress ]
                   /         \
          [Vue Service]   [API Service]
               │                │
         [Vue Pods x2]   [API Pods x2]
                                │
                          [ConfigMap: LOG_DIR]
```

### 5. Why This Matters in Production
In a cloud-native team, every engineer writes Kubernetes manifests for their service. You don't hand it to "the DevOps team." Resource limits prevent one service from starving the cluster. Liveness probes ensure Kubernetes restarts unhealthy pods automatically.

### 6. Common Mistakes to Avoid
- Not setting resource limits — one runaway pod can crash the whole node
- Using `latest` Docker tag in manifests — always pin a version tag
- Forgetting liveness and readiness probes — without them, Kubernetes sends traffic to unhealthy pods
- Putting secrets in ConfigMaps — use `kind: Secret` for sensitive values
- Not understanding the difference between a Service and an Ingress

### 7. What a Senior Would Do Differently
- Use Helm charts instead of raw YAML for parameterized deployments
- Use `kustomize` for environment-specific config (dev/staging/prod)
- Add a readiness probe separate from liveness probe
- Use `HorizontalPodAutoscaler` to scale based on CPU/memory metrics
- Use `PodDisruptionBudget` to ensure availability during rolling updates

### 8. Hints & Knowledge
- `kubectl apply -f file.yaml` — deploy a resource
- `kubectl get pods -w` — watch pods start
- `kubectl logs <pod>` — see logs
- `kubectl describe pod <pod>` — debug startup issues
- `kubectl port-forward service/api 8080:8080` — access service locally
- `minikube tunnel` — expose LoadBalancer services on Mac/Windows

### 9. Sources
- Kubernetes concepts: https://kubernetes.io/docs/concepts/
- kubectl cheatsheet: https://kubernetes.io/docs/reference/kubectl/cheatsheet/
- minikube: https://minikube.sigs.k8s.io/docs/start/
- Resource management: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/

### 10. Knowledge Gained
- ✅ Kubernetes architecture — all core components
- ✅ Pod, Deployment, Service, Ingress, ConfigMap, Secret
- ✅ Writing Kubernetes manifests (YAML)
- ✅ `kubectl` — the essential Kubernetes CLI
- ✅ Self-healing infrastructure (automatic pod restart)
- ✅ Resource management and health probes

---

# PHASE 13 — Capstone: The Full-Stack Production Service

> **Why this phase matters**
> Everything you've learned — Go, TypeScript, Vue, testing, Docker, CI/CD, Kubernetes — comes together here. This is the final challenge.

---

## Challenge 13.1 — Build and Ship a Production-Grade Full-Stack Service
### `🔴 Advanced`
**🕐 Expected duration: 40–60 hours**

### 1. Context
You are the sole engineer on a new internal product: a **real-time deployment tracker** for a small engineering team. Teams push code, deployments happen, and the tracker records what was deployed, by whom, to which environment, and whether it succeeded. Engineers need a dashboard to see the current state of all services across environments.

This is a realistic tool that real companies build internally. It touches every skill in this curriculum.

### 2. Goal
Build, test, containerize, and deploy a full-stack Go + Vue application that a real team could use. Ship it with CI/CD and document it like an open-source project.

### 3. Scope

**Backend (Go):**
- `POST /deployments` — record a new deployment
  ```json
  { "service": "api", "version": "v1.2.3", "env": "production", "deployed_by": "tom" }
  ```
- `GET /deployments` — list all deployments (with filtering: `?env=production&service=api`)
- `GET /deployments/:id` — get one deployment
- `GET /services` — list unique services and their latest deployment per env
- `GET /metrics` — basic stats: total deployments, success rate, avg per day
- In-memory store with `sync.RWMutex` (no database needed)
- Structured JSON logging with `log/slog`
- Graceful shutdown on `SIGTERM` (important for Kubernetes)
- Full unit + integration test suite, >80% coverage

**Frontend (Vue + TypeScript):**
- Dashboard page: service grid showing latest deployment per environment (prod/staging/dev)
- Color coding: green = deployed < 1h ago, yellow = < 24h, red = > 24h
- Timeline page: scrollable list of all recent deployments, filterable by service/env
- `POST /deployments` form — simulate a deployment (for demo purposes)
- Error and loading states on all views
- Component tests with Vitest

**DevOps:**
- Multi-stage Dockerfile for Go (target: < 20MB image)
- Multi-stage Dockerfile for Vue (nginx-served)
- `docker-compose.yml` — full stack runs with one command
- `.github/workflows/ci.yml`:
  - Go tests + race detection
  - Vue tests + coverage
  - Docker build on test pass
  - Image push to ghcr.io on merge to main
- `k8s/` directory with manifests for Kubernetes deployment
- `README.md` that explains: what it does, how to run locally, how to run in Docker, architecture diagram

### 4. Expected Final State
```bash
# Anyone can clone and run with:
git clone github.com/yourname/deployment-tracker
docker compose up

# Then visit:
# http://localhost:3000 — Vue dashboard
# http://localhost:8080/deployments — Go API

# CI badge on README shows green
# k8s/ manifests ready for kubectl apply
```

### 5. Why This Matters in Production
This project ties together every skill in the curriculum:
- Go backend with proper concurrency, testing, and API design
- Vue frontend with TypeScript, components, and testing
- Docker + CI/CD — fully automated build and delivery
- Kubernetes manifests — operations thinking, not just code
- YBIYRI — you built it, you can run it, you can monitor it

### 6. Common Mistakes to Avoid
- Scope creep — don't add a database, auth, or websockets; finish what's scoped first
- Not writing tests as you go — don't leave testing for the end
- Skipping graceful shutdown — Kubernetes sends SIGTERM before killing; handle it
- README as afterthought — write it as if a stranger needs to run your project
- Committing `.env` files, Docker credentials, or secrets

### 7. What a Senior Would Do Differently
- Add `GET /healthz` and `GET /readyz` — standard Kubernetes probe endpoints
- Add Prometheus metrics endpoint `GET /metrics` with `promhttp`
- Use `log/slog` with a JSON handler — structured logs for production
- Add pagination to `GET /deployments` — never return unbounded lists
- Write an ADR (Architecture Decision Record) explaining key design choices
- Add `golangci-lint` and `eslint` to CI

### 8. Architecture Diagram
```
┌─────────────────────────────────────────────────┐
│                   Docker Compose                 │
│                                                  │
│   ┌─────────────┐         ┌──────────────────┐  │
│   │  Vue + Nginx │ ──────▶│   Go API :8080   │  │
│   │   :3000      │  HTTP  │                  │  │
│   └─────────────┘         │  ┌────────────┐  │  │
│                            │  │ In-memory  │  │  │
│                            │  │   Store    │  │  │
│                            │  │ (Mutex)    │  │  │
│                            │  └────────────┘  │  │
│                            └──────────────────┘  │
└─────────────────────────────────────────────────┘

CI/CD:
push → GitHub Actions → test → build → push image → (k8s deploy)
```

### 9. Sources
- Go `log/slog`: https://pkg.go.dev/log/slog
- Graceful shutdown: https://pkg.go.dev/net/http#Server.Shutdown
- Prometheus Go client: https://github.com/prometheus/client_golang
- Vue 3 + TypeScript: https://vuejs.org/guide/typescript/overview
- GitHub Actions secrets: https://docs.github.com/en/actions/security-guides/encrypted-secrets

### 10. Knowledge Gained
- ✅ Full-stack Go + Vue application from scratch
- ✅ Production-grade API design
- ✅ Frontend + backend integration end-to-end
- ✅ Complete test coverage across the stack
- ✅ Docker multi-stage, Docker Compose
- ✅ GitHub Actions CI/CD pipeline
- ✅ Kubernetes deployment manifests
- ✅ YBIYRI — you built it, you can run it

---

# Full Curriculum Summary

| Phase | Challenge | Hours | Level | Topics |
|---|---|---|---|---|
| 1 | Go Foundations | 10h ✅ | 🟢 | syntax, slices, structs, I/O |
| 2 | Multi-Format Logger | 8–10h | 🟡 | interfaces, io.Writer, json |
| 2 | Shape Calculator Fix | 3–4h | 🟢 | type switch, interface bugs |
| 3 | Concurrent Log Scanner | 15–20h | 🟡 | goroutines, channels, WaitGroup |
| 3 | Worker Pool URL Checker | 15–20h | 🟠 | worker pool, context, http |
| 4 | Benchmark Battle | 15–20h | 🟠 | benchmarks, sync.Pool, pprof |
| 5 | Mini DevOps Dashboard (Go API) | 25–30h | 🔴 | http server, mutex, ticker |
| 6 | Code Review Gauntlet | 10h | 🔴 | spot bugs, leaks, races |
| 6 | Build Under Pressure | 8–10h | 🔴 | Pressure simulation |
| 7 | Port Go Types to TypeScript | 8–10h | 🟢 | TS types, interfaces, generics |
| 8 | Vue Dashboard for Go API | 20–25h | 🟡 | Vue 3, composition API, fetch |
| 9 | Test Go API | 12–15h | 🟠 | table tests, httptest, coverage |
| 9 | Test Vue Frontend | 8–10h | 🟡 | vitest, vue test utils, mocks |
| 10 | Containerize Full Stack | 15–18h | 🟠 | Docker, multi-stage, compose |
| 11 | CI/CD Pipeline | 15–18h | 🟠 | GitHub Actions, ghcr.io, trivy |
| 12 | Kubernetes Deploy | 15–18h | 🟠 | k8s concepts, manifests, kubectl |
| 13 | Capstone — Deployment Tracker | 40–60h | 🔴 | everything, end-to-end |
| | **Total** | **~250–290h** | | |

---

*Complete phases in order. Don't skip. Each phase builds on the previous one.*
*Phase 13 is your the final project — start the GitHub repo on day 1 and commit as you go.*