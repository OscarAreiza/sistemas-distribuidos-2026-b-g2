# Week 6 - Session 1: Orchestrate the Whole System with One `docker compose up`

> Goal: shared network, health checks gating startup, config via env, data in volumes — a single
> `docker compose up` should only report the system "up" once every service in it can actually
> serve traffic, not just once every container has started.

---

## 1. Starting point (inherited from Week 05)

Network, env-based config, and volume-backed data were already in place from the Week 05
containerization work (see `05-week/hu-status/CONTAINERIZATION.md`): one bridge network
(`lms-network`), everything configured through `.env`/`.env.example`, and one named volume per
database. What Week 05 did **not** finish was making health checks gate the *entire* startup
chain, not just the database layer.

## 2. Gap found: health checks were not gating the full chain

| Container | Had a `healthcheck`? | Was it actually awaited by its dependents? |
|---|---|---|
| `access-db`, `membership-db`, `catalog-db` | ✅ (`pg_isready`) | ✅ — their `-migrate`/`-service` containers use `condition: service_healthy` |
| `catalog-service` | ✅ (`wget` on `/health`) | ❌ — `nginx` only listed it under a plain `depends_on:` (list form), which waits for the *container to start*, not for it to be healthy |
| `access-service` | ❌ — no `healthcheck` at all, even though the code already exposes `/health` and `/health/ready` (`access-service/internal/infrastructure/http/router.go`) | ❌ |
| `membership-service` | ❌ — same: `/health`/`/health/ready` exist in code (`membership-service/internal/infrastructure/http/router.go`) but were never wired into compose | ❌ |
| `nginx` (`api-gateway`) | ❌ | `frontend` only used a plain `depends_on:` on it |

**Net effect before this fix:** `docker compose up -d` could report the whole stack "up" while
`access-service`/`membership-service` were still starting (DB pool not ready yet), and while
`nginx` could still be serving `502`s from a `catalog-service` that had started but not yet
passed its own healthcheck.

## 3. Fix applied this session

Edited `docker-compose.yml` (uncommitted in the working tree at the time of writing — see
Section 5):

- Added the same `healthcheck` pattern `catalog-service` already used (`wget --spider` against
  `http://localhost:8080/health`, `interval: 10s`, `timeout: 3s`, `retries: 5`,
  `start_period: 5s`) to **`access-service`** and **`membership-service`**, since both already
  expose `/health` in code — no application change needed, only wiring.
- Added a `healthcheck` to **`nginx`** itself (same pattern, against its own `/health` path,
  which `infra/nginx/nginx.conf` proxies to `catalog-service`).
- Changed `nginx`'s `depends_on` from a plain list (`- catalog-service`, `- access-service`,
  `- membership-service`) to the map form with `condition: service_healthy` on all three.
- Changed `frontend`'s `depends_on` on `nginx` from a plain list entry to
  `condition: service_healthy`.

Result: the dependency chain is now fully health-gated end to end —
`*-db` (healthy) → `*-migrate` (completed) → `*-service` (healthy) → `nginx` (healthy) →
`frontend`.

## 4. Verification

- `docker compose config -q` — passes; the compose file is syntactically valid with the new
  `healthcheck`/`depends_on` blocks (validated against a local `.env` copied from
  `.env.example`).
- **Not yet verified live in this sandbox:** `docker compose up -d --build` could not be run here
  — no Docker daemon is reachable in this environment
  (`/home/oscar-areiza/.docker/desktop/docker.sock` does not exist). The `wget`/`/health` pattern
  itself is not new — it is copied verbatim from `catalog-service`'s already-working healthcheck
  from Week 05 (verified then, see `CONTAINERIZATION.md` Section 4) — but the full chain must
  still be run once on a machine with Docker available before merging, to confirm timing
  (`start_period`/`retries`) is enough for `access-service`/`membership-service` to pass their
  first check.

## 5. Status and next step

- Change is made locally in the `lms-library` working tree, **not yet committed**. Suggested
  branch: `feat/HU-ORCH-01-healthcheck-gating` off `dev`, per
  `00-governance/git-conventions.md`.
- Before opening the PR: run `docker compose up -d --build` end to end on a machine with Docker,
  confirm every container reaches `healthy`, and attach that output as the PR's evidence (same
  format as `CONTAINERIZATION.md` Section 4).

## Correlations

- Prior containerization state → `05-week/hu-status/CONTAINERIZATION.md`
- Environments/config matrix this orchestration work feeds into → `MVP2-ENVIRONMENTS-PLAN.md`
  (this folder)
- Compose file → `docker-compose.yml` at the `lms-library` repo root
- Gateway routing → `infra/nginx/nginx.conf`
