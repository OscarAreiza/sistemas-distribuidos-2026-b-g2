<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       10-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 10

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Session 1 - give each service its own database, implement one saga (with a compensation) and the outbox for one critical event, and make the consumers idempotent. Session 2 is the MVP 2 release.
<!-- CONFIG-END -->

> **This week's brief (as given):**
> 1. Give each service its own database; implement one saga (with a compensation) and the
>    outbox for one critical event; make the consumers idempotent. Session 2 is the MVP 2
>    release.
> 2. Ship MVP 2: promote to main, tag v2.0.0, verify the checklist (including a
>    failure/compensation path), demo the integrated system with an injected failure, and
>    retrospect. Ensure each member's hu-status evidence is complete. On to Corte 3.

**Just opened — status below reflects day 1, not a finished week.** A meaningful share of item 1
was already underway before this week officially started (carried over from week 9's QA
promotion work), noted honestly as "doing," not claimed as "done."

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-PERSIST-01 | Give each service its own database (Anexo J: shared instance, one schema/database per domain, least-privilege login) | doing | `lms-circulation-db#4` MERGED; `lms-membership-db#5` and `lms-catalog-db#4` OPEN, unreviewed |
| HU-PERSIST-02 | One saga with a compensation path (overdue-loan penalty, per `ADR-008`) | todo | `lms-workflow` repo exists, scope documented, **zero code written** |
| HU-PERSIST-03 | Outbox for one critical event | todo — **architecture note, not just unstarted** | This project has no message broker (`AT-002`, reaffirmed by `ADR-009`) — a classic outbox (durably write, then publish to a bus) doesn't apply. `ADR-008`'s `workflow_db.saga_runs` table is this project's actual durability mechanism for the one critical event in scope (loan overdue → penalty): it durably records which step completed before/while making the direct HTTP calls, which is the same crash-safety guarantee an outbox exists to provide, without a broker to publish to. Building `saga_runs` *is* this story, not a separate one |
| HU-PERSIST-04 | Idempotent consumers | doing | HTTP-layer idempotency already real and tested: `circulation-api`'s `IdempotencyStore` (loan registration), `catalog-api`'s (book registration) — both Mongo/Postgres-backed, both verified live this period. **Not the same thing as idempotent message consumers** (no broker, so no message consumers exist) — `ADR-008`'s `loan_id`-keyed idempotency check for the penalty saga is the piece still to build |
| HU-MVP2-RELEASE | Promote to `main`, tag `v2.0.0`, verify checklist incl. failure/compensation path, demo with injected failure, retrospective | todo | Blocked behind HU-PERSIST-01..04 and the `qa` promotions still open across Circulation/Catalog/Membership (see week 9's carried-over PRs) |

## 2. My individual contribution
- Confirmed the real state of "each service its own database" across all three affected repos
  rather than assuming week 9's Circulation work generalized: `lms-circulation-db#4` is merged,
  but the equivalent PRs for `lms-membership-db` (`#5`) and `lms-catalog-db` (`#4`) are still
  open and unreviewed — same Anexo J pattern (shared instance, per-domain schema, least-privilege
  login created by infra + role granted by the domain's own migration), not yet uniform across
  the fleet.
- Read `ADR-008` closely before writing this file's HU-PERSIST-02/03 rows, specifically to avoid
  claiming "outbox" as a literal to-do when this project's own architecture decisions (`AT-002`,
  `ADR-009`) already rule out a message broker. Reframed it honestly: the saga's `saga_runs`
  state table is this project's durability mechanism for the one critical event in scope, and
  building it is how HU-PERSIST-03 actually gets satisfied — not a gap to apologize for.
- Confirmed `lms-workflow` is still exactly what its README says: scope only, zero code — so
  HU-PERSIST-02 is accurately `todo`, not understated or overstated.

## 3. Blockers and risks
- `lms-workflow` doesn't exist as code yet — HU-PERSIST-02/03 need the service scaffolded
  (hexagonal layout, `workflow_db` Postgres migration for `saga_runs`, HTTP clients to
  `circulation-api`/`membership-api`) before the saga itself can be written.
- `lms-membership-db#5` and `lms-catalog-db#4` (Anexo J) are still open — "each service its own
  database" isn't actually uniform across the fleet until these merge.
- The MVP 2 release checklist explicitly requires a demonstrated failure/compensation path — that
  can't be rehearsed until the saga exists, so this is the critical path for Session 2, not a
  parallel track.
- Week 9's `qa` promotion PRs (`lms-circulation-{db,api,portal}`, `lms-catalog-{db,portal}`) are
  still open and unreviewed — `main` can't be touched responsibly while `qa` itself isn't settled.

## 4. Plan for next week
- Scaffold `lms-workflow` (hexagonal layout matching the other `-api` repos) and the
  `workflow_db` migration for `saga_runs`, per `ADR-008`'s exact schema.
- Implement the overdue-detection scheduler and the two-step penalty call (mark late →
  suspend student), with the compensation/retry and `loan_id` idempotency check `ADR-008`
  specifies.
- Get `lms-membership-db#5` and `lms-catalog-db#4` reviewed and merged so "each service its own
  database" is actually fleet-wide, not 1-of-3.
- Only once the above is real: start the MVP 2 release checklist (promote to `main`, tag
  `v2.0.0`, injected-failure demo, retrospective).

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary` - not evaluable yet, no commits this week
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...) - same
- [x] Testable acceptance criteria - `ADR-008` already specifies the compensation path and
  idempotency key in testable terms; nothing new needed to restate them
- [ ] Tests added/updated (unit / integration) - no code written yet
- [ ] DDD / hexagonal boundaries respected (domain has no I/O) - not evaluable, `lms-workflow`
  has no code yet
- [x] No secrets; config via environment variables - carried forward from week 9, no regression

## 6. Evidence links
- `library-docs/05-architecture/decisions/records/ADR-008-circulation-saga-scope.md` - the saga's
  full scope, compensation path, idempotency key, and state schema
- `library-docs/05-architecture/decisions/records/ADR-009-worker-scheduling-model.md` - why
  there's no broker, and how the overdue scheduler runs instead
- Open Anexo J PRs: https://github.com/code-corhuila/lms-membership-db/pull/5,
  https://github.com/code-corhuila/lms-catalog-db/pull/4
- Merged Anexo J PR: https://github.com/code-corhuila/lms-circulation-db/pull/4
- `lms-workflow` repo (scope only, no code yet): https://github.com/code-corhuila/lms-workflow
