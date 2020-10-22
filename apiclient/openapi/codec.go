package openapi

import (
	"errors"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/policygroup"
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

func (c codec) decodeLicence(model openapi.Licence) (*licence.Resource, error) {
	features := make([]string, 0)
	if model.Features != nil {
		features = append(features, *model.Features...)
	}

	return &licence.Resource{
		ClusterID:            id.Cluster(model.ClusterID),
		ExpiresAt:            model.ExpiresAt,
		ClusterCapacityBytes: model.ClusterCapacityBytes,
		UsedBytes:            model.UsedBytes,
		Kind:                 model.Kind,
		CustomerName:         model.CustomerName,
		Features:             features,
		Version:              version.FromString(model.Version),
	}, nil
}

func (c codec) decodeCapacityStats(stats openapi.CapacityStats) capacity.Stats {
	return capacity.Stats{
		Total: stats.Total,
		Free:  stats.Free,
	}
}

func (c codec) decodeNode(model openapi.Node) (*node.Resource, error) {
	return &node.Resource{
		ID:       id.Node(model.Id),
		Name:     model.Name,
		Health:   health.NodeFromString(string(model.Health)),
		Capacity: c.decodeCapacityStats(model.Capacity),

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
		ID:             id.Volume(model.Id),
		Name:           model.Name,
		Description:    model.Description,
		SizeBytes:      model.SizeBytes,
		AttachedOn:     id.Node(model.AttachedOn),
		AttachmentType: volume.AttachTypeFromString(string(model.AttachmentType)),
		Nfs:            c.decodeNFSConfig(model.Nfs),

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

func (c codec) decodePolicyGroup(model openapi.PolicyGroup) (*policygroup.Resource, error) {
	users := []*policygroup.Member{}
	if model.Users != nil {
		users = make([]*policygroup.Member, 0, len(model.Users))
		for _, u := range model.Users {
			users = append(users, &policygroup.Member{
				ID:       id.User(u.Id),
				Username: u.Username,
			})
		}
	}

	specs := []*policygroup.Spec{}
	if model.Specs != nil {
		specs = make([]*policygroup.Spec, 0, len(*model.Specs))
		for _, spec := range *model.Specs {
			specs = append(specs, &policygroup.Spec{
				NamespaceID:  id.Namespace(spec.NamespaceID),
				ResourceType: spec.ResourceType,
				ReadOnly:     spec.ReadOnly,
			})
		}
	}

	return &policygroup.Resource{
		ID:        id.PolicyGroup(model.Id),
		Name:      model.Name,
		Users:     users,
		Specs:     specs,
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

func (c codec) decodeNFSConfig(model openapi.NfsConfig) volume.NFSConfig {
	cfg := volume.NFSConfig{
		Exports:         make([]volume.NFSExportConfig, 0),
		ServiceEndpoint: "",
	}

	if model.ServiceEndpoint != nil {
		cfg.ServiceEndpoint = *model.ServiceEndpoint
	}

	if model.Exports != nil {
		for _, e := range *model.Exports {
			cfg.Exports = append(cfg.Exports, c.decodeNFSExportConfig(e))
		}
	}

	return cfg
}

func (c codec) decodeNFSExportConfig(model openapi.NfsExportConfig) volume.NFSExportConfig {
	cfg := volume.NFSExportConfig{
		ExportID:   uint(model.ExportID),
		Path:       model.Path,
		PseudoPath: model.PseudoPath,
		ACLs:       make([]volume.NFSExportConfigACL, 0, len(model.Acls)),
	}

	for _, a := range model.Acls {
		cfg.ACLs = append(cfg.ACLs, volume.NFSExportConfigACL{
			Identity: volume.NFSExportConfigACLIdentity{
				IdentityType: a.Identity.IdentityType,
				Matcher:      a.Identity.Matcher,
			},
			SquashConfig: volume.NFSExportConfigACLSquashConfig{
				GID:    a.SquashConfig.Gid,
				UID:    a.SquashConfig.Uid,
				Squash: a.SquashConfig.Squash,
			},
			AccessLevel: a.AccessLevel,
		})
	}
	return cfg
}

func (c codec) encodeNFSExport(export volume.NFSExportConfig) openapi.NfsExportConfig {
	cfg := openapi.NfsExportConfig{
		ExportID:   uint64(export.ExportID),
		Path:       export.Path,
		PseudoPath: export.PseudoPath,
		Acls:       []openapi.NfsAcl{},
	}

	for _, a := range export.ACLs {
		cfg.Acls = append(cfg.Acls, openapi.NfsAcl{
			Identity: openapi.NfsAclIdentity{
				IdentityType: a.Identity.IdentityType,
				Matcher:      a.Identity.Matcher,
			},
			SquashConfig: openapi.NfsAclSquashConfig{
				Uid:    a.SquashConfig.UID,
				Gid:    a.SquashConfig.GID,
				Squash: a.SquashConfig.Squash,
			},
			AccessLevel: a.AccessLevel,
		})
	}

	return cfg
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
