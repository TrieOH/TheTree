package inbounds

import (
	"GoAuth/internal/domain/auth"
	"context"
)

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, tokenStr string) (*auth.AccessClaims, error)
	VerifyRefreshToken(ctx context.Context, tokenStr string) (*auth.RefreshClaims, error)
	VerifyVerificationToken(ctx context.Context, tokenStr string) (*auth.VerificationClaims, error)
}
