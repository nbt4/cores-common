package jwt

import (
	"context"
	"fmt"
	"net/http"
	"os"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

// Claims mirrors the cores-dashboard JWT claims structure used across all services.
type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	golangjwt.RegisteredClaims
}

// Context key for storing claims in request context.
type contextKey string

const ClaimsKey = contextKey("claims")

// JWTSecret returns the shared JWT secret. Checks CORES_JWT_SECRET first, then JWT_SECRET.
// Panics if neither is set, as the system cannot function without a shared secret.
func JWTSecret() []byte {
	if s := os.Getenv("CORES_JWT_SECRET"); s != "" {
		return []byte(s)
	}
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	panic("FATAL: CORES_JWT_SECRET or JWT_SECRET environment variable is required")
}

// ValidateToken parses and validates a JWT token string. Returns the claims and true if valid.
func ValidateToken(tokenString string) (*Claims, bool) {
	claims := &Claims{}
	token, err := golangjwt.ParseWithClaims(tokenString, claims, func(t *golangjwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*golangjwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return JWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
}

// GetClaims extracts JWT claims from the request context.
func GetClaims(r *http.Request) (*Claims, bool) {
	c, ok := r.Context().Value(ClaimsKey).(*Claims)
	return c, ok
}

// SetClaims stores JWT claims in the request context.
func SetClaims(r *http.Request, claims *Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims))
}
