package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"

	"payment-service/pkg/auth"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey       contextKey = "userID"
	EmailKey        contextKey = "email"
	NameKey         contextKey = "name"
	RoleKey         contextKey = "role"
	IsSuperAdminKey contextKey = "isSuperAdmin"
)

// AuthMiddleware validates JWT tokens from auth-service
func AuthMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "unauthorized",
						"message": "Missing authorization header",
					},
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "unauthorized",
						"message": "Invalid authorization header format",
					},
				})
				return
			}

			claims, err := auth.ValidateToken(parts[1], publicKey)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "unauthorized",
						"message": "Invalid or expired token",
					},
				})
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, NameKey, claims.Name)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			ctx = context.WithValue(ctx, IsSuperAdminKey, claims.IsSuperAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetEmailFromContext retrieves the email from request context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

// GetNameFromContext retrieves the name from request context
func GetNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(NameKey).(string)
	return name, ok
}

// GetRoleFromContext retrieves the role from request context
func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

// IsSuperAdmin checks if the user is a super admin
func IsSuperAdmin(ctx context.Context) bool {
	isSuperAdmin, ok := ctx.Value(IsSuperAdminKey).(bool)
	return ok && isSuperAdmin
}
