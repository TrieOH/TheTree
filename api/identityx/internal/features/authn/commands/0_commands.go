package commands

import (
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
	"lib/errx"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	actors        ports.ActorRepo
	platformRoles ports.PlatformRolesRepo
	cryptoKeys    ports.CryptoKeysRepo
	blacklist     ports.BlacklistRepo
	logger        *zap.Logger
	tracer        trace.Tracer
	tx            database.TxRunner
}

func NewCommands(deps ports.AuthnDeps) *Commands {
	return errx.MustProvide(&Commands{
		actors:        deps.Actors,
		platformRoles: deps.PlatformRoles,
		cryptoKeys:    deps.CryptoKeys,
		blacklist:     deps.Blacklist,
		logger:        deps.Logger,
		tracer:        deps.Tracer,
		tx:            deps.Tx,
	})
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

func (c *Commands) newIDXRefreshToken(jti, accessJTI, kid uuid.UUID, expiresAt time.Time) ([]byte, error) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSub{
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
