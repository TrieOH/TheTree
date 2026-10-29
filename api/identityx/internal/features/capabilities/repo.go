package capabilities

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.CapabilityRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrCapabilityNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapCapability(src sqlc.Capability) models.Capability {
	return models.Capability{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		Resource:  src.Resource,
		Action:    src.Action,
		CreatedBy: src.CreatedBy,
		CreatedAt: src.CreatedAt,
	}
}

func (repo *Repo) AssignToAPIKey(ctx context.Context, apiKeyID uuid.UUID, capabilityIDs []uuid.UUID, assignedBy uuid.UUID) error {
	err := repo.db(ctx).AssignCapabilitiesToApiKey(ctx, sqlc.AssignCapabilitiesToApiKeyParams{
		ApiKeyID:      apiKeyID,
		CapabilityIds: capabilityIDs,
		AssignedBy:    &assignedBy,
	})
	return repo.dbe(err)
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Capability) (*models.Capability, error) {
	capability, err := repo.db(ctx).CreateCapability(ctx, sqlc.CreateCapabilityParams{
		ProjectID: toCreate.ProjectID,
		Resource:  toCreate.Resource,
		Action:    toCreate.Action,
		CreatedBy: toCreate.CreatedBy,
	})
	return database.One(capability, err, repo.dbe, mapCapability)
}

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	capabilities, err := repo.db(ctx).ListCapabilitiesByProject(ctx, &projectID)
	return database.Many(capabilities, err, repo.dbe, mapCapability)
}

func (repo *Repo) ListByAPIKeyPrefix(ctx context.Context, prefix string) ([]models.Capability, error) {
	capabilities, err := repo.db(ctx).ListCapabilitiesByApiKeyPrefix(ctx, prefix)
	return database.Many(capabilities, err, repo.dbe, mapCapability)
}

func (repo *Repo) ValidateCapabilities(ctx context.Context, projectID *uuid.UUID, capabilities []uuid.UUID) (bool, error) {
	valid, err := repo.db(ctx).ValidateCapabilities(ctx, sqlc.ValidateCapabilitiesParams{
		CapabilityCount: len(capabilities),
		ProjectID:       projectID,
		CapabilityIds:   capabilities,
	})
	return valid, repo.dbe(err)
}
