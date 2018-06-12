package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/storageos/go-api/types"
)

const (
	defaultNamespaceQuietFormat = "{{.Name}}"
	defaultNamespaceTableFormat = "table {{.Name}}\t{{.DisplayName}}\t{{.Description}}"

	namespaceNameHeader        = "NAMESPACE"
	namespaceDisplayNameHeader = "DISPLAY"
	namespaceDescriptionHeader = "DESCRIPTION"
)

// NewNamespaceFormat returns a format for use with a namespace Context
func NewNamespaceFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultNamespaceQuietFormat
		}
		return defaultNamespaceTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\ndisplay name: {{.DisplayName}}\n`
	}
	return Format(source)
}

// NamespaceWrite writes formatted namespaces using the Context
func NamespaceWrite(ctx Context, namespaces []*types.Namespace) error {
	render := func(format func(subContext subContext) error) error {
		for _, namespace := range namespaces {
			if err := format(&namespaceContext{v: *namespace}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&namespaceContext{}, render)
}

type namespaceContext struct {
	HeaderContext
	v types.Namespace
}

func (c *namespaceContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *namespaceContext) Name() string {
	c.AddHeader(namespaceNameHeader)
	return c.v.Name
}

func (c *namespaceContext) DisplayName() string {
	c.AddHeader(namespaceDisplayNameHeader)
	return c.v.DisplayName
}

func (c *namespaceContext) Description() string {
	c.AddHeader(namespaceDescriptionHeader)
	return c.v.Description
}

func (c *namespaceContext) Labels() string {
	c.AddHeader(labelsHeader)
	if c.v.Labels == nil {
		return ""
	}

	var joinLabels []string
	for k, v := range c.v.Labels {
		joinLabels = append(joinLabels, fmt.Sprintf("%s=%s", k, v))
	}

	sort.SliceStable(joinLabels, func(i, j int) bool {
		fst := strings.Split(joinLabels[i], "=")[0]
		snd := strings.Split(joinLabels[j], "=")[0]
		return fst < snd
	})

	return strings.Join(joinLabels, ",")
}

func (c *namespaceContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}
