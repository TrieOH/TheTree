package certifications

import (
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"
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
		dbe: database.NewErrorHandler("certification"),
	}
}

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
		UserID:         src.UserID,
		RegistrationID: src.RegistrationID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}

func mapEligibleParticipant(src sqlc.ListDistinctParticipantsByProgramRow) models.CertEligibleAttendee {
	return models.CertEligibleAttendee{
		UserID:         src.UserID,
		RegistrationID: src.RegistrationID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}
