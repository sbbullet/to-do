package api

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/token"
	"github.com/sbbullet/to-do/util"
	"go.uber.org/zap"
)

type authPayloadKey string

const (
	authorizationHeaderKey                 = "authorization"
	authorizationType                      = "bearer"
	authUsernameHeaderKey                  = "auth_username"
	authorizationPayloadKey authPayloadKey = "todo_app_auth_payload"
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

// AuthMiddleware checks for authorization header and extracts payload if authorized
func AuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				util.RespondWithUauthorizedError(w, "You are not authorized to perform the action")
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				util.RespondWithUauthorizedError(w, "You are not authorized to perform the action")
				return
			}

			if authorizationType != strings.ToLower(fields[0]) {
				util.RespondWithUauthorizedError(w, "You are not authorized to perform the action")
				return
			}

			authToken := fields[1]
			payload, err := tokenMaker.VerifyToken(authToken)
			if err != nil {
				if errors.Is(err, token.ErrTokenExpired) {
					util.RespondWithUauthorizedError(w, "Your session has expired. Please, log in again to new session")
					return
				}
				util.RespondWithUauthorizedError(w, "You are not authorized to perform the action")
				return
			}

			r.Header.Set(authUsernameHeaderKey, payload.Username)

			ctx := context.WithValue(r.Context(), authorizationPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
