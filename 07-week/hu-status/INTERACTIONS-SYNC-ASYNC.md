# Week 7 - Session 1: Sync vs. Async per Interaction

> Goal: for every interaction in the system, decide sync or async and justify it; pick REST or
> gRPC for the sync ones and a topic/queue for the async ones; make at least one consumer
> idempotent.
>
> **Status: design only, no code this session** (by decision — `circulation-service` isn't
> merged to `dev`/`main` yet, it still lives on feature branches, so this session inventories the
> *real* interactions already coded across those branches and decides their shape; Session 2
> formalizes the resulting contracts, and actual broker/consumer code lands when
> `circulation-service` is merged).

---

## 1. Every interaction in the system, inventoried from the real code

| # | Caller → Callee | What it does | Where it lives today |
|---|---|---|---|
| 1 | Browser → `access-service` | Administrator login (issues JWT) | `access-service` (on `dev`) |
| 2 | Browser → `membership-service` | Student CRUD, suspend, deactivate | `membership-service` (on `dev`) |
| 3 | Browser → `catalog-service` | Book CRUD, availability | `catalog-service` (on `dev`) |
| 4 | Browser → `circulation-service` | Register loan, register return, overdue report (HU-06/07/08) | `circulation-service` (feature branch `feat/HU-08-late-return-penalty`, not yet merged) |
| 5 | `membership-service` → `circulation-service` | `CountActive(studentId)` — how many active loans a student has, to decide if `Deactivate` is allowed | `membership-service/internal/infrastructure/circulation/client.go` (on `dev`) |
| 6 | `circulation-service` → `membership-service` | `IsEligible(studentId)` — is the student currently suspended, before granting a new loan | `circulation-service/internal/infrastructure/membership/client.go` (feature branch) |
| 7 | `circulation-service` → `catalog-service` | `LoanCopy(bookId)` / `ReturnCopy(bookId)` — decrement/increment a book's available copies | `catalog-service/internal/application/usecase/adjust_book_availability.go` (on `dev`), called from `circulation-service` (feature branch) |
| 8 | `circulation-service` → `membership-service` | `Suspend(studentId, days)` — apply the flat suspension when a return is late (HU-08, INV-006) | `circulation-service/internal/infrastructure/membership/client.go` (feature branch) — **this is the one this session moves to async** |

Interactions 5–8 all reuse the same pattern today: a short-lived, self-signed JWT
(`sub: "<service-name>"`) against the shared `JWT_SECRET`, calling the *other* service's already
public REST endpoint — no separate internal-only API surface (a deliberate v1 trade-off, per the
comments in `circulation-service/internal/infrastructure/membership/client.go` and
`membership-service/internal/infrastructure/circulation/client.go`).

## 2. Decision per interaction

| # | Interaction | Sync or async? | Protocol | Why |
|---|---|---|---|---|
| 1–4 | Browser ↔ any service | **Sync** | REST (JSON over HTTP, through the NGINX gateway) | User-facing request/response — the admin's screen needs an immediate result (token, saved record, error) to render. A browser also can't natively speak gRPC without a grpc-web proxy layer NGINX doesn't have — REST is both the simplest and the only zero-extra-infrastructure option here. |
| 5 | `membership` → `circulation`: `CountActive` | **Sync** | REST | `DeactivateStudent` (`membership-service/internal/application/usecase/deactivate_student.go`) needs the count *before* it can decide whether `Student.Deactivate` succeeds — there is no useful way to answer the admin's "deactivate this student" request without it. Low volume (one admin action at a time), so REST's simplicity beats any latency gRPC would save. |
| 6 | `circulation` → `membership`: `IsEligible` | **Sync** | REST | Same reasoning as #5, mirrored: `LoanRegistrationService.RegisterLoan` cannot decide INV-003 (no loan for a suspended student) without the answer first — a blocking precondition, not a side effect. |
| 7 | `circulation` → `catalog`: `LoanCopy` / `ReturnCopy` | **Sync** | REST | Availability is a hard invariant (INV-001/INV-002 on `Book`, `library-docs/02-domain/entities-and-rules.md`) — a loan must not be confirmed unless the copy count was actually decremented first. Making this async would let the system confirm a loan for a book that turns out to have zero copies left. This is exactly the kind of interaction that stays synchronous even though everything else nearby is being pushed to events. |
| 8 | `circulation` → `membership`: `Suspend` on a late return | **Async (was sync)** | Topic/queue — see Section 3 | Unlike #7, a late-return suspension is a *consequence*, not a precondition of anything the current request needs to succeed — `RegisterReturn` (HU-07) should be able to confirm "the book is returned" even if `membership-service` is briefly down or slow, and the admin doesn't need to wait on the suspension side effect to get their response. A short delay before the suspension takes effect is an acceptable trade-off the domain doesn't forbid (no invariant says a suspension must be visible within the same request). This is also, not coincidentally, an event this system already modeled once: before the microservices split, `backend/internal/domain/event/event.go` (pre-extraction commit `faa7240`) had a `LoanReturned` domain event "consumed by ... the membership module (apply suspension if late)" — dropped only because there was no broker at the time. This session brings that same idea back, across services instead of in-process. |

**REST vs. gRPC, system-wide:** every sync interaction above stays REST. All four services
already share one HTTP/JSON stack (`chi`, the same JWT scheme, the same response envelope), and
internal call volume is low (one admin action → at most one or two internal calls). gRPC's
benefits (binary framing, generated strongly-typed stubs, HTTP/2 multiplexing) matter at a
throughput/latency point this system isn't near yet, and would mean introducing a second,
parallel toolchain (`.proto` files, `protoc` codegen in every Go module) just for 3 internal
calls. Revisit as an ADR if internal call volume or latency ever becomes a measured problem —
until then, consistency with the existing stack wins.

## 3. The async interaction, designed

**Event:** `loan.returned.v1` — published by `circulation-service` once a return is registered
(HU-07), whether or not it was late; carries the same fields the old in-process `LoanReturned`
event had, plus what a cross-process consumer needs that an in-process call didn't (an event ID,
an explicit schema version):

```json
{
  "eventId": "b3f1c2d4-...",
  "eventType": "loan.returned.v1",
  "occurredAt": "2026-09-21T14:32:00Z",
  "data": {
    "loanId": "...",
    "studentId": "...",
    "bookId": "...",
    "dueDate": "2026-09-18T00:00:00Z",
    "returnDate": "2026-09-21T14:32:00Z",
    "isLate": true,
    "suspensionDays": 3
  }
}
```

**Transport:** a topic/exchange named `lms.circulation.loan-returned` (routing key
`loan.returned`), one durable queue per consumer (broker-agnostic naming; RabbitMQ is the
proposed broker — see `MVP2-CONTRACTS-PLAN.md`, Section 1, for why). `circulation-service` is the
sole publisher; `membership-service` is the first consumer (applies the suspension when
`isLate: true`; ignores the event otherwise). Delivery is at-least-once, which is exactly why
Section 4 matters.

**Why `catalog-service` still isn't a consumer of this event:** per interaction #7's decision
above, catalog's copy count is updated synchronously, in the same request, not via this event —
publishing `loan.returned` does not replace that call, it only replaces the *suspension* side
effect.

## 4. Making the consumer idempotent

At-least-once delivery means `membership-service` must expect to see the same `loan.returned`
event more than once (broker redelivery after a slow ack, a consumer restart mid-processing,
etc.) — without idempotency, a redelivered event would extend a student's suspension a second
time for the same late return.

**Design:** a `processed_events` table in `membership_db`:

| Column | Type | Notes |
|---|---|---|
| `event_id` | `UUID PRIMARY KEY` | The event's own `eventId`, not a generated one |
| `processed_at` | `TIMESTAMPTZ NOT NULL` | For observability/cleanup, not correctness |

Consumer logic (pseudocode), inside one DB transaction:

```
on receive(event):
    if event.eventType != "loan.returned.v1" or not event.data.isLate:
        ack(event); return   # nothing to do

    begin tx
        inserted = INSERT INTO processed_events (event_id, processed_at)
                   VALUES (event.eventId, now())
                   ON CONFLICT (event_id) DO NOTHING
        if inserted == 0:
            rollback tx        # already processed — duplicate delivery
            ack(event); return

        student = SELECT ... FOR UPDATE WHERE id = event.data.studentId
        student.suspend(event.data.suspensionDays)   # same domain method suspend_student.go already uses
        UPDATE students ...
    commit tx
    ack(event)
```

The `INSERT ... ON CONFLICT DO NOTHING` inside the same transaction as the suspension write is
what makes this safe under concurrent/duplicate delivery — not a separate "check-then-act" step,
which would itself have a race.

## Correlations

- Contracts, versioning, and the consumer-driven contract test this design feeds into →
  `MVP2-CONTRACTS-PLAN.md` (this folder)
- The event this revives → `backend/internal/domain/event/event.go` (pre-extraction, commit
  `faa7240`, `lms-library` repo)
- The invariants behind each sync decision → `library-docs/02-domain/entities-and-rules.md`
  (INV-001/002 on `Book`, INV-003/004/006 on `Loan`)
- Current inter-service HTTP clients → `membership-service/internal/infrastructure/circulation/client.go`,
  `circulation-service/internal/infrastructure/membership/client.go` (feature branch)
