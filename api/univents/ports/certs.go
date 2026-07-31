package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type CertificationRepo interface {
	CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error)
	UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error)
	ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error)
	ListTemplatesForEmission(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error)

	Certify(ctx context.Context, input models.CertifyInput) (*models.Certification, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error)
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error)
	GetByHash(ctx context.Context, hash string) (*models.Certification, error)
	HasCertForProgram(ctx context.Context, userID, programID uuid.UUID) (bool, error)
	HasCertForRegistration(ctx context.Context, registrationID uuid.UUID, templateID *uuid.UUID) (bool, error)
	Invalidate(ctx context.Context, id uuid.UUID, reason *string) error
	MarkEmailSent(ctx context.Context, id uuid.UUID) error

	LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error
	UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error

	RecordEmissionError(ctx context.Context, err *models.CertEmissionError) error
	ListEmissionErrorsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error)

	ListDistinctRegistrationsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEligibleAttendee, error)
	ListDistinctParticipantsByProgram(ctx context.Context, programID uuid.UUID) ([]models.CertEligibleAttendee, error)
}
