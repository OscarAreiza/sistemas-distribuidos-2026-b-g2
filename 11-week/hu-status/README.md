<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       11-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 11

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Week 11 workshop - make HU-06 (register a loan, POST /api/v1/loans) observable end to end: trace map, log contract, instrumentation code, metrics, SLOs, burn-rate alerts, runbook and incident drill.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W11-CP1 | Critical journey HU-06 and trace map | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/19b4567 |
| W11-CP2 | Structured log contract (7 lines, one trace_id) | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/769105b |
| W11-CP3 | Instrumentation code: traceparent middleware, zap logger, propagation to lms-catalog-api, /metrics, test | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/52a19be, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/ce52809, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/6afe19d |
| W11-CP4 | Golden signals and metric catalog | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/06b2cb8 |
| W11-CP6..CP9 | SLOs, burn-rate alerts, runbook, incident drill (delivery 2) | todo | - |

## 2. My individual contribution
- Delivery 1 of the week 11 workshop (`11-week/workshop/`), applied to the Circulation context I own
  (`lms-circulation-api`, `lms-circulation-db`): HU-06 journey across `lms-membership-api`,
  `lms-catalog-api` and MongoDB `loan_db`, with the W3C `traceparent` path and an async span link
  to the `lms-workflow` overdue scanner (ADR-008/ADR-009).
- Log contract with domain fields (`loan_id`, `student_id`, `book_id`, `idempotency_key_hash`,
  `error_code`) and a "never log" list for student PII and JWTs.
- Runnable Go module `circulation-service` (`03-instrumentation/`), all tests passing locally
  (`go test ./...`), with a test proving the `trace_id` reaches `lms-catalog-api` unchanged.
- Metric catalog with NFR-001-aligned histogram buckets (0.3 s / 0.6 s boundaries).

## 3. Blockers and risks
- Found while instrumenting: the real `catalogclient`/`membershipclient` in `lms-circulation-api` do
  not forward `X-Correlation-Id`, so correlation is lost after the first hop.
- `POST /books/{id}/loan-copy` is not idempotent: a client timeout after catalog commits leaves a
  decremented copy with no Loan (shown in the CP2 example logs).

## 4. Plan for next week
- Delivery 2 of the workshop (CP6-CP10): SLOs and error budget, burn-rate alerts, runbook,
  incident drill and postmortem.
- Port the trace-context middleware and outbound propagation into `lms-circulation-api`.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Local documents attached in this same folder (`11-week/hu-status/`):
  - `figma-lms-library.txt` - link to the project's Figma design file
  - `Congreso de Robotica e IA.docx` - supporting course document
  - `release.md` - all PRs across the 12 microservice repos (`v2.0.1` release), by author
    (Oscar Areiza, Luis Meneses, Hermes Pascuas), with each repo's first/current tag commit
  - `lms-reportes.html` - local-environment dashboard/report snapshot
- Week 11 workshop (delivery 1): `11-week/workshop/` (README, 01-journey, 02-log-contract,
  03-instrumentation, 04-metrics)
