// Package oauth is the oauth feature: its repo, handlers and operations.
// Build wires it.
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/internal/providers"
	"payssage/models"
	"payssage/ports"
	idx "sdk/identityx"
	"strconv"
	"time"
	"uuid"
)

type Deps struct {
	Wallets    ports.WalletRepo
	Orgs       ports.OrganizationRepo
	OAuth      ports.OAuthStateRepo
	Collectors ports.CollectorRepo
	Sellers    ports.SellerRepo
	Authz      *authz.Service
	Registry   providers.PaymentProviders
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) Callback(ctx context.Context, providerStr, code, stateStr string) (string, error) {
	state, err := o.OAuth.Get(ctx, stateStr)
	if err != nil {
		return "", err
	}
	if state.Provider != providerStr {
		return "", errs.ErrInvalidProvider
	}

	provider, err := providers.FromString(providerStr)
	if err != nil {
		return "", err
	}

	oauthProvider := o.Registry.OAuth[provider]
	credentialData, err := oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	var credentialID string
	switch state.Flow {
	case models.OAuthFlowCollector:
		collector, err := o.Collectors.Create(ctx, models.Collector{
			OwnerID:        state.OwnerID,
			OrganizationID: state.OrganizationID,
			Provider:       providerStr,
			ProviderUserID: strconv.Itoa(credentialData.ProviderUserID),
			Credentials:    marshalCredentials(credentialData),
		})
		if err != nil {
			return "", err
		}
		credentialID = collector.ID.String()
	case models.OAuthFlowSeller:
		if state.WalletID == nil {
			return "", errs.ErrSellerFlowNeedsWallet
		}
		seller, err := o.Sellers.Create(ctx, models.Seller{
			WalletID:       *state.WalletID,
			Provider:       providerStr,
			ProviderUserID: strconv.Itoa(credentialData.ProviderUserID),
			Credentials:    marshalCredentials(credentialData),
		})
		if err != nil {
			return "", err
		}
		credentialID = seller.ID.String()
	default:
		return "", errs.ErrInvalidFlow
	}

	return fmt.Sprintf("%s?credential_id=%s&public_key=%s", state.FinalRedirectURL, credentialID, credentialData.PublicKey), nil
}

func marshalCredentials(data models.ProviderCredentialData) []byte {
	b, _ := json.Marshal(data) //nolint:gosec
	return b
}

func (o *Operations) Connect(ctx context.Context, payload models.ConnectInput) (string, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return "", err
	}

	var walletID *uuid.UUID
	var organizationID *uuid.UUID

	switch payload.Flow {
	case models.OAuthFlowSeller:
		var wallet *models.Wallet
		wallet, err = o.Wallets.GetByID(ctx, *payload.WalletID)
		if err != nil {
			return "", err
		}
		var org *models.Organization
		if wallet.OrganizationID != nil {
			org, err = o.Orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return "", err
			}
		}
		if org != nil {
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
		} else {
			_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleMember)
		}
		if err != nil {
			return "", err
		}
		walletID = &wallet.ID

	case models.OAuthFlowCollector:
		if payload.OrganizationID != nil {
			var org *models.Organization
			org, err = o.Orgs.GetByID(ctx, *payload.OrganizationID)
			if err != nil {
				return "", err
			}
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return "", err
			}
			organizationID = payload.OrganizationID
		}

	default:
		return "", errs.ErrInvalidFlow
	}

	provider, err := providers.FromString(payload.Provider)
	if err != nil {
		return "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", err
	}

	_, err = o.OAuth.Create(ctx, models.OAuthState{
		State:            stateToken,
		WalletID:         walletID,
		OrganizationID:   organizationID,
		OwnerID:          ident.Sub.ID,
		Provider:         provider.String(),
		Flow:             payload.Flow,
		FinalRedirectURL: payload.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	// The provider redirect URI is Payssage's own callback (config-owned,
	// D7) — the provider builds the auth URL with it internally.
	return o.Registry.OAuth[provider].BuildAuthURL(stateToken), nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (o *Operations) Revoke(ctx context.Context, input models.RevokeInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	switch input.Flow {
	case models.OAuthFlowCollector:
		var collector *models.Collector
		collector, err = o.Collectors.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if collector.Provider != input.Provider {
			return errs.ErrProviderMismatch
		}
		if collector.OrganizationID != nil {
			var org *models.Organization
			org, err = o.Orgs.GetByID(ctx, *collector.OrganizationID)
			if err != nil {
				return err
			}
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if collector.OwnerID != ident.Sub.ID {
			return errs.ErrNotResourceOwner
		}
		return o.Collectors.Revoke(ctx, input.ID)

	case models.OAuthFlowSeller:
		var seller *models.Seller
		seller, err = o.Sellers.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if seller.Provider != input.Provider {
			return errs.ErrProviderMismatch
		}
		var wallet *models.Wallet
		wallet, err = o.Wallets.GetByID(ctx, seller.WalletID)
		if err != nil {
			return err
		}
		if wallet.OrganizationID != nil {
			var org *models.Organization
			org, err = o.Orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return err
			}
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if wallet.OwnerID != ident.Sub.ID {
			return errs.ErrNotResourceOwner
		}
		return o.Sellers.Revoke(ctx, input.ID)

	default:
		return errs.ErrInvalidFlow
	}
}
