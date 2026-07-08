package ports

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

type CertificationRepo interface {
	CreateTemplate(ctx context.Context, input contracts.CreateCertificationTemplateInput) (*contracts.CertificationTemplate, error)
	ListTemplates(ctx context.Context, editionID uuid.UUID) ([]contracts.CertificationTemplate, error)
	GetTemplateByID(ctx context.Context, id, editionID uuid.UUID) (*contracts.CertificationTemplate, error)

	Certify(ctx context.Context, input contracts.CertifyInput) (*contracts.Certification, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]contracts.Certification, error)
	ListByTarget(ctx context.Context, targetType string, targetID uuid.UUID) ([]contracts.Certification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Certification, error)

	SetActivityTemplate(ctx context.Context, activityID uuid.UUID, templateID *uuid.UUID) error
	SetEditionTemplate(ctx context.Context, editionID uuid.UUID, templateID *uuid.UUID) error
}
