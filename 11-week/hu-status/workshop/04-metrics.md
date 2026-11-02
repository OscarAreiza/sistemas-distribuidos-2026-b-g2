# 04 — Golden Signals and Metric Catalog

Journey: HU-06 `POST /api/v1/loans` (CP1). One table per service the request touches.
Conventions: `snake_case`, service prefix (`circulation_`, `membership_`, `catalog_`), unit
suffix (`_seconds`), `_total` for counters. `/metrics` is scraped on the internal network only,
never routed by `lms-api-gateway` (`13-operations/observability.md`). The first two metrics of
`lms-circulation-api` are the ones implemented in `03-instrumentation/` (CP3).

## `lms-circulation-api` (:8083, MongoDB `loan_db`)

| Signal | Metric name | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `circulation_loan_registration_duration_seconds` | histogram | `outcome` | `histogram_quantile(0.95, sum by (le) (rate(circulation_loan_registration_duration_seconds_bucket{outcome!="rejected"}[5m])))` |
| Latency (per hop) | `circulation_downstream_request_duration_seconds` | histogram | `downstream`, `operation` | `histogram_quantile(0.95, sum by (le, downstream, operation) (rate(circulation_downstream_request_duration_seconds_bucket[5m])))` |
| Traffic | `circulation_loan_registrations_total` | counter | `outcome`, `error_code` | `sum by (outcome) (rate(circulation_loan_registrations_total[5m]))` |
| Errors | `circulation_loan_registrations_total{outcome="error"}` | counter | `outcome`, `error_code` | `sum(rate(circulation_loan_registrations_total{outcome="error"}[5m])) / sum(rate(circulation_loan_registrations_total[5m]))` |
| Errors (per hop) | `circulation_downstream_requests_total` | counter | `downstream`, `operation`, `result` | `sum by (downstream) (rate(circulation_downstream_requests_total{result=~"timeout\|error"}[5m]))` |
| Saturation | `circulation_http_in_flight_requests` | gauge | `route` | `max_over_time(circulation_http_in_flight_requests{route="/api/v1/loans"}[5m])` |
| Saturation | `circulation_mongo_pool_connections` | gauge | `state` (`in_use`, `idle`) | `circulation_mongo_pool_connections{state="in_use"} / ignoring(state) sum without(state) (circulation_mongo_pool_connections)` |

`outcome="rejected"` (422 `STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`) is
**not** an error: it is the system correctly applying INV-002/003/004. It counts as traffic
and is excluded from the error ratio and from the latency SLO.

## `lms-membership-api` (:8081, PostgreSQL `membership_db`)

| Signal | Metric name | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `membership_http_request_duration_seconds` | histogram | `method`, `route`, `status_class` | `histogram_quantile(0.95, sum by (le) (rate(membership_http_request_duration_seconds_bucket{route="/api/v1/students/{id}",method="GET"}[5m])))` |
| Traffic | `membership_http_requests_total` | counter | `method`, `route`, `status_class` | `sum by (route) (rate(membership_http_requests_total[5m]))` |
| Errors | `membership_http_requests_total{status_class="5xx"}` | counter | `method`, `route`, `status_class` | `sum(rate(membership_http_requests_total{status_class="5xx"}[5m])) / sum(rate(membership_http_requests_total[5m]))` |
| Saturation | `membership_pg_pool_acquired_connections` (+ `membership_pg_pool_max_connections`) | gauge | — | `membership_pg_pool_acquired_connections / membership_pg_pool_max_connections` |

## `lms-catalog-api` (:8082, PostgreSQL `catalog_db`)

| Signal | Metric name | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `catalog_http_request_duration_seconds` | histogram | `method`, `route`, `status_class` | `histogram_quantile(0.95, sum by (le) (rate(catalog_http_request_duration_seconds_bucket{route="/api/v1/books/{id}/loan-copy"}[5m])))` |
| Traffic | `catalog_copy_operations_total` | counter | `operation` (`loan_copy`, `return_copy`), `result` (`ok`, `no_copies`, `not_found`, `error`) | `sum by (operation) (rate(catalog_copy_operations_total[5m]))` |
| Errors | `catalog_http_requests_total{status_class="5xx"}` | counter | `method`, `route`, `status_class` | `sum(rate(catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy",status_class="5xx"}[5m])) / sum(rate(catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy"}[5m]))` |
| Saturation | `catalog_pg_pool_acquired_connections` (+ `catalog_pg_pool_max_connections`) | gauge | — | `catalog_pg_pool_acquired_connections / catalog_pg_pool_max_connections` |
| Saturation | `catalog_pg_lock_wait_seconds` | histogram | `table` (`books`) | `histogram_quantile(0.95, sum by (le) (rate(catalog_pg_lock_wait_seconds_bucket[5m])))` — row lock on `catalog.books` while decrementing `available_copies` |

## Label cardinality

| Label | Values | Estimated cardinality | Allowed / forbidden | Replacement |
|---|---|---|---|---|
| `outcome` | `registered`, `replayed`, `rejected`, `error` | 4 | allowed | — |
| `error_code` | `""`, `STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`, `VALIDATION_ERROR`, `INTERNAL_ERROR` | 6 | allowed (closed enum from the API error envelope) | — |
| `downstream` | `lms-membership-api`, `lms-catalog-api`, `mongodb.loan_db` | 3 | allowed | — |
| `operation` | `get_student`, `loan_copy`, `count_active`, `insert_loan` | 4 | allowed | — |
| `result` | `ok`, `timeout`, `error`, `no_copies`, `not_found` | 5 | allowed | — |
| `method` | GET, POST, PATCH, DELETE | 4 | allowed | — |
| `route` | templated chi/mux route, e.g. `/api/v1/books/{id}/loan-copy` | ~12 per service | allowed **only templated** | — |
| `status_class` | `2xx`, `3xx`, `4xx`, `5xx` | 4 | allowed | raw `status` not needed for SLOs |
| `state`, `table` | 2 / 1 | 2 | allowed | — |
| raw URL path (`/api/v1/books/b9c8.../loan-copy`) | one per book | unbounded (thousands of books, `01-context/scope.md`) | **forbidden** | `route` template |
| `student_id`, `book_id`, `loan_id` | one per entity, `loans` grows forever | unbounded | **forbidden** | the id goes in the CP2 log line; dashboards link to logs |
| `trace_id`, `correlation_id`, `idempotency_key` | one per request | unbounded | **forbidden** | Prometheus **exemplars** on the histogram carry `trace_id` |
| `document_id`, email, Administrator username | PII | unbounded + personal data | **forbidden** (also in the CP2 "never log" list) | nothing; not needed for any signal |

Worst case series for `lms-circulation-api`: counter 4 × 6 = 24 (≈ 8 real combinations),
latency histogram 4 outcomes × 12 series (10 buckets + `_sum` + `_count`) = 48, downstream
histogram 4 operations × 12 = 48. Well under 200 series per instance.

## Histogram buckets

- Buckets (seconds): `0.05, 0.1, 0.2, 0.3, 0.45, 0.6, 1, 2.5, 5`
  (`metrics.LatencyBuckets` in `03-instrumentation/circulation-service/internal/metrics/metrics.go`).
- Justification:
  - The user is the **Administrator at the circulation desk with the Student standing in front
    of them**: a loan must feel instant. NFR-001 makes `POST /loans` a critical endpoint with
    **P95 < 300 ms** and **P99 < 600 ms**, so `0.3` and `0.6` are exact bucket boundaries; the
    CP6 latency SLI counts `le="0.3"` directly, with no interpolation error.
  - `0.05`–`0.2`: a healthy request is 3 internal hops (membership GET, Mongo count, catalog
    loan-copy) of ~20–60 ms each, so the normal mass sits here and we can see it shift
    before it reaches the SLO.
  - `0.45`: between the two NFR targets, to see degradation coming.
  - `1` and `2.5`: "the desk notices" (> 1 s) vs "the Administrator clicks again" (≈ 2–3 s),
    which is what produces the duplicate submissions handled by `Idempotency-Key`.
  - `5`: the `http.Client` timeout of `membershipclient`/`catalogclient`; anything above it is
    a timeout and lands in `+Inf`.
