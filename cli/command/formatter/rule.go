package formatter

import (
	"fmt"
	"strings"

	"github.com/storageos/go-api/types"
)

const (
	defaultRuleQuietFormat = "{{.Name}}"
	defaultRuleTableFormat = "table {{.Name}}\t{{.Operator}}\t{{.Selector}}\t{{.RuleAction}}\t{{.Labels}}"

	ruleNameHeader     = "NAMESPACE/NAME"
	ruleSelectorHeader = "SELECTOR"
	ruleOperatorHeader = "OPERATOR"
	ruleActionHeader   = "ACTION"
	ruleLableHeader    = "LABEL"
	ruleActiveHeader   = "STATUS"
)

// NewRuleFormat returns a format for use with a rule Context
func NewRuleFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultRuleQuietFormat
		}
		return defaultRuleTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\nselectors: {{.Selector}}\noperator: {{.Operator}}\naction: {{.RuleAction}}\nlabels: {{.Labels}}\n`
	}
	return Format(source)
}

// RuleWrite writes formatted rules using the Context
func RuleWrite(ctx Context, rules []*types.Rule) error {
	render := func(format func(subContext subContext) error) error {
		for _, rule := range rules {
			if err := format(&ruleContext{v: *rule}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&ruleContext{}, render)
}

type ruleContext struct {
	HeaderContext
	v types.Rule
}

func (c *ruleContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *ruleContext) Name() string {
	c.AddHeader(ruleNameHeader)
	return fmt.Sprintf("%s/%s", c.v.Namespace, c.v.Name)
}

func (c *ruleContext) Selector() string {
	c.AddHeader(ruleSelectorHeader)
	return c.v.Selector
}

func (c *ruleContext) RuleAction() string {
	c.AddHeader(ruleActionHeader)
	return fmt.Sprintf("%s", c.v.RuleAction)
}

func (c *ruleContext) Labels() string {
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

func (c *ruleContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}

func (c *ruleContext) Active() string {
	c.AddHeader(ruleActiveHeader)
	if c.v.Active {
		return "active"
	}
	return "disabled"
}
