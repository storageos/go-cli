package apiclient

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

type mockCredentialsProvider struct {
	ReturnUsername    string
	ReturnUsernameErr error

	ReturnPassword    string
	ReturnPasswordErr error
}

func (m *mockCredentialsProvider) Username() (string, error) {
	return m.ReturnUsername, m.ReturnUsernameErr
}

func (m *mockCredentialsProvider) Password() (string, error) {
	return m.ReturnPassword, m.ReturnPasswordErr
}

func TestTransportWithReauth(t *testing.T) {
	t.Parallel()

	testWrapsFunctions := []struct {
		name string

		innerTransport *mockTransport

		doTest func(t *testing.T, inner *mockTransport)

		wantInnerTransport *mockTransport
	}{
		{
			name: "GetUser",

			innerTransport: &mockTransport{
				GetUserResource: &user.Resource{},
				GetUserError:    errors.New("user-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotUser, gotErr := wrapped.GetUser(context.Background(), "user-id")
				if gotErr != inner.GetUserError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetUserError)
				}

				if !reflect.DeepEqual(gotUser, inner.GetUserResource) {
					pretty.Ldiff(t, gotUser, inner.GetUserResource)
					t.Errorf("got user %v, want %v", gotUser, inner.GetUserResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetUserGotID:    "user-id",
				GetUserResource: &user.Resource{},
				GetUserError:    errors.New("user-error"),
			},
		},
		{
			name: "GetCluster",

			innerTransport: &mockTransport{
				GetClusterResource: &cluster.Resource{},
				GetClusterError:    errors.New("cluster-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotCluster, gotErr := wrapped.GetCluster(context.Background())
				if gotErr != inner.GetClusterError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetClusterError)
				}

				if !reflect.DeepEqual(gotCluster, inner.GetClusterResource) {
					pretty.Ldiff(t, gotCluster, inner.GetClusterResource)
					t.Errorf("got cluster config %v, want %v", gotCluster, inner.GetClusterResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetClusterResource: &cluster.Resource{},
				GetClusterError:    errors.New("cluster-error"),
			},
		},
		{
			name: "GetDiagnostics",

			innerTransport: &mockTransport{
				GetDiagnosticsReadCloser: http.Response{}.Body,
				GetDiagnosticsError:      errors.New("diagnostics-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotDiagnostics, gotErr := wrapped.GetDiagnostics(context.Background())
				if gotErr != inner.GetDiagnosticsError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetDiagnosticsError)
				}

				if !reflect.DeepEqual(gotDiagnostics, inner.GetDiagnosticsReadCloser) {
					pretty.Ldiff(t, gotDiagnostics, inner.GetDiagnosticsReadCloser)
					t.Errorf("got diagnostics %v, want %v", gotDiagnostics, inner.GetDiagnosticsReadCloser)
				}
			},

			wantInnerTransport: &mockTransport{
				GetDiagnosticsReadCloser: http.Response{}.Body,
				GetDiagnosticsError:      errors.New("diagnostics-error"),
			},
		},
		{
			name: "GetNode",

			innerTransport: &mockTransport{
				GetNodeResource: &node.Resource{},
				GetNodeError:    errors.New("node-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotNode, gotErr := wrapped.GetNode(context.Background(), "node-id")
				if gotErr != inner.GetNodeError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetNodeError)
				}

				if !reflect.DeepEqual(gotNode, inner.GetNodeResource) {
					pretty.Ldiff(t, gotNode, inner.GetNodeResource)
					t.Errorf("got node %v, want %v", gotNode, inner.GetNodeResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetNodeGotID:    "node-id",
				GetNodeResource: &node.Resource{},
				GetNodeError:    errors.New("node-error"),
			},
		},
		{
			name: "GetVolume",

			innerTransport: &mockTransport{
				GetVolumeResource: &volume.Resource{},
				GetVolumeError:    errors.New("volume-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotVolume, gotErr := wrapped.GetVolume(context.Background(), "namespace-id", "volume-id")
				if gotErr != inner.GetVolumeError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetVolumeError)
				}

				if !reflect.DeepEqual(gotVolume, inner.GetVolumeResource) {
					pretty.Ldiff(t, gotVolume, inner.GetVolumeResource)
					t.Errorf("got volume %v, want %v", gotVolume, inner.GetVolumeResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetVolumeGotNamespaceID: "namespace-id",
				GetVolumeGotVolumeID:    "volume-id",
				GetVolumeResource:       &volume.Resource{},
				GetVolumeError:          errors.New("volume-error"),
			},
		},
		{
			name: "GetNamespace",

			innerTransport: &mockTransport{
				GetNamespaceResource: &namespace.Resource{},
				GetNamespaceError:    errors.New("namespace-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotNamespace, gotErr := wrapped.GetNamespace(context.Background(), "namespace-id")
				if gotErr != inner.GetNamespaceError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetNamespaceError)
				}

				if !reflect.DeepEqual(gotNamespace, inner.GetNamespaceResource) {
					pretty.Ldiff(t, gotNamespace, inner.GetNamespaceResource)
					t.Errorf("got namespace %v, want %v", gotNamespace, inner.GetNamespaceResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetNamespaceGotID:    "namespace-id",
				GetNamespaceResource: &namespace.Resource{},
				GetNamespaceError:    errors.New("namespace-error"),
			},
		},
		{
			name: "GetPolicyGroup",

			innerTransport: &mockTransport{
				GetPolicyGroupResource: &policygroup.Resource{},
				GetPolicyGroupError:    errors.New("policygroup-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotPolicyGroup, gotErr := wrapped.GetPolicyGroup(context.Background(), "policygroup-id")
				if gotErr != inner.GetPolicyGroupError {
					t.Errorf("got error %v, want %v", gotErr, inner.GetPolicyGroupError)
				}

				if !reflect.DeepEqual(gotPolicyGroup, inner.GetPolicyGroupResource) {
					pretty.Ldiff(t, gotPolicyGroup, inner.GetPolicyGroupResource)
					t.Errorf("got policy group %v, want %v", gotPolicyGroup, inner.GetPolicyGroupResource)
				}
			},

			wantInnerTransport: &mockTransport{
				GetPolicyGroupGotID:    "policygroup-id",
				GetPolicyGroupResource: &policygroup.Resource{},
				GetPolicyGroupError:    errors.New("policygroup-error"),
			},
		},
		{
			name: "ListNodes",

			innerTransport: &mockTransport{
				ListNodesResource: []*node.Resource{},
				ListNodesError:    errors.New("node-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotNodes, gotErr := wrapped.ListNodes(context.Background())
				if gotErr != inner.ListNodesError {
					t.Errorf("got error %v, want %v", gotErr, inner.ListNodesError)
				}

				if !reflect.DeepEqual(gotNodes, inner.ListNodesResource) {
					pretty.Ldiff(t, gotNodes, inner.ListNodesResource)
					t.Errorf("got nodes %v, want %v", gotNodes, inner.ListNodesResource)
				}
			},

			wantInnerTransport: &mockTransport{
				ListNodesResource: []*node.Resource{},
				ListNodesError:    errors.New("node-error"),
			},
		},
		{
			name: "ListVolumes",

			innerTransport: &mockTransport{
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"namespace-id": []*volume.Resource{},
				},
				ListVolumesError: errors.New("volume-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotVolumes, gotErr := wrapped.ListVolumes(context.Background(), "namespace-id")
				if gotErr != inner.ListVolumesError {
					t.Errorf("got error %v, want %v", gotErr, inner.ListVolumesError)
				}

				if !reflect.DeepEqual(gotVolumes, inner.ListVolumesResource["namespace-id"]) {
					pretty.Ldiff(t, gotVolumes, inner.ListVolumesResource["namespace-id"])
					t.Errorf("got volumes %v, want %v", gotVolumes, inner.ListVolumesResource["namespace-id"])
				}
			},

			wantInnerTransport: &mockTransport{
				ListVolumesGotNamespaceIDs: []id.Namespace{"namespace-id"},
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"namespace-id": []*volume.Resource{},
				},
				ListVolumesError: errors.New("volume-error"),
			},
		},
		{
			name: "ListNamespaces",

			innerTransport: &mockTransport{
				ListNamespacesResource: []*namespace.Resource{},
				ListNamespacesError:    errors.New("namespace-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotNamespaces, gotErr := wrapped.ListNamespaces(context.Background())
				if gotErr != inner.ListNamespacesError {
					t.Errorf("got error %v, want %v", gotErr, inner.ListNamespacesError)
				}

				if !reflect.DeepEqual(gotNamespaces, inner.ListNamespacesResource) {
					pretty.Ldiff(t, gotNamespaces, inner.ListNamespacesResource)
					t.Errorf("got namespaces %v, want %v", gotNamespaces, inner.ListNamespacesResource)
				}
			},

			wantInnerTransport: &mockTransport{
				ListNamespacesResource: []*namespace.Resource{},
				ListNamespacesError:    errors.New("namespace-error"),
			},
		},
		{
			name: "ListPolicyGroups",

			innerTransport: &mockTransport{
				ListPolicyGroupsResource: []*policygroup.Resource{},
				ListPolicyGroupsError:    errors.New("policygroup-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotPolicyGroups, gotErr := wrapped.ListPolicyGroups(context.Background())
				if gotErr != inner.ListPolicyGroupsError {
					t.Errorf("got error %v, want %v", gotErr, inner.ListPolicyGroupsError)
				}

				if !reflect.DeepEqual(gotPolicyGroups, inner.ListPolicyGroupsResource) {
					pretty.Ldiff(t, gotPolicyGroups, inner.ListPolicyGroupsResource)
					t.Errorf("got policy groups %v, want %v", gotPolicyGroups, inner.ListPolicyGroupsResource)
				}
			},

			wantInnerTransport: &mockTransport{
				ListPolicyGroupsResource: []*policygroup.Resource{},
				ListPolicyGroupsError:    errors.New("policygroup-error"),
			},
		},
		{
			name: "CreateUser",

			innerTransport: &mockTransport{
				CreateUserResource: &user.Resource{},
				CreateUserError:    errors.New("user-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotUser, gotErr := wrapped.CreateUser(
					context.Background(),
					"username",
					"password",
					true,
					"group-a", "group-b",
				)
				if gotErr != inner.CreateUserError {
					t.Errorf("got error %v, want %v", gotErr, inner.CreateUserError)
				}

				if !reflect.DeepEqual(gotUser, inner.CreateUserResource) {
					pretty.Ldiff(t, gotUser, inner.CreateUserResource)
					t.Errorf("got user %v, want %v", gotUser, inner.CreateUserResource)
				}
			},

			wantInnerTransport: &mockTransport{
				CreateUserGotName:     "username",
				CreateUserGotPassword: "password",
				CreateUserGotAdmin:    true,
				CreateUserGotGroups:   []id.PolicyGroup{"group-a", "group-b"},
				CreateUserResource:    &user.Resource{},
				CreateUserError:       errors.New("user-error"),
			},
		},
		{
			name: "CreateVolume",

			innerTransport: &mockTransport{
				CreateVolumeResource: &volume.Resource{},
				CreateVolumeError:    errors.New("volume-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotVolume, gotErr := wrapped.CreateVolume(
					context.Background(),
					"namespace-id",
					"name",
					"description",
					"fs-type",
					42,
					labels.Set{},
					&CreateVolumeRequestParams{},
				)
				if gotErr != inner.CreateVolumeError {
					t.Errorf("got error %v, want %v", gotErr, inner.CreateVolumeError)
				}

				if !reflect.DeepEqual(gotVolume, inner.CreateVolumeResource) {
					pretty.Ldiff(t, gotVolume, inner.CreateVolumeResource)
					t.Errorf("got volume %v, want %v", gotVolume, inner.CreateVolumeResource)
				}
			},

			wantInnerTransport: &mockTransport{
				CreateVolumeGotNamespace:   "namespace-id",
				CreateVolumeGotName:        "name",
				CreateVolumeGotDescription: "description",
				CreateVolumeGotFs:          "fs-type",
				CreateVolumeGotSizeBytes:   42,
				CreateVolumeGotLabels:      labels.Set{},
				CreateVolumeGotParams:      &CreateVolumeRequestParams{},
				CreateVolumeResource:       &volume.Resource{},
				CreateVolumeError:          errors.New("volume-error"),
			},
		},
		{
			name: "CreatePolicyGroup",

			innerTransport: &mockTransport{
				CreatePolicyGroupResource: &policygroup.Resource{},
				CreatePolicyGroupError:    errors.New("policygroup-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotPolicyGroup, gotErr := wrapped.CreatePolicyGroup(
					context.Background(),
					"name",
					[]*policygroup.Spec{
						&policygroup.Spec{
							NamespaceID:  "namespace-id",
							ResourceType: "some-resource",
							ReadOnly:     true,
						},
						&policygroup.Spec{
							NamespaceID:  "other-namespace-id",
							ResourceType: "other-resource",
							ReadOnly:     false,
						},
					},
				)
				if gotErr != inner.CreatePolicyGroupError {
					t.Errorf("got error %v, want %v", gotErr, inner.CreatePolicyGroupError)
				}

				if !reflect.DeepEqual(gotPolicyGroup, inner.CreatePolicyGroupResource) {
					pretty.Ldiff(t, gotPolicyGroup, inner.CreatePolicyGroupResource)
					t.Errorf("got policy group %v, want %v", gotPolicyGroup, inner.CreatePolicyGroupResource)
				}
			},

			wantInnerTransport: &mockTransport{
				CreatePolicyGroupGotName: "name",
				CreatePolicyGroupGotSpecs: []*policygroup.Spec{
					&policygroup.Spec{
						NamespaceID:  "namespace-id",
						ResourceType: "some-resource",
						ReadOnly:     true,
					},
					&policygroup.Spec{
						NamespaceID:  "other-namespace-id",
						ResourceType: "other-resource",
						ReadOnly:     false,
					},
				},
				CreatePolicyGroupResource: &policygroup.Resource{},
				CreatePolicyGroupError:    errors.New("policygroup-error"),
			},
		},
		{
			name: "CreateNamespace",

			innerTransport: &mockTransport{
				CreateNamespaceResource: &namespace.Resource{},
				CreateNamespaceError:    errors.New("namespace-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotNamespace, gotErr := wrapped.CreateNamespace(
					context.Background(),
					"name",
					labels.Set{},
				)
				if gotErr != inner.CreateNamespaceError {
					t.Errorf("got error %v, want %v", gotErr, inner.CreateNamespaceError)
				}

				if !reflect.DeepEqual(gotNamespace, inner.CreateNamespaceResource) {
					pretty.Ldiff(t, gotNamespace, inner.CreateNamespaceResource)
					t.Errorf("got namespace %v, want %v", gotNamespace, inner.CreateNamespaceResource)
				}
			},

			wantInnerTransport: &mockTransport{
				CreateNamespaceGotName:   "name",
				CreateNamespaceGotLabels: labels.Set{},
				CreateNamespaceResource:  &namespace.Resource{},
				CreateNamespaceError:     errors.New("namespace-error"),
			},
		},
		{
			name: "UpdateCluster",

			innerTransport: &mockTransport{
				UpdateClusterResource: &cluster.Resource{
					ID: "out",
				},
				UpdateClusterError: errors.New("cluster-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotCluster, gotErr := wrapped.UpdateCluster(
					context.Background(),
					&cluster.Resource{
						ID: "in",
					},
					[]byte("key"),
				)
				if gotErr != inner.UpdateClusterError {
					t.Errorf("got error %v, want %v", gotErr, inner.UpdateClusterError)
				}

				if !reflect.DeepEqual(gotCluster, inner.UpdateClusterResource) {
					pretty.Ldiff(t, gotCluster, inner.UpdateClusterResource)
					t.Errorf("got cluster %v, want %v", gotCluster, inner.UpdateClusterResource)
				}
			},

			wantInnerTransport: &mockTransport{
				UpdateClusterGotResource: &cluster.Resource{
					ID: "in",
				},
				UpdateClusterGotLicenceKey: []byte("key"),
				UpdateClusterResource: &cluster.Resource{
					ID: "out",
				},
				UpdateClusterError: errors.New("cluster-error"),
			},
		},
		{
			name: "DeleteUser",

			innerTransport: &mockTransport{
				DeleteUserError: errors.New("user-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.DeleteUser(
					context.Background(),
					"user-id",
					&DeleteUserRequestParams{},
				)
				if gotErr != inner.DeleteUserError {
					t.Errorf("got error %v, want %v", gotErr, inner.DeleteUserError)
				}
			},

			wantInnerTransport: &mockTransport{
				DeleteUserGotID:     "user-id",
				DeleteUserGotParams: &DeleteUserRequestParams{},
				DeleteUserError:     errors.New("user-error"),
			},
		},
		{
			name: "DeleteVolume",

			innerTransport: &mockTransport{
				DeleteVolumeError: errors.New("volume-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.DeleteVolume(
					context.Background(),
					"namespace-id",
					"volume-id",
					&DeleteVolumeRequestParams{},
				)
				if gotErr != inner.DeleteVolumeError {
					t.Errorf("got error %v, want %v", gotErr, inner.DeleteVolumeError)
				}
			},

			wantInnerTransport: &mockTransport{
				DeleteVolumeGotNamespace: "namespace-id",
				DeleteVolumeGotVolume:    "volume-id",
				DeleteVolumeGotParams:    &DeleteVolumeRequestParams{},
				DeleteVolumeError:        errors.New("volume-error"),
			},
		},
		{
			name: "DeleteNamespace",

			innerTransport: &mockTransport{
				DeleteNamespaceError: errors.New("namespace-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.DeleteNamespace(
					context.Background(),
					"namespace-id",
					&DeleteNamespaceRequestParams{},
				)
				if gotErr != inner.DeleteNamespaceError {
					t.Errorf("got error %v, want %v", gotErr, inner.DeleteNamespaceError)
				}
			},

			wantInnerTransport: &mockTransport{
				DeleteNamespaceGotID:     "namespace-id",
				DeleteNamespaceGotParams: &DeleteNamespaceRequestParams{},
				DeleteNamespaceError:     errors.New("namespace-error"),
			},
		},
		{
			name: "DeletePolicyGroup",

			innerTransport: &mockTransport{
				DeletePolicyGroupError: errors.New("policygroup-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.DeletePolicyGroup(
					context.Background(),
					"policygroup-id",
					&DeletePolicyGroupRequestParams{},
				)
				if gotErr != inner.DeletePolicyGroupError {
					t.Errorf("got error %v, want %v", gotErr, inner.DeletePolicyGroupError)
				}
			},

			wantInnerTransport: &mockTransport{
				DeletePolicyGroupGotID:     "policygroup-id",
				DeletePolicyGroupGotParams: &DeletePolicyGroupRequestParams{},
				DeletePolicyGroupError:     errors.New("policygroup-error"),
			},
		},
		{
			name: "AttachVolume",

			innerTransport: &mockTransport{
				AttachError: errors.New("attach-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.AttachVolume(
					context.Background(),
					"namespace-id",
					"volume-id",
					"node-id",
				)
				if gotErr != inner.AttachError {
					t.Errorf("got error %v, want %v", gotErr, inner.AttachError)
				}
			},

			wantInnerTransport: &mockTransport{
				AttachGotNamespace: "namespace-id",
				AttachGotVolume:    "volume-id",
				AttachGotNode:      "node-id",
				AttachError:        errors.New("attach-error"),
			},
		},
		{
			name: "DetachVolume",

			innerTransport: &mockTransport{
				DetachError: errors.New("detach-error"),
			},

			doTest: func(t *testing.T, inner *mockTransport) {
				wrapped := NewTransportWithReauth(inner, &mockCredentialsProvider{})

				gotErr := wrapped.DetachVolume(
					context.Background(),
					"namespace-id",
					"volume-id",
					&DetachVolumeRequestParams{},
				)
				if gotErr != inner.DetachError {
					t.Errorf("got error %v, want %v", gotErr, inner.DetachError)
				}
			},

			wantInnerTransport: &mockTransport{
				DetachGotNamespace: "namespace-id",
				DetachGotVolume:    "volume-id",
				DetachGotParams:    &DetachVolumeRequestParams{},
				DetachError:        errors.New("detach-error"),
			},
		},
	}

	for _, tt := range testWrapsFunctions {
		var tt = tt
		t.Run("wraps inner transport call to "+tt.name, func(t *testing.T) {
			t.Parallel()

			tt.doTest(t, tt.innerTransport)

			if !reflect.DeepEqual(tt.innerTransport, tt.wantInnerTransport) {
				t.Errorf("got inner transport state %v, want %v", tt.innerTransport, tt.wantInnerTransport)
			}
		})
	}

	// Test doWithReauth using a mock function that does not return an auth error,
	// ensure no auth call or retry.

	t.Run("doWithReauth fn does not retry when non-auth error", func(t *testing.T) {
		var invoked int

		wantErr := errors.New("not an auth error")
		inner := &mockTransport{}
		creds := &mockCredentialsProvider{
			ReturnUsername: "username",
			ReturnPassword: "password",
		}

		wrapped := NewTransportWithReauth(inner, creds)
		gotErr := wrapped.doWithReauth(context.Background(), func() error {
			invoked++
			return wantErr
		})
		if gotErr != wantErr {
			t.Errorf("got error %v, want %v", gotErr, wantErr)
		}
		if invoked != 1 {
			t.Errorf("got invoked %v times, want %v", invoked, 1)
		}

		if inner.AuthenticateGotUsername == creds.ReturnUsername || inner.AuthenticateGotPassword == creds.ReturnPassword {
			t.Error("unexpected call to authenticate with credentials")
		}
	})

	// Test the doWithReauth using a mock function that returns an auth error,
	// ensure that an auth attempt is made with the appropriate credentials,
	// but with no more than 1 retry.

	t.Run("doWithReauth fn retries once when auth error", func(t *testing.T) {
		var invoked int

		wantErr := NewAuthenticationError("sample text")
		inner := &mockTransport{}
		creds := &mockCredentialsProvider{
			ReturnUsername: "username",
			ReturnPassword: "password",
		}

		wrapped := NewTransportWithReauth(inner, creds)
		gotErr := wrapped.doWithReauth(context.Background(), func() error {
			invoked++
			return wantErr
		})
		if gotErr != wantErr {
			t.Errorf("got error %v, want %v", gotErr, wantErr)
		}
		if invoked != 2 {
			t.Errorf("got invoked %v times, want %v", invoked, 2)
		}

		if inner.AuthenticateGotUsername != creds.ReturnUsername {
			t.Errorf("got authenticated username %q, want %q", inner.AuthenticateGotUsername, creds.ReturnUsername)
		}

		if inner.AuthenticateGotPassword != creds.ReturnPassword {
			t.Errorf("got authenticated password %q, want %q", inner.AuthenticateGotPassword, creds.ReturnPassword)
		}
	})
}
