package auth

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userIDContextKey contextKey = "authUserID"
	responseWriterContextKey contextKey = "authResponseWriter"
)

// WithUserID attaches an authenticated user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// UserIDFromContext returns the authenticated user ID when present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok && userID != ""
}

// WithResponseWriter attaches the HTTP response writer for cookie-setting resolvers.
func WithResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseWriterContextKey, w)
}

// ResponseWriterFromContext returns the response writer when available.
func ResponseWriterFromContext(ctx context.Context) (http.ResponseWriter, bool) {
	writer, ok := ctx.Value(responseWriterContextKey).(http.ResponseWriter)
	return writer, ok && writer != nil
}
