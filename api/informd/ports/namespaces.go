package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type NamespaceRepo interface {
	Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error)
	BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.Namespace, error)
	GetMember(ctx context.Context, id, namespaceID uuid.UUID) (*models.NamespaceMember, error)
	ListMembers(ctx context.Context, namespaceID uuid.UUID) ([]models.NamespaceMember, error)
	AddMember(ctx context.Context, toAdd models.NamespaceMember) error
	RemoveMember(ctx context.Context, id, namespaceID uuid.UUID) error
	ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error)
}
