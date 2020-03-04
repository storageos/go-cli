package openapi

import (
	"errors"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"

	"code.storageos.net/storageos/openapi"
)

// codec provides functionality to encode/decode openapi models, translating
// them to/from internal types.
type codec struct{}

func (c codec) decodeCluster(model openapi.Cluster) (*cluster.Resource, error) {
	return &cluster.Resource{
		ID: id.Cluster(model.Id),

		Licence: &cluster.Licence{
			ClusterID:            id.Cluster(model.Licence.ClusterID),
			ExpiresAt:            model.Licence.ExpiresAt,
			ClusterCapacityBytes: model.Licence.ClusterCapacityBytes,
			Kind:                 model.Licence.Kind,
			CustomerName:         model.Licence.CustomerName,
		},

		DisableTelemetry:      model.DisableTelemetry,
		DisableCrashReporting: model.DisableCrashReporting,
		DisableVersionCheck:   model.DisableVersionCheck,

		LogLevel:  cluster.LogLevelFromString(string(model.LogLevel)),
		LogFormat: cluster.LogFormatFromString(string(model.LogFormat)),

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}, nil
}

func (c codec) decodeNode(model openapi.Node) (*node.Resource, error) {
	return &node.Resource{
		ID:     id.Node(model.Id),
		Name:   model.Name,
		Health: health.NodeFromString(string(model.Health)),

		Labels: model.Labels,

		IOAddr:         model.IoEndpoint,
		SupervisorAddr: model.SupervisorEndpoint,
		GossipAddr:     model.GossipEndpoint,
		ClusteringAddr: model.ClusteringEndpoint,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}, nil
}

func (c codec) decodeVolume(model openapi.Volume) (*volume.Resource, error) {
	v := &volume.Resource{
		ID:          id.Volume(model.Id),
		Name:        model.Name,
		Description: model.Description,
		SizeBytes:   model.SizeBytes,

		AttachedOn: id.Node(model.AttachedOn),
		Namespace:  id.Namespace(model.NamespaceID),
		Labels:     model.Labels,
		Filesystem: volume.FsTypeFromString(string(model.FsType)),

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}

	m := model.Master
	v.Master = &volume.Deployment{
		ID:         id.Deployment(m.Id),
		Node:       id.Node(m.NodeID),
		Health:     health.MasterFromString(string(m.Health)),
		Promotable: m.Promotable,
	}

	replicas := []*volume.Deployment{}

	if model.Replicas != nil {
		replicas = make([]*volume.Deployment, len(*model.Replicas))
		for i, r := range *model.Replicas {
			replicas[i] = &volume.Deployment{
				ID:         id.Deployment(r.Id),
				Node:       id.Node(r.NodeID),
				Health:     health.ReplicaFromString(string(r.Health)),
				Promotable: r.Promotable,
			}

			p := r.SyncProgress

			if (p != openapi.SyncProgress{}) {
				replicas[i].SyncProgress = &volume.SyncProgress{
					BytesRemaining:            p.BytesRemaining,
					ThroughputBytes:           p.ThroughputBytes,
					EstimatedSecondsRemaining: p.EstimatedSecondsRemaining,
				}
			}
		}
	}

	v.Replicas = replicas

	return v, nil
}

func (c codec) decodeNamespace(model openapi.Namespace) (*namespace.Resource, error) {
	return &namespace.Resource{
		ID:     id.Namespace(model.Id),
		Name:   model.Name,
		Labels: model.Labels,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}, nil
}

func (c codec) decodeUser(model openapi.User) (*user.Resource, error) {

	groups := []id.PolicyGroup{}

	if model.Groups != nil {
		groups = make([]id.PolicyGroup, len(*model.Groups))
		for i, groupID := range *model.Groups {
			groups[i] = id.PolicyGroup(groupID)
		}
	}

	return &user.Resource{
		ID:       id.User(model.Id),
		Username: model.Username,

		IsAdmin: model.IsAdmin,
		Groups:  groups,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}, nil
}

func (c codec) encodeFsType(filesystem volume.FsType) (openapi.FsType, error) {
	v := openapi.FsType(filesystem.String())
	switch v {
	case openapi.FSTYPE_EXT2, openapi.FSTYPE_EXT3,
		openapi.FSTYPE_EXT4, openapi.FSTYPE_XFS,
		openapi.FSTYPE_BTRFS, openapi.FSTYPE_BLOCK:
		return v, nil
	default:
		return "", apiclient.NewEncodingError(
			errors.New("unknown fs type"),
			v,
			filesystem,
		)
	}
}

func (c codec) encodeLogLevel(level cluster.LogLevel) (openapi.LogLevel, error) {
	v := openapi.LogLevel(level.String())
	switch v {
	case openapi.LOGLEVEL_DEBUG, openapi.LOGLEVEL_INFO,
		openapi.LOGLEVEL_WARN, openapi.LOGLEVEL_ERROR:
		return v, nil
	default:
		return "", apiclient.NewEncodingError(
			errors.New("unknown log level"),
			v,
			level,
		)
	}
}

func (c codec) encodeLogFormat(format cluster.LogFormat) (openapi.LogFormat, error) {
	v := openapi.LogFormat(format.String())
	switch v {
	case openapi.LOGFORMAT_DEFAULT, openapi.LOGFORMAT_JSON:
		return v, nil
	default:
		return "", apiclient.NewEncodingError(
			errors.New("unknown log format"),
			v,
			format,
		)
	}
}
