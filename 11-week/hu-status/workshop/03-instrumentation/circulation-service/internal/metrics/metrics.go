// Package metrics exposes the two HU-06 journey metrics defined in
// 11-week/workshop/04-metrics.md on GET /metrics (internal network only,
// library-docs/13-operations/observability.md).
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Outcome values are a closed set, so the label stays bounded (cardinality 4).
const (
	OutcomeRegistered = "registered" // 201
	OutcomeReplayed   = "replayed"   // 200, same Idempotency-Key
	OutcomeRejected   = "rejected"   // 422, INV-002/003/004
	OutcomeError      = "error"      // 5xx, timeouts
)

// LatencyBuckets (seconds) are built around NFR-001: P95 < 300 ms and
// P99 < 600 ms for POST /loans, and the 5 s catalog/membership client timeout.
var LatencyBuckets = []float64{0.05, 0.1, 0.2, 0.3, 0.45, 0.6, 1, 2.5, 5}

type Metrics struct {
	Registry      *prometheus.Registry
	Registrations *prometheus.CounterVec
	Duration      *prometheus.HistogramVec
}

func New() *Metrics {
	m := &Metrics{
		Registry: prometheus.NewRegistry(),
		Registrations: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "circulation_loan_registrations_total",
			Help: "POST /api/v1/loans requests by outcome (HU-06).",
		}, []string{"outcome", "error_code"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "circulation_loan_registration_duration_seconds",
			Help:    "End-to-end latency of POST /api/v1/loans, including membership and catalog calls.",
			Buckets: LatencyBuckets,
		}, []string{"outcome"}),
	}
	m.Registry.MustRegister(m.Registrations, m.Duration)
	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
}
