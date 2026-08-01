package signatures

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

var (
	_ ports.SignatureRepo        = (*Repo)(nil)
	_ ports.SignatureRequestRepo = (*Repo)(nil)
)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("signature"),
	}
}

func mapSignature(src sqlc.Signature) models.Signature {
	return models.Signature{
		ID:              src.ID,
		EditionID:       src.EditionID,
		CreatedBy:       src.CreatedBy,
		SignatoryName:   src.SignatoryName,
		SignatoryTitle:  src.SignatoryTitle,
		SignatoryEmail:  src.SignatoryEmail,
		SignatoryUserID: src.SignatoryUserID,
		ImageURL:        src.ImageUrl,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
		DeletedAt:       src.DeletedAt,
	}
}

func mapSignatureRequest(src sqlc.SignatureRequest) models.SignatureRequest {
	return models.SignatureRequest{
		ID:              src.ID,
		EditionID:       src.EditionID,
		CreatedBy:       src.CreatedBy,
		SignatoryName:   src.SignatoryName,
		SignatoryTitle:  src.SignatoryTitle,
		SignatoryEmail:  src.SignatoryEmail,
		SignatoryUserID: src.SignatoryUserID,
		IdempotencyKey:  src.IdempotencyKey,
		Status:          models.SignatureRequestStatus(src.Status),
		StatusReason:    src.StatusReason,
		ExpiresAt:       src.ExpiresAt,
		SignatureID:     src.SignatureID,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
		DeletedAt:       src.DeletedAt,
	}
}
