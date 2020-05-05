package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/output"
)

func TestFileProvider(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		fileContent string

		fallback   config.Provider
		fetchValue func(p *Provider) (interface{}, error)
		wantErr    error
		wantValue  interface{}
	}{
		{
			name: "fetch noAuthCache when set",

			fileContent: "---\nnoAuthCache: true\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.AuthCacheDisabled()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch noAuthCache has invalid value",

			fileContent: "---\nnoAuthCache: notabool\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.AuthCacheDisabled()
			},

			wantValue: false,
			wantErr: &strconv.NumError{
				Func: "ParseBool",
				Num:  "notabool",
				Err:  strconv.ErrSyntax,
			},
		},
		{
			name: "fetch noAuthCache falls back when not set",

			fileContent: "---\n",
			fallback: &mockProvider{
				GetAuthCacheDisabled: true,
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.AuthCacheDisabled()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch noAuthCache fall back errors",

			fileContent: "---\n",
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.AuthCacheDisabled()
			},

			wantValue: false,
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch api endpoints when set",

			fileContent: "---\nendpoints:\n  - 1.1.1.1:5705\n  - 2.2.2.2:5705\n",
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
			name: "fetch api endpoints falls back when not set",

			fileContent: "---\n",
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
			name: "fetch api endpoints falls back when empty",

			fileContent: "---\nendpoints: []\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.APIEndpoints()
			},

			wantValue: ([]string)(nil),
			wantErr:   errMissingEndpoints,
		},
		{
			name: "fetch api endpoints fall back errors",

			fileContent: "---\n",
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

			fileContent: "---\ntimeout: 10s\n",
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
			name: "fetch cache dir when set",

			fileContent: "---\ncacheDir: my-cacheDir\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CacheDir()
			},

			wantValue: "my-cacheDir",
			wantErr:   nil,
		},
		{
			name: "fetch cache dir falls back when not set",

			fileContent: "---\n",
			fallback: &mockProvider{
				GetCacheDir: "my-cacheDir",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CacheDir()
			},

			wantValue: "my-cacheDir",
			wantErr:   nil,
		},
		{
			name: "fetch cache dir fall back errors",

			fileContent: "---\n",
			fallback: &mockProvider{
				GetError: errors.New("bananas"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CacheDir()
			},

			wantValue: "",
			wantErr:   errors.New("bananas"),
		},
		{
			name: "fetch command timeout falls back when not set",

			fileContent: "---\n",
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

			fileContent: "---\n",
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

			fileContent: "---\nusername: jeff\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Username()
			},

			wantValue: "jeff",
			wantErr:   nil,
		},
		{
			name: "fetch username falls back when not set",

			fileContent: "---\n",
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

			fileContent: "---\n",
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
			name: "fetch password returns error when password set",

			fileContent: "---\npassword: verysecret\n",
			fallback: &mockProvider{
				GetPassword: "somedefault",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "",
			wantErr:   errPasswordForbidden,
		},
		{
			name: "fetch password falls back when not set",

			fileContent: "---\n",
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

			fileContent: "---\n",
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
			name: "fetch useIds when set",

			fileContent: "---\nuseIds: true\n",
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
			name: "fetch useIds has invalid value",

			fileContent: "---\nuseIds: notabool\n",
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
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
			name: "fetch useIds falls back when not set",

			fileContent: "---\n",
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
			name: "fetch useIds fall back errors",

			fileContent: "---\n",
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

			fileContent: "---\nnamespace: my-namespace\n",
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

			fileContent: "---\n",
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

			fileContent: "---\n",
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
			name:        "fetch output format when set",
			fileContent: "---\noutput: text\n",

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
			name:        "fetch output format has invalid value",
			fileContent: "---\noutput: notAFormat\n",
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
			name:        "fetch output format falls back when not set",
			fileContent: "---\n",
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
			name:        "fetch output format fall back errors",
			fileContent: "---\n",
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
			name:        "fetch config file path falls back when not set",
			fileContent: "---\n",
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
			name:        "fetch config file path fall back errors",
			fileContent: "---\n",
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

	dir, err := ioutil.TempDir("", "storageos")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			f, err := ioutil.TempFile(dir, "config")
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(f.Name())

			_, err = f.WriteString(tt.fileContent)
			if err != nil {
				t.Fatal(err)
			}
			f.Close()

			p := NewProvider(tt.fallback)

			p.SetConfigProvider(&mockProvider{GetConfigFilePath: f.Name()})

			// Attempt to fetch the value from the config file provider
			gotValue, gotErr := tt.fetchValue(p)

			switch gotErr.(type) {
			case parseError:
				if !reflect.DeepEqual(errors.Unwrap(gotErr), tt.wantErr) {
					t.Errorf("got unwrapped error %v, want %v", errors.Unwrap(gotErr), tt.wantErr)
				}
			default:
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
				}
			}

			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("got config value %v (%T), want %v (%T)", gotValue, gotValue, tt.wantValue, tt.wantValue)
			}
		})
	}
}
