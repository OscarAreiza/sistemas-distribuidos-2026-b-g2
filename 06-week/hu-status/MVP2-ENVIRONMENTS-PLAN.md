# Week 6 - Session 2: Environments and Orchestration Plan for MVP 2

> Goal (planning): define the three environments, a documented config matrix (variable names +
> per-env values), keep secrets out of git (`.env.example`), confirm the branch↔environment
> mapping, and slice the MVP 2 orchestration user stories with testable acceptance criteria.
>
> **Status: draft.** This is this week's planning deliverable, produced to be reviewed and
> ratified with the team — environment names, exact QA/prod hosts, and the final HU IDs
> (currently placeholders `HU-ORCH-0N`, since the orchestration track doesn't exist yet in
> `03-product/product-backlog.md`) still need team sign-off before they're treated as final.

---

## 1. The three environments

| Environment | Purpose | Runs where (proposed) | Branch |
|---|---|---|---|
| **Development** | Local iteration, one dev at a time | Each developer's machine, `docker compose up -d --build` | `dev` |
| **QA / Staging** | Integration validation before release, what the team demos against | A shared VM/host running the same `docker-compose.yml`, pulled images or built from `QA` | `QA` |
| **Production** | What real users hit | A dedicated VM/host (or later, a managed platform) running the tagged release | `main` |

This mirrors the branch strategy already in force per `00-governance/git-conventions.md`
(`main -> QA -> dev -> feat/HU-XX-...`) — Section 4 makes the mapping explicit rather than
implicit.

## 2. Config matrix

Base variable set taken from the current `.env.example` (`lms-library` repo root) — no new
variables were introduced by this week's orchestration fix (`ORCHESTRATION.md`), only new
`healthcheck` blocks that read no env vars of their own.

| Variable | Used by | Development (`.env`, from `.env.example`) | QA | Production |
|---|---|---|---|---|
| `POSTGRES_USER` | `access-db`, `membership-db`, `catalog-db` | `lms_user` | `lms_user` (or QA-specific) | Distinct value, never reused from dev/QA |
| `POSTGRES_PASSWORD` | same | `lms_password` (placeholder, fine for local only) | Generated secret, injected via the QA host's secret store, not committed | Generated secret, injected via the production host's secret store, not committed |
| `JWT_SECRET` | `access-service`, `membership-service`, `catalog-service` (validates tokens `access-service` issues) | `openssl rand -base64 32`, kept in the developer's local `.env` | Its own generated secret, distinct from dev and prod | Its own generated secret, distinct from dev and QA — rotated on a schedule (policy TBD with the team) |
| `JWT_EXPIRY` | same three | `1h` | `1h` | Team to confirm — may want a shorter value in production |
| `LOG_LEVEL` | all Go services | `debug` | `info` | `warn` |
| `CORS_ORIGIN` | all Go services | `http://localhost:3000` | QA frontend origin (e.g. `https://qa.<domain>`) | Production frontend origin (e.g. `https://<domain>`) |
| `VITE_API_BASE_URL` | `frontend` build arg | `http://localhost:8080/api/v1` | QA gateway URL | Production gateway URL |

**Rule going forward (ties into HU-ORCH-02 below):** any new environment variable added to a
service must be added to this table with its dev/QA/prod value in the same PR — this is the
acceptance criterion, not just a convention.

## 3. Keeping secrets out of git

- `.env.example` (already exists at the repo root) stays the only environment file tracked in
  git — placeholder/dev-safe values only (already true today: `lms_user`/`lms_password`/a
  clearly-labeled placeholder `JWT_SECRET`).
- `.env` and `.env.*` are already `.gitignore`d (confirmed in `lms-library/.gitignore`,
  `!.env.example` is the sole exception) — no change needed there, just confirming the policy
  holds as QA/prod env files are created.
- QA and production real secret values (`POSTGRES_PASSWORD`, `JWT_SECRET`) live only on their
  respective hosts (or a secrets manager, if the team adopts one later) — never in a file
  committed to `lms-library`, and never pasted into a PR description or issue.

## 4. Branch ↔ environment mapping

| Branch | Environment | How it gets there |
|---|---|---|
| `feat/HU-XX-...` | Developer's own `development` instance | `docker compose up -d --build` locally against that branch |
| `dev` | Development (shared reference) | Merge via PR once a feature branch's checks pass |
| `QA` | QA / Staging | Promote `dev -> QA` once a batch of HUs is ready for integration testing (same practice used to ship `v1.0.0`, `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md`) |
| `main` | Production | Promote `QA -> main` and tag a release, once QA sign-off is done |

This is the same promotion flow already used for `v1.0.0` — this section just writes it down
against the three environments explicitly, closing the "not documented" gap noted in this week's
`README.md` (Section 3, Blockers and risks).

## 5. MVP 2 orchestration backlog — sliced, with acceptance criteria

| ID | Title | Acceptance criteria (Given/When/Then) | Status |
|---|---|---|---|
| **HU-ORCH-01** | Single-command bring-up gated on health, not just start order | **Given** a clean checkout with a valid `.env`, **when** `docker compose up -d --build` runs, **then** every service with a `healthcheck` reaches `healthy` and no dependent service (`nginx`, `frontend`) starts serving before its dependencies report healthy. **Given** `access-service` or `membership-service` is unhealthy, **when** `nginx` starts, **then** `nginx` does not report `healthy` either (chain, not just first hop). | doing — implemented in `docker-compose.yml` this session, live end-to-end run still pending (`ORCHESTRATION.md` Section 4) |
| **HU-ORCH-02** | Documented config matrix per environment | **Given** the config matrix (Section 2), **when** a new env var is added to any service, **then** it is added to the matrix with dev/QA/prod values in the same PR. **Given** `.env.example`, **when** compared against the matrix, **then** every matrix variable exists in `.env.example` with the development value as its default. | doing — first draft is Section 2 above, pending team review |
| **HU-ORCH-03** | Branch↔environment mapping and no-secrets-in-git policy documented | **Given** a merge to `dev`/`QA`/`main`, **when** someone asks "what environment does this deploy to", **then** Section 4's table answers it unambiguously. **Given** any commit, **when** it's reviewed, **then** it contains no real secret value (only `.env.example` placeholders) — enforced today by review + `.gitignore`, automated scanning is a stretch goal, not a blocker. | doing — Sections 3–4 above are the first draft |
| **HU-ORCH-04** | MVP 2 orchestration backlog sliced and ID'd in the product backlog | **Given** this document, **when** the team reviews it, **then** HU-ORCH-01..03 (or their renamed equivalents) are added to `03-product/product-backlog.md` with real IDs, replacing the `HU-ORCH-0N` placeholders used here. | todo — needs a team planning session, not something to finalize unilaterally |

## Correlations

- Implementation of HU-ORCH-01 → `ORCHESTRATION.md` (this folder)
- Prior branch strategy → `00-governance/git-conventions.md`
- Prior environment gap ("no staging environment exists") → `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md`, Section 2
- Current env var source of truth → `.env.example` at the `lms-library` repo root
