# How to Run — `circulation-service` instrumentation module

Go module that instruments CP1's journey (HU-06, `POST /api/v1/loans`) the way
`lms-circulation-api` should: W3C `traceparent` + `X-Correlation-Id` middleware, the CP2 JSON
logger (zap, same library as every `lms-*-api`), propagation to `lms-catalog-api`
`POST /api/v1/books/{id}/loan-copy` (CP1 hop 6) and two Prometheus metrics on `GET /metrics`.

```text
circulation-service/
├── cmd/circulation-service/main.go        wiring: POST /api/v1/loans, GET /metrics
└── internal/
    ├── tracecontext/   parse / create / child traceparent, context helpers
    ├── middleware/     TraceContext middleware (+ test)
    ├── logger/         zap JSON logger in the CP2 contract
    ├── catalogclient/  outbound loan-copy call that injects traceparent
    ├── metrics/        circulation_loan_registrations_total, ..._duration_seconds
    └── loans/          trimmed POST /api/v1/loans handler (+ propagation test)
```

Requirements: Go 1.25 (same `go 1.25.1` as `lms-circulation-api/go.mod`).

## Install

```bash
cd 11-week/workshop/03-instrumentation/circulation-service
go mod download
```

## Run

```bash
# terminal 1 — any HTTP server on :8082 that answers 200 to POST plays lms-catalog-api;
#             or point CATALOG_SERVICE_URL to a real lms-catalog-api
PORT=8083 CATALOG_SERVICE_URL=http://localhost:8082 APP_ENV=local go run ./cmd/circulation-service

# terminal 2
curl -i -X POST localhost:8083/api/v1/loans \
  -H 'traceparent: 00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-3f1a9c0e7b2d4a58-01' \
  -H 'Content-Type: application/json' \
  -d '{"studentId":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","bookId":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}'
curl -s localhost:8083/metrics | grep '^circulation_'
```

## Test

```bash
go vet ./...
go test -v ./...
```

`TestCreateLoan_PropagatesTraceIDToCatalog` starts a fake `lms-catalog-api` with `httptest`,
sends a request with a known `traceparent` and asserts that the `trace_id` received
downstream is the same one, in a **new child span**, that `X-Correlation-Id` travels too,
that `circulation_loan_registrations_total{outcome="registered"}` is 1, and that every log
line written is valid JSON with all CP2 mandatory fields and that same `trace_id`.

## Output

Real output from my machine (Ubuntu, go1.25.1 linux/amd64), 2026-10-19:

```text
$ go test -v ./...
?   	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/cmd/circulation-service	[no test files]
?   	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/catalogclient	[no test files]
=== RUN   TestCreateLoan_PropagatesTraceIDToCatalog
{"level":"INFO","timestamp":"2026-10-19T14:32:30.989Z","message":"POST /api/v1/loans received","service":"lms-circulation-api","version":"2.0.2","environment":"test","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"4608f28ed9a5b2e8","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"loan.request.received","student_id":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"level":"INFO","timestamp":"2026-10-19T14:32:30.990Z","message":"available_copies decremented","service":"lms-circulation-api","version":"2.0.2","environment":"test","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"4608f28ed9a5b2e8","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","event":"catalog.copy.loaned","downstream":"lms-catalog-api","http_status":200,"duration_ms":0,"book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"level":"INFO","timestamp":"2026-10-19T14:32:30.990Z","message":"loan registered","service":"lms-circulation-api","version":"2.0.2","environment":"test","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"4608f28ed9a5b2e8","correlation_id":"5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90","http_status":201,"duration_ms":1,"event":"loan.registered"}
--- PASS: TestCreateLoan_PropagatesTraceIDToCatalog (0.00s)
=== RUN   TestCreateLoan_NoCopiesIsBusinessRejection
--- PASS: TestCreateLoan_NoCopiesIsBusinessRejection (0.00s)
PASS
ok  	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/loans	0.011s
?   	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/logger	[no test files]
?   	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/metrics	[no test files]
=== RUN   TestTraceContext_ReusesIncomingTraceparent
--- PASS: TestTraceContext_ReusesIncomingTraceparent (0.00s)
=== RUN   TestTraceContext_CreatesWhenMissingOrInvalid
--- PASS: TestTraceContext_CreatesWhenMissingOrInvalid (0.00s)
PASS
ok  	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/middleware	0.006s
?   	github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext	[no test files]
```

Manual run (service on :18083, a 10-line Python stand-in for `lms-catalog-api` on :18082
that prints the headers it receives):

```text
$ curl -i -X POST localhost:18083/api/v1/loans -H 'traceparent: 00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-3f1a9c0e7b2d4a58-01' ...
HTTP/1.1 201 Created
Traceparent: 00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-439f0e7572ace7e7-01
X-Correlation-Id: e161f65d-4c38-4518-a7a9-d411d57abefb

# fake lms-catalog-api received (same trace_id, new child span 7502fa8f8286dcf9):
lms-catalog-api (fake) got /api/v1/books/b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e/loan-copy traceparent= 00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-7502fa8f8286dcf9-01

# service stdout
{"level":"INFO","timestamp":"2026-10-19T14:32:43.439Z","message":"circulation-service listening","service":"lms-circulation-api","version":"2.0.2","environment":"local","event":"service.started","addr":":18083"}
{"level":"INFO","timestamp":"2026-10-19T14:32:44.462Z","message":"POST /api/v1/loans received","service":"lms-circulation-api","version":"2.0.2","environment":"local","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"439f0e7572ace7e7","correlation_id":"e161f65d-4c38-4518-a7a9-d411d57abefb","event":"loan.request.received","student_id":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"level":"INFO","timestamp":"2026-10-19T14:32:44.465Z","message":"available_copies decremented","service":"lms-circulation-api","version":"2.0.2","environment":"local","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"439f0e7572ace7e7","correlation_id":"e161f65d-4c38-4518-a7a9-d411d57abefb","event":"catalog.copy.loaned","downstream":"lms-catalog-api","http_status":200,"duration_ms":3,"book_id":"b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"}
{"level":"INFO","timestamp":"2026-10-19T14:32:44.465Z","message":"loan registered","service":"lms-circulation-api","version":"2.0.2","environment":"local","trace_id":"7c2e9a41d0b84f3e9a6b15c8e2f04d71","span_id":"439f0e7572ace7e7","correlation_id":"e161f65d-4c38-4518-a7a9-d411d57abefb","http_status":201,"duration_ms":3,"event":"loan.registered"}

$ curl -s localhost:18083/metrics | grep '^circulation_'
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.05"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.1"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.2"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.3"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.45"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="0.6"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="1"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="2.5"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="5"} 1
circulation_loan_registration_duration_seconds_bucket{outcome="registered",le="+Inf"} 1
circulation_loan_registration_duration_seconds_sum{outcome="registered"} 0.003243861
circulation_loan_registration_duration_seconds_count{outcome="registered"} 1
circulation_loan_registrations_total{error_code="",outcome="registered"} 1
```

The curl did not send `X-Correlation-Id`, so the middleware generated
`e161f65d-...` and returned it; the incoming `traceparent` was kept (same `trace_id`).

## Known limitations (honest scope)

- Only CP1 hop 6 (catalog `loan-copy`) is implemented; the membership eligibility check
  (hop 4) and the MongoDB writes (hops 3, 5, 7) are not part of this module.
- Spans are not exported to a tracing backend (no OpenTelemetry SDK/Jaeger yet,
  NFR-005); the module proves the *propagation contract* that an OTel SDK would follow.
- The client-side `catalog.copy.loaned` log line uses the server `span_id`; the child span id
  is only visible in the outbound `traceparent`.
