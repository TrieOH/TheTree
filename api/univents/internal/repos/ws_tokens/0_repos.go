package ws_tokens

import (
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.WsTokenRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("ws_token"),
	}
}

func mapWsToken(src sqlc.WsToken) models.WsToken {
	return models.WsToken{
		ID:         src.ID,
		PurchaseID: src.PurchaseID,
		UserID:     src.UserID,
		TokenHash:  src.TokenHash,
		ExpiresAt:  src.ExpiresAt,
		UsedAt:     src.UsedAt,
	}
}
