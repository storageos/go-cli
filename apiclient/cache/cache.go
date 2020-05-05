package cache

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/atomicfile"
)

// ConfigProvider defines the configuration settings required by the
// SessionCache.
type ConfigProvider interface {
	APIEndpoints() ([]string, error)
	CacheDir() (string, error)
}

// SessionCache implements a lookup cache which maps usernames to auth sessions,
// loosely tracking expected expiration times.
type SessionCache struct {
	config ConfigProvider

	currentTime func() time.Time
	fuzzFactor  time.Duration

	mu *sync.RWMutex
}

// Get looks up the session cache file forUsername. If there is a valid
// session it is returned.
//
// A cached session is considered valid if the time on the system clock is
// before the estimated expiration date, after accounting for the configured
// fuzz factor.
func (c *SessionCache) Get(forUsername string) (apiclient.AuthSession, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get the expected filepath for the username
	lookupPath, err := c.getFilepath(forUsername)
	if err != nil {
		return apiclient.AuthSession{}, err
	}

	f, err := os.Open(lookupPath)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		return apiclient.AuthSession{}, errors.New("no cached session for username with current config")
	default:
		return apiclient.AuthSession{}, err
	}

	// Decode session from file
	decoder := json.NewDecoder(f)

	var session apiclient.AuthSession
	err = decoder.Decode(&session)
	if err != nil {
		return apiclient.AuthSession{}, err
	}

	// Check if the token is likely to have expired
	fuzzyExpiration := session.ExpiresAt.Add(-c.fuzzFactor)

	if c.currentTime().After(fuzzyExpiration) {
		return apiclient.AuthSession{}, errors.New("cached session has expired")
	}

	return session, nil
}

// Put stores the given session in the cache for lookup by username,
// atomically writing the cached session to disk.
func (c *SessionCache) Put(username string, session apiclient.AuthSession) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if username == "" {
		return errors.New("username required for caching")
	}

	// Get the target file path to write the cached session to
	targetPath, err := c.getFilepath(username)
	if err != nil {
		return err
	}

	txn, err := atomicfile.NewWrite(targetPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(txn)

	// Try to encode the session to the cache file
	if err = encoder.Encode(session); err != nil {
		// Error results in abort
		_ = txn.Abort()
		return err
	}

	// Otherwise commit
	return txn.Commit()
}

func (c *SessionCache) getFilepath(forUsername string) (string, error) {
	cacheDir, err := c.config.CacheDir()
	if err != nil {
		return "", err
	}

	endpoints, err := c.config.APIEndpoints()
	if err != nil {
		return "", err
	}

	// Sort the endpoints for consistency
	sort.Strings(endpoints)

	// Hash the username and the target endpoints to get cache file name
	hash := sha512.New()

	_, _ = hash.Write([]byte(forUsername))
	for _, endpoint := range endpoints {
		_, _ = hash.Write([]byte(endpoint))
	}

	return filepath.Join(
		cacheDir,                         // dir
		fmt.Sprintf("%x", hash.Sum(nil)), // file name
	), nil
}

// NewSessionCache instantiates an empty session cache.
func NewSessionCache(config ConfigProvider, timeProvider func() time.Time, fuzzFactor time.Duration) *SessionCache {
	return &SessionCache{
		config: config,

		currentTime: timeProvider,
		fuzzFactor:  fuzzFactor,

		mu: &sync.RWMutex{},
	}
}
