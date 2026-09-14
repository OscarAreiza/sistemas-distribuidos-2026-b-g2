<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       07-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 07

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Session 1 - for every interaction in the system, decide sync vs. async and justify it, pick REST/gRPC for the sync ones and a topic/queue for the async ones, and design at least one idempotent consumer. Session 2 (planning) - formalize those decisions as versioned contracts: publish OpenAPI/event-schema files, define versioning + compatibility rules, plan a consumer-driven contract test in CI, and slice the MVP 2 integration backlog with testable acceptance criteria.
<!-- CONFIG-END -->

> **This week's brief (as given):** For each interaction in your system, decide sync or async and
> justify it; pick REST or gRPC for the sync ones and topics/queues for the async ones. Make at
> least one consumer idempotent. Session 2 (planning) will formalize these as versioned
> contracts. Publish your contracts (`openapi.yaml` / `.proto` / event schemas) in the repo,
> define your versioning + compatibility rules, and add at least one consumer-driven contract
> test to CI. Slice the integration stories for MVP 2 with testable acceptance criteria.

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-INT-01 | Sync/async decision matrix, justified per interaction | doing | See individual contribution below |
| HU-INT-02 | `loan.returned` event formalized + idempotent consumer designed | doing | See individual contribution below |
| HU-INT-03 | Contracts published with versioning/compatibility rules | todo | Planned for Session 2 this week |
| HU-INT-04 | Consumer-driven contract test wired into CI | todo | Planned for Session 2 this week |

## 2. My individual contribution
- **HU-INT-01, Session 1 - inventoried every real interaction in the system**, not a
  hypothetical one: read every existing inter-service HTTP client in `lms-library` (`dev`, plus
  `circulation-service`'s not-yet-merged feature branches) to find the actual 8 interactions -
  4 browser-facing, and 4 service-to-service (`membership`→`circulation` CountActive,
  `circulation`→`membership` IsEligible, `circulation`→`catalog` LoanCopy/ReturnCopy,
  `circulation`→`membership` Suspend). Decided and justified sync vs. async for each, keeping
  REST system-wide for the sync ones (no gRPC - justified why in the doc). Full detail →
  `INTERACTIONS-SYNC-ASYNC.md` (this folder).
- **HU-INT-02 - moved exactly one interaction to async and designed its idempotent consumer:**
  the late-return suspension call (`circulation`→`membership` `Suspend`) becomes the
  `loan.returned.v1` event over a topic/queue, since it's a side effect, not a precondition
  (unlike the availability check, which stays sync because it gates a hard invariant). This
  actually revives a domain event the system already had *before* the microservices split -
  `backend/internal/domain/event/event.go`'s `LoanReturned`, dropped when everything became
  synchronous HTTP - now redesigned across services with an `eventId` and a `processed_events`
  dedupe table so a redelivered event never double-suspends a student. Full detail →
  `INTERACTIONS-SYNC-ASYNC.md`, Sections 3-4.
- **HU-INT-03/04, Session 2 (planning) - drafted the contracts plan:** which OpenAPI/event-schema
  files to publish and where, the versioning rule (path version for REST, `eventType` suffix for
  events) and compatibility rule (additive-only within a version, consumers ignore unknown
  fields), the broker choice (RabbitMQ, justified against Kafka), and a Pact-style
  consumer-driven contract test plan for CI. Full detail → `MVP2-CONTRACTS-PLAN.md` (this
  folder).

## 3. Blockers and risks
- **Scope decision this week:** by design, this session is documentation/design only - no broker
  was added to `docker-compose.yml`, no publisher/consumer code was written, and no CI workflow
  was created. Reason: `circulation-service` (the event's publisher) isn't merged to `dev`/`main`
  yet, still lives on `feat/HU-06/07/08-...` branches, so wiring real infrastructure around a
  service that isn't in the mainline yet would be built on sand.
- `lms-library` has no CI pipeline at all yet (same gap noted in
  `05-week/hu-status/MVP1-SHIPPING-CHECKLIST.md`) - HU-INT-04 (the contract test in CI) can't
  land until a first CI workflow exists, which is a bigger prerequisite than this week's scope.
- The MVP 2 integration HU IDs (`HU-INT-0N`) are placeholders, same situation as last week's
  `HU-ORCH-0N` - need team sign-off and real IDs in `03-product/product-backlog.md`.

## 4. Plan for next week
- Take `INTERACTIONS-SYNC-ASYNC.md` and `MVP2-CONTRACTS-PLAN.md` to the team for review.
- Once `circulation-service` is merged to `dev` (tracked since Week 05's "Circulation ships as a
  fast-follow" decision), add RabbitMQ to `docker-compose.yml`, implement the `loan.returned.v1`
  publisher in `circulation-service` and the idempotent consumer in `membership-service` per the
  design in Section 4 of `INTERACTIONS-SYNC-ASYNC.md` (closes HU-INT-02 for real).
- Publish the actual `openapi.v1.yaml` files and `loan-returned.v1.schema.json` in
  `library-docs/07-api/contracts/` (closes HU-INT-03).
- Stand up a first CI pipeline for `lms-library`, then add the `contract-test` job
  (closes HU-INT-04).

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration) - design-only session, no code changed
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Local documents attached in this same folder (`07-week/hu-status/`):
  - `INTERACTIONS-SYNC-ASYNC.md` - Session 1: the full interaction inventory, the sync/async
    decision per interaction, the `loan.returned.v1` event design, and the idempotent-consumer
    design.
  - `MVP2-CONTRACTS-PLAN.md` - Session 2: contracts to publish, versioning/compatibility rules,
    the consumer-driven contract test plan, and the sliced MVP 2 integration backlog.
- Real code referenced (all in `lms-library`): `membership-service/internal/infrastructure/circulation/client.go`,
  `membership-service/internal/application/usecase/deactivate_student.go` (`dev` branch);
  `circulation-service/internal/domain/service/loan_registration_service.go`,
  `circulation-service/internal/infrastructure/membership/client.go` (`feat/HU-08-late-return-penalty`
  branch, commit `7bc59c5`); `backend/internal/domain/event/event.go` (pre-extraction, commit
  `faa7240`).
