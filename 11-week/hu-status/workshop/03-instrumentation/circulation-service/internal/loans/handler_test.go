package loans

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/catalogclient"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/logger"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/metrics"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/middleware"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext"
)

const (
	incomingTraceparent = "00-7c2e9a41d0b84f3e9a6b15c8e2f04d71-3f1a9c0e7b2d4a58-01"
	incomingTraceID     = "7c2e9a41d0b84f3e9a6b15c8e2f04d71"
	bookID              = "b9c8d7e6-f5a4-4b3c-9d2e-1f0a9b8c7d6e"
)

// fakeCatalog plays lms-catalog-api and records the headers it received.
func fakeCatalog(t *testing.T, status int, got *http.Header) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/books/"+bookID+"/loan-copy" {
			t.Errorf("unexpected downstream path %s", r.URL.Path)
		}
		*got = r.Header.Clone()
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newServer(catalogURL string, logs io.Writer) (http.Handler, *metrics.Metrics) {
	m := metrics.New()
	log := logger.New(logs, "lms-circulation-api", "2.0.2", "test", "debug")
	h := NewHandler(catalogclient.New(catalogURL, 5*time.Second), log, m)
	return middleware.TraceContext(http.HandlerFunc(h.Create)), m
}

func postLoan(t *testing.T, srv http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	body := `{"studentId":"a4f1e2d3-5b6c-4d7e-8f90-1a2b3c4d5e6f","bookId":"` + bookID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", strings.NewReader(body))
	req.Header.Set("traceparent", incomingTraceparent)
	req.Header.Set("X-Correlation-Id", "5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

// The CP3 proof: the trace_id received by lms-circulation-api is the trace_id
// sent downstream to lms-catalog-api, in a new child span.
func TestCreateLoan_PropagatesTraceIDToCatalog(t *testing.T) {
	var downstream http.Header
	catalog := fakeCatalog(t, http.StatusOK, &downstream)
	var logs bytes.Buffer
	srv, m := newServer(catalog.URL, io.MultiWriter(&logs, os.Stdout))

	rec := postLoan(t, srv)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	sent, ok := tracecontext.Parse(downstream.Get("traceparent"))
	if !ok {
		t.Fatalf("catalog received no valid traceparent: %q", downstream.Get("traceparent"))
	}
	if sent.TraceID != incomingTraceID {
		t.Fatalf("trace_id sent downstream = %s, want %s", sent.TraceID, incomingTraceID)
	}
	server, _ := tracecontext.Parse(rec.Header().Get("traceparent"))
	if sent.ParentSpanID == server.ParentSpanID || sent.ParentSpanID == "3f1a9c0e7b2d4a58" {
		t.Fatalf("downstream span must be a new child span, got parent %s", sent.ParentSpanID)
	}
	if downstream.Get("X-Correlation-Id") != "5d1f6c2a-8b0e-4e57-9a3c-1f2e7d6b4a90" {
		t.Fatalf("X-Correlation-Id not propagated")
	}
	if got := testutil.ToFloat64(m.Registrations.WithLabelValues(metrics.OutcomeRegistered, "")); got != 1 {
		t.Fatalf("circulation_loan_registrations_total{outcome=registered} = %v, want 1", got)
	}
	assertLogContract(t, &logs)
}

func TestCreateLoan_NoCopiesIsBusinessRejection(t *testing.T) {
	var downstream http.Header
	catalog := fakeCatalog(t, http.StatusConflict, &downstream)
	var logs bytes.Buffer
	srv, m := newServer(catalog.URL, &logs)

	rec := postLoan(t, srv)

	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "NO_COPIES_AVAILABLE") {
		t.Fatalf("got %d %s, want 422 NO_COPIES_AVAILABLE", rec.Code, rec.Body.String())
	}
	if got := testutil.ToFloat64(m.Registrations.WithLabelValues(metrics.OutcomeRejected, "NO_COPIES_AVAILABLE")); got != 1 {
		t.Fatalf("rejected counter = %v, want 1", got)
	}
	assertLogContract(t, &logs)
}

// assertLogContract checks every line against 02-log-contract.md.
func assertLogContract(t *testing.T, logs *bytes.Buffer) {
	t.Helper()
	required := []string{"timestamp", "level", "service", "version", "environment",
		"trace_id", "span_id", "correlation_id", "event", "message"}
	sc := bufio.NewScanner(logs)
	lines := 0
	for sc.Scan() {
		lines++
		var line map[string]any
		if err := json.Unmarshal(sc.Bytes(), &line); err != nil {
			t.Fatalf("log line is not JSON: %s", sc.Text())
		}
		for _, f := range required {
			if _, ok := line[f]; !ok {
				t.Fatalf("log line missing %q: %s", f, sc.Text())
			}
		}
		if line["trace_id"] != incomingTraceID {
			t.Fatalf("log line has trace_id %v, want %s", line["trace_id"], incomingTraceID)
		}
		if _, err := time.Parse("2006-01-02T15:04:05.000Z", line["timestamp"].(string)); err != nil {
			t.Fatalf("timestamp not ISO-8601 UTC: %v", line["timestamp"])
		}
	}
	if lines == 0 {
		t.Fatal("no log lines written")
	}
}
