# Week 1 — Session 1: The Real Problem, the Backlog Seed, and Consistency/Delivery Semantics

> **Added retroactively on 2026-09-02.** Session 1's brief ("form your team, pick the real
> problem, start the backlog, write down each core operation's required consistency and
> delivery semantics") was never written down explicitly in Week 1 — the original
> `01-week/hu-status/README.md` only covers the foundations material reviewed that week. This
> document closes that gap, grounded in the system as it has actually been designed and built
> since (`library-docs`, `lms-library`), not retconned to look like it was decided differently.

---

## The team

- Oscar Mauricio Areiza Páramo — Tech Lead
- Hermes Pascuas Herrera — DevOps
- Luis Alejandro Meneses — DevOps

## The real problem

A university library needs to track its book inventory, its registered students, and the
loan/return cycle between them — including automatically penalizing late returns — without a
paper process or a spreadsheet. A single Administrator operates the whole system; students are
records the Administrator manages, not system users themselves. Full framing in
`library-docs/01-context/overview.md` and `01-context/scope.md`.

## Backlog seed

The full, sliced backlog with testable acceptance criteria lives in
`library-docs/03-product/product-backlog.md` (9 user stories, MoSCoW-prioritized, across 4
epics) and `library-docs/04-requirements/user-stories.md` (Gherkin acceptance criteria per
story). At the time this document was written, HU-01, HU-02, HU-03, and HU-04 are merged and
tagged as `v1.0.0` — see `library-docs/04-requirements/traceability-matrix.md` for the
authoritative, current status of every story.

---

## Consistency and delivery semantics per core operation

> This is the part Session 1 asked to defend in the MVP 1 design. The system's core operations
> fall into two clearly different buckets: **single-service operations** (strongly consistent,
> ACID, no delivery-semantics question at all because there is only one write) and the **one
> cross-service operation**, loan registration, where consistency and delivery semantics are a
> real, explicit design trade-off.

### Single-service operations (strong consistency, no cross-service concern)

| Operation | Owning service | Consistency | Delivery semantics |
|---|---|---|---|
| Authenticate the Administrator | `access-service` | Strong — single read against `access_db` inside one transaction | N/A — read-only, safely retryable (idempotent) |
| Register / search / edit / deactivate a student | `membership-service` | Strong — single ACID write against `membership_db` | At-most-once in practice: `documentId` has a unique constraint, so a client retry after a timeout fails with `409 DOCUMENT_ID_ALREADY_EXISTS` instead of creating a duplicate |
| Register a book | `catalog-service` (Catalog, currently inside `library-api` until extracted) | Strong — single ACID write | At-most-once in practice: `isbn` has a unique constraint, same reasoning as above |

None of these needs a saga, a broker, or a retry policy — each is a single transaction against
a single database, and PostgreSQL's own ACID guarantees are the entire consistency story.

### The cross-service operation: registering a loan (Circulation)

This is the one operation that touches three bounded contexts (Circulation, Membership,
Catalog) and is the actual design decision Session 1 asked to defend. As designed in
`02-domain/entities-and-rules.md` and implemented in `LoanRegistrationService`
(`circulation-service/internal/domain/service/loan_registration_service.go`, currently on
`feat/HU-06-loan-registration`, not yet merged):

1. Circulation calls Membership: `GET /students/by-document/{documentId}` then checks
   suspension status — **read-only**, safely retryable.
2. Circulation calls Catalog: `POST /books/{id}/loan-copy` — **mutating**, decrements
   `availableCopies`.
3. Circulation writes the `Loan` row to its own database.

**Consistency model chosen: no distributed transaction, no saga, no compensation.** Each step
commits independently. If step 3 fails after step 2 succeeded, `availableCopies` is decremented
with no corresponding `Loan` row — a real, acknowledged inconsistency window.

**Delivery semantics chosen: at-most-once per step, no automatic retry.** A failed call is
surfaced as an error to the Administrator, not silently retried — retrying a
already-succeeded `loan-copy` call would double-decrement availability, since it is not
idempotent (no idempotency key is generated per loan attempt).

**Why this trade-off, and how we'd defend it in the MVP 1 design review:** a saga (with a
compensating "give the copy back" step) or a two-phase commit would close the inconsistency
window, but at real cost — a saga orchestrator or a broker to coordinate compensating actions
is infrastructure this 3-person, one-term project cannot operate for a data-loss scenario that
is rare (a mid-request crash) and cheaply fixable by hand (an Administrator can look at
`availableCopies` vs. the loans list and correct a stuck count). This mirrors the same
reasoning already recorded for the whole system's shape in `ADR-002` and `ADR-004`: match the
consistency mechanism to the team's actual ability to operate it, not to the textbook-correct
answer regardless of cost. If this project ever needed to eliminate that window for real (e.g.
a much higher write volume where manual correction stops being cheap), the next step would be
an outbox-pattern-based saga, not a bigger single transaction — Circulation, Membership, and
Catalog are already separate databases, so a shared transaction is not on the table.

---

## Correlations

- Full domain rules → `library-docs/02-domain/entities-and-rules.md`
- The cross-service coordination this document describes → `library-docs/09-microservices/communication-patterns.md`
- Architectural style and its justification → `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`, `ADR-004-incremental-microservices-decomposition.md`
- Current backlog and story status → `library-docs/03-product/product-backlog.md`, `library-docs/04-requirements/traceability-matrix.md`
