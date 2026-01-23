package issuer

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UseCase struct{}

var _ inbounds.TokenIssuer = (*UseCase)(nil)

func NewTokenIssuer() inbounds.TokenIssuer {
	return &UseCase{}
}

func (uc *UseCase) NewAccessToken(in inbounds.NewAccessTokenInput) (string, error) {
	claims := auth.AccessClaims{
		Sub: auth.AccessSub{
			ID:         in.User.ID,
			UserType:   in.User.UserType,
			Email:      in.User.Email,
			SessionID:  in.SessionID,
			UserAgent:  in.Agent,
			UserIP:     in.IP,
			IsVerified: in.User.IsVerified,
			VerifiedAt: in.User.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        in.AccessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	accessToken.Header["kid"] = in.KeyID
	tokenStr, err := accessToken.SignedString(in.PrivateKey)
	if err != nil {
		return "", auth.ErrSigningToken{TokenType: "access", Cause: err}
	}
	return tokenStr, nil
}

func (uc *UseCase) NewRefreshToken(in inbounds.NewRefreshTokenInput) (string, error) {
	claims := auth.RefreshClaims{
		Sub: auth.RefreshSub{
			AccessJTI: in.AccessJTI,
			FamilyID:  in.FamilyID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        in.RefreshJTI.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	refreshToken.Header["kid"] = in.KeyID
	tokenStr, err := refreshToken.SignedString(in.PrivateKey)
	if err != nil {
		return "", auth.ErrSigningToken{TokenType: "refresh", Cause: err}
	}
	return tokenStr, nil
}

func (uc *UseCase) NewProjectAccessToken(in inbounds.NewProjectAccessTokenInput) (string, error) {
	claims := auth.AccessClaims{
		Sub: auth.AccessSub{
			ID:         in.User.ID,
			UserType:   in.User.UserType,
			ProjectID:  &in.User.ProjectID,
			Metadata:   in.User.Metadata,
			Email:      in.User.Email,
			SessionID:  in.SessionID,
			UserAgent:  in.Agent,
			UserIP:     in.IP,
			IsVerified: in.User.IsVerified,
			VerifiedAt: in.User.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        in.AccessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	accessToken.Header["kid"] = in.KeyID
	tokenStr, err := accessToken.SignedString(in.PrivateKey)
	if err != nil {
		return "", auth.ErrSigningToken{TokenType: "access", Cause: err}
	}
	return tokenStr, nil
}

func (uc *UseCase) NewVerificationToken(in inbounds.NewVerificationTokenInput) (string, error) {
	now := time.Now()
	claims := auth.VerificationClaims{
		Sub: auth.VerificationSub{
			Subject: in.Subject,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
			Audience:  jwt.ClaimStrings{"email-verification"},
		},
	}

	verification := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	verification.Header["kid"] = "email-verification"
	tokenStr, err := verification.SignedString(in.PrivateKey)
	if err != nil {
		return "", auth.ErrSigningToken{TokenType: "verification", Cause: err}
	}
	return tokenStr, nil
}
