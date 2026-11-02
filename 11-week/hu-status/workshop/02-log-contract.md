# 02 — Structured Log Contract

One JSON object per line, written to stdout by `zap` (the logger every `lms-*-api` already
uses, `internal/infrastructure/logger/logger.go`; NFR-005). The same contract applies to the
three services of the HU-06 journey: `lms-circulation-api`, `lms-membership-api`,
`lms-catalog-api` (and `lms-workflow` once it exists). Field names are `snake_case` so they
match the Prometheus labels of CP4 and can be grepped/queried the same way everywhere.

## Fields

| Name | Type | Required / optional | Example | Why it exists |
|---|---|---|---|---|
| `timestamp` | string, ISO-8601 UTC with ms | required | `2026-10-19T14:05:12.104Z` | Order events across 3 services and 2 DB engines; always UTC so the Mongo (`loan_db`) and Postgres timestamps line up |
| `level` | enum `DEBUG`/`INFO`/`WARN`/`ERROR` | required | `ERROR` | Filter noise; alerting reads only `ERROR` |
| `service` | string | required | `lms-circulation-api` | Repo name = deployable name (`service-catalog.md`), so a log line maps to one repo and one CODEOWNER |
| `version` | string (semver of the release tag) | required | `2.0.2` | Tells whether a failure started with a release (`release.md`: `v2.0.1` tagged, `v2.0.2` merged) |
| `environment` | enum `local`/`develop`/`qa`/`main` | required | `qa` | Same image runs in every environment (`cross-cutting.md`, centralized config) |
| `trace_id` | string, 32 hex (W3C) | required | `7c2e9a41d0b84f3e9a6b15c8e2f04d71` | Joins every line of one business intent (one loan submission, including its retry) across services |
| `span_id` | string, 16 hex (W3C) | required | `c81d4e0a92b7f356` | Says which hop of the CP1 span tree wrote the line |
| `correlation_id` | string, UUID | required | `5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90` | The `X-Correlation-Id` of one HTTP attempt; returned in the error envelope (`error.correlationId`) so the Administrator can quote it |
| `event` | string, `noun.verb_past` | required | `catalog.loan_copy.timeout` | Stable machine key to count/grep; `message` may change, `event` may not |
| `message` | string | required | `catalog-service did not answer loan-copy within 5000 ms` | Human sentence, never contains PII |
| `duration_ms` | integer | required on events that close an operation, else omitted | `5001` | Latency per hop without opening the trace UI |

## Domain fields

| Name | Type | Required / optional | Example | Why it exists |
|---|---|---|---|---|
| `loan_id` | string, UUID | optional (present once the `Loan` exists, step 7) | `e3b0c442-98fc-4c1c-a3d1-6f7b5e9a2c10` | `_id` of the `loans` document; joins logs to `loan_db.loans` and to `workflow_db.saga_runs` (ADR-008 idempotency key) |
| `student_id` | string, UUID | optional | `a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f` | Opaque id of `membership.students`, never the name/document; needed to debug INV-003/INV-004 |
| `book_id` | string, UUID | optional | `b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e` | Id of `catalog.books`; needed to debug INV-002 (`available_copies`) |
| `idempotency_key_hash` | string, first 12 hex of SHA-256 | optional (POST /loans only) | `9f86d081884c` | Shows that two attempts are the *same* intent without logging the raw `Idempotency-Key` (it is the `_id` of `loan_db.idempotency_keys`) |
| `error_code` | enum from the API error envelope | required on WARN/ERROR | `NO_COPIES_AVAILABLE` | Business code that maps 1:1 to an invariant: `STUDENT_SUSPENDED` (INV-003), `LOAN_LIMIT_REACHED` (INV-004), `NO_COPIES_AVAILABLE` (INV-002), or `INTERNAL_ERROR` |
| `downstream` | enum `lms-membership-api`/`lms-catalog-api`/`mongodb.loan_db` | optional (outbound events only) | `lms-catalog-api` | Which dependency of CP1 hops 3–7 the line is about |
| `http_status` | integer | optional | `422` | Status answered (server events) or received (client events) |
| `attempt` | integer | optional | `2` | Number of the portal attempt for the same `idempotency_key_hash`; makes retries visible |

## Never log

| Data (LMS Library domain) | Where it lives | How it is masked or replaced |
|---|---|---|
| Student full name, email, phone | `membership.students` | Never logged; only `student_id` (UUID) |
| Student `document_id` (national ID / cédula) | `membership.students.document_id` | Never logged, not even partially (unique index, a direct identifier) |
| Administrator password / `password_hash` | `access.administrators` | Never logged; login failure logs `event=auth.login_failed` without the username (INV-002 on Administrator) |
| RS256 session JWT (`Authorization: Bearer ...`) | browser → every `-api` | Header dropped by the logger; at most `token_sub` claim |
| HS256 internal service token and `INTERNAL_JWT_SECRET` | `membershipclient`/`catalogclient` | Never logged; the client logs only `downstream` and `http_status` |
| Mongo connection URI with `circulation_app` password, Postgres DSNs | `.env` / config | Logged as `mongodb://circulation_app:***@host/loan_db` only at startup, password replaced by `***` |
| Raw `Idempotency-Key` value | `loan_db.idempotency_keys._id` | Replaced by `idempotency_key_hash` (SHA-256, 12 hex) |
| Full request/response bodies | every endpoint | Never; only the ids listed in "Domain fields" |
| Payment / health data | — | LMS Library v1 has no fines payment and no health data (`01-context/scope.md`); if fines payment is added (`lms-worker`, ADR-009) card data must never reach a log |

## Example log lines (one trace_id)

Story: the Administrator registers a loan of the last copy of a book. Attempt 1 times out
on `loan-copy` (the 5 s client timeout in `catalogclient`), but `lms-catalog-api` commits the
decrement 240 ms later, so the copy is gone without a `Loan`. The portal resubmits the
same form (same `Idempotency-Key`, same `trace_id`, new `X-Correlation-Id`); the key was never
saved, so the flow runs again and is rejected by INV-002. Grep:
`grep 7c2e9a41d0b84f3e9a6b15c8e2f04d71 *.log`

```json
{"timestamp":"2026-10-19T14:05:12.104Z","level":"INFO","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"3f1a9c0e7b2d4a58","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"loan.request.received","message":"POST /api/v1/loans received","student_id":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e","idempotency_key_hash":"9f86d081884c","attempt":1}
{"timestamp":"2026-10-19T14:05:12.149Z","level":"INFO","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"9e4b27c1d05f6a3e","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"membership.student.eligible","message":"student has no active suspension","downstream":"lms-membership-api","http_status":200,"duration_ms":38,"student_id":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f"}
{"timestamp":"2026-10-19T14:05:17.158Z","level":"ERROR","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"c81d4e0a92b7f356","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"catalog.loan_copy.timeout","message":"catalog-service did not answer loan-copy within 5000 ms","downstream":"lms-catalog-api","duration_ms":5001,"error_code":"INTERNAL_ERROR","book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"timestamp":"2026-10-19T14:05:17.163Z","level":"WARN","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"3f1a9c0e7b2d4a58","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"loan.request.failed","message":"loan not registered, idempotency key not saved","http_status":500,"error_code":"INTERNAL_ERROR","duration_ms":5059,"idempotency_key_hash":"9f86d081884c","attempt":1}
{"timestamp":"2026-10-19T14:05:17.398Z","level":"INFO","service":"lms-catalog-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"5a0f3b9d17e2c864","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"book.copy.loaned","message":"available_copies decremented from 1 to 0","http_status":200,"duration_ms":5240,"book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"timestamp":"2026-10-19T14:05:19.020Z","level":"INFO","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"d29e7f1043ab5c06","correlation_id":"0b7e4c19-3d2a-4f68-b5e1-8c9d0a7f6e23","event":"loan.request.received","message":"POST /api/v1/loans received","student_id":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e","idempotency_key_hash":"9f86d081884c","attempt":2}
{"timestamp":"2026-10-19T14:05:19.087Z","level":"WARN","service":"lms-circulation-api","version":"2.0.2","environment":"qa","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"6b3c08e5f1a94d72","correlation_id":"0b7e4c19-3d2a-4f68-b5e1-8c9d0a7f6e23","event":"loan.rejected","message":"no available copies of this book (INV-002)","downstream":"lms-catalog-api","http_status":422,"error_code":"NO_COPIES_AVAILABLE","duration_ms":67,"book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e","attempt":2}
```

What the seven lines prove: line 3 is the **timeout to a downstream service**, line 7 the
**business rule rejection**; line 5 (written by a different service, same `trace_id`, its own
`span_id`) shows that the timed-out `loan-copy` actually committed, which is the real defect:
`catalog.books.available_copies` dropped to 0 with no `Loan` in `loan_db.loans`.

## Log levels

| Level | When to use | Example in LMS Library |
|---|---|---|
| `ERROR` | A request or job could not do its job because of *our* system or a dependency, and someone may need to act | `catalog.loan_copy.timeout` — `lms-catalog-api` did not answer in 5000 ms; `mongo.write.failed` on `loan_db.loans` |
| `WARN` | Expected, handled outcome that is not the happy path: business rule rejections, a request that ended in 5xx after the cause was already logged as ERROR, retries | `loan.rejected` with `error_code=LOAN_LIMIT_REACHED` (INV-004); `loan.request.failed` 500 |
| `INFO` | Business state changes and request boundaries, one line each | `loan.registered` (201, with `loan_id`); `book.copy.loaned` in `lms-catalog-api`; `loan.replayed` (200 for a repeated `Idempotency-Key`) |
| `DEBUG` | Detail for local troubleshooting, off in `qa`/`main` (`LOG_LEVEL=info`) | `mongo.loans.count` returned `active_loans=1` for the student before the INV-004 check |
