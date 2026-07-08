package contracts

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CertificationTemplate struct {
	ID        uuid.UUID       `json:"id"`
	EditionID uuid.UUID       `json:"edition_id"`
	Title     string          `json:"title"`
	Data      json.RawMessage `json:"data"`
	URL       *string         `json:"url"`
	CreatedAt time.Time       `json:"created_at"`
}

type Certification struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	TargetID    uuid.UUID `json:"target_id"`
	TargetType  string    `json:"target_type"`
	CertifiedAt time.Time `json:"certified_at"`
}

// ── Create Certification Template ─────────────────────────────────────────

type CreateCertificationTemplateRequest struct {
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
	URL   *string         `json:"url"`
}

type CreateCertificationTemplateInput struct {
	EditionID uuid.UUID       `json:"edition_id"`
	Title     string          `json:"title"`
	Data      json.RawMessage `json:"data"`
	URL       *string         `json:"url"`
}

func (r CreateCertificationTemplateRequest) ToInput(editionID uuid.UUID) CreateCertificationTemplateInput {
	return CreateCertificationTemplateInput{
		EditionID: editionID,
		Title:     r.Title,
		Data:      r.Data,
		URL:       r.URL,
	}
}

// ── Certify ───────────────────────────────────────────────────────────────

type CertifyRequest struct {
	TargetID   uuid.UUID `json:"target_id"`
	TargetType string    `json:"target_type"`
}

type CertifyInput struct {
	UserID     uuid.UUID `json:"user_id"`
	TargetID   uuid.UUID `json:"target_id"`
	TargetType string    `json:"target_type"`
}

func (r CertifyRequest) ToInput(userID uuid.UUID) CertifyInput {
	return CertifyInput{
		UserID:     userID,
		TargetID:   r.TargetID,
		TargetType: r.TargetType,
	}
}

// ── Set Certification Template on target ──────────────────────────────────

type SetCertificationTemplateRequest struct {
	CertificationTemplateID *uuid.UUID `json:"certification_template_id"`
}
