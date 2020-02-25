# Node

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | A unique identifier for a node. The format of this type is undefined and may change but the defined properties will not change.  | [optional] [readonly] 
**Name** | **string** | The hostname of the node. This value is set by the node each time it joins the StorageOS cluster.  | [optional] [readonly] 
**Health** | [**NodeHealth**](NodeHealth.md) |  | [optional] 
**Capacity** | [**CapacityStats**](CapacityStats.md) |  | [optional] 
**IoEndpoint** | **string** | Endpoint at which we operate our dataplane&#39;s dfs service. (used for IO operations) This value is set on startup by the corresponding environment variable (IO_ADVERTISE_ADDRESS)  | [optional] [readonly] 
**SupervisorEndpoint** | **string** | Endpoint at which we operate our dataplane&#39;s supervisor service (used for sync). This value is set on startup by the corresponding environment variable (SUPERVISOR_ADVERTISE_ADDRESS)  | [optional] [readonly] 
**GossipEndpoint** | **string** | Endpoint at which we operate our health checking service. This value is set on startup by the corresponding environment variable (GOSSIP_ADVERTISE_ADDRESS)  | [optional] [readonly] 
**ClusteringEndpoint** | **string** | Endpoint at which we operate our clustering GRPC API. This value is set on startup by the corresponding environment variable (INTERNAL_API_ADVERTISE_ADDRESS)  | [optional] [readonly] 
**Labels** | **map[string]string** | A set of arbitrary key value labels to apply to the entity.  | [optional] 
**CreatedAt** | [**time.Time**](time.Time.md) | The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**UpdatedAt** | [**time.Time**](time.Time.md) | The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**Version** | **string** | An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


