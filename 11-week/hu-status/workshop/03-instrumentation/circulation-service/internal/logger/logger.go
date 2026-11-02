// Package logger writes every line in the CP2 log contract
// (11-week/workshop/02-log-contract.md): timestamp, level, service, version,
// environment, trace_id, span_id, correlation_id, event, message, duration_ms.
package logger

import (
	"context"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/code-corhuila/lms-circulation-api/w11/circulation-service/internal/tracecontext"
)

type Logger struct{ base *zap.Logger }

// New builds a JSON logger bound to service/version/environment. w is stdout in
// production and a buffer in tests.
func New(w io.Writer, service, version, environment, level string) *Logger {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = zapcore.InfoLevel
	}

	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:     "timestamp",
		LevelKey:    "level",
		MessageKey:  "message",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
			e.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
		},
	})

	core := zapcore.NewCore(enc, zapcore.AddSync(w), lvl)
	return &Logger{base: zap.New(core).With(
		zap.String("service", service),
		zap.String("version", version),
		zap.String("environment", environment),
	)}
}

// NewStdout is what cmd/ uses.
func NewStdout(service, version, environment, level string) *Logger {
	return New(os.Stdout, service, version, environment, level)
}

// Ctx injects trace_id, span_id and correlation_id from the request context
// automatically, so no handler has to remember them.
func (l *Logger) Ctx(ctx context.Context) *zap.Logger {
	tc, ok := tracecontext.FromContext(ctx)
	if !ok {
		return l.base
	}
	return l.base.With(
		zap.String("trace_id", tc.TraceID),
		zap.String("span_id", tc.SpanID),
		zap.String("correlation_id", tc.CorrelationID),
	)
}

// Event is the only way handlers log: event is mandatory (CP2 "event" field).
func Event(event string) zap.Field { return zap.String("event", event) }

// Duration renders duration_ms as an integer number of milliseconds.
func Duration(d time.Duration) zap.Field { return zap.Int64("duration_ms", d.Milliseconds()) }
