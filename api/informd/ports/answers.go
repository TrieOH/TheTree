package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type AnswerRepo interface {
	GetByResponse(ctx context.Context, responseID uuid.UUID) ([]models.Answer, error)
	BatchUpsert(ctx context.Context, answers []models.Answer) error
}
