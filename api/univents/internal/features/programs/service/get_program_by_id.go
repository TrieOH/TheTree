package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProgramByID(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	return o.Programs.GetByID(ctx, id)
}
