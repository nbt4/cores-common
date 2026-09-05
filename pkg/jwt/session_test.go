package jwt

import (
	"context"
	"errors"
	"testing"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

func TestSessionUsesCurrentAccount(t *testing.T) {
	secret := []byte("session-test-secret-at-least-32-bytes")
	raw, err := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, Claims{
		UserID: 42, Username: "old-name", IsAdmin: true,
		RegisteredClaims: golangjwt.RegisteredClaims{ExpiresAt: golangjwt.NewNumericDate(time.Now().Add(time.Hour))},
	}).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}
	user := User{ID: 42, Username: "current-name", IsActive: true, IsAdmin: true}
	var lookupErr error
	lookup := func(ctx context.Context, id uint) (User, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Error("account lookup has no deadline")
		}
		if id != 42 {
			t.Fatalf("lookup id = %d", id)
		}
		return user, lookupErr
	}
	claims, ok := ValidateSession(context.Background(), raw, secret, lookup)
	if !ok || !claims.IsAdmin || claims.Username != "current-name" {
		t.Fatal("active account rejected or stale claims returned")
	}
	user.IsAdmin = false
	claims, ok = ValidateSession(context.Background(), raw, secret, lookup)
	if !ok || claims.IsAdmin {
		t.Fatal("revoked administrator role survived on the same token")
	}
	user.IsActive = false
	if _, ok := ValidateSession(context.Background(), raw, secret, lookup); ok {
		t.Fatal("disabled account accepted")
	}
	user.IsActive = true
	lookupErr = errors.New("database unavailable or account deleted")
	if _, ok := ValidateSession(context.Background(), raw, secret, lookup); ok {
		t.Fatal("lookup failure did not fail closed")
	}
}

func TestSessionRejectsInvalidTokensBeforeLookup(t *testing.T) {
	secret := []byte("session-test-secret-at-least-32-bytes")
	for _, tc := range []struct {
		name    string
		method  golangjwt.SigningMethod
		id      uint
		expires *golangjwt.NumericDate
	}{
		{"missing expiry", golangjwt.SigningMethodHS256, 42, nil},
		{"expired", golangjwt.SigningMethodHS256, 42, golangjwt.NewNumericDate(time.Now().Add(-time.Hour))},
		{"wrong algorithm", golangjwt.SigningMethodHS384, 42, golangjwt.NewNumericDate(time.Now().Add(time.Hour))},
		{"missing user", golangjwt.SigningMethodHS256, 0, golangjwt.NewNumericDate(time.Now().Add(time.Hour))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := golangjwt.NewWithClaims(tc.method, Claims{UserID: tc.id, RegisteredClaims: golangjwt.RegisteredClaims{ExpiresAt: tc.expires}}).SignedString(secret)
			if err != nil {
				t.Fatal(err)
			}
			_, ok := ValidateSession(context.Background(), raw, secret, func(context.Context, uint) (User, error) {
				t.Fatal("invalid token reached account lookup")
				return User{}, nil
			})
			if ok {
				t.Fatal("invalid token accepted")
			}
		})
	}
}
