package flags

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"

	"code.storageos.net/storageos/c2-cli/output"
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
				GetError: errors.New("don't call me"),
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
				GetError: errors.New("don't call me"),
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
				"--" + UsernameFlag, "a-username",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "a-username",
			wantErr:   nil,
		},
		{
			name: "fetch username falls back when not set",

			arguments: []string{},
			fallback: &mockProvider{
				GetUsername: "a-username",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "a-username",
			wantErr:   nil,
		},
		{
			name: "fetch username fall back errors",

			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("a-username"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "",
			wantErr:   errors.New("a-username"),
		},
		{
			name: "fetch password when set",

			arguments: []string{
				"--" + PasswordFlag, "verysecret",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "verysecret",
			wantErr:   nil,
		},
		{
			name: "fetch password falls back when not set",

			arguments: []string{},
			fallback: &mockProvider{
				GetPassword: "verysecret",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "verysecret",
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
		{
			name: "fetch use-ids when set",

			arguments: []string{
				"--" + UseIDsFlag, "true",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.UseIDs()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch use-ids falls back when not set",

			arguments: []string{},
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

			arguments: []string{},
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

			arguments: []string{
				"--" + NamespaceFlag, "my-namespace",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Namespace()
			},

			wantValue: "my-namespace",
			wantErr:   nil,
		},
		{
			name: "fetch namespace falls back when not set",

			arguments: []string{},
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

			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Namespace()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},

		{
			name: "fetch output format when set",
			arguments: []string{
				"--" + OutputFormatFlag, "text",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.OutputFormat()
			},

			wantValue: output.Text,
			wantErr:   nil,
		},
		{
			name: "fetch output format has invalid value",

			arguments: []string{
				"--" + OutputFormatFlag, "notAFormat",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.OutputFormat()
			},

			wantValue: output.Unknown,
			wantErr:   output.ErrInvalidFormat,
		},
		{
			name:      "fetch output format falls back when not set",
			arguments: []string{},
			fallback: &mockProvider{
				GetOutput: output.JSON,
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.OutputFormat()
			},

			wantValue: output.JSON,
			wantErr:   nil,
		},
		{
			name:      "fetch output format fall back errors",
			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.OutputFormat()
			},

			wantValue: output.Unknown,
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch config file path when set",
			arguments: []string{
				"--" + ConfigFileFlag, ".test_config_file",
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.ConfigFilePath()
			},

			wantValue: ".test_config_file",
			wantErr:   nil,
		},
		{
			name:      "fetch config file path falls back when not set",
			arguments: []string{},
			fallback: &mockProvider{
				GetConfigFilePath: "fallBack_path",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.ConfigFilePath()
			},

			wantValue: "fallBack_path",
			wantErr:   nil,
		},
		{
			name:      "fetch config file path fall back errors",
			arguments: []string{},
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.ConfigFilePath()
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
