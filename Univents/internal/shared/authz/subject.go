package authz

import (
	"context"
	"fmt"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

// UserSubject represents the full 'sub' claim object
type UserSubject struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type contextKey string

const UserContextKey contextKey = "user"

func WithSubject(ctx context.Context, subject *UserSubject) context.Context {
	return context.WithValue(ctx, UserContextKey, subject)
}

func RequireSubject(ctx context.Context) (*UserSubject, error) {
	val := ctx.Value(UserContextKey)
	if val == nil {
		return nil, errx.NotFound("subject").SetMessage("subject not found in context")
	}

	u, ok := val.(*UserSubject)
	if !ok {
		return nil, errx.Invalid("subject").SetMessage(fmt.Sprintf("type was %T", val))
	}

	return u, nil
}
