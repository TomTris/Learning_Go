The project is currently [`LIVE`](https://incident-handoff-latest.onrender.com/). Free tier spins down after inactivity — first load may take ~15-60s to wake.

# Incident Handoff

A backend service for on-call engineers to hand off active incidents cleanly.

When you've been firefighting an incident for hours and have to pass it to the next
person on rotation, scattered Slack messages lose the context that matters. Incident
Handoff captures timestamped actions, observations, and status during an incident and
produces a structured handoff — so the next responder sees what happened, what was
tried, and what's next.

> Status: backend complete and tested; Vue frontend in progress.

---

## Architecture

The backend is written in **Go** and exposes a REST API plus a WebSocket channel for
real-time updates. State lives in **MongoDB**, but the storage layer sits behind an
interface, so the same logic runs against an in-memory store (used in tests) or MongoDB
(used in production) without changes.

```
Client (Vue) ──REST──▶ Go API ──▶ IncidentStore (interface)
     │                    │              ├── in-memory  (tests)
     └────WebSocket───────┘              └── MongoDB     (prod)
                          │
                          └──▶ Prometheus /metrics
```

Key design choices:

- **Storage behind an interface** (`IncidentStore`, `OnCallStore`) — swappable
  in-memory and MongoDB implementations.
- **Contract tests** — one shared test suite runs against *both* implementations to
  guarantee they behave identically.
- **Optimistic locking** — updates carry an expected version to prevent lost writes
  when two responders edit the same incident.
- **Observability built in** — custom Prometheus metrics, not an afterthought.

---

## Tech stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| API | REST + WebSocket (`gorilla/websocket`) |
| Auth | JWT (`golang-jwt/v5`), bcrypt password hashing |
| Storage | MongoDB (`mongo-driver/v2`) + in-memory implementation |
| Metrics | Prometheus (`client_golang`) |
| Frontend | Vue 3 + TypeScript (Vite) — *in progress* |
| CI | GitHub Actions |
| Container | Docker |

---

## Features

- Create incidents with severity (SEV1–SEV3), service, and status lifecycle
  (`TRIGGERED` → ...).
- Timeline entries — timestamped actions, observations, and hypotheses during an
  incident.
- Structured handoff brief generated from the timeline.
- On-call registry — track who currently holds the pager.
- Real-time sync between connected clients via WebSocket.
- Feature flags for toggling behaviour without redeploys.

---

## Observability

The service exposes Prometheus metrics on `/metrics`, including:

- `handoff_http_requests_total` — request count by method, path, status code
- `handoff_http_request_duration_seconds` — latency histogram
- `handoff_incidents_total` — current incidents by status
- `handoff_entries_total` — total timeline entries created
- WebSocket connection gauge

---

## Testing

~50 tests across the backend, covering handlers, auth, middleware, config, and both
storage implementations. The store layer uses **contract tests** — a single suite
(`runStoreContractsTests`) is executed against the in-memory and MongoDB stores so any
behavioural drift between them fails the build.

```bash
make test
```

---

## Continuous Integration

GitHub Actions runs on every push and pull request to `main`:

1. **Test** — spins up a MongoDB service container and runs the full Go test suite.
2. **Build & publish** — on push, builds the Docker image and pushes it to
   GitHub Container Registry (`ghcr.io`).

> Note: this is CI (automated testing + image build/publish). Deployment to the running
> environment (currently Render) is handled separately, not through this pipeline.

---

## Running locally

Requirements: Go 1.26+, a running MongoDB instance, Docker (optional).

fill "HANDOFF_JWT_SECRET=" and
```bash
make
```


---

## Roadmap

- [ ] Complete the Vue frontend (views scaffolded, wiring to API in progress)
- [ ] Grafana dashboard for the exported metrics
- [ ] Deployment pipeline (currently CI only)

---

*Built from scratch to understand every layer — API, storage, auth, real-time,
metrics, and testing.*