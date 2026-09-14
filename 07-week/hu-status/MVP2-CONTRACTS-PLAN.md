# Week 7 - Session 2: Contracts and Integration Plan for MVP 2

> Goal (planning): publish contracts (`openapi.yaml` / `.proto` / event schemas) in the repo,
> define versioning + compatibility rules, add at least one consumer-driven contract test to CI,
> and slice the integration stories for MVP 2 with testable acceptance criteria.
>
> **Status: draft.** Builds directly on `INTERACTIONS-SYNC-ASYNC.md` (this folder). Like that
> document, this is a planning deliverable this week — no CI workflow or contract file is added
> to the `lms-library`/`library-docs` repos yet; that lands once the team reviews this plan and
> `circulation-service` is actually merged.

---

## 1. Contracts to publish

| Contract | Format | Where it lives | Covers |
|---|---|---|---|
| `access-service` API | OpenAPI 3.0 (`access.v1.yaml`) | `library-docs/07-api/contracts/openapi/` | Login, token issuance |
| `membership-service` API | OpenAPI 3.0 (`membership.v1.yaml`) | same folder | Student CRUD, suspend, the `IsEligible`/`CountActive`-backing `GET /students/{id}` |
| `catalog-service` API | OpenAPI 3.0 (`catalog.v1.yaml`) | same folder | Book CRUD, `LoanCopy`/`ReturnCopy` |
| `circulation-service` API | OpenAPI 3.0 (`circulation.v1.yaml`) | same folder | Register loan, register return, overdue report |
| `loan.returned` event | JSON Schema (`loan-returned.v1.schema.json`) | `library-docs/07-api/contracts/events/` | The async payload designed in `INTERACTIONS-SYNC-ASYNC.md`, Section 3 |

No `.proto` files are planned — Section 2 of `INTERACTIONS-SYNC-ASYNC.md` keeps every sync
interaction on REST, so there is no gRPC surface to contract yet. This closes the gap already
flagged in `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md` ("`07-api/contracts/openapi/` still has
the generic template, not real per-service contracts") — four real specs replace the one generic
placeholder.

**Broker choice for the event contract:** RabbitMQ — a single extra container in
`docker-compose.yml`, a management UI for free (useful for a course project's demo), native topic
exchanges (fits the routing-key model in `INTERACTIONS-SYNC-ASYNC.md`), and Go client libraries
(`amqp091-go`) that need no code generation step — unlike Kafka, which would need a second
container just for coordination (`ZooKeeper`/`KRaft`) that this system's volume never justifies.

## 2. Versioning rules

- **REST:** the major version lives in the URL path, which every service already does
  (`/api/v1/...`). A breaking change (removed/renamed field, changed status code semantics,
  removed endpoint) requires a new path version (`/api/v2/...`) served alongside `v1` until every
  consumer (today: only the frontend and the other three services) has migrated. Additive,
  backward-compatible changes (a new optional field, a new endpoint) ship under the existing
  `v1` with no version bump.
- **Events:** the version lives in `eventType` (`loan.returned.v1`, per
  `INTERACTIONS-SYNC-ASYNC.md`). A breaking change (removing/renaming a field in `data`, changing
  a field's meaning or type) requires a new `eventType` (`loan.returned.v2`), published
  side-by-side with `v1` for a deprecation window until `membership-service` (today's only
  consumer) has migrated its handler. Adding a new optional field to `data` does not bump the
  version.

## 3. Compatibility rules

- **Producers never remove or repurpose a field within the same version** — only add optional
  ones. This is the one rule a contract test can mechanically enforce (Section 4).
- **Consumers must ignore unknown fields** — `membership-service`'s event handler (and every
  REST client's JSON decoding) must tolerate fields it doesn't know about, so a producer can add
  fields without coordinating a simultaneous consumer deploy.
- **No required field is ever added to an existing version** — a new required field is a breaking
  change by definition (old producers wouldn't send it), so it always means a version bump, never
  an in-place edit.

## 4. Consumer-driven contract test (planned CI addition)

**Pattern:** Pact-style consumer-driven contracts (`pact-go` on both sides — everything here is
already Go).

- `membership-service` (the consumer of `loan.returned.v1`) publishes a **pact**: the minimal
  shape of the event it actually reads (`data.studentId`, `data.isLate`, `data.suspensionDays`,
  `eventId`) — not the full schema, just what its handler uses. This is the point of
  consumer-driven contracts: the consumer defines what it depends on, not the producer guessing.
- `circulation-service`'s CI job (the producer) runs a **provider verification** step against
  that pact on every PR that touches the event's construction code
  (`circulation-service/internal/domain/service/loan_registration_service.go` once merged, or its
  event-publishing equivalent) — replays the pact's expectations against the actual struct the
  service would publish, and fails the build if a field the consumer depends on is missing,
  renamed, or has changed type.
- CI wiring (once `circulation-service` is on `dev` and this is implemented, not this session):
  a new job in `lms-library`'s CI (there is no CI pipeline at all yet, per
  `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md`, Section 2 — this would be its first job) named
  `contract-test`, running after unit tests, blocking merge on failure.

This is deliberately the *smallest* contract test that satisfies "at least one" — a single
consumer (`membership-service`), a single event (`loan.returned.v1`), verified on the producer
side in CI.

## 5. MVP 2 integration backlog — sliced, with acceptance criteria

| ID | Title | Acceptance criteria (Given/When/Then) | Status |
|---|---|---|---|
| **HU-INT-01** | Sync/async decision matrix documented and justified | **Given** the system's real inter-service interactions (`INTERACTIONS-SYNC-ASYNC.md`, Section 1), **when** the matrix is reviewed, **then** each one has an explicit sync/async decision, a chosen protocol (REST or a topic/queue), and a justification tied to a concrete invariant or consistency trade-off — not a generic rule of thumb. | doing — `INTERACTIONS-SYNC-ASYNC.md`, Sections 1–2, this session |
| **HU-INT-02** | `loan.returned` event formalized and its consumer made idempotent | **Given** a `loan.returned.v1` event delivered twice (broker at-least-once semantics), **when** `membership-service` processes both deliveries, **then** the student's suspension is applied exactly once, verified via the `processed_events` dedupe table (`INTERACTIONS-SYNC-ASYNC.md`, Section 4). | doing — design complete this session; code lands once `circulation-service` is merged and the broker is added to `docker-compose.yml` |
| **HU-INT-03** | Contracts published with versioning + compatibility rules | **Given** the four `openapi.v1.yaml` files and `loan-returned.v1.schema.json`, **when** a producer adds an optional field, **then** no version bump is required and existing consumers keep working unmodified; **given** a producer removes or repurposes a field, **when** that change is reviewed, **then** it is rejected unless it ships as a new major/event version per Section 2. | todo — this document is the first draft of the rules; publishing the actual files is next |
| **HU-INT-04** | Consumer-driven contract test wired into CI | **Given** a PR that changes how `circulation-service` builds a `loan.returned` event, **when** CI runs, **then** the `contract-test` job fails if a field `membership-service`'s pact depends on is missing/renamed/retyped, and passes otherwise. | todo — needs `circulation-service` on `dev` and a first CI pipeline (neither exists yet) |

## Correlations

- Interaction inventory and the async design this formalizes → `INTERACTIONS-SYNC-ASYNC.md`
  (this folder)
- Existing contract gap → `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md`, Section 2
- Branch↔environment mapping these contracts get promoted through →
  `06-week/hu-status/MVP2-ENVIRONMENTS-PLAN.md`
- Domain invariants behind the sync decisions → `library-docs/02-domain/entities-and-rules.md`
