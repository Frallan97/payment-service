package middleware

import (
	"bytes"
	"encoding/json"
	"io"
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

// responseWriter wraps http.ResponseWriter to capture response
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// CachedResponse represents a cached idempotent response
type CachedResponse struct {
	StatusCode int             `json:"status_code"`
	Headers    http.Header     `json:"headers"`
	Body       json.RawMessage `json:"body"`
}

// readBody safely reads and restores the request body
func readBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Restore body for downstream handlers
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
