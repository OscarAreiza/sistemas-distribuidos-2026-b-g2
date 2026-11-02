# Week 11 — Session 1: Instrumenting the services

Scope: the Circulation bounded context I own (`lms-circulation-api`) and the two services it calls
synchronously (`lms-membership-api`, `lms-catalog-api`). Critical journey: **HU-06 Register a loan**
(`POST /api/v1/loans`).

The working code and the detailed checkpoints are in `11-week/workshop/`. This document summarizes
what was built and where the evidence is.

## 1. Structured logs with a propagated `trace_id`

- Every request gets a W3C `traceparent`. A **trace-context middleware** reads it from the incoming
  request, or creates one if it is missing, and stores it in the request context.
- One **structured JSON logger** (zap) writes every line with the same base fields: `ts`, `level`,
  `service`, `trace_id`, `span_id`, `route`, `msg`. The domain fields are `loan_id`, `student_id`,
  `book_id`, `idempotency_key_hash` and `error_code`.
- A **"never log" list** keeps student PII (`document_id`, email) and JWTs out of the logs.
- Contract and example lines (one HU-06 request produces 7 lines with the same `trace_id`):
  `workshop/02-log-contract.md`.

## 2. Basic metrics: rate, errors, latency

Exposed on `/metrics` on the internal network only (Prometheus format):

| Signal | Metric | Type |
|---|---|---|
| Rate | `circulation_loan_registrations_total{outcome,error_code}` | counter |
| Errors | `circulation_loan_registrations_total{outcome="error"}` | counter |
| Latency | `circulation_loan_registration_duration_seconds{outcome}` | histogram (buckets aligned to NFR-001: 0.3 s / 0.6 s) |
| Latency per hop | `circulation_downstream_request_duration_seconds{downstream,operation}` | histogram |
| Errors per hop | `circulation_downstream_requests_total{downstream,operation,result}` | counter |

Full catalog, with the PromQL for each signal and the label cardinality rules (no IDs or `trace_id`
as labels; `trace_id` goes in exemplars): `workshop/04-metrics.md`.

## 3. One distributed trace across a multi-service request

HU-06 crosses four hops: `lms-circulation-portal` → nginx gateway → `lms-circulation-api` →
`lms-membership-api` (`GET /students/{id}`) + `lms-catalog-api` (`POST /books/{id}/loan-copy`) +
MongoDB `loan_db`.

- The circulation service forwards `traceparent` on every outbound HTTP call, so the downstream
  services log the **same `trace_id`**.
- A test proves that the `trace_id` reaches `lms-catalog-api` unchanged:
  `workshop/03-instrumentation/circulation-service` (`go test ./...`, how to run it in
  `HOW-TO-RUN.md`).
- The overdue scanner in `lms-workflow` runs asynchronously, so it is linked with a **span link**,
  not as a child span (ADR-008 / ADR-009). Trace map: `workshop/01-journey.md`.

## 4. What this found

- The real `catalogclient` and `membershipclient` in `lms-circulation-api` do not forward the
  correlation header yet, so in the deployed code correlation is lost after the first hop. Porting
  the workshop middleware into the real service is story OBS-05 in
  `SLO-ALERTING-RUNBOOK.md`.
- `POST /books/{id}/loan-copy` is not idempotent. A client timeout after catalog commits leaves a
  copy decremented with no Loan. This is visible in the example logs and is the scenario for the
  runbook.

## Evidence

| Item | Commit |
|---|---|
| Journey and trace map | `19b4567` |
| Log contract | `769105b` |
| Middleware and structured logger | `52a19be` |
| Propagation downstream and `/metrics` | `ce52809` |
| Test: `trace_id` reaches `lms-catalog-api` | `6afe19d` |
| Golden signals and metric catalog | `06b2cb8` |
