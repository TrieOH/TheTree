package authn

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOAuthConnect(ctx context.Context, req openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	url, err := h.ops.OAuthConnect(ctx, string(req.Provider))
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthConnect200JSONResponse{
		Code: 200,
		Data: &struct {
			Url string `json:"url"` //nolint:revive // generated field name
		}{Url: url},
		Timestamp: time.Now(), Module: module,
	}, nil
}
