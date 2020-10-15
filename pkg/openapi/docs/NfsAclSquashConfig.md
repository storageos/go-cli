# NfsAclSquashConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uid** | **int64** |  | [optional] 
**Gid** | **int64** |  | [optional] 
**Squash** | **string** | SquashConfig defines the root squashing behaviour.  When a client creates a file, it sends the user UID from the client. If the client is running as root, this sends uid&#x3D;0. Root squashing allows the NFS administrator to prevent the client from writing as \&quot;root\&quot; to the NFS share, instead mapping the client to a new UID/GID (usually nfsnobody, -2). \&quot;none\&quot; performs no UID/GID alterations, using the values sent by the client. \&quot;root\&quot; mapps UID &amp; GID 0 to the values specified. \&quot;rootuid\&quot; maps UID 0 and a GID of any value to the value specified. \&quot;all\&quot; maps changes all UID and GID values to those specified.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


