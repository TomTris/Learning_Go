The project is [`LIVE`](https://incident-handoff-latest.onrender.com/). Free tier spins down after inactivity — first load may take ~15-60s to wake.


# Handoff `//`

A full-stack **incident handoff** tool: a Go HTTP backend plus a Vue 3 frontend that turn the messy, under-pressure moment of passing an active incident from one engineer to another into a clean, structured brief.

---

## What problem does this solve?

In teams practicing **YBIYRI** ("You Build It, You Run It"), the engineers who write the code carry
the pager. When an incident outlasts the individual — fatigue after hours of firefighting,
cross-team escalation, timezone handover — one engineer's accumulated context needs to reach another
engineer **intact, under pressure, during an active incident**. In practice it doesn't. It lands in
Slack as fragments, and the next responder burns precious minutes reconstructing *what already
happened* before they can help.

**Handoff** fixes the handoff itself:

In short, the repo turns ad-hoc incident chatter into a shareable, role-aware, structured record
optimized for the moment one person takes over from another.

- It captures **timestamped, typed actions** as the engineer works — every note is an `observation`,
  `action`, `discovery`, `open_question`, or `state_change`, so the history stays skimmable instead
  of becoming a wall of chat.
- It generates a **structured brief for the next person**: what was done, what's still open, and
  where to start — served by `GET /api/incidents/{id}/handoff` (entries, actions, open questions,
  elapsed time, and how many times the pager has changed hands). The UI's catch-up panel currently
  derives its stats client-side; wiring it to the brief endpoint is pending (Phase 7.5).
- It tracks **pager ownership** (`on_call`) per incident, so handing off is an explicit, recorded act.
- It is **role-aware**: any authenticated engineer can contribute timeline entries to any incident
  (every responder can add context), but **changing incident state** — status, severity, pager
  reassignment — is restricted to an admin or the current on-call engineer (`policy.go`).

**Design notes:** architecture in [DESIGN.md](DESIGN.md) · tradeoffs and reasoning (optimistic
locking / TOCTOU, timing-attack-resistant login, and more) in
[Decision-Documentation.md](Decision-Documentation.md).

### Progress
- [x] Phase 1 — Production Go HTTP Service
- [x] Phase 2 — Database Integration
- [x] Phase 3 — WebSocket & Real-Time
- [x] Phase 4 — Observability & Feature Flags
- [x] Phase 5 — Authentication
- [x] Phase 5.5 — Testing Backend
- [x] Phase 6 — TypeScript + Vue.js
- [x] Phase 7 — Handoff Frontend, REST flows (incident timeline, catch-up panel, on-call schedule & admin management)
- [ ] Phase 7.5 — Live timeline: WebSocket client in the SPA; wire catch-up panel to `GET /incidents/{id}/handoff`
- [ ] Phase 8 — Testing Frontend (unit + e2e scaffolding in place, coverage being expanded)
- [x] Phase 9 — Ship It

The project is [`LIVE`](https://incident-handoff-latest.onrender.com/)

### How the pieces fit together

- **Layered HTTP.** Every request passes through `RequestID → Observability → Timeout` middleware.
  Routes are grouped into three mounts in `router.go`:
  - `/api/*` — authenticated app routes (incidents, entries, handoff brief, auth/me, incident WebSocket)
  - `/admin/*` — admin-only routes (feature flags, on-call shifts)
  - `/*` — public routes (`/login`, `/registration`, `/healthz`, `/readyz`, `/metrics`) **and the SPA fallback**
- **Auth.** `POST /login` verifies credentials and sets an `HttpOnly`, `Secure`, `SameSite=Strict`
  **JWT cookie** (`access_token`, HS256). `AuthMiddleware` parses it and injects a `UserContext`
  (id, username, role) into the request; `AuthAdminOnlyMiddleware` gates the admin mux.
- **Storage is pluggable.** Incidents, users, and the on-call roster each have **in-memory** and
  **MongoDB** implementations behind store interfaces, chosen at boot from
  `HANDOFF_CONNECT_STRING` (empty → in-memory); contract tests run against both. An instrumented
  decorator records incident-store metrics. Feature flags use an in-memory store, created empty at
  boot and managed via `/admin/flags`.
- **Real-time (backend side).** Each incident exposes `GET /api/incidents/{id}/ws`; a `Registry`/`Hub`
  fans new-entry and incident-update events out to connected clients over WebSocket. The frontend
  does not subscribe yet — the UI currently updates on refetch; the WebSocket client is the main
  remaining piece of Phase 7.
- **Observability.** Prometheus metrics at `/metrics`, structured `slog` request logs with request IDs,
  and `/healthz` / `/readyz` (the latter pings MongoDB) for probes.
- **Frontend.** The SPA talks only to the JSON API; in dev, Vite proxies `/api`, `/admin`, `/login`,
  `/registration` to `:8080`, and in production the Go server serves `frontend-vue/dist` with an SPA
  fallback so client-side routes resolve.

### Tech stack

| Layer        | Tech                                                                 |
| ------------ | -------------------------------------------------------------------- |
| Backend      | Go 1.26, `net/http` (std mux), `golang-jwt`, `gorilla/websocket`     |
| Database     | MongoDB 7 (replica set) — or in-memory fallback                      |
| Observability| Prometheus client, `log/slog`                                        |
| Frontend     | Vue 3 + TypeScript, Vite, vue-router, Pinia                          |
| Tests        | Go `testing` (unit + contract, both store implementations); frontend: Vitest (components) + Playwright (e2e), minimal — expanding in Phase 8 |

---

## Running locally

**Prerequisites:** Go 1.26+, Node 20.19+/22.12+, and Docker (for MongoDB; optional — see below).

```sh
# 1. Configure environment
cp .env.example .env          # then set HANDOFF_JWT_SECRET to any non-empty value

# 2. Run everything (starts MongoDB via docker compose, builds the frontend,
#    serves the SPA at http://localhost:8080)
make run

# — or, without Docker: leave HANDOFF_CONNECT_STRING empty for the in-memory store and run
#   cd frontend-vue && npm install && npm run build && cd ../backend-go && go run .
```

For frontend development with hot reload, run the API (`go run .`) and, in another terminal,
`cd frontend-vue && npm install && npm run dev` (Vite proxies API calls to `:8080`).

### Configuration

| Env var                  | Default            | Purpose                                              |
| ------------------------ | ------------------ | ---------------------------------------------------- |
| `HANDOFF_JWT_SECRET`     | — (**required**)   | HMAC secret for signing auth tokens                  |
| `HANDOFF_PORT`           | `8080`             | HTTP listen port                                     |
| `HANDOFF_CONNECT_STRING` | `""`               | MongoDB URI; empty → in-memory store                 |
| `HANDOFF_DB`             | `incident_tracker` | MongoDB database name                                |
| `HANDOFF_LOG_LEVEL`      | `info`             | Log level                                            |
| `HANDOFF_ENV`            | `development`      | Environment label                                    |

### Tests

```sh
make test        # go test with HTML coverage report
make test-race   # race detector + coverage

cd frontend-vue
npm run test:unit -- --run   # component tests (Vitest)
npm run test:e2e             # e2e (Playwright; first run: npx playwright install)
```

### Accounts

There are no seeded accounts — a fresh boot starts with zero users. Create one via the
registration page (`/registration`): choose the `engineer` role, or `admin` with the admin token
(a hardcoded constant in this demo — see `DefaultAdminToken`; a real deployment would move this
to configuration). Note: with the in-memory store, accounts and data reset on every restart —
including the free-tier live deployment when it spins down.