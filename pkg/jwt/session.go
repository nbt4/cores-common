package jwt

import (
	"context"
	"database/sql"
	"errors"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

// User contains only the current account data needed to authorize a request.
type User struct {
	ID       uint
	Username string
	IsActive bool
	IsAdmin  bool
}

type UserLookup func(context.Context, uint) (User, error)

// DatabaseUserLookup deliberately does not cache accounts: disabling an account
// or changing its role must take effect on the next authenticated request.
func DatabaseUserLookup(db *sql.DB) UserLookup {
	return func(ctx context.Context, id uint) (User, error) {
		if db == nil {
			return User{}, errors.New("account database unavailable")
		}
		var user User
		err := db.QueryRowContext(ctx, `SELECT userid, username, is_active, is_admin FROM users WHERE userid = $1`, id).
			Scan(&user.ID, &user.Username, &user.IsActive, &user.IsAdmin)
		return user, err
	}
}

// ValidateSession verifies the suite token and replaces mutable claims with the
// current database values. Database failures, deleted and inactive users fail closed.
func ValidateSession(ctx context.Context, raw string, secret []byte, lookup UserLookup) (*Claims, bool) {
	if len(secret) == 0 || lookup == nil {
		return nil, false
	}
	claims := &Claims{}
	token, err := golangjwt.ParseWithClaims(raw, claims, func(*golangjwt.Token) (any, error) {
		return secret, nil
	}, golangjwt.WithValidMethods([]string{"HS256"}), golangjwt.WithExpirationRequired())
	if err != nil || !token.Valid || claims.UserID == 0 {
		return nil, false
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	user, err := lookup(ctx, claims.UserID)
	if err != nil || !user.IsActive || user.ID != claims.UserID {
		return nil, false
	}
	claims.Username, claims.IsAdmin = user.Username, user.IsAdmin
	return claims, true
}
