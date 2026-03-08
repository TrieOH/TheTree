package providers

import (
	"TriePayments/internal/core/domain"
	"context"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
)

const mpAuthURL = "https://auth.mercadopago.com/authorization"

type MercadoPagoProvider struct {
	clientID    string
	redirectURI string
	oauthClient oauth.Client
}

func NewMercadoPagoProvider(clientID, accessToken, redirectURI string) (domain.OAuthProvider, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &MercadoPagoProvider{
		clientID:    clientID,
		redirectURI: redirectURI,
		oauthClient: oauth.NewClient(cfg),
	}, nil
}

func (p *MercadoPagoProvider) BuildAuthURL(state string) string {
	return p.oauthClient.GetAuthorizationURL(p.clientID, p.redirectURI, state)
}

func (p *MercadoPagoProvider) ExchangeCode(ctx context.Context, code string) (domain.ProviderCredentialData, error) {
	resource, err := p.oauthClient.Create(ctx, code, p.redirectURI)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}

	return domain.ProviderCredentialData{
		AccessToken:    resource.AccessToken,
		RefreshToken:   resource.RefreshToken,
		ProviderUserID: resource.PublicKey,
	}, nil
}
