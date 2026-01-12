package middleware

import (
	"context"
	"fmt"
	"net/http"
	"payment-service/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	redis  *redis.Client
	limit  int           // requests per window
	window time.Duration // time window
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		limit:  limit,
		window: window,
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (if authenticated)
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			// For unauthenticated requests, use IP address
			key := fmt.Sprintf("ratelimit:ip:%s", r.RemoteAddr)
			if allowed, remaining, resetTime := rl.checkLimit(r.Context(), key); !allowed {
				rl.writeRateLimitResponse(w, remaining, resetTime)
				return
			}
		} else {
			// For authenticated requests, use user ID
			key := fmt.Sprintf("ratelimit:user:%s", userID.String())
			if allowed, remaining, resetTime := rl.checkLimit(r.Context(), key); !allowed {
				rl.writeRateLimitResponse(w, remaining, resetTime)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// checkLimit checks if the request is allowed under the rate limit
func (rl *RateLimiter) checkLimit(ctx context.Context, key string) (allowed bool, remaining int, resetTime time.Time) {
	// Use Redis sliding window log algorithm
	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove old entries
	rl.redis.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixNano(), 10))

	// Count requests in current window
	count, err := rl.redis.ZCard(ctx, key).Result()
	if err != nil {
		// On error, allow the request (fail open)
		return true, rl.limit, now.Add(rl.window)
	}

	if count >= int64(rl.limit) {
		// Rate limit exceeded
		// Get the oldest entry to calculate reset time
		oldestEntries, err := rl.redis.ZRangeWithScores(ctx, key, 0, 0).Result()
		if err == nil && len(oldestEntries) > 0 {
			oldestTime := time.Unix(0, int64(oldestEntries[0].Score))
			resetTime = oldestTime.Add(rl.window)
		} else {
			resetTime = now.Add(rl.window)
		}
		return false, 0, resetTime
	}

	// Add current request
	member := now.UnixNano()
	rl.redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(member),
		Member: member,
	})

	// Set expiry on the key
	rl.redis.Expire(ctx, key, rl.window*2)

	remaining = rl.limit - int(count) - 1
	resetTime = now.Add(rl.window)

	return true, remaining, resetTime
}

// writeRateLimitResponse writes a rate limit exceeded response
func (rl *RateLimiter) writeRateLimitResponse(w http.ResponseWriter, remaining int, resetTime time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
	w.Header().Set("Retry-After", strconv.FormatInt(int64(time.Until(resetTime).Seconds()), 10))

	w.WriteHeader(http.StatusTooManyRequests)

	err := models.NewAPIError(
		models.ErrCodeInvalidRequest,
		"Rate limit exceeded. Please try again later.",
		http.StatusTooManyRequests,
	)

	_, _ = fmt.Fprintf(w, `{"error":{"code":"%s","message":"%s"},"retry_after":%d}`,
		err.Code, err.Message, int64(time.Until(resetTime).Seconds()))
}
