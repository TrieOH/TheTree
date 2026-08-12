package ws

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

// GetWsToken issues the one-time WS handshake token (split 6): the caller
// must be the purchase's purchaser — a non-owner gets 404 via the service
// (no existence leak). The response carries the raw token (never stored —
// only its hash) and its expiry so the front can re-request before
// reconnecting.
func (h *Handlers) GetWsToken(ctx context.Context, req openapi.GetWsTokenRequestObject) (openapi.GetWsTokenResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	token, expiresAt, err := h.ops.IssueToken(ctx, req.Params.PurchaseId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.GetWsToken200JSONResponse{
		Code:      200,
		Data:      &openapi.WsToken{Token: token, ExpiresAt: expiresAt},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}
