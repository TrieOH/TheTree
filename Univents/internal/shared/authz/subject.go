package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"univents/internal/shared/errx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserMetadata represents the nested object inside 'sub'
type UserMetadata struct {
	// Add metadata fields here
}

// UserSubject represents the full 'sub' claim object
type UserSubject struct {
	ID       uuid.UUID    `json:"id"`
	Email    string       `json:"email"`
	Metadata UserMetadata `json:"metadata"`
}

// GetSubjectFromToken extracts and parses the complex 'sub' claim
func GetSubjectFromToken(token *jwt.Token) (*UserSubject, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errx.Invalid("access claims")
	}

	subInterface, exists := claims["sub"]
	if !exists {
		return nil, errx.NotFound("sub claims")
	}

	subBytes, err := json.Marshal(subInterface)
	if err != nil {
		return nil, errx.Internal("sub").SetMessage("marshal failed")
	}

	var userSubject UserSubject
	if err := json.Unmarshal(subBytes, &userSubject); err != nil {
		return nil, errx.Internal("sub").SetMessage("unmarshal failed")
	}

	return &userSubject, nil
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
