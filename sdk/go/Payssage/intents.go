package payssage

import (
	"context"
	"encoding/json"
	"fmt"
)

type InitiateCheckoutRequest struct {
	Amount               int64           `json:"amount"`
	Currency             string          `json:"currency"`
	Provider             string          `json:"provider"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	PaymentMethodID      string          `json:"payment_method_id,omitempty"`
	Installments         int             `json:"installments,omitempty"`
	CardToken            string          `json:"card_token,omitempty"`
	PaymentMethodType    string          `json:"payment_method_type,omitempty"`
	SellerCredentialID   string          `json:"seller_credential_id,omitempty"`
	PayerEmail           string          `json:"payer_email,omitempty"`
	IdentificationNumber string          `json:"identification_number,omitempty"`
	IdentificationType   string          `json:"identification_type,omitempty"`
}

func (c *Client) InitiateCheckout(ctx context.Context, req InitiateCheckoutRequest) (*Intent, error) {
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

type CancelPixRequest struct {
	Provider           string `json:"provider"`
	SellerCredentialID string `json:"seller_credential_id"`
}

func (c *Client) CancelPixIntent(ctx context.Context, intentID string, req CancelPixRequest) (*Intent, error) {
	var out Intent
	if err := c.do(ctx, "POST", fmt.Sprintf("/intents/%s/cancel-pix", intentID), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type ChargeRequest struct {
	SellerCredentialID string `json:"seller_credential_id"`
}

func (c *Client) Charge(ctx context.Context, intentID string, req ChargeRequest) (*Intent, error) {
	var out Intent
	if err := c.do(ctx, "POST", fmt.Sprintf("/intents/%s/charge", intentID), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
