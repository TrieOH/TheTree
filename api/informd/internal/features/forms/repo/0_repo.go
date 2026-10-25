package repo

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.FormsRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrFormNotFound),
	}
}

func mapForm(src sqlc.Form) models.Form {
	return models.Form{
		ID:          src.ID,
		NamespaceID: src.NamespaceID,
		OwnerID:     src.OwnerID,
		Title:       src.Name,
		Status:      models.FormStatus(src.Status),
		OpenedAt:    src.OpenedAt,
		ClosedAt:    src.ClosedAt,
		ArchivedAt:  src.ArchivedAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func mapFormMember(src sqlc.FormMember) models.FormMember {
	return models.FormMember{
		UserID:  src.UserID,
		FormID:  src.FormID,
		Role:    models.FormMemberRole(src.Role),
		AddedAt: src.AddedAt,
		AddedBy: src.AddedBy,
	}
}
