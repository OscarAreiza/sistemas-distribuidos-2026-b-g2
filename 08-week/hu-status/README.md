<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       08-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 08

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Session 1 - run the sprint like a pro: a prioritized backlog, a WIP limit, a PR for every change, a daily sync, and throughput tracked. Session 2 (planning) - build a story map of the product, estimate the MVP 2 backlog (the repo-per-context decomposition) with planning poker, map and sequence cross-service dependencies contract-first, and commit a realistic MVP 2 scope tied to velocity.
<!-- CONFIG-END -->

> **This week's brief (as given):** Run your sprint like a pro: a prioritized backlog with
> testable stories, a WIP limit, PRs for every change, and a daily sync. Track throughput.
> Session 2 (planning) refines and estimates the MVP 2 backlog: build a story map of your
> product, estimate the MVP 2 stories with planning poker, map and sequence cross-service
> dependencies (contract-first + mocks), and commit a realistic MVP 2 scope tied to your
> velocity and the sprint goal.

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-SPRINT-01 | Prioritized backlog + WIP visibility for the MVP 2 migration | doing | `SPRINT-EXECUTION.md` §1-2 |
| HU-SPRINT-02 | PR-per-change discipline, throughput tracked | done | 12 PRs opened / 9 merged this period — `SPRINT-EXECUTION.md` §3, §5 |
| HU-SPRINT-03 | Daily sync | todo | Not held this week (async solo session) — `SPRINT-EXECUTION.md` §4 |
| HU-MVP2-01 | Story map of the MVP 2 backlog (18-repo decomposition) | done | `MVP2-STORY-MAP-AND-SCOPE.md` §1 |
| HU-MVP2-02 | Planning-poker estimates for the MVP 2 backlog | doing | Draft proposed, pending live team round — `MVP2-STORY-MAP-AND-SCOPE.md` §2 |
| HU-MVP2-03 | Cross-service dependency sequencing (contract-first) | done | Already applied building Circulation against Membership/Catalog's contracts — `MVP2-STORY-MAP-AND-SCOPE.md` §3 |
| HU-MVP2-04 | Committed MVP 2 scope tied to velocity | done | `MVP2-STORY-MAP-AND-SCOPE.md` §4 |

## 2. My individual contribution
- **Session 1 - sprint execution:** defined a WIP limit (2 open PRs/person) and measured real
  compliance against it — the visible board (2 open PRs: `lms-circulation-api#2`,
  `lms-circulation-db#2`) looks compliant, but a large share of finished MVP 2 work
  (`lms-membership-api`, `lms-catalog-api` incl. HU-05/HU-09, all three portal repos) has no
  PR open at all, so the limit isn't really being tested yet. Tracked real throughput across the
  5 repos touched this period: 12 PRs opened, 9 merged, 2 still open. No real daily sync was
  held - flagged honestly rather than reported as done. Full detail -> `SPRINT-EXECUTION.md`.
- **Session 2 - MVP 2 planning:** built a story map crossing the product's backbone
  (Authenticate / Manage Students / Manage Catalog / Manage Loans) against the MVP 1 (shipped)
  and MVP 2 (in-progress, 18 repos) releases. Proposed Fibonacci estimates for all 18 MVP 2
  repos (79 points, draft - not yet run as a live planning-poker round), grounded in the actual
  complexity difference observed between a straight port (Access, 5 pts) and new work
  (Circulation's Postgres->MongoDB rewrite, 8 pts). Documented the cross-service dependency
  sequencing actually used this week - `lms-catalog-api`'s and `lms-membership-api`'s endpoints
  existed before `lms-circulation-api`'s HTTP clients were written against them, i.e.
  contract-first in practice, not just in principle. Committed a scope for the next review tied
  to that velocity: land the 2 open PRs and open PRs for the invisible-WIP work; explicitly not
  committing to `lms-infra`/`lms-api-gateway`/`lms-front`/`lms-workflow`/`lms-worker` this cycle,
  since each needs its own scope-defining ADR first. Full detail -> `MVP2-STORY-MAP-AND-SCOPE.md`.

## 3. Blockers and risks
- `lms-circulation-api#2` can't merge as-is: `go.sum` has no verifiable hash for
  `go.mongodb.org/mongo-driver` (no Go toolchain/network access in this environment to run
  `go mod tidy`) - real throughput bottleneck for next week.
- The WIP limit only covers what has a PR open; a large amount of finished work has none yet
  (see Session 1 above) - the board currently understates real progress and real WIP both.
- No CI pipeline exists yet (same gap `07-week/hu-status` already flagged) - can't automate any
  of this tracking.
- MVP 2 story points are a draft, not team-ratified - same caveat prior weeks' first-draft
  planning docs (`06-week`, `07-week`) already carried forward.

## 4. Plan for next week
- Fix `lms-circulation-api#2`'s `go.sum` (run `go mod tidy` with a real Go toolchain) and merge
  both open PRs.
- Open PRs for `lms-membership-api`, `lms-catalog-api`, and the three portal repos so they count
  against the WIP limit and get reviewed instead of sitting local.
- Take the planning-poker draft in `MVP2-STORY-MAP-AND-SCOPE.md` to the team for a real
  estimation round; register `HU-SPRINT-0N`/`HU-MVP2-0N` with real IDs once ratified.
- Hold the first real daily sync, even if async/written, before starting `lms-infra`.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration) - planning/process session, no application code changed
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Local documents attached in this same folder (`08-week/hu-status/`):
  - `SPRINT-EXECUTION.md` - Session 1: backlog, WIP limit, PR discipline, daily sync, throughput
  - `MVP2-STORY-MAP-AND-SCOPE.md` - Session 2: story map, planning-poker draft, dependency
    sequencing, committed MVP 2 scope
  - `ADR-010-liquibase-for-database-migrations.md` - local copy of the ADR (`library-docs`
    original: `05-architecture/decisions/records/ADR-010-...md`) committed this same week,
    deciding every `-db` repo's migration tool (Liquibase, not golang-migrate) ahead of the
    MVP 2 repo-per-context split this sprint's planning covers
- Open PRs: https://github.com/code-corhuila/lms-circulation-api/pull/2,
  https://github.com/code-corhuila/lms-circulation-db/pull/2
- Merged PRs this period: https://github.com/code-corhuila/lms-access-api/pull/2,
  .../pull/3, .../pull/4, .../pull/5
- Backlogs referenced: `library-docs/03-product/product-backlog.md` (MVP 1),
  `library-docs/09-microservices/repo-migration-map.md` (MVP 2)
