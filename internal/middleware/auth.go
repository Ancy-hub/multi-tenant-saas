package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
)

// contextKey is the type for context keys.
type contextKey string

// UserIDKey is the context key for storing user ID.
const UserIDKey contextKey = "user_id"

// AuthMiddleware validates JWT tokens and attaches user ID to request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// attach user_id to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
