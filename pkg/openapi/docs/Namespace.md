# Namespace

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | A unique identifier for a namespace. The format of this type is undefined and may change but the defined properties will not change..  | [optional] 
**Name** | **string** |  | [optional] [readonly] 
**Labels** | **map[string]string** | A set of arbitrary key value labels to apply to the entity.  | [optional] 
**CreatedAt** | [**time.Time**](time.Time.md) | The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**UpdatedAt** | [**time.Time**](time.Time.md) | The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node&#39;s local clock was skewed. This value is for the user&#39;s informative purposes only, and correctness is not required. String format is RFC3339.  | [optional] [readonly] 
**Version** | **string** | An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


