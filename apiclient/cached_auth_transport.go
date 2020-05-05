package apiclient

import (
	"context"
)

// AuthCache describes a type that provides a key-value cache mapping usernames
// to cached auth sessions.
type AuthCache interface {
	Get(forUsername string) (AuthSession, error)
	Put(username string, session AuthSession) error
}

// AuthCachedTransport decorates the embedded transport with cached auth
// capabilities, modifying the Authenticate() behaviour. An instance will
// try to authenticate from the cache if the requested username is not the
// previous username requested for auth. If that fails, it will continue as
// normal, caching the result.
type AuthCachedTransport struct {
	Transport

	cache              AuthCache
	lastAuthedUsername string
}

// Authenticate attempts to fetch an authentication session for username from
// its cache, unless it is a repeated attempt. If the cache lookup fails or the
// call is a re-attempted auth then this authenticates as normal but caches
// the result.
func (co *AuthCachedTransport) Authenticate(ctx context.Context, username, password string) (AuthSession, error) {
	if username != co.lastAuthedUsername {
		session, err := co.authenticateFromCache(ctx, username)
		if err == nil {
			return session, nil
		}
	}

	session, err := co.Transport.Authenticate(
		ctx,
		username,
		password,
	)
	if err != nil {
		return session, err
	}

	_ = co.cache.Put(username, session)
	co.lastAuthedUsername = username

	return session, nil
}

// authenticateFromCache sets the inner transport to use the auth session found
// for username, if present.
func (co *AuthCachedTransport) authenticateFromCache(ctx context.Context, username string) (AuthSession, error) {
	session, err := co.cache.Get(username)
	if err != nil {
		return AuthSession{}, err
	}

	return session, co.UseAuthSession(ctx, session)
}

// NewAuthCachedTransport decorates inner with use of an authentication session
// cache.
func NewAuthCachedTransport(inner Transport, cache AuthCache) *AuthCachedTransport {
	return &AuthCachedTransport{
		Transport: inner,
		cache:     cache,
	}
}
