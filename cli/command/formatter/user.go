package formatter

import (
	"github.com/storageos/go-api/types"
	"strings"
)

const (
	defaultUserQuietFormat = "{{.Username}}"
	defaultUserTableFormat = "table {{.UUID}}\t{{.Username}}\t{{.Groups}}\t{{.Role}}"

	userUUIDHeader     = "ID"
	userUsernameHeader = "Username"
	userGroupsHeader   = "Groups"
	userRoleHeader     = "Role"
)

func NewUserFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultUserQuietFormat
		}
		return defaultUserTableFormat
	case RawFormatKey:
		if quiet {
			return `username: {{.Username}}`
		}
		return `id: {{.UUID}}\nusername: {{.Username}}\ngroups: {{.Groups}}\nrole: {{.Role}}\n`
	}
	return Format(source)
}

func UserWrite(ctx Context, users []*types.User) error {
	render := func(format func(subContext subContext) error) error {
		for _, user := range users {
			if err := format(&userContext{v: *user}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&userContext{}, render)
}

type userContext struct {
	HeaderContext
	v types.User
}

func (c *userContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *userContext) UUID() string {
	c.AddHeader(userUUIDHeader)
	return c.v.UUID
}

func (c *userContext) Username() string {
	c.AddHeader(userUsernameHeader)
	return c.v.Username
}

func (c *userContext) Groups() string {
	c.AddHeader(userGroupsHeader)
	return strings.Join(c.v.Groups, ",")
}

func (c *userContext) Role() string {
	c.AddHeader(userRoleHeader)
	return c.v.Role
}
