package commands

import (
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/crypto"
	"lib/database"
	"lib/errx"
	"os"
	"time"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	actors             ports.ActorRepo
	platformRoles      ports.PlatformRolesRepo
	cryptoKeys         ports.CryptoKeysRepo
	blacklist          ports.BlacklistRepo
	externalIdentities ports.ExternalIdentitiesRepo
	logger             *zap.Logger
	tracer             trace.Tracer
	tx                 database.TxRunner
}

func NewCommands(deps ports.AuthnDeps) *Commands {
	return errx.MustProvide(&Commands{
		actors:             deps.Actors,
		platformRoles:      deps.PlatformRoles,
		cryptoKeys:         deps.CryptoKeys,
		blacklist:          deps.Blacklist,
		externalIdentities: deps.ExternalIdentities,
		logger:             deps.Logger,
		tracer:             deps.Tracer,
		tx:                 deps.Tx,
	})
}

func (c *Commands) cryptoKeyFromToken(ctx context.Context, token *jwt.Token) (*models.CryptoKey, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fun.ErrUnauthorized("missing kid")
	}
	keyID, err := uuid.Parse(kid)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid kid")
	}
	cryptoKey, err := c.cryptoKeys.GetByID(ctx, keyID)
	if err != nil && fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("outdated token")
	}
	if err != nil {
		return nil, err
	}
	if cryptoKey.Status == "revoked" {
		return nil, fun.ErrUnauthorized("token signing key revoked")
	}
	return cryptoKey, nil
}

func (c *Commands) issueTokens(ctx context.Context, actor *models.Actor) (*models.UserTokensOutput, error) {
	activeKeyPair, err := c.cryptoKeys.GetActive(ctx, models.SigningCryptoKeyType, nil)
	if err != nil {
		return nil, err
	}
	accessJTI := uuid.New()
	refreshJTI := uuid.New()
	accessExpiresAt := time.Now().Add(errx.Env[time.Duration]("ACCESS_TOKEN_EXPIRATION", time.ParseDuration, 15*time.Minute))
	refreshExpiresAt := time.Now().Add(errx.Env[time.Duration]("REFRESH_TOKEN_EXPIRATION", time.ParseDuration, 7*24*time.Hour))
	accessPayload, err := c.newIDXAccessToken(*actor, accessJTI, activeKeyPair.ID, accessExpiresAt)
	if err != nil {
		return nil, err
	}
	refreshPayload, err := c.newIDXRefreshToken(actor, refreshJTI, accessJTI, activeKeyPair.ID, refreshExpiresAt)
	if err != nil {
		return nil, err
	}
	kp := activeKeyPair.ToKeyPair()
	accessToken, err := crypto.SignToken(accessPayload, kp)
	if err != nil {
		return nil, err
	}
	refreshToken, err := crypto.SignToken(refreshPayload, kp)
	if err != nil {
		return nil, err
	}
	return &models.UserTokensOutput{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Domain:           os.Getenv("ISSUER"),
	}, nil
}

func (c *Commands) newIDXAccessToken(actor models.Actor, jti, kid uuid.UUID, expiresAt time.Time) ([]byte, error) {
	claims := models.AccessClaims{
		Sub: models.AccessSub{
			ID:           actor.ID,
			ProjectID:    nil,
			Email:        actor.Email,
			Type:         actor.Type,
			Capabilities: nil,
			Metadata:     nil,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    os.Getenv("ISSUER"),
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func (c *Commands) newIDXRefreshToken(actor *models.Actor, jti, accessJTI, kid uuid.UUID, expiresAt time.Time) ([]byte, error) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSub{
			ID:        actor.ID,
			ProjectID: actor.ProjectID,
			AccessJTI: accessJTI,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    os.Getenv("ISSUER"),
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}
