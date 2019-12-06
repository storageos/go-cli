# UpdateClusterData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LicenceKey** | **string** | A StorageOS product licence key, used to register a cluster. The format of this type is opaque and may change. If this field is empty we assume that the called does not want to update the cluster&#39;s licence, and do not perform any operation.  | [optional] 
**DisableTelemetry** | **bool** | Disables collection of telemetry data across the cluster.  | [optional] [default to false]
**DisableCrashReporting** | **bool** | Disables collection of reports for any fatal crashes across the cluster.  | [optional] [default to false]
**DisableVersionCheck** | **bool** | Disables the mechanism responsible for checking if there is an updated version of StorageOS available for installation.  | [optional] [default to false]
**LogLevel** | [**LogLevel**](LogLevel.md) |  | [optional] 
**LogFormat** | [**LogFormat**](LogFormat.md) |  | [optional] 
**Version** | **string** | An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


