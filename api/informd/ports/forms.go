package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type FormsRepo interface {
	Create(ctx context.Context, toCreate models.Form) (*models.Form, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Form, error)
	GetMember(ctx context.Context, userID, formID uuid.UUID) (*models.FormMember, error)
	AddMember(ctx context.Context, toCreate models.FormMember) error
	RemoveMember(ctx context.Context, userID, formID uuid.UUID) error
	ListMine(ctx context.Context, userID uuid.UUID) ([]models.Form, error)
	ListMineArchived(ctx context.Context, userID uuid.UUID) ([]models.Form, error)
	ListFromNamespace(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error)
	ListFromNamespaceArchived(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error)
	ListDirectMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error)
	Open(ctx context.Context, formID uuid.UUID) (*models.Form, error)
	Close(ctx context.Context, formID uuid.UUID) (*models.Form, error)
	Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error)
	ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error)
	ResponsesCount(ctx context.Context, formID uuid.UUID) (int, error)
}
