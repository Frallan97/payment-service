package middleware

import (
	"net/http"
)

// IdempotencyMiddleware handles idempotency keys for safe retries
// This is a simplified version - a production version would use Redis or database
func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to POST/PUT/PATCH/DELETE
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			// No idempotency key, proceed normally
			next.ServeHTTP(w, r)
			return
		}

		// TODO: Implement proper idempotency checking with Redis/database
		// For now, just pass through
		// In production:
		// 1. Check if key exists in cache/db
		// 2. If exists, return cached response
		// 3. If not, execute request and cache response
		// 4. Set TTL of 24 hours on cached response

		next.ServeHTTP(w, r)
	})
}
