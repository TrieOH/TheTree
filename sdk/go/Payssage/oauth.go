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

// ConnectProviderPayload is what actually goes over the wire. The provider
// redirect URI is Payssage's own business (D7) — its backend builds the
// callback URL from its own config — so callers only supply the flow and the
// final redirect.
type ConnectProviderPayload struct {
	Flow             OAuthFlow  `json:"flow"`
	WalletID         *uuid.UUID `json:"wallet_id,omitempty"`
	OrganizationID   *uuid.UUID `json:"organization_id,omitempty"`
	FinalRedirectURL string     `json:"final_redirect_url"`
}

// ConnectProvider starts a payment-provider OAuth connection (flow
// `collector` or `seller`). It returns the provider consent URL — the caller
// is responsible for redirecting the user to it. The provider redirect URI is
// Payssage's own `/callback/{provider}` route (built from its own config, D7),
// and after the flow completes Payssage redirects the browser to
// FinalRedirectURL.
func (c *Client) ConnectProvider(ctx context.Context, provider string, req ConnectProviderRequest) (string, error) {
	var out string
	err := c.do(ctx, "POST", fmt.Sprintf("/providers/%s/connect", provider), ConnectProviderPayload(req), &out)
	if err != nil {
		return "", err
	}
	return out, nil
}
