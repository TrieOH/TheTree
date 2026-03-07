package paymentsSDK

import (
	"context"
	"encoding/json"
	"fmt"
)

type CreateIntentRequest struct {
	Amount   int64           `json:"amount"`
	Currency string          `json:"currency"`
	Provider string          `json:"provider"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

func (c *Client) CreateIntent(ctx context.Context, req CreateIntentRequest) (*Intent, error) {
	var out Intent
	if err := c.do(ctx, "POST", "/intents", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetIntent(ctx context.Context, intentID string) (*Intent, error) {
	var out Intent
	if err := c.do(ctx, "GET", fmt.Sprintf("/intents/%s", intentID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListIntents(ctx context.Context) ([]Intent, error) {
	var out []Intent
	if err := c.do(ctx, "GET", "/intents", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CancelIntent(ctx context.Context, intentID string) (*Intent, error) {
	var out Intent
	if err := c.do(ctx, "POST", fmt.Sprintf("/intents/%s/cancel", intentID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
