package apiclient

import (
	"errors"
	"time"
)

// ErrAuthSessionNoToken is an error value indicating that a given auth session
// does not have an associated token.
var ErrAuthSessionNoToken = errors.New("auth session has no token")

// AuthSession encapsulates re-usable auth session.
type AuthSession struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// NewAuthSession initialises a re-usable auth session.
func NewAuthSession(token string, expiresAt time.Time) AuthSession {
	return AuthSession{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
