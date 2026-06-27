package jobs

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"payssage/ports"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type DeliverWebhookArgs struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
	EndpointID uuid.UUID `json:"endpoint_id"`
	URL        string    `json:"url"`
	Secret     string    `json:"secret"`
	Payload    []byte    `json:"payload"`
}

func (DeliverWebhookArgs) Kind() string { return "webhook.deliver" }

type DeliverWebhookWorker struct {
	river.WorkerDefaults[DeliverWebhookArgs]
	deliveries ports.WebhookDeliveryRepo
}

func NewDeliverWebhookWorker(deliveries ports.WebhookDeliveryRepo) *DeliverWebhookWorker {
	return &DeliverWebhookWorker{deliveries: deliveries}
}

func (w *DeliverWebhookWorker) Work(ctx context.Context, job *river.Job[DeliverWebhookArgs]) error {
	args := job.Args

	mac := hmac.New(sha256.New, []byte(args.Secret))
	mac.Write(args.Payload)
	sig := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, args.URL, bytes.NewReader(args.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Payssage-Signature", sig)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if _, updateErr := w.deliveries.IncrementAttempt(context.Background(), args.DeliveryID); updateErr != nil {
			return updateErr
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if _, err := w.deliveries.MarkDelivered(context.Background(), args.DeliveryID); err != nil {
			return err
		}
		return nil
	}

	if _, err := w.deliveries.IncrementAttempt(context.Background(), args.DeliveryID); err != nil {
		return err
	}

	if job.Attempt >= job.MaxAttempts {
		if _, err := w.deliveries.MarkFailed(context.Background(), args.DeliveryID); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("endpoint returned %d", resp.StatusCode)
}
