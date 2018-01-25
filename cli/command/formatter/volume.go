package formatter

import (
	"fmt"
	"strconv"
	"strings"

	units "github.com/docker/go-units"
	"github.com/storageos/go-cli/api/types"
	cliconfig "github.com/storageos/go-cli/cli/config"
)

const (
	defaultVolumeQuietFormat = "{{.Name}}"
	defaultVolumeTableFormat = "table {{.Name}}\t{{.Size}}\t{{.MountedBy}}\t{{.NodeSelector}}\t{{.Status}}\t{{.Replicas}}\t{{.Location}}"

	volumeNameHeader         = "NAMESPACE/NAME"
	volumeMountedByHeader    = "MOUNTED BY"
	volumeNodeSelectorHeader = "NODE SELECTOR"
	volumeStatusHeader       = "STATUS"
	volumeReplicasHeader     = "REPLICAS"
	volumeLocationHeader     = "LOCATION"
)

// NewVolumeFormat returns a format for use with a volume Context
func NewVolumeFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultVolumeQuietFormat
		}
		return defaultVolumeTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\ndriver: {{.Driver}}\n`
	}
	return Format(source)
}

// VolumeWrite writes formatted volumes using the Context
func VolumeWrite(ctx Context, volumes []*types.Volume, nodes []*types.Controller) error {
	render := func(format func(subContext subContext) error) error {
		for _, volume := range volumes {
			if err := format(&volumeContext{v: *volume, nodes: nodes}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&volumeContext{}, render)
}

type volumeContext struct {
	HeaderContext
	v     types.Volume
	nodes []*types.Controller
}

func (c *volumeContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *volumeContext) Name() string {
	c.AddHeader(volumeNameHeader)
	return fmt.Sprintf("%s/%s", c.v.Namespace, c.v.Name)
}

func (c *volumeContext) Labels() string {
	c.AddHeader(labelsHeader)
	if c.v.Labels == nil {
		return ""
	}

	var joinLabels []string
	for k, v := range c.v.Labels {
		joinLabels = append(joinLabels, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(joinLabels, ",")
}

func (c *volumeContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}

func (c *volumeContext) MountedBy() string {
	c.AddHeader(volumeMountedByHeader)
	if c.v.MountedBy == "" {
		return ""
	}
	return fmt.Sprintf("%s", c.v.MountedBy)
}

func (c *volumeContext) NodeSelector() string {
	c.AddHeader(volumeNodeSelectorHeader)
	if c.v.NodeSelector == "" {
		return ""
	}
	return fmt.Sprintf("%s", c.v.NodeSelector)
}

func (c *volumeContext) Size() string {
	c.AddHeader(sizeHeader)
	return units.HumanSize(float64(c.v.Size * 1000000000))
}

func (c *volumeContext) Status() string {
	c.AddHeader(volumeStatusHeader)
	return c.v.Status
}

func (c *volumeContext) Replicas() string {
	c.AddHeader(volumeReplicasHeader)

	// desired
	desired := getDesiredReplicas(&c.v)
	activeReplicas := activeReplicas(&c.v)

	return fmt.Sprintf("%d/%d", activeReplicas, desired)
}

func (c *volumeContext) Location() string {
	c.AddHeader(volumeLocationHeader)
	if c.v.Master != nil {
		master, err := c.nodeByID(c.v.Master.Controller)
		if err != nil {
			return "-"
		}

		return fmt.Sprintf("%s (%s)", master.Name, master.Health)
	}

	return "-"
}

func (c *volumeContext) nodeByID(id string) (*types.Controller, error) {
	for _, node := range c.nodes {
		if node.ID == id {
			return node, nil
		}
	}
	return nil, fmt.Errorf("node not found")
}

func activeReplicas(volume *types.Volume) int {
	found := 0
	for _, replica := range volume.Replicas {
		// looking for active replicas
		if replica.Status == "active" && replica.Health == "healthy" {
			found++
		}
	}

	return found
}

// GetDesiredReplicas - get desired replicas.
// If the value is invalid (i.e. storageos.feature.replicas="hi") - desired
// replicas will be set to 0. Only valid values will be tolerated.
func getDesiredReplicas(volume *types.Volume) int {
	r, ok := volume.Labels[cliconfig.FeatureReplicas]
	// if replication label is missing - do nothing
	if !ok {
		return 0
	}

	desiredReplicas, err := strconv.Atoi(r)
	if err != nil {
		return 0
	}

	return desiredReplicas
}
