package payssage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateWalletRequest struct {
	Name           string     `json:"name"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
}

// CreateWallet creates a wallet owned by the authenticated (service) actor,
// optionally scoped to an organization.
func (c *Client) CreateWallet(ctx context.Context, req CreateWalletRequest) (*Wallet, error) {
	var out Wallet
	err := c.do(ctx, "POST", "/wallets", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWallet fetches a wallet by id. Used for fail-fast boot checks (split
// 2) and direct lookups.
func (c *Client) GetWallet(ctx context.Context, walletID uuid.UUID) (*Wallet, error) {
	var out Wallet
	err := c.do(ctx, "GET", fmt.Sprintf("/wallets/%s", walletID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// SetWalletFee sets a wallet's marketplace fee in basis points (bps;
// 100 bps = 1%).
func (c *Client) SetWalletFee(ctx context.Context, walletID uuid.UUID, feeBps int) error {
	return c.do(ctx, "PATCH", fmt.Sprintf("/wallets/%s/fee", walletID), map[string]int{"fee_bps": feeBps}, nil)
}

// ListWalletSellers lists the sellers (provider accounts) bound to a wallet.
func (c *Client) ListWalletSellers(ctx context.Context, walletID uuid.UUID) ([]Seller, error) {
	var out []Seller
	err := c.do(ctx, "GET", fmt.Sprintf("/wallets/%s/sellers", walletID), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
