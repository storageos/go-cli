package openapi

import (
	"errors"

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

// codec provides functionality to "decode" openapi models, translating them to
// internal types.
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
		Health: health.FromString(string(model.Health)),

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
		Inode:      model.Inode,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   version.FromString(model.Version),
	}

	m := model.Master
	v.Master = &volume.Deployment{
		ID:      id.Deployment(m.Id),
		Node:    id.Node(m.NodeID),
		Inode:   m.Inode,
		Health:  health.FromString(string(m.Health)),
		Syncing: m.Syncing,
	}

	v.Replicas = make([]*volume.Deployment, len(v.Replicas))
	for i, r := range *model.Replicas {
		v.Replicas[i] = &volume.Deployment{
			ID:      id.Deployment(r.Id),
			Node:    id.Node(r.NodeID),
			Inode:   r.Inode,
			Health:  health.FromString(string(r.Health)),
			Syncing: r.Syncing,
		}
	}

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
	groups := make([]id.PolicyGroup, len(*model.Groups))
	for i, groupID := range *model.Groups {
		groups[i] = id.PolicyGroup(groupID)
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
	switch filesystem.String() {
	case string(openapi.EXT2):
		return openapi.EXT2, nil
	case string(openapi.EXT3):
		return openapi.EXT3, nil
	case string(openapi.EXT4):
		return openapi.EXT4, nil
	case string(openapi.XFS):
		return openapi.XFS, nil
	case string(openapi.BTRFS):
		return openapi.BTRFS, nil
	case string(openapi.BLOCK):
		return openapi.BLOCK, nil
	default:
		return "", errors.New("unknown fs type")
	}
}
