package issuer

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/ports/inbounds"
	"encoding/base64"
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

func (uc *UseCase) NewAccessToken(in inbounds.NewAccessTokenInput) ([]byte, error) {
	claims := auth.AccessClaims{
		Sub: auth.AccessSub{
			ID:         in.User.ID,
			UserType:   in.User.UserType,
			Email:      in.User.Email,
			SessionID:  in.SessionID,
			UserAgent:  in.Agent,
			UserIP:     in.IP,
			IsVerified: in.User.IsVerified,
			FamilyID:   in.FamilyID,
			VerifiedAt: in.User.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        in.AccessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (uc *UseCase) NewRefreshToken(in inbounds.NewRefreshTokenInput) ([]byte, error) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (uc *UseCase) NewProjectAccessToken(in inbounds.NewProjectAccessTokenInput) ([]byte, error) {
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
			FamilyID:   in.FamilyID,
			VerifiedAt: in.User.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    in.User.ProjectID.String(),
			ID:        in.AccessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (uc *UseCase) NewVerificationToken(in inbounds.NewVerificationTokenInput) ([]byte, error) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (uc *UseCase) NewResetPasswordToken(in inbounds.NewResetPasswordInput) ([]byte, error) {
	now := time.Now()
	claims := auth.ResetPasswordClaims{
		Sub: auth.ResetPasswordSub{
			Subject:   in.Subject,
			ProjectID: in.ProjectID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    viper.GetString("ISSUER"),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
			Audience:  jwt.ClaimStrings{"password-reset"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (uc *UseCase) AssembleJWT(payload []byte, sig []byte) string {
	return string(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
}
