# CreateUserData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** |  | 
**Password** | **string** | If not present, the existing password is not changed | [default to unchanged]
**IsAdmin** | **bool** | If true, this user is an administrator of the cluster. Administrators bypass the usual authentication checks and are granted access to all resources. Some actions (such as adding a new user) can only be performed by an administrator.  | [optional] [default to false]
**Groups** | Pointer to **[]string** | Defines a set of policy group IDs this user is a member of. Policy groups can be used to logically group users and apply authorisation  policies to all members.  | [optional] [default to []]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


