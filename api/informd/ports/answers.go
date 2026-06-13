package ports

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

type AnswerRepo interface {
	GetByResponse(ctx context.Context, responseID uuid.UUID) ([]models.Answer, error)
	GetByField(ctx context.Context, fieldID uuid.UUID) ([]models.Answer, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Answer, error)
	BatchUpsert(ctx context.Context, answers []models.Answer) error
}
