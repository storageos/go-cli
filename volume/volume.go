package volume

import (
	"errors"
	"strings"
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
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
