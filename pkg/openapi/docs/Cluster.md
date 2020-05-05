# Cluster

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | A unique identifier for a cluster. The format of this type is undefined and may change but the defined properties will not change.  | [optional] [readonly] 
**DisableTelemetry** | **bool** | Disables collection of telemetry data across the cluster. | [optional] [default to false]
**DisableCrashReporting** | **bool** | Disables collection of reports for any fatal crashes across the cluster.  | [optional] [default to false]
**DisableVersionCheck** | **bool** | Disables the mechanism responsible for checking if there is an updated version of StorageOS available for installation.  | [optional] [default to false]
**LogLevel** | [**LogLevel**](LogLevel.md) |  | [optional] 
**LogFormat** | [**LogFormat**](LogFormat.md) |  | [optional] 
**CreatedAt** | [**time.Time**](time.Time.md) | The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**UpdatedAt** | [**time.Time**](time.Time.md) | The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**Version** | **string** | An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


