package paymentsSDK

import (
	"context"
	"fmt"
)

type BeginOAuthRequest struct {
	FinalRedirectURL string `json:"final_redirect_url"`
}

type BeginOAuthResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// BeginOAuth returns the provider's authorization URL.
// The caller is responsible for redirecting the user to this URL.
func (c *Client) BeginOAuth(ctx context.Context, workspaceName, provider string, req BeginOAuthRequest) (string, error) {
	var out BeginOAuthResponse
	if err := c.do(ctx, "POST", fmt.Sprintf("/workspaces/%s/oauth/%s/begin", workspaceName, provider), req, &out); err != nil {
		return "", err
	}
	return out.RedirectURL, nil
}

func (c *Client) ListCredentials(ctx context.Context, workspaceName string) ([]ProviderCredential, error) {
	var out []ProviderCredential
	if err := c.do(ctx, "GET", fmt.Sprintf("/workspaces/%s/credentials", workspaceName), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RevokeCredential(ctx context.Context, workspaceName, credentialID string) error {
	return c.do(ctx, "DELETE", fmt.Sprintf("/workspaces/%s/credentials/%s", workspaceName, credentialID), nil, nil)
}
