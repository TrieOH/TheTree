package tokens

import (
	"IdentityX/internal/features/keys"
	"IdentityX/internal/shared/contracts"
	errx2 "IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/validation"
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CommandService struct {
	keys keys.CommandService
}

func NewCommandService(
	keys keys.CommandService,
) *CommandService {
	return &CommandService{
		keys: keys,
	}
}

func (uc *CommandService) NewAccessToken(in contracts.NewAccessTokenInput) ([]byte, error) {
	claims := contracts.AccessClaims{
		Sub: contracts.AccessSub{
			ID:         in.User.ID,
			UserType:   string(in.User.UserType),
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

func (uc *CommandService) NewRefreshToken(in contracts.NewRefreshTokenInput) ([]byte, error) {
	claims := contracts.RefreshClaims{
		Sub: contracts.RefreshSub{
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

func (uc *CommandService) NewProjectAccessToken(in contracts.NewProjectAccessTokenInput) ([]byte, error) {
	claims := contracts.AccessClaims{
		Sub: contracts.AccessSub{
			ID:         in.User.ID,
			UserType:   string(in.User.UserType),
			ProjectID:  in.User.ProjectID,
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

func (uc *CommandService) NewVerificationToken(in contracts.NewVerificationTokenInput) ([]byte, error) {
	now := time.Now()
	claims := contracts.VerificationClaims{
		Sub: contracts.VerificationSub{
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

func (uc *CommandService) NewResetPasswordToken(in contracts.NewResetPasswordInput) ([]byte, error) {
	now := time.Now()
	claims := contracts.ResetPasswordClaims{
		Sub: contracts.ResetPasswordSub{
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

func (uc *CommandService) AssembleJWT(payload []byte, sig []byte) string {
	return string(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func (uc *CommandService) VerifyAccessToken(
	ctx context.Context,
	tokenStr string,
) (*contracts.AccessClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"access",
		tokenStr,
		&contracts.AccessClaims{},
	)
}

func (uc *CommandService) VerifyRefreshToken(
	ctx context.Context,
	tokenStr string,
) (*contracts.RefreshClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"refresh",
		tokenStr,
		&contracts.RefreshClaims{},
	)
}

func (uc *CommandService) VerifyVerificationToken(
	ctx context.Context,
	tokenStr string,
) (*contracts.VerificationClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"verification",
		tokenStr,
		&contracts.VerificationClaims{},
	)
}

func (uc *CommandService) VerifyResetPasswordToken(
	ctx context.Context,
	tokenStr string,
) (*contracts.ResetPasswordClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"reset password",
		tokenStr,
		&contracts.ResetPasswordClaims{},
	)
}

func verifyToken[T jwt.Claims](
	ctx context.Context,
	uc *CommandService,
	tokenType string,
	tokenStr string,
	claims T,
) (T, error) {
	token, err := parseJWTUnverified(tokenStr, claims)
	if err != nil {
		return claims, errx2.FromJWTError(err, tokenType)
	}

	alg, _ := token.Header["alg"].(string)
	if alg != jwt.SigningMethodEdDSA.Alg() {
		return claims, fail.New(errx2.TokenInvalidAlg).WithArgs(tokenType, jwt.SigningMethodEdDSA.Alg(), alg).RecordCtx(ctx)
	}

	if token.Method == nil || token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
		methodAlg := ""
		if token.Method != nil {
			methodAlg = token.Method.Alg()
		}
		return claims, fail.New(errx2.TokenInvalidAlg).WithArgs(tokenType, jwt.SigningMethodEdDSA.Alg(), methodAlg).RecordCtx(ctx)
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return claims, fail.New(errx2.TokenMissingKid).WithArgs(tokenType).RecordCtx(ctx)
	}

	payload, sig, err := splitJWT(ctx, tokenType, tokenStr)
	if err != nil {
		return claims, err
	}

	switch {
	case strings.HasPrefix(kid, "goauth:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 2 {
			return claims, fail.New(errx2.TokenInvalidKid).WithArgs(tokenType).RecordCtx(ctx)
		}

		if err := uc.keys.VerifyGoAuth(ctx, kid, payload, sig); err != nil {
			if fail.Is(err, errx2.SQLNotFound) {
				return claims, fail.New(errx2.TokenUntrusted).WithArgs(tokenType).RecordCtx(ctx)
			}
			return claims, err
		}

	case strings.HasPrefix(kid, "project:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 3 {
			return claims, fail.New(errx2.TokenInvalidKid).WithArgs(tokenType).RecordCtx(ctx)
		}

		projectID, err := validation.ParseUUID(parts[1], "project_id")
		if err != nil {
			return claims, err
		}

		if err := uc.keys.VerifyProject(ctx, projectID, kid, payload, sig); err != nil {
			if fail.Is(err, errx2.SQLNotFound) {
				return claims, fail.New(errx2.TokenUntrusted).WithArgs(tokenType).RecordCtx(ctx)
			}
			return claims, err
		}

	default:
		return claims, fail.New(errx2.TokenUnknownKid).WithArgs(tokenType).RecordCtx(ctx)
	}

	return claims, nil
}

func parseJWTUnverified[T jwt.Claims](tokenStr string, claims T) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	return token, err
}

func splitJWT(ctx context.Context, tokenType, tokenStr string) (signingInput, sig []byte, err error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, nil, fail.New(errx2.TokenInvalidFormat).WithArgs(tokenType).RecordCtx(ctx)
	}

	signingInput = []byte(parts[0] + "." + parts[1])

	sig, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, err
	}

	return signingInput, sig, nil
}
