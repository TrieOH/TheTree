package inbounds

import (
	"GoAuth/internal/domain/authz"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, in RegisterUserInput) error
	Login(ctx context.Context, in LoginUserInput) (*UserTokensOutput, error)
	Logout(ctx context.Context, snapshot authz.ServiceSnapshot) error
	Refresh(ctx context.Context, in RefreshInput) (*UserTokensOutput, error)
	RegisterProjectUser(ctx context.Context, in ProjectRegisterInput) error
	LoginProjectUser(ctx context.Context, in ProjectLoginInput) (*UserTokensOutput, error)
	LogoutProjectUser(ctx context.Context, in ProjectLogoutInput) error
	Verify(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context) error
	ForgotPassword(ctx context.Context, in ForgotPasswordInput) error
	ResetPassword(ctx context.Context, in ResetPasswordInput) error
	GetJWKS(ctx context.Context) (map[string]any, error)
	Exchange(ctx context.Context, globalAccess string) (*ExchangeOutput, error)
}
