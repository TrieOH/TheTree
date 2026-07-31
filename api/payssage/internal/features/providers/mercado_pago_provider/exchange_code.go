package mercado_pago_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payssage/models"

	"lib/telemetry"

	"go.uber.org/zap"
)

func (p *Provider) ExchangeCode(ctx context.Context, code, redirectURI string) (models.ProviderCredentialData, error) {
	body, err := json.Marshal(map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     p.cfg.MpClientID,
		"client_secret": p.cfg.MpClientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	})
	if err != nil {
		telemetry.Log().Error("error marshaling MP exchange code request body", zap.Error(err))
		return models.ProviderCredentialData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.mercadopago.com/oauth/token",
		bytes.NewReader(body),
	)
	if err != nil {
		telemetry.Log().Error("error creating MP exchange code request", zap.Error(err))
		return models.ProviderCredentialData{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		telemetry.Log().Error("error executing MP exchange code request", zap.Error(err))
		return models.ProviderCredentialData{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		telemetry.Log().Error("error reading MP exchange code response body", zap.Error(err))
		return models.ProviderCredentialData{}, err
	}

	telemetry.Log().Info("MP exchange response",
		zap.Int("status", resp.StatusCode),
		zap.String("body", string(rawBody)),
	)

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		PublicKey    string `json:"public_key"`
		UserID       int    `json:"user_id"`
	}
	if err := json.Unmarshal(rawBody, &result); err != nil {
		telemetry.Log().Error("error unmarshaling MP exchange code response body", zap.Error(err))
		return models.ProviderCredentialData{}, err
	}

	if result.AccessToken == "" {
		telemetry.Log().Error("MP exchange code response had empty access token", zap.Any("result struct", result), zap.Any("rawBody", rawBody))
		return models.ProviderCredentialData{}, fmt.Errorf("MP token exchange failed: %s", string(rawBody))
	}

	return models.ProviderCredentialData{
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
		ProviderUserID: result.UserID,
		PublicKey:      result.PublicKey,
	}, nil
}
