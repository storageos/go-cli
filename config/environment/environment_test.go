package environment

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestEnvironmentProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		setenv     func(t *testing.T)
		fallback   *mockProvider
		fetchValue func(p *Provider) (interface{}, error)

		wantValue interface{}
		wantErr   error
	}{
		{
			name: "fetch api endpoints when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(APIEndpointsVar, "1.1.1.1:5705,2.2.2.2:5705")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.APIEndpoints()
			},

			wantValue: []string{"1.1.1.1:5705", "2.2.2.2:5705"},
			wantErr:   nil,
		},
		{
			name: "fetch api endpoints falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetAPIEndpoints: []string{"1.1.1.1:5705", "2.2.2.2:5705"},
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.APIEndpoints()
			},

			wantValue: []string{"1.1.1.1:5705", "2.2.2.2:5705"},
			wantErr:   nil,
		},
		{
			name: "fetch api endpoints fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.APIEndpoints()
			},

			wantValue: ([]string)(nil),
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch command timeout when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(CommandTimeoutVar, "10s")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CommandTimeout()
			},

			wantValue: 10 * time.Second,
			wantErr:   nil,
		},
		{
			name: "fetch command timeout falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetCommandTimeout: time.Second,
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CommandTimeout()
			},

			wantValue: time.Second,
			wantErr:   nil,
		},
		{
			name: "fetch command timeout fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CommandTimeout()
			},

			wantValue: time.Duration(0),
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch username when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(UsernameVar, "jeff")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "jeff",
			wantErr:   nil,
		},
		{
			name: "fetch username falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetUsername: "username",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "username",
			wantErr:   nil,
		},
		{
			name: "fetch username fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch password when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(PasswordVar, "verysecret")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "verysecret",
			wantErr:   nil,
		},
		{
			name: "fetch password falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetPassword: "somedefault",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "somedefault",
			wantErr:   nil,
		},
		{
			name: "fetch password fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch use-ids when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(UseIDsVar, "true")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.UseIDs()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch use-ids has invalid value",

			setenv: func(t *testing.T) {
				err := os.Setenv(UseIDsVar, "notabool")
				if err != nil {
					t.Fatalf("got unexpected error setting up environment: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.UseIDs()
			},

			wantValue: false,
			wantErr: &strconv.NumError{
				Func: "ParseBool",
				Num:  "notabool",
				Err:  strconv.ErrSyntax,
			},
		},
		{
			name: "fetch use-ids falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetUseIDs: true,
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.UseIDs()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch use-ids fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.UseIDs()
			},

			wantValue: false,
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch namespace when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(NamespaceVar, "my-namespace")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Namespace()
			},

			wantValue: "my-namespace",
			wantErr:   nil,
		},
		{
			name: "fetch namespace falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetNamespace: "my-namespace",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Namespace()
			},

			wantValue: "my-namespace",
			wantErr:   nil,
		},
		{
			name: "fetch namespace fall back errors",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Namespace()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			// Handle the environment setup tear down
			os.Clearenv()
			defer os.Clearenv()
			tt.setenv(t)

			p := NewProvider(tt.fallback)

			// Attempt to fetch the value from the flag provider
			gotValue, gotErr := tt.fetchValue(p)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("got config value %v (%T), want %v (%T)", gotValue, gotValue, tt.wantValue, tt.wantValue)
			}
		})
	}
}
