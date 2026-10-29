package repo

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var (
	_ ports.SignatureRepo        = (*Repo)(nil)
	_ ports.SignatureRequestRepo = (*Repo)(nil)
)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrSignatureNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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

func (repo *Repo) CancelRequest(ctx context.Context, id uuid.UUID, reason *string) error {
	rows, err := repo.db(ctx).CancelSignatureRequest(ctx, sqlc.CancelSignatureRequestParams{
		ID:           id,
		StatusReason: reason,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return &database.NotFoundError{Named: errs.ErrSignatureRequestHandled}
	}
	return nil
}

func (repo *Repo) CompleteRequest(ctx context.Context, id, signatureID uuid.UUID) error {
	rows, err := repo.db(ctx).CompleteSignatureRequest(ctx, sqlc.CompleteSignatureRequestParams{
		ID:          id,
		SignatureID: &signatureID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if rows == 0 {
		return &database.NotFoundError{Named: errs.ErrSignatureRequestCompleted}
	}
	return nil
}

func (repo *Repo) Create(ctx context.Context, toCreate *models.Signature) (*models.Signature, error) {
	result, err := repo.db(ctx).CreateSignature(ctx, sqlc.CreateSignatureParams{
		EditionID:       toCreate.EditionID,
		CreatedBy:       toCreate.CreatedBy,
		SignatoryName:   toCreate.SignatoryName,
		SignatoryTitle:  toCreate.SignatoryTitle,
		SignatoryEmail:  toCreate.SignatoryEmail,
		SignatoryUserID: toCreate.SignatoryUserID,
		ImageUrl:        toCreate.ImageURL,
	})
	return database.One(result, err, repo.dbe, mapSignature)
}

func (repo *Repo) CreateRequest(ctx context.Context, toCreate *models.SignatureRequest) (*models.SignatureRequest, error) {
	result, err := repo.db(ctx).CreateSignatureRequest(ctx, sqlc.CreateSignatureRequestParams{
		EditionID:       toCreate.EditionID,
		CreatedBy:       toCreate.CreatedBy,
		SignatoryName:   toCreate.SignatoryName,
		SignatoryTitle:  toCreate.SignatoryTitle,
		SignatoryEmail:  toCreate.SignatoryEmail,
		SignatoryUserID: toCreate.SignatoryUserID,
		IdempotencyKey:  toCreate.IdempotencyKey,
		ExpiresAt:       toCreate.ExpiresAt,
	})
	return database.One(result, err, repo.dbe, mapSignatureRequest)
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteSignature(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) ExpireStaleRequests(ctx context.Context) error {
	err := repo.db(ctx).ExpireStaleSignatureRequests(ctx)
	return repo.dbe(err)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	result, err := repo.db(ctx).GetSignatureByID(ctx, id)
	return database.One(result, err, repo.dbe, mapSignature)
}

func (repo *Repo) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	result, err := repo.db(ctx).GetSignatureRequestByID(ctx, id)
	return database.One(result, err, repo.dbe, mapSignatureRequest)
}

func (repo *Repo) GetRequestByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	result, err := repo.db(ctx).GetSignatureRequestByIdempotencyKey(ctx, idempotencyKey)
	return database.One(result, err, repo.dbe, mapSignatureRequest)
}

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	results, err := repo.db(ctx).ListSignaturesByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapSignature)
}

func (repo *Repo) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	err := repo.ExpireStaleRequests(ctx)
	if err != nil {
		return nil, err
	}
	results, err := repo.db(ctx).ListSignatureRequestsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapSignatureRequest)
}
