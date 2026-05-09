package authz

import (
	"context"
	"encoding/json"

	"github.com/MintzyG/fun"
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
		return nil, fun.ErrBadRequest("invalid token claims")
	}

	subInterface, exists := claims["sub"]
	if !exists {
		return nil, fun.ErrNotFound("sub claims not found")
	}

	subBytes, err := json.Marshal(subInterface)
	if err != nil {
		return nil, fun.Err("sub marshal failed").WithErr(err).Internal()
	}

	var userSubject UserSubject
	if err = json.Unmarshal(subBytes, &userSubject); err != nil {
		return nil, fun.Err("sub unmarshal failed").WithErr(err).Internal()
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
		return nil, fun.ErrNotFound("subject not found in context")
	}

	u, ok := val.(*UserSubject)
	if !ok {
		return nil, fun.Errf("invalid subject type, was: %T", val).Internal()
	}

	return u, nil
}
