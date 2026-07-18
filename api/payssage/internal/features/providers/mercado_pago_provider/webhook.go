package mercado_pago_provider

//func (p *Provider) VerifySignature(xSignature, xRequestID, dataID string) bool {
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
//	mac := hmac.New(sha256.New, []byte(p.cfg.MpWebhookSecret))
//	mac.Write([]byte(manifest))
//	computed := hex.EncodeToString(mac.Sum(nil))
//
//	return hmac.Equal([]byte(computed), []byte(hash))
//}
