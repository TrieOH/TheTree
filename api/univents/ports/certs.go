package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type CertificationRepo interface {
	CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error)
	ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error)

	Certify(ctx context.Context, input models.CertifyInput) (*models.Certification, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error)
	ListByTarget(ctx context.Context, targetType string, targetID uuid.UUID) ([]models.Certification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error)
	GetByHash(ctx context.Context, hash string) (*models.Certification, error)

	SetActivityTemplate(ctx context.Context, activityID uuid.UUID, templateID *uuid.UUID) error
	SetEditionTemplate(ctx context.Context, editionID uuid.UUID, templateID *uuid.UUID) error
}
