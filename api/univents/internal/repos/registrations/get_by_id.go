package registrations

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.GetByID")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetRegistrationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(row)), nil
}
