package payssage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RegisterWebhookEndpointRequest struct {
	URL string `json:"url"`
}

func (c *Client) RegisterWebhookEndpoint(ctx context.Context, workspaceName, url string) (*CreateWebhookEndpointResponse, error) {
	var out CreateWebhookEndpointResponse
	if err := c.do(ctx, "POST", fmt.Sprintf("/workspaces/%s/webhooks", workspaceName), RegisterWebhookEndpointRequest{URL: url}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListWebhookEndpoints(ctx context.Context, workspaceName string) ([]WebhookEndpoint, error) {
	var out []WebhookEndpoint
	if err := c.do(ctx, "GET", fmt.Sprintf("/workspaces/%s/webhooks", workspaceName), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) DeleteWebhookEndpoint(ctx context.Context, workspaceName, endpointID string) error {
	return c.do(ctx, "DELETE", fmt.Sprintf("/workspaces/%s/webhooks/%s", workspaceName, endpointID), nil, nil)
}

// VerifyWebhookSignature verifies the X-Payssage-Signature header on an inbound delivery.
// Call this in your webhook receiver handler.
func VerifyWebhookSignature(r *http.Request, secret string) (*WebhookPayload, error) {
	sig := r.Header.Get("X-Payssage-Signature")
	if sig == "" {
		return nil, errors.New("payssage: missing X-Payssage-Signature header")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("payssage: read body: %w", err)
	}

	// decode signature from hex
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, fmt.Errorf("payssage: invalid signature encoding")
	}

	// compute expected HMAC
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)

	// constant-time comparison
	if !hmac.Equal(expected, sigBytes) {
		return nil, errors.New("payssage: invalid signature")
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("payssage: decode payload: %w", err)
	}

	return &payload, nil
}
