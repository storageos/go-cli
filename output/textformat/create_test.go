package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

func TestDisplayer_CreateUser(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		user    *output.User
		wantW   string
		wantErr error
	}{
		{
			name: "display created user with single group ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS           
banana-name  admin  donkeys years ago  policy-group-name
`,
			wantErr: nil,
		},
		{
			name: "display created user with multiple groups ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
					{
						ID:   "policy-group-id-2",
						Name: "policy-group-name-2",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS                               
banana-name  admin  donkeys years ago  policy-group-name,policy-group-name-2
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.CreateUser(context.Background(), w, tt.user)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_CreateNamespace(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		resource *namespace.Resource
		wantW    string
		wantErr  error
	}{
		{
			name: "create namespace",
			resource: &namespace.Resource{
				ID:        "bananaID",
				Name:      "bananaName",
				Labels:    map[string]string{"bananaKey": "bananaValue"},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME        AGE              
bananaName  donkeys years ago
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.CreateNamespace(context.Background(), w, output.NewNamespace(tt.resource))
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_CreatePolicyGroup(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		resource   *policygroup.Resource
		namespaces []*namespace.Resource
		wantW      string
		wantErr    error
	}{
		{
			name: "create policy group",
			resource: &policygroup.Resource{
				ID:    "bananaID",
				Name:  "bananaName",
				Users: []*policygroup.Member{},
				Specs: []*policygroup.Spec{
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "volume",
						ReadOnly:     false,
					},
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "*",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			namespaces: []*namespace.Resource{
				{
					ID:        "banana-namespace-id",
					Name:      "banana-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "kiwi-namespace-id",
					Name:      "kiwi-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "pineapple-namespace-id",
					Name:      "pineapple-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME        USERS  SPECS  AGE              
bananaName  0      2      donkeys years ago
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.CreatePolicyGroup(context.Background(), w, output.NewPolicyGroup(tt.resource, tt.namespaces))
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}
