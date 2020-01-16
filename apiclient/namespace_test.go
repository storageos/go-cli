package apiclient

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"code.storageos.net/storageos/c2-cli/namespace"
)

func TestGetNamespaceByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		configProvider *mockConfigProvider
		transport      *mockTransport

		namespaceName string

		wantResource *namespace.Resource
		wantErr      error
	}{
		{
			name: "ok",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				ListNamespacesResource: []*namespace.Resource{
					&namespace.Resource{
						Name: "possibly-dave",
					},
					&namespace.Resource{
						Name: "definitely-steve",
					},
				},
			},

			namespaceName: "definitely-steve",

			wantResource: &namespace.Resource{
				Name: "definitely-steve",
			},
			wantErr: nil,
		},
		{
			name: "namespace with name does not exist",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				ListNamespacesResource: []*namespace.Resource{
					&namespace.Resource{
						Name: "possibly-dave",
					},
					&namespace.Resource{
						Name: "not-steve",
					},
				},
			},

			namespaceName: "definitely-steve",

			wantResource: nil,
			wantErr: NamespaceNotFoundError{
				name: "definitely-steve",
			},
		},
		{
			name: "error getting list of namespaces",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				ListNamespacesError: errors.New("bananas"),
			},

			namespaceName: "a-namespace",

			wantResource: nil,
			wantErr:      errors.New("bananas"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := New(tt.configProvider)
			if err := client.ConfigureTransport(tt.transport); err != nil {
				t.Fatalf("got error configuring client transport: %v", err)
			}

			gotResource, gotErr := client.GetNamespaceByName(context.Background(), tt.namespaceName)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				t.Errorf("got namespace resource %v, want %v", gotResource, tt.wantResource)
			}
		})
	}
}
