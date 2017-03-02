package formatter

import (
	"fmt"
	"strings"

	units "github.com/docker/go-units"
	"github.com/storageos/go-api/types"
)

const (
	defaultVolumeQuietFormat = "{{.Name}}"
	defaultVolumeTableFormat = "table {{.Name}}\t{{.Size}}\t{{.MountedBy}}\t{{.Status}}"

	volumeNameHeader      = "NAMESPACE/NAME"
	volumeMountedByHeader = "MOUNTED BY"
	volumeStatusHeader    = "STATUS"
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
func VolumeWrite(ctx Context, volumes []*types.Volume) error {
	render := func(format func(subContext subContext) error) error {
		for _, volume := range volumes {
			if err := format(&volumeContext{v: *volume}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&volumeContext{}, render)
}

type volumeContext struct {
	HeaderContext
	v types.Volume
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
	return fmt.Sprintf("%d", c.v.MountedBy)
}

func (c *volumeContext) Size() string {
	c.AddHeader(sizeHeader)
	return units.HumanSize(float64(c.v.Size * 1000000000))
}

func (c *volumeContext) Status() string {
	c.AddHeader(volumeStatusHeader)
	return c.v.Status
}
