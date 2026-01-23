package inbounds

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, in RegisterUserInput) error
	Login(ctx context.Context, in LoginUserInput) (*UserTokensOutput, error)
	Logout(ctx context.Context) error
	Refresh(ctx context.Context, in RefreshInput) (*UserTokensOutput, error)
	RegisterProjectUser(ctx context.Context, in ProjectRegisterInput) error
	LoginProjectUser(ctx context.Context, in ProjectLoginInput) (*UserTokensOutput, error)
	Verify(ctx context.Context, token string) error
}
