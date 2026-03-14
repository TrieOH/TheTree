package inbounds

import (
	"GoAuth/internal/domain/authz"
	"context"
)

type RequestAuthenticator interface {
	AuthenticateRequest(ctx context.Context, in AuthenticateRequestInput) (*authz.Principal, error)
}
