<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       05-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 05

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Containerize every service (multi-stage Dockerfile each, docker-compose.yml on one network with env-based config and volume-backed data) as the runtime base for the MVP 1 release; ship MVP 1 (promote to main, tag v1.0.0, verify the DoD, demo, retrospective).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01 | Administrator Authentication | done | https://github.com/OscarAreiza/lms-library/pull/8 (merged to `dev`) |
| HU-02 | Student Registration | done | https://github.com/OscarAreiza/lms-library/pull/9 (merged to `dev`) |
| HU-04 | Book Registration | done | https://github.com/OscarAreiza/lms-library/commit/d390674 (merged to `dev`) |
| HU-03, HU-05, HU-06, HU-07, HU-08, HU-09 | Student search/edit/deactivate, book search/edit, loans, returns, penalties | doing | Built and tested on their own `feat/HU-XX-*` branches; not yet merged to `dev` (see `04-week`/`05-week` blockers below) |

## 2. My individual contribution
- **Containerization (Session 1):** every current service (`access-service`, `membership-service`, `backend`) has its own multi-stage `Dockerfile` (`golang:1.25-alpine` builder → `alpine:3.21` runtime), and the frontend has its own (`node:20-alpine` builder → `nginx:alpine` runtime). `docker-compose.yml` at the repo root brings up all of them plus one PostgreSQL instance per service (`access-db`, `membership-db`, `db`) and an NGINX API gateway, all on one bridge network (`lms-network`), configured entirely through `.env` (see `.env.example`), with each database's data on its own named volume (`access-db-data`, `membership-db-data`, `lms-db-data`). See `library-docs/05-architecture/deployment.md` for the up-to-date infrastructure diagram.
- **Access and Membership extraction:** while updating the container topology, also finished extracting the Access (HU-01) and Membership (HU-02) bounded contexts into their own services (`access-service`, `membership-service`), each with its own database — see PR [#8](https://github.com/OscarAreiza/lms-library/pull/8) and [#9](https://github.com/OscarAreiza/lms-library/pull/9). Recorded the decision to decompose incrementally (one bounded context at a time, not a big-bang cutover) in `ADR-004-incremental-microservices-decomposition.md`.
- **Documentation coherence pass:** with the architecture now partially split, brought `library-docs` back in sync with the real system — updated `02-domain/domain-map.md`, `05-architecture/overview.md` and `deployment.md`, `06-data/models.md`, `09-microservices/service-catalog.md`, and the per-service docs under `09-microservices/services/`, and corrected the traceability matrix (`04-requirements/traceability-matrix.md`), which incorrectly stated no HU had been implemented yet.
- **Session 2 (Ship MVP 1):** not completed this week — see Blockers below. The Docker Compose stack (this week's Session 1 output) is the runtime base MVP 1 will ship on, but promoting to `main`/tagging `v1.0.0` is deliberately left for a manual step once the remaining Must-Have stories land.

## 3. Blockers and risks
- **MVP 1 scope is not fully merged yet.** The original Cut 1 (Must Have) backlog is HU-01, HU-02, HU-04, HU-06, HU-07 (`03-product/product-backlog.md`) — HU-06 (Loan Registration) and HU-07 (Return Registration) are the Circulation bounded context, which has **no implementation at all yet** (only an unused domain skeleton inside `backend`, see `ADR-004`). Tagging `v1.0.0` today would ship login + student management + book catalog, without the loan/return flow that is the system's Core Domain.
- **Resolved this week:** the team decided on option (b) — Circulation (HU-06/07/08) will **not** block the `v1.0.0` tag. It ships as a fast-follow **after** MVP 1, once the login + student management + book catalog stories already merged/mergeable (HU-01, HU-02, HU-03, HU-04, HU-05) are on `dev`. Promoting `dev -> QA -> main` and tagging `v1.0.0` proceeds on that reduced scope; Circulation's implementation starts only after the tag.
- No CI pipeline exists yet, so the DoD's "CI/CD green on the branch" and "deployed to staging" criteria (`00-governance/definition-of-done.md`) cannot be verified automatically — checked manually instead (see `MVP1-SHIPPING-CHECKLIST.md` in this folder).

## 4. Plan for next week
- Merge the remaining `feat/HU-XX` branches into `dev` (HU-03, HU-05) following the same incremental, one-bounded-context-at-a-time pattern already used for Access and Membership (`ADR-004`).
- Promote `dev -> QA -> main`, tag `v1.0.0` on the resolved MVP 1 scope (HU-01/02/03/04/05), run the demo, and hold the retrospective.
- Start Circulation (HU-06/07/08) only after the `v1.0.0` tag, per the decision recorded in Blockers and risks above.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Base Docker setup: https://github.com/OscarAreiza/lms-library/commit/e5940af
- Compose wiring (migrations, healthchecks, JWT env vars): https://github.com/OscarAreiza/lms-library/commit/74703f6
- Access extraction (own Dockerfile, own DB, gateway wiring): https://github.com/OscarAreiza/lms-library/pull/8
- Membership extraction (own Dockerfile, own DB, gateway wiring): https://github.com/OscarAreiza/lms-library/pull/9
- `library-docs` coherence pass: https://github.com/code-corhuila/library-docs/commit/94e8951, https://github.com/code-corhuila/library-docs/commit/a667bb9, https://github.com/code-corhuila/library-docs/commit/c6e2c6f
- Local copies attached in this same folder (`05-week/hu-status/`):
  - `CONTAINERIZATION.md` — Dockerfile-per-service, `.dockerignore` status, and the compose topology (Session 1)
  - `MVP1-SHIPPING-CHECKLIST.md` — DoD verification against the actual repo state, and what is still pending before the `v1.0.0` tag (Session 2)
- MVP snapshot (`mvp.zip`): uploaded manually by the team, not part of this automated pass
