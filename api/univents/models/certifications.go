package models

import (
	"time"

	"github.com/google/uuid"
)

type CertificationTemplateKind string

const (
	CertificationTemplateKindEditionAttendance CertificationTemplateKind = "edition_attendance"
	CertificationTemplateKindProgramAttendance CertificationTemplateKind = "program_attendance"
)

type CertificationTemplate struct {
	ID          uuid.UUID                 `json:"id"`
	EditionID   uuid.UUID                 `json:"edition_id"`
	Kind        CertificationTemplateKind `json:"kind"`
	Name        string                    `json:"name"`
	Description *string                   `json:"description"`
	DesignData  []byte                    `json:"design_data"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   *time.Time                `json:"updated_at"`
	DeletedAt   *time.Time                `json:"deleted_at"`
}

type CertificationTemplateProgram struct {
	ID         uuid.UUID `json:"id"`
	TemplateID uuid.UUID `json:"template_id"`
	ProgramID  uuid.UUID `json:"program_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Certification struct {
	ID               uuid.UUID  `json:"id"`
	EditionID        uuid.UUID  `json:"edition_id"`
	TemplateID       *uuid.UUID `json:"template_id"`
	RegistrationID   uuid.UUID  `json:"registration_id"`
	UserID           uuid.UUID  `json:"user_id"`
	ProgramID        *uuid.UUID `json:"program_id"`
	VerificationHash string     `json:"verification_hash"`
	Valid            bool       `json:"valid"`
	InvalidReason    *string    `json:"invalid_reason"`
	EmailSent        bool       `json:"email_sent"`
	IssuedAt         time.Time  `json:"issued_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type CertEmissionError struct {
	ID           uuid.UUID  `json:"id"`
	EditionID    uuid.UUID  `json:"edition_id"`
	UserID       uuid.UUID  `json:"user_id"`
	TemplateID   *uuid.UUID `json:"template_id"`
	ProgramID    *uuid.UUID `json:"program_id"`
	ErrorMessage string     `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CreateCertificationTemplateRequest struct {
	Kind        string  `json:"kind"        validate:"required,oneof=edition_attendance program_attendance"`
	Name        string  `json:"name"        validate:"required,min=2,max=256"`
	Description *string `json:"description"`
	DesignData  []byte  `json:"design_data"`
}

func (r CreateCertificationTemplateRequest) ToInput(editionID uuid.UUID) CreateCertificationTemplateInput {
	return CreateCertificationTemplateInput{
		EditionID:   editionID,
		Kind:        CertificationTemplateKind(r.Kind),
		Name:        r.Name,
		Description: r.Description,
		DesignData:  r.DesignData,
	}
}

type CreateCertificationTemplateInput struct {
	EditionID   uuid.UUID
	Kind        CertificationTemplateKind
	Name        string
	Description *string
	DesignData  []byte
}

type UpdateCertificationTemplateRequest struct {
	Kind        string  `json:"kind"        validate:"required,oneof=edition_attendance program_attendance"`
	Name        string  `json:"name"        validate:"required,min=2,max=256"`
	Description *string `json:"description"`
	DesignData  []byte  `json:"design_data"`
}

func (r UpdateCertificationTemplateRequest) ToInput(templateID uuid.UUID) UpdateCertificationTemplateInput {
	return UpdateCertificationTemplateInput{
		TemplateID:  templateID,
		Kind:        CertificationTemplateKind(r.Kind),
		Name:        r.Name,
		Description: r.Description,
		DesignData:  r.DesignData,
	}
}

type UpdateCertificationTemplateInput struct {
	TemplateID  uuid.UUID
	Kind        CertificationTemplateKind
	Name        string
	Description *string
	DesignData  []byte
}

type CertifyInput struct {
	EditionID        uuid.UUID
	TemplateID       *uuid.UUID
	RegistrationID   uuid.UUID
	UserID           uuid.UUID
	ProgramID        *uuid.UUID
	VerificationHash string
}

type InvalidCertReason struct {
	Reason string `json:"reason" validate:"required"`
}

type VerifyCertResponse struct {
	Valid      bool           `json:"valid"`
	TemplateID *uuid.UUID     `json:"template_id"`
	Cert       *Certification `json:"cert"`
}

type CertEligibleAttendee struct {
	UserID         uuid.UUID
	RegistrationID uuid.UUID
	AttendeeEmail  string
	AttendeeName   string
}
