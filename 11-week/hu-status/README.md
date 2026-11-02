<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       11-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 11

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Make HU-06 (register a loan, POST /api/v1/loans) observable end to end - structured logs with a propagated trace_id, rate/errors/latency metrics and one distributed trace across circulation, membership and catalog (Session 1) - then define its SLOs, error-budget policy, symptom-based alerts with owner and runbook, and the observability stories (Session 2).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W11-S1 | Instrument the services: structured logs with propagated trace_id, rate/errors/latency metrics, one distributed trace across services | done | `OBSERVABILITY-INSTRUMENTATION.md` (this folder); https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/52a19be, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/ce52809, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/6afe19d |
| W11-S2 | Define 3 SLOs (availability, latency, freshness), error-budget policy, symptom-based alerts with owner + runbook, runbook, observability stories | done | `SLO-ALERTING-RUNBOOK.md` (this folder) - not yet committed/pushed |
| W11-CP1 | Critical journey HU-06 and trace map | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/19b4567 |
| W11-CP2 | Structured log contract (7 lines, one trace_id) | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/769105b |
| W11-CP3 | Instrumentation code: traceparent middleware, zap logger, propagation to lms-catalog-api, /metrics, test | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/52a19be, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/ce52809, https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/6afe19d |
| W11-CP4 | Golden signals and metric catalog | done | https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/commit/06b2cb8 |
| W11-CP6..CP9 | Workshop files `05-slo.md` to `08-incident.md` | doing | Their content (SLOs, budget, alerts, runbook) is written in `SLO-ALERTING-RUNBOOK.md`; the workshop stub files themselves are still `TODO` and the incident drill (CP9) has not been run |

## 2. My individual contribution
- Applied the whole week to the Circulation context I own (`lms-circulation-api`), on HU-06 across
  `lms-membership-api`, `lms-catalog-api` and MongoDB `loan_db`, with the W3C `traceparent` path and
  an async span link to the `lms-workflow` overdue scanner (ADR-008/ADR-009).
- Session 1: log contract with domain fields and a "never log" list for student PII and JWTs; a
  runnable Go module with the trace-context middleware, zap logger, outbound propagation and
  `/metrics`; a test proving the `trace_id` reaches `lms-catalog-api` unchanged; a metric catalog
  with NFR-001-aligned histogram buckets.
- Session 2: 3 SLOs on a 28-day window (availability 99.5 %, latency 99 % <= 0.6 s, overdue-status
  freshness 99 %) with their PromQL; error budget from an explicit traffic assumption; a 50/75/100 %
  budget policy; multi-window burn-rate alerts (fast 14.4x page, slow 6x ticket) with a minimum
  request count for low traffic, owner and runbook link; the runbook for the fast-burn alert; and
  5 observability stories (OBS-01..OBS-05) with testable acceptance criteria.

## 3. Blockers and risks
- The real `catalogclient`/`membershipclient` in `lms-circulation-api` do not forward the
  correlation header yet, so correlation is lost after the first hop in the deployed code (story
  OBS-05).
- `POST /books/{id}/loan-copy` is not idempotent: a client timeout after catalog commits leaves a
  decremented copy with no Loan (runbook mitigation 4).
- The traffic estimate (about 400 loans/day) is an assumption; there is no production traffic yet.
- No Prometheus is running in the local pod yet, so the alert rules are written but not loaded
  (story OBS-02).

## 4. Plan for next week
- Week 12: deploy with a zero-downtime strategy and write the MVP 3 deployment plan.
- Carry over: run the incident drill (CP9) and fill the workshop stub files from
  `SLO-ALERTING-RUNBOOK.md`; port the trace-context middleware into `lms-circulation-api` (OBS-05).

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Session documents in this same folder (`11-week/hu-status/`):
  - `OBSERVABILITY-INSTRUMENTATION.md` - Session 1: logs with trace_id, metrics, distributed trace
  - `SLO-ALERTING-RUNBOOK.md` - Session 2: SLOs, error budget and policy, alerts, runbook, stories
- Week 11 workshop: `11-week/workshop/` (README, 01-journey, 02-log-contract,
  03-instrumentation, 04-metrics)
- Other documents attached in this folder:
  - `figma-lms-library.txt` - link to the project's Figma design file
  - `Congreso de Robotica e IA.docx` - supporting course document
  - `release.md` - all PRs across the 12 microservice repos (`v2.0.1` release), by author
    (Oscar Areiza, Luis Meneses, Hermes Pascuas), with each repo's first/current tag commit
  - `lms-reportes.html` - local-environment dashboard/report snapshot
