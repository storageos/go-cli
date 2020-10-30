# NfsAclIdentity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IdentityType** | **string** | The identity type used to identify the nfs client.  | [optional] 
**Matcher** | **string** | NFS identity matcher. For \&quot;cidr\&quot;, this should be a valid CIDR block string such as \&quot;10.0.0.0/8\&quot;. For \&quot;hostname\&quot;, this must be the hostname sent by the client, with ? and * wildcard characters. For netgroup, this must be in the form of \&quot;@netgroup\&quot; with ? and * wildcard characters.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


