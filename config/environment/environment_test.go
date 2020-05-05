package environment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"code.storageos.net/storageos/c2-cli/output"
)

func TestEnvironmentProvider(t *testing.T) {
	t.Parallel()

	createShellScript := func(cmd string) string {
		file, err := ioutil.TempFile("", "storageos_cli_pwd_cmd.*.sh")
		if err != nil {
			t.Fatalf("error creating an empty file for password sourcing command")
		}
		defer file.Close()

		_, err = file.WriteString("#!/bin/sh\n" + cmd)
		if err != nil {
			t.Fatalf("error writing content in password sourcing command")
		}

		err = file.Chmod(0777)
		if err != nil {
			t.Fatalf("error on changing permissions on password sourcing command")
		}

		return file.Name()
	}

	defer func() {
		files, err := filepath.Glob(os.TempDir() + "/storageos_cli_pwd_cmd.*")
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}
	}()

	tests := []struct {
		name string

		setenv     func(t *testing.T)
		fallback   *mockProvider
		fetchValue func(p *Provider) (interface{}, error)

		wantValue interface{}
		wantErr   error
	}{
		{
			name: "fetch auth cache disabled when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(AuthCacheDisabledVar, "true")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.AuthCacheDisabled()
			},

			wantValue: true,
			wantErr:   nil,
		},
		{
			name: "fetch auth cache disabled has invalid value",

			setenv: func(t *testing.T) {
				err := os.Setenv(AuthCacheDisabledVar, "notabool")
				if err != nil {
					t.Fatalf("got unexpected error setting up environment: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
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
			name: "fetch auth cache disabled falls back when not set",

			setenv: func(t *testing.T) {},
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
			name: "fetch auth cache disabled fall back errors",

			setenv: func(t *testing.T) {},
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

			setenv: func(t *testing.T) {
				err := os.Setenv(APIEndpointsVar, "1.1.1.1:5705,2.2.2.2:5705")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
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
			name: "fetch cache dir when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(CacheDirVar, "/tmp/.cache")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("dont call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CacheDir()
			},

			wantValue: "/tmp/.cache",
			wantErr:   nil,
		},
		{
			name: "fetch cache dir falls back when not set",

			setenv: func(t *testing.T) {},
			fallback: &mockProvider{
				GetCacheDir: "/tmp/.cache",
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.CacheDir()
			},

			wantValue: "/tmp/.cache",
			wantErr:   nil,
		},
		{
			name: "fetch cache dir fall back errors",

			setenv: func(t *testing.T) {},
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
			name: "fetch command timeout when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(CommandTimeoutVar, "10s")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
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
			name: "fetch password with password command",

			setenv: func(t *testing.T) {
				path := createShellScript("echo 'verysecret'")

				err := os.Setenv(PasswordCommandVar, path)
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
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
			name: "fetch password with failing password command",

			setenv: func(t *testing.T) {
				path := createShellScript("exit 42")

				err := os.Setenv(PasswordCommandVar, path)
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
			},
			fallback: &mockProvider{
				GetError: errors.New("don't call me"),
			},
			fetchValue: func(p *Provider) (interface{}, error) {
				return p.Password()
			},

			wantValue: "",
			wantErr:   fmt.Errorf("password command exited with error code 42"),
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
				GetError: errors.New("don't call me"),
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

		{
			name: "fetch output format when set",

			setenv: func(t *testing.T) {
				err := os.Setenv(OutputFormatVar, "text")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
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

			setenv: func(t *testing.T) {
				err := os.Setenv(OutputFormatVar, "notAFormat")
				if err != nil {
					t.Fatalf("got unexpected error setting up environment: %v", err)
				}
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
			name: "fetch output format falls back when not set",

			setenv: func(t *testing.T) {},
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
			name: "fetch output format fall back errors",

			setenv: func(t *testing.T) {},
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

			setenv: func(t *testing.T) {
				err := os.Setenv(ConfigFilePathVar, ".test_config_file")
				if err != nil {
					t.Fatalf("got unexpected error setting up env: %v", err)
				}
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
			name: "fetch config file path falls back when not set",

			setenv: func(t *testing.T) {},
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
			name: "fetch config file path fall back errors",

			setenv: func(t *testing.T) {},
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
			// Handle the environment setup tear down
			os.Clearenv()
			defer os.Clearenv()
			tt.setenv(t)

			p := NewProvider(tt.fallback)

			// Attempt to fetch the value from the env provider
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
