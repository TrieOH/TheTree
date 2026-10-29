package certifications

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var (
	_ ports.CertificationRepo = (*Repo)(nil)
)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrCertificationNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapCertTemplate(src sqlc.CertificationTemplate) models.CertificationTemplate {
	return models.CertificationTemplate{
		ID:          src.ID,
		EditionID:   src.EditionID,
		Kind:        models.CertificationTemplateKind(src.Kind),
		Name:        src.Name,
		Description: src.Description,
		DesignData:  src.DesignData,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func mapCertTemplateProgram(src sqlc.CertificationTemplateProgram) models.CertificationTemplateProgram {
	return models.CertificationTemplateProgram{
		ID:         src.ID,
		TemplateID: src.TemplateID,
		ProgramID:  src.ProgramID,
		CreatedAt:  src.CreatedAt,
	}
}

func mapCertification(src sqlc.Certification) models.Certification {
	return models.Certification{
		ID:               src.ID,
		EditionID:        src.EditionID,
		TemplateID:       src.TemplateID,
		RegistrationID:   src.RegistrationID,
		UserID:           src.UserID,
		ProgramID:        src.ProgramID,
		VerificationHash: src.VerificationHash,
		Valid:            src.Valid,
		InvalidReason:    src.InvalidReason,
		EmailSent:        src.EmailSent,
		IssuedAt:         src.IssuedAt,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
	}
}

func mapCertEmissionError(src sqlc.CertEmissionError) models.CertEmissionError {
	return models.CertEmissionError{
		ID:           src.ID,
		EditionID:    src.EditionID,
		UserID:       src.UserID,
		TemplateID:   src.TemplateID,
		ProgramID:    src.ProgramID,
		ErrorMessage: src.ErrorMessage,
		CreatedAt:    src.CreatedAt,
	}
}

func mapEligibleAttendee(src sqlc.ListDistinctRegistrationsByEditionRow) models.CertEligibleAttendee {
	return models.CertEligibleAttendee{
		UserID:         userIDOrNil(src.UserID),
		RegistrationID: src.RegistrationID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}

func mapEligibleParticipant(src sqlc.ListDistinctParticipantsByProgramRow) models.CertEligibleAttendee {
	return models.CertEligibleAttendee{
		UserID:         userIDOrNil(src.UserID),
		RegistrationID: src.RegistrationID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}

// userIDOrNil dereferences the nullable attendee user id, falling back to
// uuid.Nil for accountless (email-only) attendees — the sentinel the cert
// grant job skips (a cert must be owned by an account).
func userIDOrNil(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

func (repo *Repo) Certify(ctx context.Context, input models.CertifyInput) (*models.Certification, error) {
	result, err := repo.db(ctx).CreateCertification(ctx, sqlc.CreateCertificationParams{
		EditionID:        input.EditionID,
		TemplateID:       input.TemplateID,
		RegistrationID:   input.RegistrationID,
		UserID:           input.UserID,
		ProgramID:        input.ProgramID,
		VerificationHash: input.VerificationHash,
	})
	return database.One(result, err, repo.dbe, mapCertification)
}

func (repo *Repo) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	result, err := repo.db(ctx).CreateCertificationTemplate(ctx, sqlc.CreateCertificationTemplateParams{
		EditionID:   input.EditionID,
		Kind:        string(input.Kind),
		Name:        input.Name,
		Description: input.Description,
		DesignData:  input.DesignData,
	})
	return database.One(result, err, repo.dbe, mapCertTemplate)
}

func (repo *Repo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteCertificationTemplate(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) GetByHash(ctx context.Context, hash string) (*models.Certification, error) {
	result, err := repo.db(ctx).GetCertificationByHash(ctx, hash)
	return database.One(result, err, repo.dbe, mapCertification)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	result, err := repo.db(ctx).GetCertificationByID(ctx, id)
	return database.One(result, err, repo.dbe, mapCertification)
}

func (repo *Repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	result, err := repo.db(ctx).GetCertificationTemplateByID(ctx, id)
	return database.One(result, err, repo.dbe, mapCertTemplate)
}

func (repo *Repo) HasCertForProgram(ctx context.Context, userID, programID uuid.UUID) (bool, error) {
	has, err := repo.db(ctx).HasCertForProgram(ctx, sqlc.HasCertForProgramParams{
		UserID:    userID,
		ProgramID: &programID,
	})
	return has, repo.dbe(err)
}

func (repo *Repo) HasCertForRegistration(ctx context.Context, registrationID uuid.UUID, templateID *uuid.UUID) (bool, error) {
	has, err := repo.db(ctx).HasCertForRegistration(ctx, sqlc.HasCertForRegistrationParams{
		RegistrationID: registrationID,
		TemplateID:     templateID,
	})
	return has, repo.dbe(err)
}

func (repo *Repo) Invalidate(ctx context.Context, id uuid.UUID, reason *string) error {
	err := repo.db(ctx).InvalidateCertification(ctx, sqlc.InvalidateCertificationParams{
		ID:            id,
		InvalidReason: reason,
	})
	return repo.dbe(err)
}

func (repo *Repo) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	err := repo.db(ctx).CreateCertificationTemplateProgram(ctx, sqlc.CreateCertificationTemplateProgramParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
	return repo.dbe(err)
}

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	results, err := repo.db(ctx).ListCertificationsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapCertification)
}

func (repo *Repo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	results, err := repo.db(ctx).ListCertificationsByUser(ctx, userID)
	return database.Many(results, err, repo.dbe, mapCertification)
}

func (repo *Repo) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	results, err := repo.db(ctx).ListCertificationTemplateProgramsByTemplate(ctx, templateID)
	return database.Many(results, err, repo.dbe, mapCertTemplateProgram)
}

func (repo *Repo) ListDistinctParticipantsByProgram(ctx context.Context, programID uuid.UUID) ([]models.CertEligibleAttendee, error) {
	results, err := repo.db(ctx).ListDistinctParticipantsByProgram(ctx, programID)
	return database.Many(results, err, repo.dbe, mapEligibleParticipant)
}

func (repo *Repo) ListDistinctRegistrationsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEligibleAttendee, error) {
	results, err := repo.db(ctx).ListDistinctRegistrationsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapEligibleAttendee)
}

func (repo *Repo) ListEmissionErrorsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	results, err := repo.db(ctx).ListCertEmissionErrorsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapCertEmissionError)
}

func (repo *Repo) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	results, err := repo.db(ctx).ListCertificationTemplatesByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapCertTemplate)
}

func (repo *Repo) ListTemplatesForEmission(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	results, err := repo.db(ctx).ListCertificationTemplatesByEditionForEmission(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapCertTemplate)
}

func (repo *Repo) MarkEmailSent(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).MarkCertificationEmailSent(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) RecordEmissionError(ctx context.Context, input *models.CertEmissionError) error {
	err := repo.db(ctx).RecordCertEmissionError(ctx, sqlc.RecordCertEmissionErrorParams{
		EditionID:    input.EditionID,
		UserID:       input.UserID,
		TemplateID:   input.TemplateID,
		ProgramID:    input.ProgramID,
		ErrorMessage: input.ErrorMessage,
	})
	return repo.dbe(err)
}

func (repo *Repo) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	err := repo.db(ctx).DeleteCertificationTemplatePrograms(ctx, sqlc.DeleteCertificationTemplateProgramsParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
	return repo.dbe(err)
}

func (repo *Repo) UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	result, err := repo.db(ctx).UpdateCertificationTemplate(ctx, sqlc.UpdateCertificationTemplateParams{
		ID:          input.TemplateID,
		Kind:        string(input.Kind),
		Name:        input.Name,
		Description: input.Description,
		DesignData:  input.DesignData,
	})
	return database.One(result, err, repo.dbe, mapCertTemplate)
}
