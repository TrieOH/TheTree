package mercado_pago_provider

//
//func VerifyMercadoPagoSignature(xSignature, xRequestID, dataID, secret string) bool {
//	var ts, hash string
//	for _, part := range strings.Split(xSignature, ",") {
//		kv := strings.SplitN(part, "=", 2)
//		if len(kv) != 2 {
//			continue
//		}
//		key := strings.TrimSpace(kv[0])
//		val := strings.TrimSpace(kv[1])
//		switch key {
//		case "ts":
//			ts = val
//		case "v1":
//			hash = val
//		}
//	}
//
//	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, ts)
//
//	mac := hmac.New(sha256.New, []byte(secret))
//	mac.Write([]byte(manifest))
//	computed := hex.EncodeToString(mac.Sum(nil))
//
//	return hmac.Equal([]byte(computed), []byte(hash))
//}
//
//// -- PAL methods --
//
//
//func (p *MercadoPagoImpl) CancelPixCode(ctx context.Context, paymentID string, sellerToken string) error {
//	body := map[string]any{
//		"status": "cancelled",
//	}
//
//	bodyBytes, err := json.Marshal(body)
//	if err != nil {
//		return err
//	}
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
//		fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID),
//		bytes.NewReader(bodyBytes),
//	)
//	if err != nil {
//		return err
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+sellerToken)
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	rawBody, _ := io.ReadAll(resp.Body)
//
//	telemetry.Log().Info("MP Cancel Pix Payment Response",
//		zap.String("payment_id", paymentID),
//		zap.String("body", string(rawBody)),
//	)
//
//	if resp.StatusCode >= 400 {
//		return fmt.Errorf("mercadopago cancel pix payment error %d: %s", resp.StatusCode, string(rawBody))
//	}
//
//	return nil
//}
//
//func (p *MercadoPagoImpl) Charge(ctx context.Context, request *ports.ChargeRequest) (*models.Intent, error) {
//	chargeIDKey, err := uuid.NewV7()
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
//		"https://api.mercadopago.com/v1/orders/"+request.Intent.MercadoPagoData.OrderID+"/process",
//		nil,
//	)
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+request.MPSellerToken)
//	req.Header.Set("X-Idempotency-Key", chargeIDKey.String())
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//	defer resp.Body.Close()
//
//	rawBody, _ := io.ReadAll(resp.Body)
//
//	if resp.StatusCode >= 400 {
//		return nil, wrapMPError(fmt.Errorf("mercadopago process order error %d: %s", resp.StatusCode, string(rawBody)))
//	}
//
//	var mpResp struct {
//		ID           string `json:"id"`
//		Status       string `json:"status"`
//		StatusDetail string `json:"status_detail"`
//		Transactions struct {
//			Payments []struct {
//				ID           string `json:"id"`
//				Status       string `json:"status"`
//				StatusDetail string `json:"status_detail"`
//				PaidAmount   string `json:"paid_amount"`
//			} `json:"payments"`
//		} `json:"transactions"`
//	}
//
//	if err := json.Unmarshal(rawBody, &mpResp); err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	intent := request.Intent
//	intent.MercadoPagoData.OrderStatus = mpResp.Status
//	intent.MercadoPagoData.OrderStatusDetail = mpResp.StatusDetail
//
//	if len(mpResp.Transactions.Payments) > 0 {
//		tx := mpResp.Transactions.Payments[0]
//		intent.MercadoPagoData.TransactionID = tx.ID
//		intent.MercadoPagoData.TransactionStatus = tx.Status
//		intent.MercadoPagoData.TransactionStatusDetail = tx.StatusDetail
//	}
//
//	telemetry.Log().Info("MP Process Order Response",
//		zap.String("order_id", mpResp.ID),
//		zap.String("status", mpResp.Status),
//		zap.String("status_detail", mpResp.StatusDetail),
//	)
//
//	return &intent, nil
//}
//
//func (p *MercadoPagoImpl) Refund(ctx context.Context, request *ports.RefundRequest) (*models.Intent, error) {
//	return nil, wrapMPError(errors.New("not implemented"))
//}
//
//// -- MercadoPagoImpl internal methods --
//
//func (p *MercadoPagoImpl) CreatePixOrder(ctx context.Context, req ports.ChargeRequest) (*models.MercadoPagoIntentData, error) {
//	return nil, wrapMPError(errors.New("not implemented"))
//}
//
//func (p *MercadoPagoImpl) NormalizeStatus(status, statusDetail string) models.IntentStatus {
//	switch status {
//	case "processed":
//		return models.IntentStatusSucceeded
//	case "processing":
//		return models.IntentStatusPending
//	case "action_required":
//		// statusDetail disambiguates — for now treat as pending
//		return models.IntentStatusPending
//	case "canceled":
//		return models.IntentStatusCancelled
//	default:
//		return models.IntentStatusPending
//	}
//}
//
//// -- helpers --
//
//// formatAmount converts int64 centavos to MP's decimal string: 1050 → "10.50"
//func formatAmount(centavos int64) string {
//	return fmt.Sprintf("%d.%02d", centavos/100, centavos%100)
//}
//
//// parseAmount is the inverse: "10.50" → 1050
//func parseAmount(s string) int64 {
//	f, _ := strconv.ParseFloat(s, 64)
//	return int64(f * 100)
//}
//
//func nullableString(s string) *string {
//	if s == "" {
//		return nil
//	}
//	return &s
//}
//
//func wrapMPError(err error) error {
//	return fmt.Errorf("mercadopago: %w", err)
//}
//
//func extractOrderID(payload []byte) string       { /* parse JSON "id" field */ return "" }
//func extractTransactionID(payload []byte) string { /* parse JSON transaction id */ return "" }
//func extractExternalRef(payload []byte) string   { /* parse JSON "external_reference" */ return "" }
//
//func calcApplicationFee(amountCents int64, feeBps int) int64 {
//	return (amountCents*int64(feeBps) + 5000) / 10000
//}
