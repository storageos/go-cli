package output

import (
	"code.storageos.net/storageos/c2-cli/volume"
)

// NFSConfig contains all NFS info for a volume.
type NFSConfig struct {
	Exports         []NFSExportConfig `json:"exports"`
	ServiceEndpoint string            `json:"serviceEndpoint"`
}

// NFSExportConfig contains a single export configuration for NFS attaching.
type NFSExportConfig struct {
	ExportID   uint                 `json:"exportID"`
	Path       string               `json:"path"`
	PseudoPath string               `json:"pseudoPath"`
	ACLs       []NFSExportConfigACL `json:"acls"`
}

// NFSExportConfigACL contains a single ACL policy for NFS attaching export
// configuration.
type NFSExportConfigACL struct {
	Identity     NFSExportConfigACLIdentity     `json:"identity"`
	SquashConfig NFSExportConfigACLSquashConfig `json:"squashConfig"`
	AccessLevel  string                         `json:"accessLevel"`
}

// NFSExportConfigACLIdentity contains identity info for an ACL in a NFS export
// config.
type NFSExportConfigACLIdentity struct {
	IdentityType string `json:"identityType"`
	Matcher      string `json:"matcher"`
}

// NFSExportConfigACLSquashConfig contains squash info for an ACL in a NFS
// export config.
type NFSExportConfigACLSquashConfig struct {
	GID    int64  `json:"gid"`
	UID    int64  `json:"uid"`
	Squash string `json:"squash"`
}

// NewNFSConfig returns a new NFSConfig object that contains all the info needed
// to be outputted.
func NewNFSConfig(c volume.NFSConfig) NFSConfig {
	return NFSConfig{
		Exports:         NewNFSExportConfigs(c.Exports),
		ServiceEndpoint: c.ServiceEndpoint,
	}
}

// NewNFSExportConfig returns a new NFSExportConfig object that contains all the
// info needed to be outputted.
func NewNFSExportConfig(n volume.NFSExportConfig) NFSExportConfig {
	exp := NFSExportConfig{
		ExportID:   n.ExportID,
		Path:       n.Path,
		PseudoPath: n.PseudoPath,
		ACLs:       []NFSExportConfigACL{},
	}

	for _, a := range n.ACLs {
		exp.ACLs = append(exp.ACLs, NFSExportConfigACL{
			Identity: NFSExportConfigACLIdentity{
				IdentityType: a.Identity.IdentityType,
				Matcher:      a.Identity.Matcher,
			},
			SquashConfig: NFSExportConfigACLSquashConfig{
				UID:    a.SquashConfig.UID,
				GID:    a.SquashConfig.GID,
				Squash: a.SquashConfig.Squash,
			},
			AccessLevel: a.AccessLevel,
		})
	}

	return exp
}

// NewNFSExportConfigs returns a list of NFSExportConfig objects that contains
// all the info needed to be outputted.
func NewNFSExportConfigs(exps []volume.NFSExportConfig) []NFSExportConfig {
	configs := make([]NFSExportConfig, 0, len(exps))
	for _, e := range exps {
		configs = append(configs, NewNFSExportConfig(e))
	}
	return configs
}
