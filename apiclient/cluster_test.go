package apiclient

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

func TestUpdateLicence(t *testing.T) {
	t.Parallel()

	var (
		mockExpiry     = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
		newMockLicence = func() *cluster.Licence {
			return &cluster.Licence{
				ClusterID:            "some-cluster-id",
				ExpiresAt:            mockExpiry,
				ClusterCapacityBytes: 42,
				Kind:                 "mock-licence",
				CustomerName:         "go-test",
			}
		}
	)

	tests := []struct {
		name string

		configProvider *mockConfigProvider
		transport      *mockTransport

		opaqueKey []byte

		wantLicence *cluster.Licence
		wantErr     error
	}{
		{
			name: "ok",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				GetClusterResource: &cluster.Resource{
					ID:                    "some-cluster-id",
					Licence:               newMockLicence(),
					DisableTelemetry:      true,
					DisableCrashReporting: true,
					DisableVersionCheck:   true,

					LogLevel:  "info",
					LogFormat: "default",

					CreatedAt: mockExpiry,
					UpdatedAt: mockExpiry,
					Version:   version.Version("42"),
				},
				UpdateClusterResource: &cluster.Resource{
					ID:                    "some-cluster-id",
					Licence:               newMockLicence(),
					DisableTelemetry:      true,
					DisableCrashReporting: true,
					DisableVersionCheck:   true,

					LogLevel:  "info",
					LogFormat: "default",

					CreatedAt: mockExpiry,
					UpdatedAt: mockExpiry,
					Version:   version.Version("42"),
				},
			},

			opaqueKey: []byte("a licence key"),

			wantLicence: newMockLicence(),
			wantErr:     nil,
		},
		{
			name: "fails to fetch current cluster configuration",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				GetClusterError: errors.New("oh bananas"),
			},

			opaqueKey: nil, // Shouldn't attempt to update without knowing current settings

			wantLicence: nil,
			wantErr:     errors.New("oh bananas"),
		},
		{
			name: "fails to apply new licence key",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				GetClusterResource: &cluster.Resource{
					ID:                    "some-cluster-id",
					Licence:               newMockLicence(),
					DisableTelemetry:      true,
					DisableCrashReporting: true,
					DisableVersionCheck:   true,

					LogLevel:  "info",
					LogFormat: "default",

					CreatedAt: mockExpiry,
					UpdatedAt: mockExpiry,
					Version:   version.Version("42"),
				},

				UpdateClusterError: errors.New("oh bananas"),
			},

			opaqueKey: []byte("a licence key"), // Fail updating the cluster

			wantLicence: nil,
			wantErr:     errors.New("oh bananas"),
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

			gotLicence, gotErr := client.UpdateLicence(context.Background(), tt.opaqueKey)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotLicence, tt.wantLicence) {
				pretty.Ldiff(t, gotLicence, tt.wantLicence)
				t.Errorf("got returned licence %v, want %v", pretty.Sprint(gotLicence), pretty.Sprint(tt.wantLicence))
			}

			if !reflect.DeepEqual(tt.transport.UpdateClusterGotResource, tt.transport.GetClusterResource) {
				pretty.Ldiff(t, tt.transport.UpdateClusterGotResource, tt.transport.GetClusterResource)
				t.Errorf("gave cluster config %v, wanted %v", tt.transport.UpdateClusterGotResource, tt.transport.GetClusterResource)
			}

			// Should pass in the key
			if !reflect.DeepEqual(tt.transport.UpdateClusterGotLicenceKey, tt.opaqueKey) {
				t.Errorf("transport got key %s, want %s", string(tt.transport.UpdateClusterGotLicenceKey), string(tt.opaqueKey))
			}
		})
	}
}
