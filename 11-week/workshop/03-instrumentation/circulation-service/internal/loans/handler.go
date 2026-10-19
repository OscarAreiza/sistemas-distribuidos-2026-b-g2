// Package loans is a trimmed-down POST /api/v1/loans (HU-06): only the
// catalog hop is real here, enough to show propagation, metrics and logs.
// The membership check and Mongo writes of the real
// LoanRegistrationService are out of scope for this workshop module.
package loans

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/catalogclient"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/logger"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/metrics"
)

type Handler struct {
	catalog *catalogclient.Client
	log     *logger.Logger
	metrics *metrics.Metrics
}

func NewHandler(catalog *catalogclient.Client, log *logger.Logger, m *metrics.Metrics) *Handler {
	return &Handler{catalog: catalog, log: log, metrics: m}
}

type createLoanRequest struct {
	StudentID string `json:"studentId"`
	BookID    string `json:"bookId"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log := h.log.Ctx(r.Context())

	var req createLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookID == "" || req.StudentID == "" {
		h.finish(w, log, start, http.StatusBadRequest, metrics.OutcomeRejected, "VALIDATION_ERROR")
		return
	}
	log.Info("POST /api/v1/loans received", logger.Event("loan.request.received"),
		zap.String("student_id", req.StudentID), zap.String("book_id", req.BookID))

	callStart := time.Now()
	err := h.catalog.LoanCopy(r.Context(), req.BookID)
	switch {
	case errors.Is(err, catalogclient.ErrNoCopiesAvailable):
		h.finish(w, log, start, http.StatusUnprocessableEntity, metrics.OutcomeRejected, "NO_COPIES_AVAILABLE")
		return
	case errors.Is(err, catalogclient.ErrBookNotFound):
		h.finish(w, log, start, http.StatusUnprocessableEntity, metrics.OutcomeRejected, "VALIDATION_ERROR")
		return
	case err != nil:
		log.Error("catalog-service loan-copy failed", logger.Event("catalog.loan_copy.failed"),
			zap.String("downstream", "lms-catalog-api"), logger.Duration(time.Since(callStart)),
			zap.String("error_code", "INTERNAL_ERROR"), zap.String("book_id", req.BookID))
		h.finish(w, log, start, http.StatusInternalServerError, metrics.OutcomeError, "INTERNAL_ERROR")
		return
	}
	log.Info("available_copies decremented", logger.Event("catalog.copy.loaned"),
		zap.String("downstream", "lms-catalog-api"), zap.Int("http_status", http.StatusOK),
		logger.Duration(time.Since(callStart)), zap.String("book_id", req.BookID))

	h.finish(w, log, start, http.StatusCreated, metrics.OutcomeRegistered, "")
}

// finish writes the response, the closing log line and both metrics in one place.
func (h *Handler) finish(w http.ResponseWriter, log *zap.Logger, start time.Time, status int, outcome, errorCode string) {
	elapsed := time.Since(start)
	h.metrics.Registrations.WithLabelValues(outcome, errorCode).Inc()
	h.metrics.Duration.WithLabelValues(outcome).Observe(elapsed.Seconds())

	fields := []zap.Field{zap.Int("http_status", status), logger.Duration(elapsed)}
	switch outcome {
	case metrics.OutcomeRegistered:
		log.Info("loan registered", append(fields, logger.Event("loan.registered"))...)
	case metrics.OutcomeRejected:
		log.Warn("loan rejected by a business rule", append(fields, logger.Event("loan.rejected"), zap.String("error_code", errorCode))...)
	default:
		log.Warn("loan not registered", append(fields, logger.Event("loan.request.failed"), zap.String("error_code", errorCode))...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if errorCode != "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{
			"code": errorCode, "correlationId": w.Header().Get("X-Correlation-Id"),
		}})
	}
}
