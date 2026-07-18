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
//// -- MercadoPagoImpl internal methods --
//
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
