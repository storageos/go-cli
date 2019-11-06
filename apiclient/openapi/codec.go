package openapi

import (
	"code.storageos.net/storageos/c2-cli/pkg/entity"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"

	"code.storageos.net/storageos/openapi"
)

type codec struct{}

func (c codec) decodeGetNode(model openapi.Node) (*node.Resource, error) {
	node := &node.Resource{
		ID:     id.Node(model.Id),
		Name:   model.Name,
		Health: entity.HealthFromString(string(model.Health)),

		Labels: map[string]string(model.Labels),

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   entity.VersionFromString(model.Version),
	}

	return node, nil
}

func (c codec) decodeDescribeNode(model openapi.Node) (*node.Resource, error) {
	n, err := c.decodeGetNode(model)
	if err != nil {
		return nil, err
	}

	n.Configuration = &node.Configuration{
		IOAddr:         model.IoEndpoint,
		SupervisorAddr: model.SupervisorEndpoint,
		GossipAddr:     model.GossipEndpoint,
		ClusteringAddr: model.ClusteringEndpoint,
	}

	return n, nil
}

func (c codec) decodeGetVolume(model openapi.Volume) (*volume.Resource, error) {

	volume := &volume.Resource{
		ID:          id.Volume(model.Id),
		Name:        model.Name,
		Description: model.Description,
		SizeBytes:   model.SizeBytes,

		AttachedOn: id.Node(model.AttachedOn),
		Namespace:  id.Namespace(model.NamespaceID),
		Labels:     map[string]string(model.Labels),
		Filesystem: volume.FsTypeFromString(string(model.FsType)),
		Inode:      model.Inode,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Version:   entity.VersionFromString(model.Version),
	}

	return volume, nil
}

func (c codec) decodeDescribeVolume(model openapi.Volume) (*volume.Resource, error) {
	v, err := c.decodeGetVolume(model)
	if err != nil {
		return nil, err
	}

	m := model.Master
	v.Master = &volume.Deployment{
		ID:      id.Deployment(m.Id),
		Node:    id.Node(m.NodeID),
		Inode:   m.Inode,
		Health:  entity.HealthFromString(string(m.Health)),
		Syncing: m.Syncing,
	}

	v.Replicas = make([]*volume.Deployment, len(v.Replicas))
	for i, r := range *model.Replicas {
		v.Replicas[i] = &volume.Deployment{
			ID:      id.Deployment(r.Id),
			Node:    id.Node(r.NodeID),
			Inode:   r.Inode,
			Health:  entity.HealthFromString(string(r.Health)),
			Syncing: r.Syncing,
		}
	}

	return v, nil
}
