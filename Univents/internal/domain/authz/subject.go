package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"univents/internal/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
)

// UserMetadata represents the nested object inside 'sub'
type UserMetadata struct {
	// Add metadata fields here
}

// UserSubject represents the full 'sub' claim object
type UserSubject struct {
	ID       string       `json:"id"`
	Email    string       `json:"email"`
	Metadata UserMetadata `json:"metadata"`
}

// GetSubjectFromToken extracts and parses the complex 'sub' claim
func GetSubjectFromToken(token *jwt.Token) (*UserSubject, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fail.New(errx.TokenInvalidAccessClaims).WithArgs("access")
	}

	subInterface, exists := claims["sub"]
	if !exists {
		return nil, fail.New(errx.TokenMissingSubClaim)
	}

	subBytes, err := json.Marshal(subInterface)
	if err != nil {
		return nil, fail.New(errx.TokenSubMarshalFailed)
	}

	var userSubject UserSubject
	if err := json.Unmarshal(subBytes, &userSubject); err != nil {
		return nil, fail.New(errx.TokenSubUnmarshallingFailed)
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
		return nil, fail.New(errx.AuthSubjectNotInContext).RecordCtx(ctx)
	}

	u, ok := val.(*UserSubject)
	if !ok {
		return nil, fail.New(errx.AuthInvalidSubject).
			WithArgs(fmt.Sprintf("type was %T", val)).
			RecordCtx(ctx)
	}

	return u, nil
}
