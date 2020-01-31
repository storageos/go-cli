package volume

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

const (
	// LabelNoCache is a StorageOS volume label which when enabled disables the
	// caching of volume data.
	LabelNoCache = "storageos.com/nocache"
	// LabelNoCompress is a StorageOS volume label which when enabled disables the
	// compression of volume data (both at rest and during transit).
	LabelNoCompress = "storageos.com/nocompress"
	// LabelReplicas is a StorageOS volume label which decides how many replicas
	// must be provisioned for that volume.
	LabelReplicas = "storageos.com/replicas"
	// LabelThrottle is a StorageOS volume label which when enabled deprioritises
	// the volume's traffic by reducing disk I/O rate.
	LabelThrottle = "storageos.com/throttle"
	// LabelHintMaster is a StorageOS volume label holding a list of nodes. When set,
	// placement of the volume master on one of the nodes present in the list is
	// preferred.
	LabelHintMaster = "storageos.com/hint.master"
	// LabelHintReplicas is a StorageOS volume label holding a list of nodes. When
	// set, placement of volume replicas on nodes present in the list is
	// preferred.
	LabelHintReplicas = "storageos.com/hint.replicas"
)

// ErrNoNamespace is an error stating that an ID based volume reference
// string is missing its namespace.
var ErrNoNamespace = errors.New("namespace not specified for id format")

// FsType indicates the kind of filesystem which a volume has been given.
type FsType string

// String returns the name string for fs.
func (fs FsType) String() string {
	return string(fs)
}

// FsTypeFromString wraps name as an FsType. It doesn't perform validity
// checks.
func FsTypeFromString(name string) FsType {
	return FsType(name)
}

// Resource encapsulates a StorageOS volume API resource as a data type.
type Resource struct {
	ID          id.Volume `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AttachedOn  id.Node   `json:"attachedOn"`

	Namespace  id.Namespace `json:"namespaceID"`
	Labels     labels.Set   `json:"labels"`
	Filesystem FsType       `json:"filesystem"`
	Inode      uint32       `json:"inode"`
	SizeBytes  uint64       `json:"sizeBytes"`

	Master   *Deployment   `json:"master"`
	Replicas []*Deployment `json:"replicas"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// Deployment encapsulates a deployment instance for a
// volume as a data type.
type Deployment struct {
	ID      id.Deployment `json:"id"`
	Node    id.Node       `json:"nodeID"`
	Inode   uint32        `json:"inode"`
	Health  health.State  `json:"health"`
	Syncing bool          `json:"syncing"`
}

// IsCachingDisabled returns if the volume resource is configured to disable
// caching of data.
func (r *Resource) IsCachingDisabled() (bool, error) {
	value, exists := r.Labels[LabelNoCache]
	if !exists {
		return false, nil
	}

	return strconv.ParseBool(value)
}

// IsCompressionDisabled returns if the volume resource is configured to disable
// compression of data at rest and during transit.
func (r *Resource) IsCompressionDisabled() (bool, error) {
	value, exists := r.Labels[LabelNoCompress]
	if !exists {
		return false, nil
	}

	return strconv.ParseBool(value)
}

// IsThrottleEnabled returns if the volume resource is configured to have its
// traffic deprioritised by reducing its disk I/O rate.
func (r *Resource) IsThrottleEnabled() (bool, error) {
	value, exists := r.Labels[LabelThrottle]
	if !exists {
		return false, nil
	}

	return strconv.ParseBool(value)
}

// ParseReferenceName will parse a volume reference string built of
// a namespace name and a volume name.
//
// if no namespace name is present then "default" is returned for the
// namespace.
func ParseReferenceName(ref string) (namespace string, volume string, err error) {
	parts := strings.Split(ref, "/")

	switch len(parts) {
	case 2:
		return parts[0], parts[1], nil
	case 1:
		return "default", parts[0], nil
	default:
		return "", "", errors.New("invalid volume reference string")
	}
}

// ParseReferenceID will parse a volume reference string built of a namespace
// ID and a volume ID.
//
// if the reference string does not contain a namespace then the volume id
// is returned along with an ErrNoNamespace, so that the caller can check
// for the value and decide on using the default namespace (as this is not
// free for ID usecases)
func ParseReferenceID(ref string) (id.Namespace, id.Volume, error) {
	parts := strings.Split(ref, "/")

	switch len(parts) {
	case 2:
		return id.Namespace(parts[0]), id.Volume(parts[1]), nil
	case 1:
		return "", id.Volume(parts[0]), ErrNoNamespace
	default:
		return "", "", errors.New("invalid volume reference string")
	}
}
