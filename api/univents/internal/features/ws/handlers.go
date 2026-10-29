package ws

import (
	"context"
	"net/http"

	idx "sdk/identityx"
	"univents/internal/features/ws/service"
	"univents/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

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
	return openapi.GetWsToken200JSONResponse(openapi.WsToken{Token: token, ExpiresAt: expiresAt}), nil
}

// ServeWS is the raw `WS /ws?token=...` route (D9) — mounted by the app
// router outside the strict handler, because the handshake carries the
// one-time token in the query string (not an Authorization header). The
// token proves prior REST auth for this purchase; the service consumes it
// atomically (one-time), sends the snapshot frame, then relays live frames
// until a terminal event closes the socket.
func (h *Handlers) ServeWS(w http.ResponseWriter, r *http.Request) {
	h.ops.ServeWS(w, r)
}
