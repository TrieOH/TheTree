package mercado_pago_provider

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"payssage/models"

	"lib/telemetry"
)

// tokenRequestBody is the MP token-exchange body. The redirect URI is
// Payssage's own configured callback (D7) — extracted for testability.
func (p *Provider) tokenRequestBody(code string) map[string]any {
	return map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     p.cfg.MpClientID,
		"client_secret": p.cfg.MpClientSecret,
		"code":          code,
		"redirect_uri":  p.cfg.MpRedirectURI,
	}
}

func (p *Provider) ExchangeCode(ctx context.Context, code string) (models.ProviderCredentialData, error) {
	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		PublicKey    string `json:"public_key"`
		UserID       int    `json:"user_id"`
		Error        string `json:"error"`
		Message      string `json:"message"`
	}

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetBody(p.tokenRequestBody(code)).
		SetResult(&result).
		SetResultError(&result).
		Post("https://api.mercadopago.com/oauth/token")
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "error executing MP exchange code request", telemetry.Err(err))
		return models.ProviderCredentialData{}, err
	}

	if result.AccessToken == "" {
		telemetry.Log().ErrorContext(ctx, "MP exchange code response had empty access token",
			slog.Int("status", resp.StatusCode()),
			slog.String("mp_error", result.Error),
			slog.String("mp_message", result.Message),
		)
		return models.ProviderCredentialData{}, fmt.Errorf("MP token exchange failed: status %d: %s", resp.StatusCode(), cmp.Or(result.Message, result.Error, "no access token"))
	}

	return models.ProviderCredentialData{
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
		ProviderUserID: result.UserID,
		PublicKey:      result.PublicKey,
	}, nil
}
