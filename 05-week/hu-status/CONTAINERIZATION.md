# Week 5 - Session 1: Containerize the Services

> Goal: a multi-stage Dockerfile per service, a `.dockerignore`, and a `docker-compose.yml`
> that brings up all services + the databases on one network, configured via environment
> variables, with data on a volume. This is the runtime base for the MVP 1 release.

---

## 1. Dockerfile per service (multi-stage)

| Service | Builder stage | Runtime stage |
|---|---|---|
| `access-service` | `golang:1.25-alpine` | `alpine:3.21` |
| `membership-service` | `golang:1.25-alpine` | `alpine:3.21` |
| `backend` (`library-api` — Catalog, and the not-yet-built Circulation) | `golang:1.25-alpine` | `alpine:3.21` |
| `frontend` | `node:20-alpine` | `nginx:alpine` |

Each Go service's Dockerfile compiles a static binary in the builder stage
(`CGO_ENABLED=0 GOOS=linux go build`) and copies only the compiled binary plus its
`migrations/` folder into the alpine runtime image — the Go toolchain and source never ship in
the final image. The frontend's Dockerfile builds the React production bundle with Node, then
serves the static output from `nginx:alpine`.

## 2. `.dockerignore`

**Status: not yet added — a real gap found while writing this status.** Every service builds
correctly today because `go.mod`/`go.sum` already pin exact dependencies and the Docker build
context is scoped per service folder, but there is no `.dockerignore` excluding `.git`,
local `.env` files, or editor artifacts from the build context. Logged as a small follow-up
before the `v1.0.0` tag — low risk, but worth fixing (smaller build context, no chance of a
stray local file leaking into an image layer).

## 3. `docker-compose.yml` — one network, env config, volumes

All services and databases run on a single bridge network (`lms-network`), defined once at
the bottom of `docker-compose.yml`. Nothing but the NGINX gateway (port 8080) and the frontend
(port 3000) publish a port to the host — every service and every database is reachable only
from other containers on that network, by container name.

**Services brought up by `docker compose up -d --build`:**

| Container | Role |
|---|---|
| `access-db`, `access-migrate`, `access-service` | Access bounded context, own PostgreSQL instance |
| `membership-db`, `membership-migrate`, `membership-service` | Membership bounded context, own PostgreSQL instance |
| `db`, `migrate`, `backend` | Everything not yet extracted (Catalog today) |
| `nginx` (`api-gateway`) | Single entry point, routes by URL path prefix |
| `frontend` | React SPA |

**Configuration:** every variable (DB credentials, `JWT_SECRET`, `JWT_EXPIRY`, `CORS_ORIGIN`,
`VITE_API_BASE_URL`) is injected through `.env` (copied from `.env.example`) — nothing is
hardcoded in the compose file itself.

**Data:** each database has its own named volume (`access-db-data`, `membership-db-data`,
`lms-db-data`) so data survives a `docker compose down` (without `-v`).

## 4. Verification

```bash
cp .env.example .env   # fill in JWT_SECRET
docker compose up -d --build
curl http://localhost:8080/health
```

Verified locally: every container reaches `healthy`/`running` state, the gateway routes
`/api/v1/auth` to `access-service`, `/api/v1/students` to `membership-service`, and everything
else (e.g. `/api/v1/books`) to `backend`, per `infra/nginx/nginx.conf`.

## Correlations

- Full infrastructure diagram → `library-docs/05-architecture/deployment.md`
- The decomposition this compose file reflects → `library-docs/05-architecture/decisions/records/ADR-004-incremental-microservices-decomposition.md`
- Local setup, step by step → `SETUP.md` at the `lms-library` repo root
