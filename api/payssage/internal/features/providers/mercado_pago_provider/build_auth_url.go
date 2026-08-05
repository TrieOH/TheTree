package mercado_pago_provider

import "net/url"

// BuildAuthURL builds the MercadoPago authorization URL for the marketplace
// OAuth flow (redirect the seller to grant permissions, then exchange the
// returned code for tokens via the token endpoint).
func (p *Provider) BuildAuthURL(state, redirectURI string) string {
	params := url.Values{}
	params.Set("client_id", p.cfg.MpClientID)
	params.Set("response_type", "code")
	params.Set("platform_id", "mp")
	params.Set("state", state)
	params.Set("redirect_uri", redirectURI)
	return "https://auth.mercadopago.com/authorization?" + params.Encode()
}
