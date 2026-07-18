package mercado_pago_provider

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
