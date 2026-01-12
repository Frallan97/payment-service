package auth

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims from auth-service
type Claims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// PublicKeyCache caches the public key from auth-service
type PublicKeyCache struct {
	Key       *rsa.PublicKey
	FetchedAt time.Time
	mu        sync.RWMutex
}

var keyCache = &PublicKeyCache{}

// FetchPublicKey fetches the RSA public key from auth-service
// Returns cached key if it was fetched less than 1 hour ago
func FetchPublicKey(authServiceURL string) (*rsa.PublicKey, error) {
	keyCache.mu.RLock()
	if keyCache.Key != nil && time.Since(keyCache.FetchedAt) < 1*time.Hour {
		key := keyCache.Key
		keyCache.mu.RUnlock()
		return key, nil
	}
	keyCache.mu.RUnlock()

	// Fetch new key
	url := fmt.Sprintf("%s/api/public-key", authServiceURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth-service returned status %d", resp.StatusCode)
	}

	keyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Update cache
	keyCache.mu.Lock()
	keyCache.Key = publicKey
	keyCache.FetchedAt = time.Now()
	keyCache.mu.Unlock()

	return publicKey, nil
}

// ValidateToken validates a JWT token using the RSA public key
func ValidateToken(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
