# ADR-009 — `lms-worker`'s Scheduling Model (Polling, Not a Broker)

- **ID:** ADR-009
- **Date:** 2026-09-15
- **Status:** Accepted
- **Authors:** Oscar Areiza — Tech Lead
- **Reviewers:** Hermes Pascuas, Luis Alejandro Meneses — Development team

---

## Context

`lms-worker`'s README describes it as "Asynchronous jobs: due-date notices, fines accrual,
catalog reindexing. **Event-driven**, no HTTP surface of its own." Read literally, "event-driven"
implies a message broker or event bus — components a broker/consumer architecture would need.

That directly contradicts the system's own, already-accepted architectural debt: `AT-002` in
`05-architecture/overview.md` states plainly "No message broker — within `library-api`,
cross-module events would be in-process; already-extracted services coordinate via plain
synchronous HTTP instead of adopting a broker," and `09-microservices/communication-patterns.md`
is equally explicit: "There is no message broker, no event bus, and no asynchronous messaging
anywhere in this system." Adopting a broker just so `lms-worker` has something to consume would
silently overturn `AT-002` as a side effect of a repo-split ADR, not as its own deliberate,
separately-justified decision.

**Known constraints:**
- Same 3-person academic team, no budget for a managed broker (`01-context/scope.md`)
- `lms-worker`'s three responsibilities (due-date notices, fines accrual, catalog reindexing) are
  all things that can legitimately run on a schedule — none of them require sub-second reaction
  to a single event the moment it happens
- The "never read another service's database directly" rule
  (`09-microservices/data-ownership-matrix.md`) applies to `lms-worker` exactly like any other
  client — it can only reach `circulation-api`/`membership-api`/`catalog-api` through their public
  HTTP APIs

---

## Decision

**We decided:** `lms-worker` is **polling/cron-triggered**, not event-driven in the message-broker
sense. It runs on an internal scheduler (a Go ticker, or the container's own cron), and on each
tick:
1. Calls `lms-circulation-api`'s public API to find loans past `dueDate` (due-date notices) and
   loans eligible for a fines/penalty pass
2. Calls `lms-membership-api` where a notification or suspension follow-up is needed
3. Calls `lms-catalog-api` to trigger a reindex of search data

No message broker, no event bus, no new infrastructure is adopted. `AT-002` remains true and
unrevisited — this ADR reaffirms it rather than overturning it. `lms-worker`'s own README's
"event-driven" wording is treated as generic scaffold language describing the general shape
("things that happen without a user clicking a button"), not a literal requirement to build
pub/sub infrastructure.

**Relationship to `lms-workflow` (`ADR-008`):** the overdue-detection and penalty-application
work described in `ADR-008` and the due-date-notice/reindexing work here are adjacent but
distinct — `lms-workflow` owns the saga that changes state across services with compensation;
`lms-worker` owns simpler, non-compensable background tasks. If the two turn out to overlap in
practice (e.g. both wanting to scan for overdue loans), that overlap is resolved when
`lms-workflow` and `lms-worker` are actually built, not speculated on here.

**Justification:** every one of `lms-worker`'s three responsibilities tolerates the latency of a
polling interval (minutes, not milliseconds) — none of them is "notify the instant this happens."
A scheduler calling existing public APIs costs nothing new to operate; a broker would be a
standing piece of infrastructure a 3-person team has to run, monitor, and understand, for a
workload that doesn't need it.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| **Polling/cron-triggered, calling existing public APIs (CHOSEN)** | Zero new infrastructure; reuses `circulation-api`/`membership-api`/`catalog-api`'s existing public endpoints, same pattern every other client already uses | Not "real-time" — a due-date notice could be up to one polling interval late | — (chosen) |
| Adopt a message broker (RabbitMQ/Kafka/similar) and make `lms-worker` a real consumer | Matches the README's "event-driven" wording literally; lower latency between an event and the worker reacting | Overturns `AT-002`, a deliberately accepted piece of architectural debt, as an unplanned side effect; new infrastructure to operate with no CI/ops maturity to support it yet (`11-quality/testing-strategy.md` confirms no CI exists) | This project's own accepted scope explicitly defers this; revisit only if a real, separately-justified need for sub-second reactivity appears |
| Outbox pattern (services write an events table, `lms-worker` polls that table instead of the public API) | Slightly more efficient than polling a full API response | Requires `lms-worker` to read another service's database table directly, violating the ownership rule (`data-ownership-matrix.md`) unless each service also exposes an `/events` endpoint — extra API surface for no clear benefit at this scale | Public API polling already satisfies the ownership rule for free |

---

## Consequences

**Positive:**
- No new infrastructure category (broker/queue) enters the system's operational surface
- `lms-worker` is testable the same way every other client in this system is: mock the HTTP
  responses of the services it calls
- `AT-002` stays a single, consistently-true statement about the whole system, not "true except
  for `lms-worker`"

**Negative / Trade-offs:**
- Polling has inherent latency (bounded by the tick interval) and wastes some requests when
  nothing has changed — acceptable at this project's traffic scale, called out explicitly so it
  isn't mistaken for an oversight
- If a genuine low-latency requirement appears later (e.g. instant due-date push notifications),
  this ADR needs to be revisited/superseded, not silently worked around

**Impact on the system:**
- Affected repo: `lms-worker` only, once it receives its first commit
- Documents to update: none beyond this ADR and its cross-references — `AT-002` in
  `05-architecture/overview.md` is reaffirmed, not changed

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Someone builds `lms-worker` per its README's literal "event-driven" wording and starts standing up a broker | Medium | Medium — introduces infrastructure this ADR deliberately avoided, with no ops maturity to run it | This ADR is the canonical answer; link it from `lms-worker`'s own `decisions.md` once that repo exists |
| Polling interval is too coarse for a due-date notice to feel timely | Low | Low — single-admin system, not user-facing real-time chat | Tune the interval based on real usage once built; not a correctness issue, just a UX tuning knob |

---

## References

- Part of the repo decomposition this ADR's target (`lms-worker`) belongs to →
  `ADR-006-repo-per-context-decomposition.md`
- The architectural debt this ADR reaffirms rather than resolves → `05-architecture/overview.md`,
  `AT-002`
- No broker anywhere in this system, already established → `09-microservices/communication-patterns.md`
- Never read another service's database directly → `09-microservices/data-ownership-matrix.md`
- Related, adjacent scope → `ADR-008-circulation-saga-scope.md`
