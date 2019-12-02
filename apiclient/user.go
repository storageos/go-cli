package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
)

// CreateUser requests the creation of a new StorageOS user account from the
// provided fields.
func (c *Client) CreateUser(
	ctx context.Context,
	username, password string,
	withAdmin bool,
	groups ...id.PolicyGroup,
) (*user.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.CreateUser(
		ctx,
		username,
		password,
		withAdmin,
		groups...,
	)
}
