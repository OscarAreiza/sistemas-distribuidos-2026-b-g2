package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext"
)

func TestTraceContext_ReusesIncomingTraceparent(t *testing.T) {
	const in = "00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-3f1a9c0e7b2d4a58-01"
	var got tracecontext.TraceContext
	h := TraceContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = tracecontext.FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", nil)
	req.Header.Set("traceparent", in)
	req.Header.Set("X-Correlation-Id", "5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got.TraceID != "7c2e9a41d0b84f3e9a6b15c8e2f04d71" || got.ParentSpanID != "3f1a9c0e7b2d4a58" {
		t.Fatalf("trace context not reused: %+v", got)
	}
	if !strings.Contains(rec.Header().Get("traceparent"), got.TraceID) {
		t.Fatalf("response traceparent missing trace_id: %q", rec.Header().Get("traceparent"))
	}
	if rec.Header().Get("X-Correlation-Id") != "5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90" {
		t.Fatalf("correlation id not echoed")
	}
}

func TestTraceContext_CreatesWhenMissingOrInvalid(t *testing.T) {
	h := TraceContext(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", nil)
	req.Header.Set("traceparent", "garbage")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if _, ok := tracecontext.Parse(rec.Header().Get("traceparent")); !ok {
		t.Fatalf("expected a valid new traceparent, got %q", rec.Header().Get("traceparent"))
	}
	if rec.Header().Get("X-Correlation-Id") == "" {
		t.Fatalf("expected a generated X-Correlation-Id")
	}
}
