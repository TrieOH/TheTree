package domain

import (
	"TriePayments/internal/plataform/telemetry"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"github.com/mercadopago/sdk-go/pkg/order"
	"go.uber.org/zap"
)

const mpAuthURL = "https://auth.mercadopago.com/authorization"

type MercadoPagoImpl struct {
	clientID      string
	accessToken   string
	clientSecret  string
	redirectURI   string
	webhookSecret string
	oauthClient   oauth.Client
}

func MustMP(clientID, accessToken, clientSecret, redirectURI, webhookSecret string) *MercadoPagoImpl {
	if clientID == "" || accessToken == "" || clientSecret == "" || redirectURI == "" || webhookSecret == "" {
		log.Fatal("clientID, accessToken, clientSecret, redirectURI and webhookSecret are required")
	}
	mp, err := NewMercadoPagoProvider(clientID, accessToken, clientSecret, redirectURI, webhookSecret)
	if err != nil {
		log.Fatal("NewMercadoPagoProvider: ", err)
	}

	return mp
}

func NewMercadoPagoProvider(clientID, accessToken, clientSecret, redirectURI, webhookSecret string) (*MercadoPagoImpl, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &MercadoPagoImpl{
		clientID:      clientID,
		accessToken:   accessToken,
		clientSecret:  clientSecret,
		redirectURI:   redirectURI,
		webhookSecret: webhookSecret,
		oauthClient:   oauth.NewClient(cfg),
	}, nil
}

func (p *MercadoPagoImpl) BuildAuthURL(state, redirectURI string) string {
	return p.oauthClient.GetAuthorizationURL(p.clientID, redirectURI, state)
}

func (p *MercadoPagoImpl) ExchangeCode(ctx context.Context, code, redirectURI string) (ProviderCredentialData, error) {
	body, err := json.Marshal(map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
	})
	if err != nil {
		return ProviderCredentialData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.mercadopago.com/oauth/token",
		bytes.NewReader(body),
	)
	if err != nil {
		return ProviderCredentialData{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ProviderCredentialData{}, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProviderCredentialData{}, err
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
		return ProviderCredentialData{}, err
	}

	if result.AccessToken == "" {
		return ProviderCredentialData{}, fmt.Errorf("MP token exchange failed: %s", string(rawBody))
	}

	return ProviderCredentialData{
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
		ProviderUserID: result.UserID,
		Nickname:       result.Nickname,
	}, nil
}

func VerifyMercadoPagoSignature(xSignature, xRequestID, dataID, secret string) bool {
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

// -- PAL methods --

func (p *MercadoPagoImpl) InitiateCheckout(ctx context.Context, request *InitiateCheckoutRequest) (*Intent, error) {
	intent, err := NewIntent(request.WorkspaceID, request.Amount, request.Currency, request.Provider, request.Metadata)
	if err != nil {
		return nil, err
	}

	cfg, err := config.New(request.MPSellerToken)
	if err != nil {
		return nil, wrapMPError(err)
	}

	client := order.NewClient(cfg)

	mpOrder, err := client.Create(ctx, order.Request{
		Type:              "online",
		TotalAmount:       formatAmount(request.Amount),
		ExternalReference: "tp_" + intent.ID.String(),
		ProcessingMode:    "automatic",
		MarketPlaceFee:    formatAmount(calcApplicationFee(request.Amount, request.MPMarketplaceFeeBPS)),
		Currency:          request.Currency,
		Transactions: &order.TransactionRequest{
			Payments: []order.PaymentRequest{
				{
					Amount:         formatAmount(request.Amount),
					ExpirationTime: "",
					PaymentMethod: &order.PaymentMethodRequest{
						ID:           request.MPPaymentMethodID,
						Type:         request.MPPaymentMethodType,
						Token:        request.MPPayerToken,
						Installments: request.Installments,
					},
					AutomaticPayments: nil,
					StoredCredential:  nil,
					SubscriptionData:  nil,
				},
			},
		},
		Payer: &order.PayerRequest{
			Email: request.Payer.Email,
		},
	})
	if err != nil {
		return nil, wrapMPError(err)
	}

	intent.MercadoPagoData = &MercadoPagoIntentData{
		OrderID:           mpOrder.ID,
		OrderStatus:       mpOrder.Status,
		OrderStatusDetail: mpOrder.StatusDetail,
	}

	return intent, nil
}

func (p *MercadoPagoImpl) Charge(ctx context.Context, request *ChargeRequest) (*Intent, error) {
	return nil, wrapMPError(errors.New("not implemented"))
}

func (p *MercadoPagoImpl) Refund(ctx context.Context, request *RefundRequest) (*Intent, error) {
	return nil, wrapMPError(errors.New("not implemented"))
}

// -- MercadoPagoImpl internal methods --

func (p *MercadoPagoImpl) CreatePixOrder(ctx context.Context, req ChargeRequest) (*MercadoPagoIntentData, error) {
	return nil, wrapMPError(errors.New("not implemented"))
}

func (p *MercadoPagoImpl) NormalizeStatus(status, statusDetail string) IntentStatus {
	switch status {
	case "processed":
		return IntentStatusSucceeded
	case "processing":
		return IntentStatusPending
	case "action_required":
		// statusDetail disambiguates — for now treat as pending
		return IntentStatusPending
	case "canceled":
		return IntentStatusCancelled
	default:
		return IntentStatusPending
	}
}

// -- helpers --

// formatAmount converts int64 centavos to MP's decimal string: 1050 → "10.50"
func formatAmount(centavos int64) string {
	return fmt.Sprintf("%d.%02d", centavos/100, centavos%100)
}

// parseAmount is the inverse: "10.50" → 1050
func parseAmount(s string) int64 {
	f, _ := strconv.ParseFloat(s, 64)
	return int64(f * 100)
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func wrapMPError(err error) error {
	return fmt.Errorf("mercadopago: %w", err)
}

func extractOrderID(payload []byte) string       { /* parse JSON "id" field */ return "" }
func extractTransactionID(payload []byte) string { /* parse JSON transaction id */ return "" }
func extractExternalRef(payload []byte) string   { /* parse JSON "external_reference" */ return "" }

func calcApplicationFee(amountCents int64, feeBps int) int64 {
	return (amountCents*int64(feeBps) + 5000) / 10000
}
