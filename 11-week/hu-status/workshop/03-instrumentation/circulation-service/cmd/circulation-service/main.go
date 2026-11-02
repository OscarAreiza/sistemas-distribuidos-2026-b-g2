// Command circulation-service runs the instrumented POST /api/v1/loans and
// GET /metrics on PORT (default 8083, same as lms-circulation-api).
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/catalogclient"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/loans"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/logger"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/metrics"
	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/middleware"
)

func main() {
	log := logger.NewStdout("lms-circulation-api", "2.0.2", env("APP_ENV", "local"), env("LOG_LEVEL", "info"))
	m := metrics.New()
	catalog := catalogclient.New(env("CATALOG_SERVICE_URL", "http://localhost:8082"), 5*time.Second)

	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/loans", middleware.TraceContext(http.HandlerFunc(loans.NewHandler(catalog, log, m).Create)))
	mux.Handle("GET /metrics", m.Handler())

	addr := ":" + env("PORT", "8083")
	log.Ctx(context.Background()).Info("circulation-service listening", logger.Event("service.started"), zap.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Ctx(context.Background()).Error("server stopped", logger.Event("service.stopped"), zap.Error(err))
		os.Exit(1)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
