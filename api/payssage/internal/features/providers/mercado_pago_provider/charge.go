package mercado_pago_provider

//func (p *Provider) Charge(ctx context.Context, intent *models.Intent, sellerToken string) (*models.Intent, error) {
//	var data models.MercadoPagoIntentData
//	if err := json.Unmarshal(intent.ProviderData, &data); err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	chargeIDKey, err := uuid.NewV7()
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
//		"https://api.mercadopago.com/v1/orders/"+data.OrderID+"/process",
//		nil,
//	)
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+sellerToken)
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
//	data.OrderStatus = mpResp.Status
//	data.OrderStatusDetail = mpResp.StatusDetail
//
//	if len(mpResp.Transactions.Payments) > 0 {
//		tx := mpResp.Transactions.Payments[0]
//		data.TransactionID = tx.ID
//		data.TransactionStatus = tx.Status
//		data.TransactionStatusDetail = tx.StatusDetail
//	}
//
//	dataBytes, err := json.Marshal(data)
//	if err != nil {
//		return nil, wrapMPError(err)
//	}
//
//	intent.ProviderData = dataBytes
//
//	telemetry.Log().Info("MP Process Order Response",
//		zap.String("order_id", mpResp.ID),
//		zap.String("status", mpResp.Status),
//		zap.String("status_detail", mpResp.StatusDetail),
//	)
//
//	return intent, nil
//}
