package authz

import (
	"context"
	"encoding/json"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
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
		return nil, fun.NewError("invalid token claims").BadRequest()
	}

	subInterface, exists := claims["sub"]
	if !exists {
		return nil, fun.NewError("sub claims not found").NotFound()
	}

	subBytes, err := json.Marshal(subInterface)
	if err != nil {
		return nil, fun.NewError("sub marshal failed").WithErr(err).Internal()
	}

	var userSubject UserSubject
	if err = json.Unmarshal(subBytes, &userSubject); err != nil {
		return nil, fun.NewError("sub unmarshal failed").WithErr(err).Internal()
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
		return nil, fun.NewError("subject not found in context").NotFound()
	}

	u, ok := val.(*UserSubject)
	if !ok {
		return nil, fun.NewErrorf("invalid subject type, was: %T", val).Internal()
	}

	return u, nil
}
