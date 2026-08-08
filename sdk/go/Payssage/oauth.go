package payssage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OAuthFlow string

const (
	OAuthFlowCollector OAuthFlow = "collector"
	OAuthFlowSeller    OAuthFlow = "seller"
)

type ConnectProviderRequest struct {
	Flow             OAuthFlow  `json:"flow"`
	WalletID         *uuid.UUID `json:"wallet_id,omitempty"`
	OrganizationID   *uuid.UUID `json:"organization_id,omitempty"`
	FinalRedirectURL string     `json:"final_redirect_url"`
}

// ConnectProviderPayload is what actually goes over the wire: the provider
// redirect URI is derived from the client's AppURL, so callers only supply
// the flow and the final redirect.
type ConnectProviderPayload struct {
	Flow                OAuthFlow  `json:"flow"`
	WalletID            *uuid.UUID `json:"wallet_id,omitempty"`
	OrganizationID      *uuid.UUID `json:"organization_id,omitempty"`
	ProviderRedirectURL string     `json:"provider_redirect_url"`
	FinalRedirectURL    string     `json:"final_redirect_url"`
}

// ConnectProvider starts a payment-provider OAuth connection (flow
// `collector` or `seller`). It returns the provider consent URL — the caller
// is responsible for redirecting the user to it. The provider redirect URI is
// the Payssage app's own `/callback/{provider}` route (built from the
// client's AppURL), and after the flow completes Payssage redirects the
// browser to FinalRedirectURL.
func (c *Client) ConnectProvider(ctx context.Context, provider string, req ConnectProviderRequest) (string, error) {
	var out string
	err := c.do(ctx, "POST", fmt.Sprintf("/providers/%s/connect", provider), ConnectProviderPayload{
		Flow:                req.Flow,
		WalletID:            req.WalletID,
		OrganizationID:      req.OrganizationID,
		ProviderRedirectURL: fmt.Sprintf("%s/callback/%s", c.appURL, provider),
		FinalRedirectURL:    req.FinalRedirectURL,
	}, &out)
	if err != nil {
		return "", err
	}
	return out, nil
}
