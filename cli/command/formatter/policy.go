package formatter

import (
	"github.com/storageos/go-api/types"
)

const (
	defaultPolicyTableFormat = "table {{.ID}}\t{{.User}}\t{{.Group}}\t{{.Namespace}}"

	policyIDHeader              = "ID"
	policyUserHeader            = "USER"
	policyGroupHeader           = "GROUP"
	policyReadonlyHeader        = "READONLY"
	policyAPIGroupHeader        = "API_GROUP"
	policyResourceHeader        = "RESOURCE"
	policyNamespaceHeader       = "NAMESPACE"
	policyNonResourcePathHeader = "NON_RESOURCE_PATH"
)

func NewPolicyFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultPolicyTableFormat
	case RawFormatKey:
		return "id: {{.ID}}\nuser: {{.User}}\ngroup: {{.Group}}\nnamespace: {{.Namespace}}"
	}
	return Format(source)
}

func PolicyWrite(ctx Context, policies []*types.PolicyWithID) error {
	render := func(format func(subContext subContext) error) error {
		for _, policy := range policies {
			if err := format(&policyContext{v: *policy}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&policyContext{}, render)
}

type policyContext struct {
	HeaderContext
	v types.PolicyWithID
}

func (c *policyContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *policyContext) ID() string {
	c.AddHeader(policyIDHeader)
	return c.v.ID
}

func (c *policyContext) User() string {
	c.AddHeader(policyUserHeader)
	return c.v.Spec.User
}

func (c *policyContext) Group() string {
	c.AddHeader(policyGroupHeader)
	return c.v.Spec.Group
}

func (c *policyContext) Readonly() string {
	c.AddHeader(policyReadonlyHeader)
	if c.v.Spec.Readonly {
		return "true"
	}
	return "false"
}

func (c *policyContext) APIGroup() string {
	c.AddHeader(policyAPIGroupHeader)
	return c.v.Spec.APIGroup
}

func (c *policyContext) Resource() string {
	c.AddHeader(policyResourceHeader)
	return c.v.Spec.Resource
}

func (c *policyContext) Namespace() string {
	c.AddHeader(policyNamespaceHeader)
	return c.v.Spec.Namespace
}

func (c *policyContext) NonResourcePath() string {
	c.AddHeader(policyNonResourcePathHeader)
	return c.v.Spec.NonResourcePath
}
