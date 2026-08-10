package mercado_pago_provider

import (
	"net/url"
	"testing"

	"payssage/internal/config"
)

func newTestProvider() *Provider {
	return &Provider{cfg: config.MercadoPagoConfig{
		MpClientID:     "TEST_CLIENT_ID",
		MpClientSecret: "TEST_CLIENT_SECRET",
		MpAccessToken:  "TEST_ACCESS_TOKEN",
		MpRedirectURI:  "https://pay.example/callback/mercado_pago",
	}}
}

// TestBuildAuthURL_UsesConfiguredRedirectURI pins D7: the auth URL carries
// Payssage's own configured redirect URI — callers never supply one.
func TestBuildAuthURL_UsesConfiguredRedirectURI(t *testing.T) {
	p := newTestProvider()

	authURL := p.BuildAuthURL("state-token")

	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("auth url not parseable: %v", err)
	}
	q := parsed.Query()
	if got := q.Get("redirect_uri"); got != p.cfg.MpRedirectURI {
		t.Errorf("redirect_uri = %q, want %q", got, p.cfg.MpRedirectURI)
	}
	if got := q.Get("state"); got != "state-token" {
		t.Errorf("state = %q, want state-token", got)
	}
	if got := q.Get("client_id"); got != "TEST_CLIENT_ID" {
		t.Errorf("client_id = %q, want TEST_CLIENT_ID", got)
	}
}

// TestTokenRequestBody_UsesConfiguredRedirectURI pins the same D7 wiring on
// the token exchange: the body carries the configured redirect URI.
func TestTokenRequestBody_UsesConfiguredRedirectURI(t *testing.T) {
	p := newTestProvider()

	body := p.tokenRequestBody("auth-code")

	if got := body["redirect_uri"]; got != p.cfg.MpRedirectURI {
		t.Errorf("redirect_uri = %v, want %q", got, p.cfg.MpRedirectURI)
	}
	if got := body["code"]; got != "auth-code" {
		t.Errorf("code = %v, want auth-code", got)
	}
	if got := body["grant_type"]; got != "authorization_code" {
		t.Errorf("grant_type = %v, want authorization_code", got)
	}
}
