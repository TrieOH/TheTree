package providers

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/telemetry"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"github.com/mercadopago/sdk-go/pkg/order"
	"github.com/mercadopago/sdk-go/pkg/user"
	"go.uber.org/zap"
)

const mpAuthURL = "https://auth.mercadopago.com/authorization"

type MercadoPagoProvider struct {
	clientID     string
	accessToken  string
	clientSecret string
	redirectURI  string
	oauthClient  oauth.Client
}

func NewMercadoPagoProvider(clientID, accessToken, clientSecret, redirectURI string) (*MercadoPagoProvider, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &MercadoPagoProvider{
		clientID:     clientID,
		accessToken:  accessToken,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		oauthClient:  oauth.NewClient(cfg),
	}, nil
}

func (p *MercadoPagoProvider) BuildAuthURL(state, redirectURI string) string {
	return p.oauthClient.GetAuthorizationURL(p.clientID, redirectURI, state)
}

func (p *MercadoPagoProvider) ExchangeCode(ctx context.Context, code, redirectURI string) (domain.ProviderCredentialData, error) {
	body, err := json.Marshal(map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	})
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.mercadopago.com/oauth/token",
		bytes.NewReader(body),
	)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}

	telemetry.Log().Info("MP exchange response",
		zap.Int("status", resp.StatusCode),
		zap.String("body", string(rawBody)),
	)

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserID       int    `json:"user_id"`
		Nickname     string `json:"nickname"`
	}
	if err := json.Unmarshal(rawBody, &result); err != nil {
		return domain.ProviderCredentialData{}, err
	}

	if result.AccessToken == "" {
		return domain.ProviderCredentialData{}, fmt.Errorf("MP token exchange failed: %s", string(rawBody))
	}

	return domain.ProviderCredentialData{
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
		ProviderUserID: result.UserID,
		Nickname:       result.Nickname,
	}, nil
}

func (p *MercadoPagoProvider) Charge(ctx context.Context, req domain.ChargeRequest) (*domain.ChargeResult, error) {
	cfg, err := config.New(req.SellerToken)
	if err != nil {
		return nil, err
	}

	client := order.NewClient(cfg)
	processedOrder, err := client.Process(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	return &domain.ChargeResult{
		OrderID: processedOrder.ID,
		Status:  mapMPStatus(processedOrder.StatusDetail),
	}, nil
}

func mapMPStatus(status string) domain.IntentStatus {
	switch status {
	case "accredited":
		return domain.IntentStatusSucceeded
	case "waiting_capture":
		return domain.IntentStatusPending
	case "pending_review_manual":
		return domain.IntentStatusPending
	case "in_process":
		return domain.IntentStatusPending
	default:
		return domain.IntentStatusPending
	}
}

func (p *MercadoPagoProvider) MeID(ctx context.Context, accessToken string) (int, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return 0, err
	}

	userClient := user.NewClient(cfg)
	me, err := userClient.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get MP platform user ID: %w", err)
	}
	return me.ID, nil
}

func (p *MercadoPagoProvider) MeName(ctx context.Context, accessToken string) (string, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return "", err
	}

	userClient := user.NewClient(cfg)
	me, err := userClient.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get MP platform user ID: %w", err)
	}
	return me.Nickname, nil
}

func VerifyMercadoPagoSignature(r *http.Request, secret string) bool {
	xSignature := r.Header.Get("x-signature")
	xRequestID := r.Header.Get("x-request-id")
	dataID := r.URL.Query().Get("data.id")

	var ts, hash string
	for _, part := range strings.Split(xSignature, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "ts":
			ts = val
		case "v1":
			hash = val
		}
	}

	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, ts)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(manifest))
	computed := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computed), []byte(hash))
}
