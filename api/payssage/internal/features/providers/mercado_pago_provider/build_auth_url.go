package mercado_pago_provider

func (p *Provider) BuildAuthURL(state, redirectURI string) string {
	return p.oauthClient.GetAuthorizationURL(p.cfg.MpClientID, redirectURI, state)
}
