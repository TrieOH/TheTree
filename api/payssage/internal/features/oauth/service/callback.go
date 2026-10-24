package service

import (
	"context"
	"encoding/json"
	"fmt"
	"payssage/internal/providers"
	"payssage/models"
	"strconv"

	"github.com/MintzyG/fun"
)

func (o *Operations) Callback(ctx context.Context, providerStr, code, stateStr string) (string, error) {
	state, err := o.OAuth.Get(ctx, stateStr)
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

	oauthProvider := o.Registry.OAuth[provider]
	credentialData, err := oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	var credentialID string
	switch state.Flow {
	case models.OAuthFlowCollector:
		collector, err := o.Collectors.Create(ctx, models.Collector{
			OwnerID:        state.OwnerID,
			OrganizationID: state.OrganizationID,
			Provider:       providerStr,
			ProviderUserID: strconv.Itoa(credentialData.ProviderUserID),
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
		seller, err := o.Sellers.Create(ctx, models.Seller{
			WalletID:       *state.WalletID,
			Provider:       providerStr,
			ProviderUserID: strconv.Itoa(credentialData.ProviderUserID),
			Credentials:    marshalCredentials(credentialData),
		})
		if err != nil {
			return "", err
		}
		credentialID = seller.ID.String()
	default:
		return "", fun.ErrBadRequest("invalid flow")
	}

	return fmt.Sprintf("%s?credential_id=%s&public_key=%s", state.FinalRedirectURL, credentialID, credentialData.PublicKey), nil
}

func marshalCredentials(data models.ProviderCredentialData) []byte {
	b, _ := json.Marshal(data) //nolint:gosec
	return b
}
