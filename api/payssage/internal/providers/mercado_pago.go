package providers

import (
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
	"payssage/ports"
	"strconv"
	"strings"
	"time"

	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"github.com/spf13/viper"
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

func (p *MercadoPagoImpl) ExchangeCode(ctx context.Context, code, redirectURI string) (models.ProviderCredentialData, error) {
	body, err := json.Marshal(map[string]any{
		"grant_type":    "authorization_code",
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
		"test_token":    viper.GetString("TEST_MODE"),
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
	defer resp.Body.Close()

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

func (p *MercadoPagoImpl) InitiateCheckout(ctx context.Context, request *ports.InitiateCheckoutRequest) (*models.Intent, error) {
	intent, err := models.NewIntent(request.WorkspaceID, request.Amount, request.Currency, request.Provider, request.Metadata)
	if err != nil {
		return nil, err
	}

	intent.SellerCredentialID = &request.SellerCredentialID

	body := map[string]any{
		"transaction_amount":   json.Number(formatAmount(request.Amount)),
		"application_fee":      json.Number(formatAmount(calcApplicationFee(request.Amount, request.MPMarketplaceFeeBPS))),
		"installments":         request.Installments,
		"token":                request.MPCardToken,
		"payment_method_id":    request.MPPaymentMethodID,
		"external_reference":   intent.ID.String(),
		"statement_descriptor": "payssage",
		"payer": map[string]any{
			"email": request.Payer.Email,
			"identification": map[string]any{
				"type":   request.IdentificationType,
				"number": request.IdentificationNumber,
			},
		},
		"additional_info": map[string]any{
			"items": []map[string]any{
				{
					"title":      "Online Purchase",
					"quantity":   1,
					"unit_price": json.Number(formatAmount(request.Amount)),
				},
			},
		},
	}

	telemetry.Log().Info("MP Create Payment Request", zap.Any("body", body))

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	idempotencyKey, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.mercadopago.com/v1/payments", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+request.MPSellerToken)
	req.Header.Set("X-Idempotency-Key", idempotencyKey.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	telemetry.Log().Info("MP Create Payment Request Raw Body", zap.Any("body", string(rawBody)))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mercadopago create payment error %d: %s", resp.StatusCode, string(rawBody))
	}

	var mpResp struct {
		ID           int    `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	if err := json.Unmarshal(rawBody, &mpResp); err != nil {
		return nil, err
	}

	intent.MercadoPagoData = &models.MercadoPagoIntentData{
		OrderID:           strconv.Itoa(mpResp.ID),
		TransactionID:     strconv.Itoa(mpResp.ID),
		OrderStatus:       mpResp.Status,
		OrderStatusDetail: mpResp.StatusDetail,
		PaymentMethodID:   request.MPPaymentMethodID,
		PaymentMethodType: request.MPPaymentMethodType,
	}

	return intent, nil
}

func (p *MercadoPagoImpl) InitiatePixCheckout(ctx context.Context, request *ports.InitiateCheckoutRequest) (*models.Intent, error) {
	intent, err := models.NewIntent(request.WorkspaceID, request.Amount, request.Currency, request.Provider, request.Metadata)
	if err != nil {
		return nil, err
	}

	intent.SellerCredentialID = &request.SellerCredentialID

	loc := time.FixedZone("BRT", -3*60*60)
	expirationTime := time.Now().In(loc).Add(30 * time.Minute).Format("2006-01-02T15:04:05.000-07:00")

	body := map[string]any{
		"transaction_amount":   json.Number(formatAmount(request.Amount)),
		"application_fee":      json.Number(formatAmount(calcApplicationFee(request.Amount, request.MPMarketplaceFeeBPS))),
		"payment_method_id":    "pix",
		"external_reference":   intent.ID.String(),
		"date_of_expiration":   expirationTime,
		"statement_descriptor": "payssage",
		"payer": map[string]any{
			"email": request.Payer.Email,
			"identification": map[string]any{
				"type":   request.IdentificationType,
				"number": request.IdentificationNumber,
			},
		},
		"additional_info": map[string]any{
			"items": []map[string]any{
				{
					"title":      "Online Purchase",
					"quantity":   1,
					"unit_price": json.Number(formatAmount(request.Amount)),
				},
			},
		},
	}

	telemetry.Log().Info("MP Create Pix Payment Request", zap.Any("body", body))

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.mercadopago.com/v1/payments", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+request.MPSellerToken)
	req.Header.Set("X-Idempotency-Key", intent.ID.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	telemetry.Log().Info("MP Create Pix Payment Response", zap.String("body", string(rawBody)))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mercadopago create pix payment error %d: %s", resp.StatusCode, string(rawBody))
	}

	var mpResp struct {
		ID                 int64  `json:"id"`
		Status             string `json:"status"`
		StatusDetail       string `json:"status_detail"`
		PointOfInteraction struct {
			TransactionData struct {
				QRCode       string `json:"qr_code"`
				QRCodeBase64 string `json:"qr_code_base64"`
			} `json:"transaction_data"`
		} `json:"point_of_interaction"`
	}

	if err := json.Unmarshal(rawBody, &mpResp); err != nil {
		return nil, err
	}

	paymentID := fmt.Sprintf("%d", mpResp.ID)

	intent.MercadoPagoData = &models.MercadoPagoIntentData{
		OrderID:           paymentID,
		TransactionID:     paymentID,
		OrderStatus:       mpResp.Status,
		OrderStatusDetail: mpResp.StatusDetail,
		PaymentMethodID:   "pix",
		PaymentMethodType: "bank_transfer",
		PixQRCode:         mpResp.PointOfInteraction.TransactionData.QRCode,
		PixQRCodeB64:      mpResp.PointOfInteraction.TransactionData.QRCodeBase64,
	}

	return intent, nil
}

func (p *MercadoPagoImpl) CancelPixCode(ctx context.Context, paymentID string, sellerToken string) error {
	body := map[string]any{
		"status": "cancelled",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID),
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sellerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	telemetry.Log().Info("MP Cancel Pix Payment Response",
		zap.String("payment_id", paymentID),
		zap.String("body", string(rawBody)),
	)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("mercadopago cancel pix payment error %d: %s", resp.StatusCode, string(rawBody))
	}

	return nil
}

func (p *MercadoPagoImpl) Charge(ctx context.Context, request *ports.ChargeRequest) (*models.Intent, error) {
	chargeIDKey, err := uuid.NewV7()
	if err != nil {
		return nil, wrapMPError(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.mercadopago.com/v1/orders/"+request.Intent.MercadoPagoData.OrderID+"/process",
		nil,
	)
	if err != nil {
		return nil, wrapMPError(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+request.MPSellerToken)
	req.Header.Set("X-Idempotency-Key", chargeIDKey.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, wrapMPError(err)
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, wrapMPError(fmt.Errorf("mercadopago process order error %d: %s", resp.StatusCode, string(rawBody)))
	}

	var mpResp struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
		Transactions struct {
			Payments []struct {
				ID           string `json:"id"`
				Status       string `json:"status"`
				StatusDetail string `json:"status_detail"`
				PaidAmount   string `json:"paid_amount"`
			} `json:"payments"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal(rawBody, &mpResp); err != nil {
		return nil, wrapMPError(err)
	}

	intent := request.Intent
	intent.MercadoPagoData.OrderStatus = mpResp.Status
	intent.MercadoPagoData.OrderStatusDetail = mpResp.StatusDetail

	if len(mpResp.Transactions.Payments) > 0 {
		tx := mpResp.Transactions.Payments[0]
		intent.MercadoPagoData.TransactionID = tx.ID
		intent.MercadoPagoData.TransactionStatus = tx.Status
		intent.MercadoPagoData.TransactionStatusDetail = tx.StatusDetail
	}

	telemetry.Log().Info("MP Process Order Response",
		zap.String("order_id", mpResp.ID),
		zap.String("status", mpResp.Status),
		zap.String("status_detail", mpResp.StatusDetail),
	)

	return &intent, nil
}

func (p *MercadoPagoImpl) Refund(ctx context.Context, request *ports.RefundRequest) (*models.Intent, error) {
	return nil, wrapMPError(errors.New("not implemented"))
}

// -- MercadoPagoImpl internal methods --

func (p *MercadoPagoImpl) CreatePixOrder(ctx context.Context, req ports.ChargeRequest) (*models.MercadoPagoIntentData, error) {
	return nil, wrapMPError(errors.New("not implemented"))
}

func (p *MercadoPagoImpl) NormalizeStatus(status, statusDetail string) models.IntentStatus {
	switch status {
	case "processed":
		return models.IntentStatusSucceeded
	case "processing":
		return models.IntentStatusPending
	case "action_required":
		// statusDetail disambiguates — for now treat as pending
		return models.IntentStatusPending
	case "canceled":
		return models.IntentStatusCancelled
	default:
		return models.IntentStatusPending
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
