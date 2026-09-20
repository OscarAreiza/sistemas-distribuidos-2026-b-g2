# ADR-008 — Scope of the `lms-workflow` Saga for Circulation

- **ID:** ADR-008
- **Date:** 2026-09-15
- **Status:** Accepted
- **Authors:** Oscar Areiza — Tech Lead
- **Reviewers:** Hermes Pascuas, Luis Alejandro Meneses — Development team

---

## Context

`lms-workflow`'s README describes it as orchestrating "loan → renewal → return → penalty (saga)"
across `lms-circulation-api`, `lms-membership-api`, and `lms-catalog-api`. Read as "every
cross-service interaction in the loan lifecycle goes through the saga," this would replace code
that already exists and already works in `lms-library` today:

- `membership-service`'s `ActiveLoansChecker` port + HTTP client already calls
  `circulation-service` directly (`GET /api/v1/loans?studentId=...&status=ACTIVE`) to check active
  loans before a student can be deactivated
- `catalog-service` already exposes `POST /books/{id}/loan-copy` and `POST /books/{id}/return-copy`
  specifically "needed by circulation-service (HU-06/HU-07)," per its own code comments

Both are single-request/response calls: ask a question, get an answer, act on it in the same
call. Neither has a multi-step, time-spanning, or compensatable shape — the two properties that
justify a saga/orchestrator in the first place. Routing them through `lms-workflow` would wrap an
already-correct synchronous call in orchestration machinery for no benefit, adding a new failure
mode (the orchestrator itself) to something that doesn't need one.

Two parts of the loan lifecycle genuinely do have that shape and have **no existing code at
all**: detecting a loan has become overdue (which must run on its own schedule, not in response to
a request) and applying the resulting penalty across `circulation-api` (mark the loan late) and
`membership-api` (suspend the student) — two writes to two different services for one logical
event, exactly the kind of multi-step, partially-recoverable operation a saga exists for.

**Known constraints:**
- Same 3-person team; no message broker adopted (`AT-002`, reaffirmed by `ADR-009`) — the saga is
  triggered by direct HTTP calls / a schedule, not by consuming events off a bus
- `ADR-005`'s MongoDB choice for Circulation and every existing invariant on `Loan`
  (`02-domain/entities-and-rules.md`, INV-001–INV-005) are unaffected by this ADR — it only
  decides who calls whom, not what the rules are

---

## Decision

**We decided:** `lms-workflow` owns only the **time-spanning, compensable** parts of the loan
lifecycle:
1. **Overdue detection** — runs on its own schedule (not triggered by a user request), scanning
   `lms-circulation-api` for active loans past `dueDate`
2. **Penalty application** — for each overdue loan found, calls `lms-circulation-api` to mark it
   late and `lms-membership-api` to apply the suspension, with retry if either call fails, so a
   partial failure (loan marked late, suspension call failed) doesn't get silently lost

Everything else — the eligibility check before registering a loan, the availability
decrement/restore when a loan is registered or returned — **stays exactly where it already is**:
a direct, synchronous HTTP call made by `circulation-api` (once built) to `membership-api`/
`catalog-api`, using the same client/port pattern already proven in `membership-service`'s
`circulation` package. `lms-workflow` is not in that call path at all.

**Justification:** a saga earns its complexity by handling failure across multiple steps over
time — that's exactly the overdue/penalty flow (nothing forces those two writes to happen in the
same instant, and one can legitimately fail while the other succeeds). It buys nothing for a
same-request eligibility check, which either succeeds or fails atomically as a single HTTP
round-trip today. Keeping the existing pattern for those calls means zero rework of code that
already exists and is already correct.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| **`lms-workflow` owns only overdue detection + penalty application (CHOSEN)** | Reuses the existing eligibility/availability client code unchanged; the saga only exists where its failure-handling is actually needed | `lms-workflow`'s own README's "loan → renewal → return → penalty" wording overstates its scope — must be documented clearly | — (chosen) |
| `lms-workflow` orchestrates the entire lifecycle, including loan registration and return | Single, consistent orchestration story for the whole lifecycle | Rewrites already-working synchronous calls into saga steps for a single-request operation with no partial-failure surface to protect against; more code, more moving parts, no corresponding benefit | Solves a problem that doesn't exist at the cost of real rework |
| No saga at all; a cron job directly applies penalties as one combined operation, no compensation | Simpler than a saga | If the suspension call fails after the loan is already marked late, nothing retries or reconciles it — silent data drift between two services | The whole reason `lms-workflow` exists is to handle exactly this multi-service partial-failure case; skipping it defeats the purpose of building the repo at all |

---

## Consequences

**Positive:**
- Zero changes to `membership-service`'s existing `ActiveLoansChecker`/`circulation` client code
  or to the planned `catalog-service` loan-copy/return-copy calls
- `lms-workflow` has a small, well-justified scope instead of an open-ended "orchestrates
  everything" mandate that would be hard to test and reason about
- Matches `pattern-guide.md`'s general Saga guidance ("use it when the transaction fits in a
  single service, don't; use it when it genuinely spans services over time, do") applied
  correctly instead of by default

**Negative / Trade-offs:**
- Two different coordination patterns now coexist in the same domain (direct sync calls for
  eligibility/availability, saga for overdue/penalty) — must be documented clearly so a future
  contributor doesn't assume one pattern covers everything
- "Renewal" (mentioned in `lms-workflow`'s README title) has no corresponding HU or domain rule
  anywhere in `02-domain/entities-and-rules.md` today — out of scope until a renewal HU actually
  exists; not designed for in this ADR

**Impact on the system:**
- Affected repos: `lms-workflow` (scope, once it receives its first commit), `lms-circulation-api`
  (the eligibility/availability calls it makes stay direct, not saga steps)
- Documents to update: `05-architecture/pattern-guide.md`'s Saga row (already updated in this
  same change to point here instead of a blanket "Yes")

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Someone builds `lms-workflow` per its README's literal wording and duplicates the eligibility/availability calls as saga steps | Medium | Medium — two code paths doing the same check, possible inconsistency | This ADR is the canonical scope; link it from `lms-workflow`'s own `decisions.md` once that repo exists |
| Penalty application's retry logic is under-designed (e.g., no idempotency check before re-suspending an already-suspended student) | Medium | Medium — a retried penalty step could double-apply a side effect | Design the retry step to be idempotent (check current state before writing) as part of `lms-workflow`'s own implementation, not deferred |
| A future "renewal" HU is added without revisiting this ADR's scope | Low | Low — scope creep into an ADR that didn't design for it | Revisit/supersede this ADR when a renewal HU is actually written, rather than silently expanding `lms-workflow`'s scope |

---

## References

- Part of the repo decomposition this ADR's target (`lms-workflow`) belongs to →
  `ADR-006-repo-per-context-decomposition.md`
- General Saga guidance this ADR applies → `05-architecture/pattern-guide.md`
- The existing, unchanged eligibility check this ADR keeps outside the saga →
  `09-microservices/services/03-membership-service/decisions.md` (once created;
  currently in `lms-library`'s `membership-service/internal/domain/membership/circulation_port.go`)
- Loan invariants this ADR does not change → `02-domain/entities-and-rules.md`, Loan (INV-001–INV-005)
- No message broker — how the saga is triggered instead → `ADR-009-worker-scheduling-model.md`
