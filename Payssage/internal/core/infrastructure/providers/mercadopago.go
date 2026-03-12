package providers

import (
	"TriePayments/internal/core/domain"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/user"
)

const mpAuthURL = "https://auth.mercadopago.com/authorization"

type MercadoPagoProvider struct {
	clientID    string
	accessToken string // add this
	redirectURI string
	oauthClient oauth.Client
}

func NewMercadoPagoProvider(clientID, accessToken, redirectURI string) (*MercadoPagoProvider, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &MercadoPagoProvider{
		clientID:    clientID,
		accessToken: accessToken, // add this
		redirectURI: redirectURI,
		oauthClient: oauth.NewClient(cfg),
	}, nil
}

func (p *MercadoPagoProvider) BuildAuthURL(state, redirectURI string) string {
	return p.oauthClient.GetAuthorizationURL(p.clientID, redirectURI, state)
}

func (p *MercadoPagoProvider) ExchangeCode(ctx context.Context, code, redirectURI string) (domain.ProviderCredentialData, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.accessToken)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.mercadopago.com/oauth/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.ProviderCredentialData{}, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserID       int    `json:"user_id"`
		Nickname     string `json:"nickname"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.ProviderCredentialData{}, err
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

	client := payment.NewClient(cfg)

	amountInUnits := float64(req.Intent.Amount) / 100.0

	resource, err := client.Create(ctx, payment.Request{
		TransactionAmount: amountInUnits,
		Token:             req.CardToken,
		PaymentMethodID:   req.PaymentMethodID,
		Installments:      req.Installments,
		Payer: &payment.PayerRequest{
			Email: req.PayerEmail,
		},
		ApplicationFee: req.ApplicationFee,
	})
	if err != nil {
		return nil, err
	}

	return &domain.ChargeResult{
		ProviderPaymentID: fmt.Sprintf("%d", resource.ID),
		Status:            mapMPStatus(resource.Status),
	}, nil
}

func mapMPStatus(status string) domain.IntentStatus {
	switch status {
	case "approved":
		return domain.IntentStatusSucceeded
	case "rejected":
		return domain.IntentStatusFailed
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
