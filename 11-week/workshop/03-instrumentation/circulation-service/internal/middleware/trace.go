// Package middleware replaces lms-circulation-api's current CorrelationID
// middleware (internal/adapter/in/httpapi/middleware/correlation_id.go), which
// only knows X-Correlation-Id, with one that also speaks W3C traceparent.
package middleware

import (
	"net/http"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext"
)

// TraceContext reads traceparent and X-Correlation-Id from the incoming
// request, creates them when missing, stores them in the request context and
// returns them in the response headers, so the Administrator's error envelope
// and the portal can quote them.
func TraceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tc, ok := tracecontext.Parse(r.Header.Get(tracecontext.TraceparentHeader))
		if !ok {
			tc = tracecontext.New()
		}

		tc.CorrelationID = r.Header.Get(tracecontext.CorrelationIDHeader)
		if tc.CorrelationID == "" {
			tc.CorrelationID = tracecontext.NewCorrelationID()
		}

		w.Header().Set(tracecontext.TraceparentHeader, tc.Traceparent())
		w.Header().Set(tracecontext.CorrelationIDHeader, tc.CorrelationID)

		next.ServeHTTP(w, r.WithContext(tracecontext.WithContext(r.Context(), tc)))
	})
}
