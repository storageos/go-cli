package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/apiclient"
)

type mockConfigProvider struct {
	APIEndpointsGiveValue []string
	APIEndpointsGiveErr   error

	CacheDirGiveValue string
	CacheDirGiveErr   error
}

func (m *mockConfigProvider) APIEndpoints() ([]string, error) {
	return m.APIEndpointsGiveValue, m.APIEndpointsGiveErr
}

func (m *mockConfigProvider) CacheDir() (string, error) {
	return m.CacheDirGiveValue, m.CacheDirGiveErr
}

func TestSessionCacheGet(t *testing.T) {
	t.Parallel()

	var (
		mockTime     = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		timeProvider = func() time.Time {
			return mockTime
		}
		fuzzFactor = time.Second

		newTmpDir = func() string {
			dir, err := ioutil.TempDir("", "storageos-cli-test-*")
			if err != nil {
				t.Fatalf("failed to set up test case: %v", err)
			}
			return dir
		}
	)

	tests := []struct {
		name string

		config        *mockConfigProvider
		setupCacheDir func(t *testing.T, forCache *SessionCache) error

		lookupUsername string

		wantSession apiclient.AuthSession
		wantErr     error
	}{
		{
			name: "valid session exists for username",

			config: &mockConfigProvider{
				APIEndpointsGiveValue: []string{"endpoint-a", "endpoint-b"},
				CacheDirGiveValue:     newTmpDir(),
			},
			setupCacheDir: func(t *testing.T, forCache *SessionCache) error {
				session := apiclient.NewAuthSession(
					"i am token",
					mockTime.Add(2*fuzzFactor),
				)

				path, err := forCache.getFilepath("i-am-user")
				if err != nil {
					return err
				}

				f, err := os.Create(path)
				if err != nil {
					return err
				}
				defer f.Close()

				enc := json.NewEncoder(f)
				return enc.Encode(&session)
			},

			lookupUsername: "i-am-user",

			wantSession: apiclient.NewAuthSession(
				"i am token",
				mockTime.Add(2*fuzzFactor),
			),
			wantErr: nil,
		},
		{
			name: "expired session exists for username",

			config: &mockConfigProvider{
				APIEndpointsGiveValue: []string{"endpoint-a", "endpoint-b"},
				CacheDirGiveValue:     newTmpDir(),
			},
			setupCacheDir: func(t *testing.T, forCache *SessionCache) error {
				session := apiclient.NewAuthSession(
					"i am token",
					mockTime,
				)

				path, err := forCache.getFilepath("i-am-user")
				if err != nil {
					return err
				}

				f, err := os.Create(path)
				if err != nil {
					return err
				}
				defer f.Close()

				enc := json.NewEncoder(f)
				return enc.Encode(&session)
			},

			lookupUsername: "i-am-user",

			wantSession: apiclient.AuthSession{},
			wantErr:     errors.New("cached session has expired"),
		},
		{
			name: "no session for username",

			config: &mockConfigProvider{
				APIEndpointsGiveValue: []string{"endpoint-a", "endpoint-b"},
				CacheDirGiveValue:     newTmpDir(),
			},
			setupCacheDir: func(t *testing.T, forCache *SessionCache) error {
				return nil
			},

			lookupUsername: "i-am-not-user",

			wantSession: apiclient.AuthSession{},
			wantErr:     errors.New("no cached session for username with current config"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if strings.HasPrefix(tt.config.CacheDirGiveValue, os.TempDir()) {
				defer os.RemoveAll(tt.config.CacheDirGiveValue)
			}

			cache := NewSessionCache(tt.config, timeProvider, fuzzFactor)

			err := tt.setupCacheDir(t, cache)
			if err != nil {
				t.Fatalf("failed to set-up test case: %v", err)
			}

			gotToken, gotErr := cache.Get(tt.lookupUsername)
			if gotToken != tt.wantSession {
				t.Errorf("got token %v, want %v", gotToken, tt.wantSession)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got err %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestSessionCachePut(t *testing.T) {
	t.Parallel()

	var (
		mockTime     = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		timeProvider = func() time.Time {
			return mockTime
		}
		fuzzFactor = time.Second

		newTmpDir = func() string {
			dir, err := ioutil.TempDir("", "storageos-cli-test-*")
			if err != nil {
				t.Fatalf("failed to set up test case: %v", err)
			}
			return dir
		}
	)

	tests := []struct {
		name string

		config *mockConfigProvider

		checkUsername string
		session       apiclient.AuthSession

		wantAuthSession apiclient.AuthSession
		wantErr         error
	}{
		{
			name: "ok",

			config: &mockConfigProvider{
				APIEndpointsGiveValue: []string{"endpoint-a", "endpoint-b"},
				CacheDirGiveValue:     newTmpDir(),
			},

			checkUsername: "i-am-user",
			session: apiclient.NewAuthSession(
				"i am token",
				mockTime.Add(5*time.Second),
			),

			wantAuthSession: apiclient.NewAuthSession(
				"i am token",
				mockTime.Add(5*time.Second),
			),
			wantErr: nil,
		},
		{
			name: "empty username rejected",

			config: &mockConfigProvider{
				APIEndpointsGiveValue: []string{"endpoint-a", "endpoint-b"},
				CacheDirGiveValue:     newTmpDir(),
			},

			checkUsername: "",
			session: apiclient.NewAuthSession(
				"i am token",
				mockTime.Add(5*time.Second),
			),

			wantAuthSession: apiclient.AuthSession{},
			wantErr:         errors.New("username required for caching"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := NewSessionCache(tt.config, timeProvider, fuzzFactor)

			gotErr := cache.Put(tt.checkUsername, tt.session)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if tt.wantAuthSession != (apiclient.AuthSession{}) {
				path, err := cache.getFilepath(tt.checkUsername)
				if err != nil {
					t.Fatalf("unexpected failure checking results: %v", err)
				}

				f, err := os.Open(path)
				if err != nil {
					t.Fatalf("unexpected failure checking results: %v", err)
				}
				defer os.RemoveAll(tt.config.CacheDirGiveValue)

				dec := json.NewDecoder(f)

				var gotAuthSession apiclient.AuthSession
				dec.Decode(&gotAuthSession)

				if !reflect.DeepEqual(gotAuthSession, tt.wantAuthSession) {
					pretty.Ldiff(t, gotAuthSession, tt.wantAuthSession)
					t.Errorf("for username %v got session %v, want %v", tt.checkUsername, gotAuthSession, tt.wantAuthSession)
				}
			}
		})
	}
}
