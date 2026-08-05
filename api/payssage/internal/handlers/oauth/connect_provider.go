package oauth

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) ConnectProvider(ctx context.Context, req openapi.ConnectProviderRequestObject) (openapi.ConnectProviderResponseObject, error) {
	url, err := h.ops.Connect(ctx, models.ConnectInput{
		Provider:            req.Provider,
		Flow:                req.Body.Flow,
		ProviderRedirectURL: req.Body.ProviderRedirectUrl,
		FinalRedirectURL:    req.Body.FinalRedirectUrl,
		WalletID:            req.Body.WalletId,
		OrganizationID:      req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.ConnectProvider200JSONResponse{
		Code: 200, Data: &url, Timestamp: time.Now(), Module: module,
	}, nil
}
