// Package tracecontext holds the W3C Trace Context (traceparent) of the
// request currently handled by circulation-service, plus the X-Correlation-Id
// the lms-front apiClient already sends on every call (CP1, hops 1-2).
package tracecontext

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	TraceparentHeader   = "traceparent"
	CorrelationIDHeader = "X-Correlation-Id"
)

// TraceContext is what travels between lms-circulation-api, lms-membership-api
// and lms-catalog-api. SpanID is the span of *this* service; ParentSpanID is
// the span id received from the caller ("" when this service started the trace).
type TraceContext struct {
	TraceID       string
	SpanID        string
	ParentSpanID  string
	Flags         string
	CorrelationID string
}

type ctxKey struct{}

// Parse reads a traceparent header ("00-<32 hex>-<16 hex>-<2 hex>") and opens a
// new server span under it. ok=false means the header is missing or invalid
// and the caller must start a new trace with New.
func Parse(traceparent string) (tc TraceContext, ok bool) {
	parts := strings.Split(strings.TrimSpace(traceparent), "-")
	if len(parts) != 4 || parts[0] != "00" ||
		!isHex(parts[1], 32) || !isHex(parts[2], 16) || !isHex(parts[3], 2) ||
		parts[1] == strings.Repeat("0", 32) || parts[2] == strings.Repeat("0", 16) {
		return TraceContext{}, false
	}
	return TraceContext{
		TraceID:      parts[1],
		SpanID:       randomHex(8),
		ParentSpanID: parts[2],
		Flags:        parts[3],
	}, true
}

// New starts a brand-new trace (this service is the root).
func New() TraceContext {
	return TraceContext{TraceID: randomHex(16), SpanID: randomHex(8), Flags: "01"}
}

// Child opens a client span for an outbound call (membership/catalog/mongo):
// same trace_id, new span_id, parent = the current span.
func (tc TraceContext) Child() TraceContext {
	return TraceContext{
		TraceID:       tc.TraceID,
		SpanID:        randomHex(8),
		ParentSpanID:  tc.SpanID,
		Flags:         tc.Flags,
		CorrelationID: tc.CorrelationID,
	}
}

// Traceparent renders the header value that identifies *this* span to the next hop.
func (tc TraceContext) Traceparent() string {
	return fmt.Sprintf("00-%s-%s-%s", tc.TraceID, tc.SpanID, tc.Flags)
}

// NewCorrelationID is used when the request arrives without X-Correlation-Id
// (e.g. a call that bypassed the lms-front apiClient).
func NewCorrelationID() string { return uuid.NewString() }

func WithContext(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, tc)
}

func FromContext(ctx context.Context) (TraceContext, bool) {
	tc, ok := ctx.Value(ctxKey{}).(TraceContext)
	return tc, ok
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func isHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil && s == strings.ToLower(s)
}
