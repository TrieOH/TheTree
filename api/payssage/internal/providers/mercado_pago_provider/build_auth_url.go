package mercado_pago_provider

import "net/url"

// BuildAuthURL builds the MercadoPago authorization URL for the marketplace
// OAuth flow (redirect the seller to grant permissions, then exchange the
// returned code for tokens via the token endpoint). The redirect URI is
// Payssage's own `/callback/{provider}` route from config (MP validates it
// at code exchange), so callers never supply one (D7).
func (p *Provider) BuildAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", p.cfg.MpClientID)
	params.Set("response_type", "code")
	params.Set("platform_id", "mp")
	params.Set("state", state)
	params.Set("redirect_uri", p.cfg.MpRedirectURI)
	return "https://auth.mercadopago.com/authorization?" + params.Encode()
}
