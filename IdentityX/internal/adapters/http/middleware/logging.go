package middleware

import (
	"GoAuth/internal/adapters/observability/logs"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func Logs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		reqID := GetRequestID(r.Context())
		userID := GetUserID(r.Context())

		routePattern := "not_found"
		if rctx := chi.RouteContext(r.Context()); rctx != nil {
			if pattern := rctx.RoutePattern(); pattern != "" {
				routePattern = pattern
			}
		}

		logs.L().Info("http_request",
			zap.String("request_id", reqID),
			zap.String("user_id", userID),
			zap.String("method", r.Method),
			zap.String("path", routePattern),
			zap.Int("status", ww.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

type ctxKey string

const requestIDKey ctxKey = "requestID"
const userIDKey ctxKey = "userID"

var (
	GoAuthMiddlewareTracer = otel.Tracer("GoAuthMiddlewareTracer")
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := GoAuthMiddlewareTracer.Start(r.Context(), "Middleware.RequestID")
		defer span.End()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		span.SetAttributes(attribute.String("request_id", reqID))

		ctx = context.WithValue(ctx, requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		userID := r.Header.Get("X-User-ID")

		span.SetAttributes(attribute.String("user_id", userID))

		ctx = context.WithValue(ctx, userIDKey, userID)
		w.Header().Set("X-User-ID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(userIDKey).(string); ok {
		return v
	}
	return ""
}
