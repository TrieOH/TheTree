package async

import (
	"TriePayments/internal/core/domain"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
)

func (h *AsynqHandlers) HandleDeliverWebhook(ctx context.Context, t *asynq.Task) error {
	var p domain.DeliverWebhookPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	// sign payload with HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(p.Secret))
	mac.Write(p.Payload)
	sig := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, bytes.NewReader(p.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TriePayments-Signature", sig)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = h.deliveries.IncrementAttempt(context.Background(), p.DeliveryID)
		return err // asynq retries
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_, _ = h.deliveries.MarkDelivered(context.Background(), p.DeliveryID)
		return nil
	}

	_, _ = h.deliveries.IncrementAttempt(context.Background(), p.DeliveryID)

	// if max retries exhausted asynq will stop — mark as failed
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		_, _ = h.deliveries.MarkFailed(context.Background(), p.DeliveryID)
		return nil
	}

	return fmt.Errorf("endpoint returned %d", resp.StatusCode)
}
