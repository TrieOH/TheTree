package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignatureRequestStatus string

const (
	SignatureRequestStatusPending   SignatureRequestStatus = "pending"
	SignatureRequestStatusCompleted SignatureRequestStatus = "completed"
	SignatureRequestStatusExpired   SignatureRequestStatus = "expired"
	SignatureRequestStatusCancelled SignatureRequestStatus = "cancelled"
)

type Signature struct {
	ID              uuid.UUID  `json:"id"`
	EditionID       uuid.UUID  `json:"edition_id"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	SignatoryName   string     `json:"signatory_name"`
	SignatoryTitle  *string    `json:"signatory_title"`
	SignatoryEmail  *string    `json:"signatory_email"`
	SignatoryUserID *uuid.UUID `json:"signatory_user_id"`
	ImageURL        string     `json:"image_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type SignatureRequest struct {
	ID              uuid.UUID              `json:"id"`
	EditionID       uuid.UUID              `json:"edition_id"`
	CreatedBy       uuid.UUID              `json:"created_by"`
	SignatoryName   string                 `json:"signatory_name"`
	SignatoryTitle  *string                `json:"signatory_title"`
	SignatoryEmail  *string                `json:"signatory_email"`
	SignatoryUserID *uuid.UUID             `json:"signatory_user_id"`
	IdempotencyKey  uuid.UUID              `json:"idempotency_key"`
	Status          SignatureRequestStatus `json:"status"`
	StatusReason    *string                `json:"status_reason"`
	ExpiresAt       time.Time              `json:"expires_at"`
	SignatureID     *uuid.UUID             `json:"signature_id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at"`
}

type AddSignatureRequest struct {
	SignatoryName   string     `json:"signatory_name"   validate:"required,min=2,max=256"`
	SignatoryTitle  *string    `json:"signatory_title"`
	SignatoryEmail  *string    `json:"signatory_email"`
	SignatoryUserID *uuid.UUID `json:"signatory_user_id"`
	ImageURL        string     `json:"image_url"         validate:"required,url"`
}

func (r AddSignatureRequest) ToInput(editionID uuid.UUID) AddSignatureInput {
	return AddSignatureInput{
		EditionID:       editionID,
		SignatoryName:   r.SignatoryName,
		SignatoryTitle:  r.SignatoryTitle,
		SignatoryEmail:  r.SignatoryEmail,
		SignatoryUserID: r.SignatoryUserID,
		ImageURL:        r.ImageURL,
	}
}

type AddSignatureInput struct {
	EditionID       uuid.UUID  `json:"edition_id"`
	SignatoryName   string     `json:"signatory_name"`
	SignatoryTitle  *string    `json:"signatory_title"`
	SignatoryEmail  *string    `json:"signatory_email"`
	SignatoryUserID *uuid.UUID `json:"signatory_user_id"`
	ImageURL        string     `json:"image_url"`
}

type CreateSignatureRequestRequest struct {
	SignatoryName   string     `json:"signatory_name"   validate:"required,min=2,max=256"`
	SignatoryTitle  *string    `json:"signatory_title"`
	SignatoryEmail  *string    `json:"signatory_email"   validate:"required,email"`
	SignatoryUserID *uuid.UUID `json:"signatory_user_id"`
	ExpiresInDays   int        `json:"expires_in_days"   validate:"omitempty,gte=1,lte=365"`
}

func (r CreateSignatureRequestRequest) ToInput(editionID uuid.UUID) CreateSignatureRequestInput {
	days := r.ExpiresInDays
	if days == 0 {
		days = 7
	}
	return CreateSignatureRequestInput{
		EditionID:       editionID,
		SignatoryName:   r.SignatoryName,
		SignatoryTitle:  r.SignatoryTitle,
		SignatoryEmail:  r.SignatoryEmail,
		SignatoryUserID: r.SignatoryUserID,
		ExpiresInDays:   days,
	}
}

type CreateSignatureRequestInput struct {
	EditionID       uuid.UUID  `json:"edition_id"`
	SignatoryName   string     `json:"signatory_name"`
	SignatoryTitle  *string    `json:"signatory_title"`
	SignatoryEmail  *string    `json:"signatory_email"`
	SignatoryUserID *uuid.UUID `json:"signatory_user_id"`
	ExpiresInDays   int        `json:"expires_in_days"`
}

type SignatureRequestClaims struct {
	RequestID uuid.UUID `json:"request_id"`
	EditionID uuid.UUID `json:"edition_id"`
	jwt.RegisteredClaims
}

type SignatureRevocationClaims struct {
	SignatureID uuid.UUID `json:"signature_id"`
	EditionID   uuid.UUID `json:"edition_id"`
	jwt.RegisteredClaims
}
