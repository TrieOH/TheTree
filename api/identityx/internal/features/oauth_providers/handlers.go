package oauth_providers

import (
	"context"

	"IdentityX/errs"
	"IdentityX/internal/features/oauth_providers/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CreateProjectOAuthProvider(ctx context.Context, req openapi.CreateProjectOAuthProviderRequestObject) (openapi.CreateProjectOAuthProviderResponseObject, error) {
	created, err := h.ops.Create(ctx, models.CreateOAuthProviderInput{
		ProjectID:    req.ProjectId,
		Provider:     req.Body.Provider,
		ClientID:     req.Body.ClientId,
		ClientSecret: req.Body.ClientSecret,
		CallbackURL:  req.Body.CallbackUrl,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProjectOAuthProvider201JSONResponse(*created), nil
}

func (h *Handlers) DeleteOAuthProvider(ctx context.Context, req openapi.DeleteOAuthProviderRequestObject) (openapi.DeleteOAuthProviderResponseObject, error) {
	err := h.ops.Delete(ctx, req.OauthProviderId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOAuthProvider204Response{}, nil
}

func (h *Handlers) DisableOAuthProvider(ctx context.Context, req openapi.DisableOAuthProviderRequestObject) (openapi.DisableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, false)
	if err != nil {
		return nil, err
	}
	return openapi.DisableOAuthProvider200JSONResponse(*provider), nil
}

func (h *Handlers) EnableOAuthProvider(ctx context.Context, req openapi.EnableOAuthProviderRequestObject) (openapi.EnableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, true)
	if err != nil {
		return nil, err
	}
	return openapi.EnableOAuthProvider200JSONResponse(*provider), nil
}

func (h *Handlers) GetOAuthCallback(ctx context.Context, req openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	if req.Params.Code == "" {
		return nil, errs.ErrMissingOAuthCode
	}
	if req.Params.State == "" {
		return nil, errs.ErrMissingOAuthState
	}
	tokens, err := h.ops.Callback(ctx, string(req.Provider), req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthCallback201JSONResponse(*tokens), nil
}

func (h *Handlers) GetOAuthConnect(ctx context.Context, req openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	acceptedTos := false
	if req.Params.AcceptedTos != nil {
		acceptedTos = *req.Params.AcceptedTos
	}
	url, err := h.ops.Connect(ctx, string(req.Provider), req.Params.ProjectId, acceptedTos)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthConnect200JSONResponse{Url: url}, nil
}

func (h *Handlers) GetOAuthProviders(ctx context.Context, req openapi.GetOAuthProvidersRequestObject) (openapi.GetOAuthProvidersResponseObject, error) {
	providers, err := h.ops.ListEnabledProviders(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	items := make([]openapi.OAuthProviderDiscoveryItem, 0, len(providers))
	for _, p := range providers {
		items = append(items, openapi.OAuthProviderDiscoveryItem{Provider: p})
	}
	return openapi.GetOAuthProviders200JSONResponse{Items: items}, nil
}

func (h *Handlers) ListProjectOAuthProviders(ctx context.Context, req openapi.ListProjectOAuthProvidersRequestObject) (openapi.ListProjectOAuthProvidersResponseObject, error) {
	providers, err := h.ops.ListByProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectOAuthProviders200JSONResponse{Items: providers}, nil
}

func (h *Handlers) UpdateOAuthProvider(ctx context.Context, req openapi.UpdateOAuthProviderRequestObject) (openapi.UpdateOAuthProviderResponseObject, error) {
	updated, err := h.ops.Update(ctx, models.UpdateOAuthProviderInput{
		ID:           req.OauthProviderId,
		ClientID:     req.Body.ClientId,
		ClientSecret: req.Body.ClientSecret,
		CallbackURL:  req.Body.CallbackUrl,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpdateOAuthProvider200JSONResponse(*updated), nil
}
