package api

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sbbullet/to-do/logger"
	"go.uber.org/zap"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error(
						"Request",
						zap.Any("err", err),
						zap.String("trace", string(debug.Stack())),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Info(
				"Request",
				zap.String("Method", r.Method),
				zap.String("Endpoint", r.RequestURI),
				zap.Int("Status", wrapped.status),
				zap.Duration("Duration", time.Since(start)),
			)
		}

		return http.HandlerFunc(fn)
	}
}
