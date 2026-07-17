package commands

import (
	"context"
	"fmt"
	"lib/telemetry"
	"payssage/internal/providers"
	"payssage/models"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
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

	switch state.Flow {
	case models.OAuthFlowCollector:
		// Create Collector
	case models.OAuthFlowSeller:
		// Create Seller tied to Wallet
	default:
		return "", fun.ErrBadRequest("invalid flow")
	}

	telemetry.Log().Info("credential data", zap.Any("credentials", credentialData))

	return fmt.Sprintf("%s&credential_id=%s&public_key=%s", state.FinalRedirectUrl, "bogus", credentialData.PublicKey), nil
}
