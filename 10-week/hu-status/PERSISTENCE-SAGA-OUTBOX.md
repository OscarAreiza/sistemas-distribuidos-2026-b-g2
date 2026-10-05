# Session 1 — Persistence, Saga, Outbox, Idempotent Consumers

> Brief: "Give each service its own database; implement one saga (with a compensation) and the
> outbox for one critical event; make the consumers idempotent."

Written day 1 of this week — a plan grounded in real current state, not a report of finished
work. Each section says exactly what exists today vs. what's still to build.

---

## 1. Each service its own database (Anexo J)

Real status per repo, checked directly against GitHub (not assumed):

| Domain | PR | State | What it does |
|---|---|---|---|
| Circulation | `lms-circulation-db#4` | **MERGED** | Migrates to the shared `lms-infra-mongo` instance: own database (`loan_db`) inside it, own login (`circulation_app`, created by infra with no roles), own role grant (`circulation_writer`) done by this repo's own migration |
| Membership | `lms-membership-db#5` | OPEN, unreviewed | Same pattern, shared `lms-infra` Postgres instance, own schema |
| Catalog | `lms-catalog-db#4` | OPEN, unreviewed | Same pattern, own schema |

**What "own database" means here, precisely** (Anexo J, `rules/3-Anexo-J`): not a separate
container/instance per service — one shared Postgres instance and one shared MongoDB instance for
the whole project, each environment (`dev`/`qa`/`main`) its own instance. Within that shared
instance, each domain gets its own schema (Postgres) or its own database (Mongo), its own
suffixed migration-tracking table, and a dedicated least-privilege login user that only that
domain's own migration grants a role to. The infra repo creates the login; the domain's `-db`
repo grants the role — never the reverse, and never one repo touching another's schema.

**To close this story:** get `lms-membership-db#5` and `lms-catalog-db#4` reviewed and merged.
Nothing new to build — the pattern is proven (Circulation already works), this is a review
bottleneck, not an implementation one.

## 2. One saga with a compensation path

Per `ADR-008-circulation-saga-scope.md` (already accepted, reviewed by Hermes and Luis) — the one
event in this project's entire domain that actually has the multi-step, time-spanning,
compensable shape a saga exists for: **a loan becomes overdue → apply the penalty.**

**Everything else stays exactly where it already is** (ADR-008's own words) — eligibility checks,
availability decrement/restore on loan registration — those are single-request/response calls,
already correct, not saga candidates. Routing them through a saga would add a failure mode for no
benefit. This is why only one saga is in scope, not "a saga per event."

**What it does, concretely:**
1. **Overdue detection** — runs on its own schedule (not request-triggered), scans
   `lms-circulation-api` for active loans past `dueDate`.
2. **Penalty application** — two writes, two services: mark the loan late in `circulation-api`,
   suspend the student in `membership-api`. If the second call fails after the first succeeds,
   retry — don't let it drift silently.

**Compensation path** (ADR-008 §"Compensation path"): if marking the loan late succeeds but the
suspension call fails, the saga retries the suspension call on its next scheduled pass — it does
**not** undo the "marked late" step, because that step is already factually true (the loan really
is late) and undoing it would be incorrect, not a safety net. Compensation here means "retry the
failed half," not "roll back the succeeded half" — worth stating explicitly since "compensation"
defaults to meaning rollback in most saga literature, and that default is wrong for this specific
event.

**Nothing of this exists as code yet.** `lms-workflow`'s repo is still exactly its scope README,
zero implementation.

**To close this story:**
1. Scaffold `lms-workflow` with the same hexagonal layout every other `-api` repo uses.
2. Add the `workflow_db` Postgres migration for `saga_runs` (schema below, §4).
3. Implement the overdue-detection scheduler (see `ADR-009` for the scheduling model — polling,
   not a broker).
4. Implement the two-step penalty call with the retry-on-failure compensation path above.

## 3. The outbox for one critical event

**Architecture note before anything else:** this project has no message broker.
`library-docs/11-quality/testing-strategy.md`'s own Tier 2 section says so directly ("There is
also no message broker to integration-test against"), and `ADR-009-worker-scheduling-model.md`
reaffirms it as a deliberate choice (`AT-002`), not an oversight. A classic outbox pattern exists
to solve one specific problem: durably recording "this needs to be published" in the same
transaction as the business write, so a crash between the write and the publish can't silently
lose the event. Without a broker to publish to, there's no publish step to protect.

**What actually serves the same purpose here:** `ADR-008`'s `saga_runs` table (see §4) **is** this
project's outbox-equivalent for the one critical event in scope. It durably records, in the same
database the saga writes to, which step of the penalty flow has completed — `PENDING` →
`MARKED_LATE` → `COMPLETED` or `FAILED` — before/while making the direct HTTP calls that are this
project's substitute for publishing to a bus. A crash mid-flow leaves a `saga_runs` row the next
scheduled pass can pick up and resume, exactly the crash-safety guarantee a classic outbox exists
to provide.

**This means HU-PERSIST-03 isn't a separate deliverable from HU-PERSIST-02 — building
`saga_runs` correctly satisfies both at once.** Treating "outbox" as a literal checkbox needing
its own broker-publishing code would be building something this project's own architecture
decisions already ruled out.

## 4. Saga state schema (shared by §2 and §3)

Exactly as `ADR-008` specifies — not redesigned here, just restated as the concrete thing to
build:

```sql
-- workflow_db, PostgreSQL (per stack.md's default engine — no document-shape
-- argument applies to a flat, read/written-by-status row like this one)
CREATE TABLE saga_runs (
  loan_id        UUID PRIMARY KEY,
  status         TEXT NOT NULL CHECK (status IN ('PENDING','MARKED_LATE','COMPLETED','FAILED')),
  attempt_count  INT NOT NULL DEFAULT 0,
  last_error     TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**Idempotency key:** `loan_id`, scoped to one overdue cycle (`ADR-008` §"Idempotency"). Before
acting on a loan, the scheduler checks whether a `saga_runs` row already exists for that
`loan_id` in a non-`FAILED` state — if so, skip it. This is what prevents a retried pass from
re-suspending an already-suspended student.

## 5. Idempotent consumers

**Already real, already tested** — HTTP-layer idempotency, not message-consumer idempotency
(there are no message consumers in this project, per §3's no-broker note):

| Service | Mechanism | Verified |
|---|---|---|
| `lms-circulation-api` | `IdempotencyStore`, Mongo-backed, keyed on `Idempotency-Key` header | Yes — retried loan registration confirmed not to double-decrement available copies |
| `lms-catalog-api` | `IdempotencyStore`, Postgres-backed, same pattern | Yes — same mechanism, book registration |

**What's still missing:** the `loan_id`-keyed idempotency check for the penalty saga itself (§4)
— a different idempotency key for a different operation, not a gap in the existing HTTP-layer
mechanism. Building §4's check closes this.

## Correlations

- Full saga scope, evaluated alternatives, and risk register → `library-docs/05-architecture/decisions/records/ADR-008-circulation-saga-scope.md`
- Why no broker, and the scheduling model instead → `library-docs/05-architecture/decisions/records/ADR-009-worker-scheduling-model.md`
- MVP 2 release checklist this unlocks → `MVP2-RELEASE-CHECKLIST.md`, this same folder
