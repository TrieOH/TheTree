package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"payssage/internal/providers"
	"payssage/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) Callback(ctx context.Context, providerStr, code, stateStr string) (string, error) {
	ctx, span := c.tracer.Start(ctx, "Callback")
	defer span.End()

	state, err := c.oauth.Get(ctx, stateStr)
	if err != nil {
		return "", err
	}
	if state.Provider != providerStr {
		return "", fun.ErrBadRequest("invalid provider")
	}

	provider, err := providers.FromString(providerStr)
	if err != nil {
		return "", err
	}

	oauthProvider := providers.PayssageProviders.OAuth[provider]
	credentialData, err := oauthProvider.ExchangeCode(ctx, code, state.ProviderRedirectUrl)
	if err != nil {
		return "", err
	}

	var credentialID string
	switch state.Flow {
	case models.OAuthFlowCollector:
		collector, err := c.collectors.Create(ctx, models.Collector{
			OwnerID:        state.OwnerID,
			OrganizationID: state.OrganizationID,
			Provider:       providerStr,
			ProviderUserID: fmt.Sprintf("%d", credentialData.ProviderUserID),
			Credentials:    marshalCredentials(credentialData),
		})
		if err != nil {
			return "", err
		}
		credentialID = collector.ID.String()
	case models.OAuthFlowSeller:
		if state.WalletID == nil {
			return "", fun.ErrBadRequest("wallet_id is required for seller flow")
		}
		seller, err := c.sellers.Create(ctx, models.Seller{
			WalletID:       *state.WalletID,
			Provider:       providerStr,
			ProviderUserID: fmt.Sprintf("%d", credentialData.ProviderUserID),
			Credentials:    marshalCredentials(credentialData),
		})
		if err != nil {
			return "", err
		}
		credentialID = seller.ID.String()
	default:
		return "", fun.ErrBadRequest("invalid flow")
	}

	return fmt.Sprintf("%s&credential_id=%s&public_key=%s", state.FinalRedirectUrl, credentialID, credentialData.PublicKey), nil
}

func marshalCredentials(data models.ProviderCredentialData) []byte {
	b, _ := json.Marshal(data)
	return b
}
