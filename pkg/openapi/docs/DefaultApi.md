# \DefaultApi

All URIs are relative to *http://localhost/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AttachNFSVolume**](DefaultApi.md#AttachNFSVolume) | **Post** /namespaces/{namespaceID}/volumes/{id}/nfs/attach | attach and share the volume using NFS
[**AttachVolume**](DefaultApi.md#AttachVolume) | **Post** /namespaces/{namespaceID}/volumes/{id}/attach | Attach a volume to the given node
[**AuthenticateUser**](DefaultApi.md#AuthenticateUser) | **Post** /auth/login | Authenticate a user
[**CreateNamespace**](DefaultApi.md#CreateNamespace) | **Post** /namespaces | Create a new namespace
[**CreatePolicyGroup**](DefaultApi.md#CreatePolicyGroup) | **Post** /policies | Create a new policy group
[**CreateUser**](DefaultApi.md#CreateUser) | **Post** /users | Create a new user
[**CreateVolume**](DefaultApi.md#CreateVolume) | **Post** /namespaces/{namespaceID}/volumes | Create a new Volume in the specified namespace
[**DeleteAuthenticatedUser**](DefaultApi.md#DeleteAuthenticatedUser) | **Delete** /users/self | Delete the authenticated user
[**DeleteAuthenticatedUserSessions**](DefaultApi.md#DeleteAuthenticatedUserSessions) | **Delete** /users/self/sessions | Invalidate the logged in user&#39;s sessions
[**DeleteNamespace**](DefaultApi.md#DeleteNamespace) | **Delete** /namespaces/{id} | Delete a namespace
[**DeleteNode**](DefaultApi.md#DeleteNode) | **Delete** /nodes/{id} | Delete a node
[**DeletePolicyGroup**](DefaultApi.md#DeletePolicyGroup) | **Delete** /policies/{id} | Delete a policy group
[**DeleteSessions**](DefaultApi.md#DeleteSessions) | **Delete** /users/{id}/sessions | Invalidate login sessions
[**DeleteUser**](DefaultApi.md#DeleteUser) | **Delete** /users/{id} | Delete a user
[**DeleteVolume**](DefaultApi.md#DeleteVolume) | **Delete** /namespaces/{namespaceID}/volumes/{id} | Delete a volume
[**DetachVolume**](DefaultApi.md#DetachVolume) | **Delete** /namespaces/{namespaceID}/volumes/{id}/attach | Detach the given volume
[**GetAuthenticatedUser**](DefaultApi.md#GetAuthenticatedUser) | **Get** /users/self | Get the currently authenticated user&#39;s information
[**GetCluster**](DefaultApi.md#GetCluster) | **Get** /cluster | Retrieves the cluster&#39;s global configuration settings
[**GetDiagnostics**](DefaultApi.md#GetDiagnostics) | **Get** /diagnostics | Retrieves a diagnostics bundle from the target node
[**GetLicence**](DefaultApi.md#GetLicence) | **Get** /cluster/licence | Retrieves the cluster&#39;s licence information
[**GetNamespace**](DefaultApi.md#GetNamespace) | **Get** /namespaces/{id} | Fetch a namespace
[**GetNode**](DefaultApi.md#GetNode) | **Get** /nodes/{id} | Fetch a node
[**GetPolicyGroup**](DefaultApi.md#GetPolicyGroup) | **Get** /policies/{id} | Fetch a policy group
[**GetUser**](DefaultApi.md#GetUser) | **Get** /users/{id} | Fetch a user
[**GetVolume**](DefaultApi.md#GetVolume) | **Get** /namespaces/{namespaceID}/volumes/{id} | Fetch a volume
[**ListNamespaces**](DefaultApi.md#ListNamespaces) | **Get** /namespaces | Fetch the list of namespaces
[**ListNodes**](DefaultApi.md#ListNodes) | **Get** /nodes | Fetch the list of nodes
[**ListPolicyGroups**](DefaultApi.md#ListPolicyGroups) | **Get** /policies | Fetch the list of policy groups
[**ListUsers**](DefaultApi.md#ListUsers) | **Get** /users | Fetch the list of users
[**ListVolumes**](DefaultApi.md#ListVolumes) | **Get** /namespaces/{namespaceID}/volumes | Fetch the list of volumes in the given namespace
[**RefreshJwt**](DefaultApi.md#RefreshJwt) | **Post** /auth/refresh | Refresh the JWT
[**ResizeVolume**](DefaultApi.md#ResizeVolume) | **Put** /namespaces/{namespaceID}/volumes/{id}/size | Increase the size of a volume.
[**SetReplicas**](DefaultApi.md#SetReplicas) | **Put** /namespaces/{namespaceID}/volumes/{id}/replicas | Set the number of replicas to maintain for the volume.
[**Spec**](DefaultApi.md#Spec) | **Get** /openapi | Serves this openapi spec file
[**UpdateAuthenticatedUser**](DefaultApi.md#UpdateAuthenticatedUser) | **Put** /users/self | Update the authenticated user&#39;s information
[**UpdateCluster**](DefaultApi.md#UpdateCluster) | **Put** /cluster | Update the cluster&#39;s global configuration settings
[**UpdateLicence**](DefaultApi.md#UpdateLicence) | **Put** /cluster/licence | Update the licence global configuration settings
[**UpdateNFSVolumeExports**](DefaultApi.md#UpdateNFSVolumeExports) | **Put** /namespaces/{namespaceID}/volumes/{id}/nfs/export-config | Update an nfs volume&#39;s export configuration
[**UpdateNFSVolumeMountEndpoint**](DefaultApi.md#UpdateNFSVolumeMountEndpoint) | **Put** /namespaces/{namespaceID}/volumes/{id}/nfs/mount-endpoint | Update an nfs volume&#39;s mount endpoint
[**UpdateNamespace**](DefaultApi.md#UpdateNamespace) | **Put** /namespaces/{id} | Update a namespace
[**UpdateNode**](DefaultApi.md#UpdateNode) | **Put** /nodes/{id} | Update a node
[**UpdatePolicyGroup**](DefaultApi.md#UpdatePolicyGroup) | **Put** /policies/{id} | Update a policy group
[**UpdateUser**](DefaultApi.md#UpdateUser) | **Put** /users/{id} | Update a user
[**UpdateVolume**](DefaultApi.md#UpdateVolume) | **Put** /namespaces/{namespaceID}/volumes/{id} | Update a volume



## AttachNFSVolume

> AttachNFSVolume(ctx, namespaceID, id, attachNfsVolumeData, optional)

attach and share the volume using NFS

Attach the given volume as an NFS volume. If no export configuration has been set via the /nfs/export-config endpoint, the nfs service will start with defaults settings (sharing the volume at its root). 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**attachNfsVolumeData** | [**AttachNfsVolumeData**](AttachNfsVolumeData.md)|  | 
 **optional** | ***AttachNFSVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a AttachNFSVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachVolume

> AttachVolume(ctx, namespaceID, id, attachVolumeData)

Attach a volume to the given node

Attach the volume identified by id to the node identified in the request's body. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**attachVolumeData** | [**AttachVolumeData**](AttachVolumeData.md)|  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthenticateUser

> UserSession AuthenticateUser(ctx, authUserData)

Authenticate a user

Generate a new JWT token for a user.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**authUserData** | [**AuthUserData**](AuthUserData.md)|  | 

### Return type

[**UserSession**](UserSession.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateNamespace

> Namespace CreateNamespace(ctx, createNamespaceData)

Create a new namespace

Create a new namespace in the cluster - only administrators can create new namespaces. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**createNamespaceData** | [**CreateNamespaceData**](CreateNamespaceData.md)|  | 

### Return type

[**Namespace**](Namespace.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreatePolicyGroup

> PolicyGroup CreatePolicyGroup(ctx, createPolicyGroupData)

Create a new policy group

Create a new policy group in the cluster - only administrators can create new policy groups. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**createPolicyGroupData** | [**CreatePolicyGroupData**](CreatePolicyGroupData.md)|  | 

### Return type

[**PolicyGroup**](PolicyGroup.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateUser

> User CreateUser(ctx, createUserData)

Create a new user

Create a new user in the cluster - only administrators can create new users. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**createUserData** | [**CreateUserData**](CreateUserData.md)|  | 

### Return type

[**User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateVolume

> Volume CreateVolume(ctx, namespaceID, createVolumeData, optional)

Create a new Volume in the specified namespace

Create a new volume in the given namespace

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**createVolumeData** | [**CreateVolumeData**](CreateVolumeData.md)|  | 
 **optional** | ***CreateVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a CreateVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAuthenticatedUser

> DeleteAuthenticatedUser(ctx, version, optional)

Delete the authenticated user

Remove the authenticated user from the cluster.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeleteAuthenticatedUserOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeleteAuthenticatedUserOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAuthenticatedUserSessions

> DeleteAuthenticatedUserSessions(ctx, )

Invalidate the logged in user's sessions

Invalidates logged in user's active JWTs.

### Required Parameters

This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNamespace

> DeleteNamespace(ctx, id, version, optional)

Delete a namespace

Remove the namespace identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a namespace | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeleteNamespaceOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeleteNamespaceOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNode

> DeleteNode(ctx, id, version, optional)

Delete a node

Remove the node identified by id. A node can only be deleted if it is currently offline and does not host any master deployments. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a node | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeleteNodeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeleteNodeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePolicyGroup

> DeletePolicyGroup(ctx, id, version, optional)

Delete a policy group

Remove the policy group identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a policy group | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeletePolicyGroupOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeletePolicyGroupOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSessions

> DeleteSessions(ctx, id)

Invalidate login sessions

Invalidates active JWTs on a per-user basis, specified by id. This request will not succeed if the target account is the currently authenticated account. Use the separate users/self endpoint for this purpose. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a user | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUser

> DeleteUser(ctx, id, version, optional)

Delete a user

Remove the user identified by id. This request will not succeed if the target account is the currently authenticated account. Use the separate users/self endpoint for this purpose. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a user | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeleteUserOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeleteUserOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVolume

> DeleteVolume(ctx, namespaceID, id, version, optional)

Delete a volume

Remove the volume identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DeleteVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DeleteVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 
 **offlineDelete** | **optional.Bool**| If set to true, enables deletion of a volume when all  deployments are offline, bypassing the host nodes which cannot be reached. An offline delete request will be rejected when either a) there are online deployments for the target volume or b) there is evidence that an unreachable node still has the volume master  | [default to false]

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DetachVolume

> DetachVolume(ctx, namespaceID, id, version, optional)

Detach the given volume

Detach the volume identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**version** | **string**| This value is used to perform a conditional delete or update of the entity. If the entity has been modified since the version token was obtained, the request will fail with a HTTP 409 Conflict.  | 
 **optional** | ***DetachVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a DetachVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAuthenticatedUser

> User GetAuthenticatedUser(ctx, )

Get the currently authenticated user's information

Fetch authenticated user's information.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCluster

> Cluster GetCluster(ctx, )

Retrieves the cluster's global configuration settings

Retrieves the current global configuration settings in use by the cluster. 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Cluster**](Cluster.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDiagnostics

> *os.File GetDiagnostics(ctx, )

Retrieves a diagnostics bundle from the target node

Requests that the target node gathers detailed information about the state of the cluster, using it to then build and return a bundle which can be used for troubleshooting. The request will only be served when the authenticated user is an administrator. The node will attempt to gather information about its local state, cluster-wide state and local state of other nodes in the cluster. If the cluster is unhealthy this may cause a slower response. 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[***os.File**](*os.File.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/gzip, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLicence

> Licence GetLicence(ctx, )

Retrieves the cluster's licence information

Retrieves the cluster's current licence information 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Licence**](Licence.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNamespace

> Namespace GetNamespace(ctx, id)

Fetch a namespace

Fetch the namespace identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a namespace | 

### Return type

[**Namespace**](Namespace.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNode

> Node GetNode(ctx, id)

Fetch a node

Fetch the node identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a node | 

### Return type

[**Node**](Node.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPolicyGroup

> PolicyGroup GetPolicyGroup(ctx, id)

Fetch a policy group

Fetch the policy group identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a policy group | 

### Return type

[**PolicyGroup**](PolicyGroup.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUser

> User GetUser(ctx, id)

Fetch a user

Fetch the user identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a user | 

### Return type

[**User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolume

> Volume GetVolume(ctx, namespaceID, id)

Fetch a volume

Fetch the volume identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 

### Return type

[**Volume**](Volume.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNamespaces

> []Namespace ListNamespaces(ctx, )

Fetch the list of namespaces

Fetch the list of namespaces in the cluster.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]Namespace**](Namespace.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodes

> []Node ListNodes(ctx, )

Fetch the list of nodes

Fetch the list of nodes of the cluster.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]Node**](Node.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPolicyGroups

> []PolicyGroup ListPolicyGroups(ctx, )

Fetch the list of policy groups

Fetch the list of policy groups in the cluster.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]PolicyGroup**](PolicyGroup.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUsers

> []User ListUsers(ctx, )

Fetch the list of users

Fetch the list of users of the cluster.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVolumes

> []Volume ListVolumes(ctx, namespaceID)

Fetch the list of volumes in the given namespace

Fetch the list of volumes in the cluster.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 

### Return type

[**[]Volume**](Volume.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshJwt

> UserSession RefreshJwt(ctx, )

Refresh the JWT

Obtain a fresh token with an updated expiry deadline.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**UserSession**](UserSession.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResizeVolume

> Volume ResizeVolume(ctx, namespaceID, id, resizeVolumeRequest, optional)

Increase the size of a volume.

Resize the volume identified by id in the namespace identified by namespaceID. A volume's size cannot be reduced. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**resizeVolumeRequest** | [**ResizeVolumeRequest**](ResizeVolumeRequest.md)|  | 
 **optional** | ***ResizeVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ResizeVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 
 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**Volume**](Volume.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetReplicas

> AcceptedMessage SetReplicas(ctx, namespaceID, id, setReplicasRequest, optional)

Set the number of replicas to maintain for the volume.

Set the number of replicas for the volume identified by id to the number specified in the request's body. This modifies the protected StorageOS system label \"storageos.com/replicas\". This request changes the desired replica count, and returns an error if changing the desired replica count failed. StorageOS satisfies the new replica configuration asynchronously. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**setReplicasRequest** | [**SetReplicasRequest**](SetReplicasRequest.md)|  | 
 **optional** | ***SetReplicasOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a SetReplicasOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**AcceptedMessage**](AcceptedMessage.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Spec

> string Spec(ctx, )

Serves this openapi spec file

Serves this openapi spec file

### Required Parameters

This endpoint does not need any parameter.

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/yaml

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAuthenticatedUser

> User UpdateAuthenticatedUser(ctx, updateAuthenticatedUserData, optional)

Update the authenticated user's information

Update the authenticated user.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**updateAuthenticatedUserData** | [**UpdateAuthenticatedUserData**](UpdateAuthenticatedUserData.md)|  | 
 **optional** | ***UpdateAuthenticatedUserOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateAuthenticatedUserOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCluster

> Cluster UpdateCluster(ctx, updateClusterData, optional)

Update the cluster's global configuration settings

Update the global configuration settings to use for the cluster.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**updateClusterData** | [**UpdateClusterData**](UpdateClusterData.md)|  | 
 **optional** | ***UpdateClusterOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateClusterOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**Cluster**](Cluster.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLicence

> Licence UpdateLicence(ctx, updateLicence, optional)

Update the licence global configuration settings

Update the cluster's licence.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**updateLicence** | [**UpdateLicence**](UpdateLicence.md)|  | 
 **optional** | ***UpdateLicenceOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateLicenceOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**Licence**](Licence.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNFSVolumeExports

> UpdateNFSVolumeExports(ctx, namespaceID, id, nfsVolumeExports, optional)

Update an nfs volume's export configuration

Update the NFS volume's export configuration 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**nfsVolumeExports** | [**NfsVolumeExports**](NfsVolumeExports.md)|  | 
 **optional** | ***UpdateNFSVolumeExportsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateNFSVolumeExportsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNFSVolumeMountEndpoint

> UpdateNFSVolumeMountEndpoint(ctx, namespaceID, id, nfsVolumeMountEndpoint, optional)

Update an nfs volume's mount endpoint

Update the NFS volume's mount endpoint 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**nfsVolumeMountEndpoint** | [**NfsVolumeMountEndpoint**](NfsVolumeMountEndpoint.md)|  | 
 **optional** | ***UpdateNFSVolumeMountEndpointOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateNFSVolumeMountEndpointOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

 (empty response body)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNamespace

> Namespace UpdateNamespace(ctx, id, updateNamespaceData, optional)

Update a namespace

Update the namespace identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a namespace | 
**updateNamespaceData** | [**UpdateNamespaceData**](UpdateNamespaceData.md)|  | 
 **optional** | ***UpdateNamespaceOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateNamespaceOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**Namespace**](Namespace.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNode

> Node UpdateNode(ctx, id, updateNodeData, optional)

Update a node

Update the node identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a node | 
**updateNodeData** | [**UpdateNodeData**](UpdateNodeData.md)|  | 
 **optional** | ***UpdateNodeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateNodeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**Node**](Node.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePolicyGroup

> PolicyGroup UpdatePolicyGroup(ctx, id, updatePolicyGroupData, optional)

Update a policy group

Update the policy group identified by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a policy group | 
**updatePolicyGroupData** | [**UpdatePolicyGroupData**](UpdatePolicyGroupData.md)|  | 
 **optional** | ***UpdatePolicyGroupOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdatePolicyGroupOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**PolicyGroup**](PolicyGroup.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUser

> User UpdateUser(ctx, id, updateUserData, optional)

Update a user

Update the user identified by id. This request will not succeed if the target account is the currently authenticated account. Use the separate users/self endpoint for this purpose. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| ID of a user | 
**updateUserData** | [**UpdateUserData**](UpdateUserData.md)|  | 
 **optional** | ***UpdateUserOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateUserOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]

### Return type

[**User**](User.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVolume

> Volume UpdateVolume(ctx, namespaceID, id, updateVolumeData, optional)

Update a volume

Update the description and non-storageos labels configured for the volume identified by id. 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**namespaceID** | **string**| ID of a Namespace | 
**id** | **string**| ID of a Volume | 
**updateVolumeData** | [**UpdateVolumeData**](UpdateVolumeData.md)|  | 
 **optional** | ***UpdateVolumeOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a UpdateVolumeOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **ignoreVersion** | **optional.Bool**| If set to true this value indicates that the user wants to ignore entity version constraints, thereby \&quot;forcing\&quot; the operation.  | [default to false]
 **asyncMax** | **optional.String**| Optional parameter which will make the api request asynchronous. The operation will not be cancelled even if the client disconnect. The URL parameter value overrides the \&quot;async-max\&quot; header value, if any. The value of this header defines the timeout duration for the request, it must be set to a valid duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \&quot;300ms\&quot;, or \&quot;2h45m\&quot;. Valid time units are \&quot;ns\&quot;, \&quot;us\&quot; (or \&quot;µs\&quot;), \&quot;ms\&quot;, \&quot;s\&quot;, \&quot;m\&quot;, \&quot;h\&quot;. We reject negative or nil duration values.  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[jwt](../README.md#jwt)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

