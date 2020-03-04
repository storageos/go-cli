# SyncProgress

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BytesRemaining** | **uint64** | Number of bytes left remaining to complete the sync.  | [optional] 
**ThroughputBytes** | **uint64** | The average throughput of the sync given as bytes per  second.  | [optional] 
**EstimatedSecondsRemaining** | **uint64** | The estimated time left for the sync to complete, given in seconds. When this field has a value of 0 either the  sync is complete or no duration estimate could be made. The values reported for bytesRemaining and  throughputBytes provide the client with the information needed to choose what to display.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


