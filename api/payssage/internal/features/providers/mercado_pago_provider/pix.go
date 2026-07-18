package mercado_pago_provider

//func (p *Provider) InitiatePixCheckout(ctx context.Context, intent *models.Intent, sellerToken string, marketplaceFeeBPS int, payerEmail, identType, identNum string) (*models.Intent, error) {
//	loc := time.FixedZone("BRT", -3*60*60)
//	expirationTime := time.Now().In(loc).Add(30 * time.Minute).Format("2006-01-02T15:04:05.000-07:00")
//
//	body := map[string]any{
//		"transaction_amount":   json.Number(formatAmount(intent.AmountCents)),
//		"application_fee":      json.Number(formatAmount(calcApplicationFee(intent.AmountCents, marketplaceFeeBPS))),
//		"payment_method_id":    "pix",
//		"external_reference":   intent.ID.String(),
//		"date_of_expiration":   expirationTime,
//		"statement_descriptor": "payssage",
//		"payer": map[string]any{
//			"email": payerEmail,
//			"identification": map[string]any{
//				"type":   identType,
//				"number": identNum,
//			},
//		},
//		"additional_info": map[string]any{
//			"items": []map[string]any{
//				{
//					"title":      "Online Purchase",
//					"quantity":   1,
//					"unit_price": json.Number(formatAmount(intent.AmountCents)),
//				},
//			},
//		},
//	}
//
//	telemetry.Log().Info("MP Create Pix Payment Request", zap.Any("body", body))
//
//	bodyBytes, err := json.Marshal(body)
//	if err != nil {
//		return nil, err
//	}
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.mercadopago.com/v1/payments", bytes.NewReader(bodyBytes))
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+sellerToken)
//	req.Header.Set("X-Idempotency-Key", intent.ID.String())
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	rawBody, _ := io.ReadAll(resp.Body)
//
//	telemetry.Log().Info("MP Create Pix Payment Response", zap.String("body", string(rawBody)))
//
//	if resp.StatusCode >= 400 {
//		return nil, fmt.Errorf("mercadopago create pix payment error %d: %s", resp.StatusCode, string(rawBody))
//	}
//
//	var mpResp struct {
//		ID                 int64  `json:"id"`
//		Status             string `json:"status"`
//		StatusDetail       string `json:"status_detail"`
//		PointOfInteraction struct {
//			TransactionData struct {
//				QRCode       string `json:"qr_code"`
//				QRCodeBase64 string `json:"qr_code_base64"`
//			} `json:"transaction_data"`
//		} `json:"point_of_interaction"`
//	}
//
//	if err := json.Unmarshal(rawBody, &mpResp); err != nil {
//		return nil, err
//	}
//
//	paymentID := fmt.Sprintf("%d", mpResp.ID)
//
//	data := &models.MercadoPagoIntentData{
//		OrderID:           paymentID,
//		TransactionID:     paymentID,
//		OrderStatus:       mpResp.Status,
//		OrderStatusDetail: mpResp.StatusDetail,
//		PaymentMethodID:   "pix",
//		PaymentMethodType: "bank_transfer",
//		PixQRCode:         mpResp.PointOfInteraction.TransactionData.QRCode,
//		PixQRCodeB64:      mpResp.PointOfInteraction.TransactionData.QRCodeBase64,
//	}
//
//	dataBytes, err := json.Marshal(data)
//	if err != nil {
//		return nil, err
//	}
//
//	intent.ProviderData = dataBytes
//
//	return intent, nil
//}

//func (p *Provider) CancelPixCode(ctx context.Context, paymentID, sellerToken string) error {
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
