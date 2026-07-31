package payssage

import (
	"context"
	"fmt"
)

func (c *Client) DisconnectProvider(ctx context.Context, workspaceName, credentialID string) (*ProviderCredential, error) {
	var out ProviderCredential
	err := c.do(ctx, "DELETE", fmt.Sprintf("/workspaces/%s/providers/%s/disconnect", workspaceName, credentialID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
