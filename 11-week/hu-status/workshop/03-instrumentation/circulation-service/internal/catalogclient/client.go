// Package catalogclient is the CP1 hop 6 outbound call:
// lms-circulation-api -> lms-catalog-api POST /api/v1/books/{id}/loan-copy.
// Unlike the current lms-circulation-api catalogclient, it propagates the
// trace context (traceparent + X-Correlation-Id) to the downstream service.
package catalogclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrNoCopiesAvailable = errors.New("no copies available") // INV-002, catalog answers 409
)

type Client struct {
	baseURL string
	http    *http.Client
}

// New uses the same 5 s timeout as the production catalogclient.
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: timeout}}
}

// LoanCopy decrements catalog.books.available_copies for bookID. The outbound
// request carries a *child* span of the current request: same trace_id, new
// span_id, so lms-catalog-api's server span is parented to this call.
func (c *Client) LoanCopy(ctx context.Context, bookID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/api/v1/books/%s/loan-copy", c.baseURL, bookID), nil)
	if err != nil {
		return err
	}

	if tc, ok := tracecontext.FromContext(ctx); ok {
		child := tc.Child()
		req.Header.Set(tracecontext.TraceparentHeader, child.Traceparent())
		req.Header.Set(tracecontext.CorrelationIDHeader, child.CorrelationID)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("calling lms-catalog-api: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrBookNotFound
	case http.StatusConflict:
		return ErrNoCopiesAvailable
	default:
		return fmt.Errorf("lms-catalog-api returned %d", resp.StatusCode)
	}
}
