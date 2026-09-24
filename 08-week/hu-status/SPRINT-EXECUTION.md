# Session 1 — Run the Sprint Like a Pro

> **Brief (as given):** A prioritized backlog with testable stories, a WIP limit, PRs for every
> change, and a daily sync. Track throughput.

## 1. Prioritized backlog

Two backlogs are in play, at two different levels:

- **MVP 1 backlog** (`03-product/product-backlog.md`) — 9 user stories, 33 story points, already
  fully pointed and prioritized (Must Have: HU-01/02/04/06/07; Should Have: HU-03/05/08/09).
  Shipped as `v1.0.0` (Week 05).
- **MVP 2 backlog** — the repo-per-context decomposition tracked in
  `09-microservices/repo-migration-map.md`: 18 repositories (4 bounded contexts × api/db/portal,
  plus `lms-api-gateway`, `lms-front`, `lms-workflow`, `lms-worker`, `lms-infra`). This backlog
  has **no story points yet** — that gap is exactly what Session 2 (planning) closes this week;
  see `MVP2-STORY-MAP-AND-SCOPE.md`.

The migration order in `repo-migration-map.md` already gives MVP 2 its priority ranking (Access →
Membership/Catalog → `lms-infra` → gateway/front → Circulation → workflow/worker) — the same
"no dependency blocks the next slice" reasoning `ADR-004`/`ADR-006` used for the original
incremental split.

## 2. WIP limit

**Policy set this week:** at most 2 open PRs per person at a time, so review attention isn't
split across more in-flight branches than can realistically be reviewed.

**Actual compliance — mixed.** Measured against the *visible* board (open PRs), the limit holds:
only `lms-circulation-api#2` and `lms-circulation-db#2` are open right now. But that undercounts
the real WIP: `lms-membership-api` (domain/usecase/infra), `lms-catalog-api` (domain/usecase/infra
plus HU-05/HU-09), all three portal repos (`lms-membership-portal`, `lms-catalog-portal`,
`lms-circulation-portal`), and `lms-front` are **fully built but have zero open PRs** — that work
exists only in local working trees, invisible to anyone looking at the repos' PR lists. A WIP
limit that only counts open PRs doesn't actually limit this kind of work; the real fix for next
week is opening PRs for that work (even as drafts) so it counts against the limit and gets
reviewed, instead of accumulating unseen.

## 3. PRs for every change

Every change this week went through a branch → PR, never a direct commit to `develop`
(`00-governance/branching-policy.md`). Two granularities were used on purpose:

- **Coarse, layer-per-PR** — `lms-access-api`'s PRs #2–#5 (entry point/tooling, domain,
  application, infrastructure as four separate PRs against `develop`), each independently
  buildable-or-explicitly-marked-not-yet.
- **Fine, file-per-commit inside one PR** — `lms-circulation-api#2` (26 commits, one file each:
  `go.mod`, then domain, then application, then config/infrastructure, in dependency order) and
  `lms-circulation-db#2` (2 commits: compose service, then README). Every commit carries
  `Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>`.

## 4. Daily sync

**Not held this week as a real synchronous ceremony** — this session's work was done solo,
async, against the team's shared repos. Flagging this honestly rather than reporting a sync that
didn't happen. Proposed format for next week: a 10-minute async standup (what merged, what's
blocked, what's next) posted in the team channel before each session, since a live daily sync
isn't realistic for three people on different schedules.

## 5. Throughput (this period)

| Repo | PRs opened | PRs merged | PRs open | Notes |
|---|---|---|---|---|
| `lms-access-api` | 5 | 5 | 0 | Fully migrated; #2–#5 merged same-day, fast cycle time |
| `lms-membership-api` | 2 | 1 | 0 | #2 (entry point) opened, then closed unmerged — see Blockers |
| `lms-catalog-api` | 1 | 1 | 0 | Only the docs PR is remote; full implementation is local, no PR yet |
| `lms-circulation-api` | 2 | 1 | 1 | #2 open since 2026-09-22, still open — blocked on `go.sum`/`go.mod` (see Blockers) |
| `lms-circulation-db` | 2 | 1 | 1 | #2 open since 2026-09-23, still open |
| **Total** | **12** | **9** | **2** | |

**Reading it:** cycle time on `lms-access-api` (open → merged same day, sometimes within
minutes) is the healthy baseline — that's what "reviewed promptly" looks like. The two currently
open PRs (`circulation-api#2`, `circulation-db#2`) are the actual bottleneck to watch next week,
not new WIP to start.

## Blockers and risks (Session 1 specific)

- `lms-circulation-api#2` cannot merge as-is: `go.sum` has no hash for `go.mongodb.org/mongo-driver`
  (no Go toolchain/network access in this environment to run `go mod tidy`) — this is the actual
  throughput bottleneck for next week, not a new story.
- `lms-membership-api#2` was opened, then closed unmerged (an unauthorized push corrected mid-week)
  — not counted as delivered throughput; the same work now exists correctly-scoped locally,
  not yet re-opened as a PR.
- A large share of finished MVP 2 work (see WIP section) has no PR at all — the backlog/board
  understates how much is actually done.
