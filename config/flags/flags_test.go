package flags

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func TestFlagProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		arguments  []string
		fallback   *mockProvider
		fetchValue func(p *Provider) (interface{}, error)

		wantValue interface{}
		wantErr   error
	}{
		{
			name: "fetch api endpoints config when set",

			arguments: []string{
				"--" + APIEndpointsFlag, "1.1.1.1:5705",
				"--" + APIEndpointsFlag, "2.2.2.2:5705",
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
			name: "fetch api endpoints config falls back when not set",

			arguments: []string{},
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
			name: "fetch api endpoints config fall back errors",

			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.APIEndpoints()
			},

			wantValue: ([]string)(nil), // Typed nils are a stupid thing.
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch command timeout when set",

			arguments: []string{
				"--" + CommandTimeoutFlag, "10s",
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

			arguments: []string{},
			fallback: &mockProvider{
				GetCommandTimeout: 5 * time.Second,
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CommandTimeout()
			},

			wantValue: 5 * time.Second,
			wantErr:   nil,
		},
		{
			name: "fetch command timeout fall back errors",

			arguments: []string{},
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

			arguments: []string{
				"--" + UsernameFlag, "mange",
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "mange",
			wantErr:   nil,
		},
		{
			name: "fetch username falls back when not set",

			arguments: []string{},
			fallback: &mockProvider{
				GetUsername: "mange",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "mange",
			wantErr:   nil,
		},
		{
			name: "fetch username fall back errors",

			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("fest hos mange"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "",
			wantErr:   errors.New("fest hos mange"),
		},
		{
			name: "fetch password when set",

			arguments: []string{
				"--" + PasswordFlag, "mangemakers",
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "mangemakers",
			wantErr:   nil,
		},
		{
			name: "fetch password falls back when not set",

			arguments: []string{},
			fallback: &mockProvider{
				GetPassword: "mangemakers",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "mangemakers",
			wantErr:   nil,
		},
		{
			name: "fetch password falls back errors",

			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flagset := pflag.NewFlagSet(tt.name, pflag.ContinueOnError)

			p := NewProvider(
				flagset,
				tt.fallback,
			)

			// Parse the provided flags
			parseErr := flagset.Parse(tt.arguments)
			if parseErr != nil {
				t.Fatalf("got parse error %v", parseErr)
			}

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
