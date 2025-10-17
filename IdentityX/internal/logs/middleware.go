package logs

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"GoAuth/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LogsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		reqID := GetRequestID(r.Context())
		userID := GetUserID(r.Context())

		L().Info("http_request",
			zap.String("request_id", reqID),
			zap.String("user_id", userID),
			zap.String("method", r.Method),
			zap.String("path", NormalizePath(r)),
			zap.Int("status", ww.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (ww *statusWriter) WriteHeader(statusCode int) {
	ww.status = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

type ctxKey string

const requestIDKey ctxKey = "requestID"
const userIDKey ctxKey = "userID"

func RequestIDMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			if access_token_cookie, err := r.Cookie("access_token"); err == nil {
				if uid := utils.ParseAccessTokenUserIDUnsafe(access_token_cookie.Value, viper.GetString("JWT_SECRET")); uid != nil {
					userID = *uid
				}
			}
		}

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
