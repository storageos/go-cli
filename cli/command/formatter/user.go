package formatter

import (
	"strings"

	"github.com/storageos/go-api/types"
)

const (
	defaultUserQuietFormat = "{{.Username}}"
	defaultUserTableFormat = "table {{.UUID}}\t{{.Username}}\t{{.Groups}}\t{{.Role}}"

	userUUIDHeader     = "ID"
	userUsernameHeader = "USERNAME"
	userGroupsHeader   = "GROUPS"
	userRoleHeader     = "ROLE"
)

// NewUserFormat returns a format string for user list operations
// corresponding to the format key passed to it. If the key given
// is not supported it will be converted to a format string and returned.
// If the quiet parameter is set and the format key is supported,
// only the username will be displayed.
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

// UserWrite writes the given usuers to the provided context,
// using the format specified within the context.
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
