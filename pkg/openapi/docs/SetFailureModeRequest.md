# SetFailureModeRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FailureThreshold** | **uint64** | The minimum number of replicas required to be online and receiving writes in order for the volume to remain read-writable. This value replaces any previously set failure threshold or intent-based failure mode.  | [optional] 
**Version** | **string** | An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.  | [optional] 
**Mode** | [**FailureModeIntent**](FailureModeIntent.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


