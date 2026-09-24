# Session 2 (Planning) — MVP 2 Story Map, Estimation & Scope

> **Brief (as given):** Build a story map of the product, estimate the MVP 2 stories with
> planning poker, map and sequence cross-service dependencies (contract-first + mocks), and
> commit a realistic MVP 2 scope tied to velocity and the sprint goal.

## 1. Story map

Backbone = what the Administrator does end to end (`02-domain/domain-map.md`); each column below
is a walking skeleton slice, already fully working as one monolith (`v1.0.0`) and now being
re-sliced into independent services per `ADR-006`. Rows are the release: MVP 1 (shipped, single
repo) and MVP 2 (in progress, 18 repos).

| Backbone activity → | Authenticate | Manage Students | Manage Catalog | Manage Loans |
|---|---|---|---|---|
| **MVP 1** (`lms-library`, `v1.0.0`) | HU-01 ✅ | HU-02/03 ✅ | HU-04/05/09 ✅ | HU-06/07/08 — not shipped in v1.0.0 |
| **MVP 2** (repo-per-context) | `lms-access-api` 🟢, `lms-access-db` 🔴, `lms-access-portal` 🔴 | `lms-membership-api` 🟡(local), `lms-membership-db` 🔴, `lms-membership-portal` 🟡(local) | `lms-catalog-api` 🟡(local), `lms-catalog-db` 🔴, `lms-catalog-portal` 🟡(local) | `lms-circulation-api` 🟡(PR open), `lms-circulation-db` 🟡(PR open), `lms-circulation-portal` 🟡(local) |

Cross-cutting row (not tied to one backbone column — needed once at least two api+portal pairs
exist): `lms-infra` 🔴, `lms-api-gateway` 🔴, `lms-front` 🟡(local, shell only — portals not yet
composed as remotes), `lms-workflow` 🔴, `lms-worker` 🔴.

Legend: 🟢 merged to `develop` · 🟡 built, not yet merged (PR open or still local) · 🔴 not started.

## 2. Planning poker estimates — MVP 2 backlog (draft)

**Not yet run as a live team session** — no synchronous poker round happened this week. What
follows is a *proposed* estimate per repo, grounded in the actual complexity already observed
building the first four (not guessed from the ticket title), for the team to re-estimate for
real next sync. Scale: Fibonacci (1/2/3/5/8/13).

| Repo | Proposed points | Why |
|---|---|---|
| `lms-access-api` | 5 | Reference case — done: 5 PRs, ~24 files, straight port, no new dependency |
| `lms-access-db` | 2 | Compose service + 1 migration set — same shape as `lms-circulation-db`'s compose commit |
| `lms-access-portal` | 2 | One screen (`LoginPage`), no list/search state |
| `lms-membership-api` | 5 | Same shape as access-api, more usecases (6 vs 1) |
| `lms-membership-db` | 2 | Same shape as access-db |
| `lms-membership-portal` | 3 | Two screens with list/search/inline-edit state (`StudentsListPage`, `StudentFormPage`) |
| `lms-catalog-api` | 5 | Same shape, plus HU-05/HU-09 found and ported this session (search + edit, not just register) |
| `lms-catalog-db` | 2 | Same shape |
| `lms-catalog-portal` | 3 | Same shape as membership-portal |
| `lms-circulation-api` | **8** | Not a port — new persistence layer (Postgres→MongoDB rewrite), new dependency with no verifiable `go.sum`, 26 files |
| `lms-circulation-db` | 3 | New engine (Mongo, not Postgres) — one extra design decision (`ADR-005`), same compose shape otherwise |
| `lms-circulation-portal` | 5 | Three screens, one (`LoanFormPage`) built new rather than ported — no placeholder existed for it |
| `lms-infra` | 8 | Assembles all of the above via sibling-directory compose — highest integration risk, nothing to copy from |
| `lms-api-gateway` | 5 | Path-based routing already exists in `infra/nginx/nginx.conf`; scope (auth/rate-limiting) still needs its own ADR first |
| `lms-front` | 5 | Shell exists; real remote composition mechanism still undecided (no ADR) |
| `lms-workflow` | 8 | From scratch, scope not yet defined by its own ADR |
| `lms-worker` | 8 | From scratch, scope not yet defined by its own ADR |
| **Total (draft)** | **79** | vs. MVP 1's 33 — expected, this is infrastructure decomposition, not new user value |

## 3. Cross-service dependency sequencing (contract-first)

This was already applied in practice this week, not just planned. The real order followed:

```
lms-access-api  (no dependency — issues the JWT every other service validates)
      │
      ├──► lms-membership-api  (validates access's JWT locally; no call to access-api)
      │         ▲
      │         │  GET /students/{id}, POST /students/{id}/suspend
      │
      └──► lms-catalog-api     (validates access's JWT locally; no call to access-api)
                ▲
                │  POST /books/{id}/loan-copy, POST /books/{id}/return-copy
                │
         lms-circulation-api  ─┘  (built LAST, once both contracts above already existed)
```

**Contract-first, concretely:** `lms-catalog-api`'s `LoanCopy`/`ReturnCopy` endpoints and
`lms-membership-api`'s `GetStudent`/`SuspendStudent` endpoints were built and already present
*before* `lms-circulation-api`'s HTTP clients (`internal/infrastructure/catalog/client.go`,
`internal/infrastructure/membership/client.go`) were written against them — the clients were
written against real, already-defined response shapes, not mocks against a hoped-for contract.
No consumer-driven contract test exists yet in CI (same gap `07-week/hu-status` already flagged —
no CI pipeline exists at all); that remains a real risk, not a solved one.

**Sequencing implication for what's left:** `lms-infra` cannot be meaningfully built until at
least two `api`+`db` pairs are mergeable (already true — Access, Membership, Catalog qualify);
`lms-api-gateway`/`lms-front` need the portal repos to exist first; `lms-workflow`/`lms-worker`
need all three domain services reachable, which puts them last regardless of their own points.

## 4. Committed MVP 2 scope (tied to velocity)

Observed velocity this week: one context fully migrated end-to-end (Access, reference case) plus
three more contexts' application code fully written (Membership, Catalog, Circulation), in the
time it took to also design and rewrite Circulation's persistence layer for MongoDB. That's
roughly "one full context per working session" once a pattern exists to copy, and closer to two
sessions for the one context (Circulation) that wasn't a straight port.

**Committing to, for the next review:**
- Land the 2 currently-open PRs (`lms-circulation-api#2`, `lms-circulation-db#2`) — the
  `go.sum`/`go.mod` gap is the only known blocker on the first.
- Open PRs for the work that's already done but invisible on the board: `lms-membership-api`,
  `lms-catalog-api` (incl. HU-05/HU-09), and the three portal repos.
- **Not** committing to `lms-infra`, `lms-api-gateway`, `lms-front`'s real composition,
  `lms-workflow`, or `lms-worker` this cycle — each needs its own scope-defining ADR first
  (`lms-api-gateway`, `lms-workflow`, `lms-worker` per their own READMEs; `lms-front`'s remote
  composition mechanism has no ADR at all yet), and per the draft estimates above they're the
  most expensive items in the backlog. Starting them without that decision made is the same
  "built on sand" risk `07-week/hu-status` already called out for wiring infrastructure around
  an unmerged service.

This keeps MVP 2's scope tied to what velocity actually supports, not to what's left on the map.
