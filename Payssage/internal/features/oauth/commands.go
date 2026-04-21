package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"payssage/internal/platform/database"
	"payssage/internal/platform/telemetry"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/ports"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	intents     ports.IntentRepository
	workspaces  ports.WorkspaceRepo
	oauthStates ports.OAuthStateRepo
	credentials ports.ProviderCredentialRepo
	marketplace ports.MarketplaceConfigRepo
	providers   map[string]ports.OAuthProvider
	tx          database.TxRunner
	tracer      trace.Tracer
}

func NewCommandService(
	intents ports.IntentRepository,
	workspaces ports.WorkspaceRepo,
	oauthStates ports.OAuthStateRepo,
	credentials ports.ProviderCredentialRepo,
	marketplace ports.MarketplaceConfigRepo,
	providers map[string]ports.OAuthProvider,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		intents:     intents,
		workspaces:  workspaces,
		oauthStates: oauthStates,
		credentials: credentials,
		marketplace: marketplace,
		providers:   providers,
		tx:          tx,
		tracer:      tracer,
	}
}

func (uc *CommandService) CompleteOAuth(ctx context.Context, provider, stateToken, code, redirectURI string) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CompleteOAuth")
	defer span.End()

	oauthState, err := uc.oauthStates.Get(ctx, stateToken)
	if err != nil {
		return "", errx.Invalid("oauth_state").SetMessage("invalid or expired state")
	}

	if oauthState.Provider != provider {
		return "", errx.Invalid("oauth_state").SetMessage("provider mismatch")
	}

	p, err := uc.getProvider(provider)
	if err != nil {
		return "", err
	}

	credData, err := p.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		telemetry.Log().Info("Error exchanging codes", zap.Error(err))
		return "", errx.Internal("oauth").SetMessage(fmt.Sprintf("failed to exchange code: %s", err.Error()))
	}

	cred, err := uc.credentials.Create(ctx, contracts.ProviderCredential{
		WorkspaceID: oauthState.WorkspaceID,
		Provider:    provider,
		Credentials: credData,
	})
	if err != nil {
		return "", err
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("redirect_url", oauthState.FinalRedirectURL)
	u.RawQuery = q.Encode()

	FinalRedirectURL := u.String()

	telemetry.Log().Info("Exchange result",
		zap.String("access_token_prefix", credData.AccessToken[:20]),
		zap.Int("user_id", credData.ProviderUserID),
		zap.String("provider", provider),
		zap.String("flow", oauthState.Flow),
		zap.String("credential_id", cred.ID.String()),
		zap.String("url", oauthState.FinalRedirectURL),
	)

	// if setup flow + marketplace, auto-create marketplace config
	if oauthState.Flow == contracts.OAuthFlowSetup && oauthState.IsMarketplace {
		existing, err := uc.marketplace.Get(ctx, oauthState.WorkspaceID, cred.ID)
		if err != nil && !errx.IsKind(err, "not_found") {
			return "", err
		}
		if existing != nil {
			if provider != existing.Provider {
				return "", errx.Invalid("marketplace_config").SetMessage("cannot change provider of a config through OAuth")
			}
			_, err = uc.marketplace.Update(ctx, contracts.MarketplaceConfig{
				WorkspaceID:  oauthState.WorkspaceID,
				CredentialID: cred.ID,
				FeeBps:       oauthState.FeeBps,
			})
		} else {
			_, err = uc.marketplace.Create(ctx, contracts.MarketplaceConfig{
				WorkspaceID:  oauthState.WorkspaceID,
				Provider:     provider,
				CredentialID: cred.ID,
				FeeBps:       oauthState.FeeBps,
			})
		}
		if err != nil {
			return "", err
		}
	} else {

	}

	_ = uc.oauthStates.Delete(ctx, stateToken)

	switch oauthState.Flow {
	case contracts.OAuthFlowSetup:
		return fmt.Sprintf("%s&provider=%s&status=success", FinalRedirectURL, provider), nil
	case contracts.OAuthFlowConnect:
		return fmt.Sprintf("%s&credential_id=%s&provider=%s&public_key=%s", FinalRedirectURL, cred.ID, provider, cred.Credentials.PublicKey), nil
	default:
		return FinalRedirectURL, nil
	}
}

type ConnectSellerInput struct {
	WorkspaceName       string
	Provider            string
	ProviderRedirectURL string
	FinalRedirectURL    string
}

func (uc *CommandService) ConnectSeller(ctx context.Context, req ConnectSellerInput) (string, string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.ConnectSeller")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return "", "", err
	}

	workspace, err := uc.workspaces.GetByID(ctx, ws.ID)
	if err != nil {
		return "", "", err
	}

	_, err = uc.marketplace.GetByProvider(ctx, workspace.ID, req.Provider)
	if err != nil {
		return "", "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", "", errx.Internal("oauth_state").SetCause(err)
	}

	_, err = uc.oauthStates.Create(ctx, contracts.OAuthState{
		State:            stateToken,
		WorkspaceID:      workspace.ID,
		Provider:         req.Provider,
		Flow:             contracts.OAuthFlowConnect,
		IsMarketplace:    false,
		FeeBps:           0,
		FinalRedirectURL: req.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", "", err
	}

	provider, _ := uc.getProvider(req.Provider)
	return provider.BuildAuthURL(stateToken, req.ProviderRedirectURL), req.FinalRedirectURL, nil
}

func (uc *CommandService) DeleteMarketplaceConfig(ctx context.Context, workspaceName string, credentialID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DeleteMarketplaceConfig")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	// verify the credential belongs to this workspace before deleting
	cred, err := uc.credentials.GetByID(ctx, credentialID)
	if err != nil {
		return err
	}
	if cred.WorkspaceID != workspace.ID {
		return errx.Forbidden("credential").SetMessage("credential does not belong to this workspace")
	}

	// If this config was backing a marketplace revoke it
	_, err = uc.credentials.Revoke(ctx, credentialID, workspace.ID)
	if err != nil {
		return err
	}

	return uc.marketplace.Delete(ctx, workspace.ID, credentialID)
}

func (uc *CommandService) DisconnectCredential(ctx context.Context, credentialID uuid.UUID) (*contracts.ProviderCredential, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DisconnectCredential")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	return uc.credentials.Revoke(ctx, credentialID, workspace.ID)
}

func (uc *CommandService) RevokeCredential(ctx context.Context, workspaceName string, credentialID uuid.UUID) (*contracts.ProviderCredential, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RevokeCredential")
	defer span.End()

	subject, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, subject.ID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.credentials.Revoke(ctx, credentialID, workspace.ID)
	if err != nil {
		return nil, err
	}

	// if this credential was backing a marketplace config, remove it
	_ = uc.marketplace.Delete(ctx, workspace.ID, credentialID)

	return cred, nil
}

type SetMarketplaceConfigInput struct {
	WorkspaceName string
	CredentialID  uuid.UUID
	FeeBps        int
}

func (uc *CommandService) SetMarketplaceConfig(ctx context.Context, req SetMarketplaceConfigInput) (*contracts.MarketplaceConfig, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.SetMarketplaceConfig")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	// verify credential belongs to this workspace
	cred, err := uc.credentials.GetByID(ctx, req.CredentialID)
	if err != nil {
		return nil, err
	}
	if cred.WorkspaceID != workspace.ID {
		return nil, errx.Forbidden("credential").SetMessage("credential does not belong to this workspace")
	}

	existing, err := uc.marketplace.Get(ctx, workspace.ID, req.CredentialID)
	if err != nil && !errx.IsKind(err, "not_found") {
		return nil, err
	}

	if existing != nil {
		return uc.marketplace.Update(ctx, contracts.MarketplaceConfig{
			WorkspaceID:  workspace.ID,
			CredentialID: req.CredentialID,
			FeeBps:       req.FeeBps,
		})
	}

	return uc.marketplace.Create(ctx, contracts.MarketplaceConfig{
		WorkspaceID:  workspace.ID,
		CredentialID: req.CredentialID,
		FeeBps:       req.FeeBps,
	})
}

type SetupProviderInput struct {
	WorkspaceName       string
	Provider            string
	IsMarketplace       bool
	FeeBps              int
	ProviderRedirectURL string
	FinalRedirectURL    string
}

func (uc *CommandService) SetupProvider(ctx context.Context, req SetupProviderInput) (string, string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.SetupProvider")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return "", "", err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return "", "", err
	}

	if _, err := uc.getProvider(req.Provider); err != nil {
		return "", "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", "", errx.Internal("oauth_state").SetCause(err)
	}

	_, err = uc.oauthStates.Create(ctx, contracts.OAuthState{
		State:            stateToken,
		WorkspaceID:      workspace.ID,
		Provider:         req.Provider,
		Flow:             contracts.OAuthFlowSetup,
		IsMarketplace:    req.IsMarketplace,
		FeeBps:           req.FeeBps,
		FinalRedirectURL: req.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", "", err
	}

	provider, _ := uc.getProvider(req.Provider)
	return provider.BuildAuthURL(stateToken, req.ProviderRedirectURL), req.FinalRedirectURL, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (uc *CommandService) getProvider(name string) (ports.OAuthProvider, error) {
	p, ok := uc.providers[name]
	if !ok {
		return nil, errx.Invalid("provider").SetMessage(fmt.Sprintf("unsupported provider: %s", name))
	}
	return p, nil
}
