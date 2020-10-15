# NfsExportConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExportID** | **uint64** | ID for this export  | [optional] 
**Path** | **string** | The path relative to the volume root to serve as the export root  | [optional] [default to ]
**PseudoPath** | **string** | The configured pseudo path in the NFS virtual filesystem. This is the path clients will see when traversing to this export on the NFS share.  | [optional] [default to ]
**Acls** | [**[]NfsAcl**](NfsAcl.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


