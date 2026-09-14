// Package tos implements the project terms-of-service endpoints: the
// versioned document (get/post/put) and the consent ledger (accept/list).
package tos

import (
	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

type Handlers struct {
	ops *services.Tos
}

func New(ops *services.Tos) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"

func toOpenTos(t models.TermsOfService) openapi.ProjectTos {
	return openapi.ProjectTos{
		Id:          t.ID,
		ProjectId:   t.ProjectID,
		Version:     t.Version,
		Content:     t.Content,
		EffectiveAt: t.EffectiveAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toOpenAcceptance(a models.TosAcceptanceWithActor) openapi.TosAcceptance {
	return openapi.TosAcceptance{
		Id:          a.ID,
		ActorId:     a.ActorID,
		ProjectId:   a.ProjectID,
		TosVersion:  a.TosVersion,
		ContentHash: a.ContentHash,
		Source:      openapi.TosAcceptanceSource(a.Source),
		AcceptedAt:  a.AcceptedAt,
		ActorEmail:  a.ActorEmail,
	}
}
