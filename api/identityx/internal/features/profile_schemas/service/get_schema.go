package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// GetSchema returns the active profile schema for the scope. Public read:
// schemas shape public profiles, so anyone may fetch them.
func (o *Operations) GetSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	return o.Schemas.Get(ctx, projectID)
}
