package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) RevokeByRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	err := database.Queries(ctx, repo.q).RevokeBadgeEmissionByRegistration(ctx, sqlc.RevokeBadgeEmissionByRegistrationParams{
		RegistrationID: new(registrationID),
		Reason:         new(reason),
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
