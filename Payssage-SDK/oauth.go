package paymentsSDK

import (
	"context"
	"fmt"
)

type SetupProviderRequest struct {
	IsMarketplace    bool   `json:"is_marketplace"`
	FeeBps           int    `json:"fee_bps"`
	FinalRedirectURL string `json:"final_redirect_url"`
}

type ConnectSellerRequest struct {
	FinalRedirectURL string `json:"final_redirect_url"`
}

type OAuthRedirectResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// SetupProvider begins the OAuth flow for a workspace owner to connect a payment provider.
// Returns a redirect URL — the caller is responsible for redirecting the user to it.
func (c *Client) SetupProvider(ctx context.Context, workspaceName, provider string, req SetupProviderRequest) (string, error) {
	var out OAuthRedirectResponse
	if err := c.do(ctx, "POST", fmt.Sprintf("/workspaces/%s/providers/%s/setup", workspaceName, provider), req, &out); err != nil {
		return "", err
	}
	return out.RedirectURL, nil
}

// ConnectSeller begins the OAuth flow for a seller to connect their account for split payments.
// Returns a redirect URL — the caller is responsible for redirecting the user to it.
func (c *Client) ConnectSeller(ctx context.Context, workspaceName, provider string, req ConnectSellerRequest) (string, error) {
	var out OAuthRedirectResponse
	if err := c.do(ctx, "POST", fmt.Sprintf("/workspaces/%s/providers/%s/connect", workspaceName, provider), req, &out); err != nil {
		return "", err
	}
	return out.RedirectURL, nil
}

// SetMarketplaceConfig sets the marketplace configuration for a workspace.
func (c *Client) SetMarketplaceConfig(ctx context.Context, workspaceName string, req SetMarketplaceConfigRequest) (*MarketplaceConfig, error) {
	var out MarketplaceConfig
	if err := c.do(ctx, "PUT", fmt.Sprintf("/workspaces/%s/marketplace", workspaceName), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMarketplaceConfig removes the marketplace configuration for a workspace.
func (c *Client) DeleteMarketplaceConfig(ctx context.Context, workspaceName string) error {
	return c.do(ctx, "DELETE", fmt.Sprintf("/workspaces/%s/marketplace", workspaceName), nil, nil)
}
