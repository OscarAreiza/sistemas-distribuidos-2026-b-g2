<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       06-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 06

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Session 1 - bring the whole system up with a single `docker compose up`: shared network, health checks gating startup, config via env, data in volumes (harden the MVP 1 compose setup, closing the health-check gating still missing on `access-service`/`membership-service` and on the services that depend on them). Session 2 (planning) - define the three environments and the orchestration plan for MVP 2: a documented config matrix (variable names + per-env values), keep secrets out of git via `.env.example`, confirm the branch<->environment mapping, and slice the MVP 2 orchestration user stories with testable acceptance criteria.
<!-- CONFIG-END -->

> **This week's brief (as given):** Bring your whole system up with a single `docker compose up`:
> shared network, health checks gating startup, config via env, data in volumes. In Session 2
> (planning) you will define your environments and the orchestration plan for MVP 2. Define your
> three environments and a documented config matrix (variable names + per-env values), keep
> secrets out of git (with `.env.example`), confirm the branch<->environment mapping, and slice
> the orchestration stories for MVP 2 with testable acceptance criteria.

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-ORCH-01 | Single-command bring-up with startup gated on health checks | doing | Started this week - see individual contribution below |
| HU-ORCH-02 | Define the three environments + documented config matrix | doing | Started this week - see individual contribution below |
| HU-ORCH-03 | Branch<->environment mapping and secrets-out-of-git policy | todo | Planned for Session 2 this week |
| HU-ORCH-04 | Slice MVP 2 orchestration backlog with testable acceptance criteria | todo | Planned for Session 2 this week |

## 2. My individual contribution
- **Starting point (carried over from Week 05):** `lms-library`'s `docker-compose.yml` already
  has one shared bridge network (`lms-network`), env-based config (`.env` / `.env.example`), and
  one named volume per database (`access-db-data`, `membership-db-data`, `catalog-db-data`) from
  the Week 05 containerization work.
- **HU-ORCH-01, Session 1 (this week) - implemented:** found and closed the health-check gating
  gap. `access-service` and `membership-service` already exposed `/health`/`/health/ready` in
  code but were never wired into `docker-compose.yml`; `nginx` had no `healthcheck` of its own;
  and both `nginx` (on `catalog-service`/`access-service`/`membership-service`) and `frontend`
  (on `nginx`) used a plain `depends_on` list, which only waits for a container to *start*, not
  to become *healthy*. Fixed by adding the same `wget`-on-`/health` `healthcheck` pattern
  `catalog-service` already used to `access-service`, `membership-service`, and `nginx`, and by
  switching every `depends_on` in that chain to `condition: service_healthy`. Full detail →
  `ORCHESTRATION.md` (this folder). Validated with `docker compose config -q`; a live
  `docker compose up` end-to-end run is still pending (no Docker daemon reachable in this
  session's environment) and is the first thing to do before opening the PR.
- **HU-ORCH-02/03, Session 2 (planning) - first draft produced:** drafted the three environments
  (development/QA/production), the config matrix for every variable already in `.env.example`
  (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `JWT_SECRET`, `JWT_EXPIRY`, `LOG_LEVEL`, `CORS_ORIGIN`,
  `VITE_API_BASE_URL`) with per-environment values, the secrets-out-of-git policy, and the
  branch↔environment mapping (`dev`/`QA`/`main`, already used informally per
  `00-governance/git-conventions.md`, now written down explicitly). Full detail →
  `MVP2-ENVIRONMENTS-PLAN.md` (this folder).
- **HU-ORCH-04:** sliced HU-ORCH-01..04 with Gherkin-style acceptance criteria in
  `MVP2-ENVIRONMENTS-PLAN.md`, Section 5 - as placeholder IDs, since the orchestration track
  doesn't exist yet in `03-product/product-backlog.md`.

## 3. Blockers and risks
- **Live verification pending:** the health-check fix (HU-ORCH-01) was validated statically
  (`docker compose config -q`) but not run end to end - this sandbox has no reachable Docker
  daemon. Needs one `docker compose up -d --build` run on a machine with Docker before the PR
  can be opened with real evidence (same bar `CONTAINERIZATION.md` used in Week 05).
- The environments/config-matrix/branch-mapping documents are a first draft, not yet reviewed or
  ratified by the team - values for QA/production hosts are proposed, not confirmed.
- The MVP 2 orchestration HU IDs (`HU-ORCH-0N`) are placeholders - they need to be reconciled
  with the team and given real IDs in `03-product/product-backlog.md` (HU-ORCH-04's own
  acceptance criterion).

## 4. Plan for next week
- Run `docker compose up -d --build` end to end on a machine with Docker, confirm every service
  reaches `healthy`, and open the PR for HU-ORCH-01 (suggested branch:
  `feat/HU-ORCH-01-healthcheck-gating`) with that evidence attached.
- Take `MVP2-ENVIRONMENTS-PLAN.md` to the team for review; once ratified, register HU-ORCH-01..04
  with real IDs in `03-product/product-backlog.md` (closes HU-ORCH-04).
- Start acting on the config matrix: create the QA and production `.env` files (never committed)
  from the matrix's values.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration) - infra-only change this week, no application code
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- `docker-compose.yml` health-check/`depends_on` fix: applied locally in the `lms-library`
  working tree (based on `dev` @ `58c472d`), **not yet committed** - pending the live verification
  run noted in Section 3/4.
- Local documents attached in this same folder (`06-week/hu-status/`):
  - `ORCHESTRATION.md` - Session 1: the health-check gating gap, the fix, and its verification
    status.
  - `MVP2-ENVIRONMENTS-PLAN.md` - Session 2: the three environments, config matrix, secrets
    policy, branch↔environment mapping, and the sliced MVP 2 orchestration backlog.
