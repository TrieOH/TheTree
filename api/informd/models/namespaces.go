package models

import (
	"time"

	"github.com/google/uuid"
)

type Namespace struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"   validate:"required"`
	Name      string    `json:"name"       validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNamespace(ownerID uuid.UUID, name string) (*Namespace, error) {
	p := &Namespace{
		OwnerID: ownerID,
		Name:    name,
	}
	err := validate.Struct(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

type NamespaceMemberRole string

const (
	NamespaceMemberRoleMember NamespaceMemberRole = "member"
	NamespaceMemberRoleAdmin  NamespaceMemberRole = "admin"
	NamespaceMemberRoleOwner  NamespaceMemberRole = "owner"
)

func (r NamespaceMemberRole) Rank() int {
	switch r {
	case NamespaceMemberRoleMember:
		return 0
	case NamespaceMemberRoleAdmin:
		return 1
	case NamespaceMemberRoleOwner:
		return 2
	default:
		return 0
	}
}

func (r NamespaceMemberRole) String() string { return string(r) }

type NamespaceMember struct {
	UserID      uuid.UUID           `json:"user_id"`
	NamespaceID uuid.UUID           `json:"namespace_id"`
	Role        NamespaceMemberRole `json:"role"`
	AddedAt     time.Time           `json:"added_at"`
	AddedBy     uuid.UUID           `json:"added_by"`
}

type CreateNamespaceRequest struct {
	Name string `json:"name"`
}

type AddNamespaceMemberRequest struct {
	UserID uuid.UUID           `json:"user_id"`
	Role   NamespaceMemberRole `json:"role"`
}

type AddNamespaceMemberInput struct {
	UserID      uuid.UUID           `json:"user_id"`
	Role        NamespaceMemberRole `json:"role"`
	NamespaceID uuid.UUID           `json:"namespace_id"`
}

func (r *AddNamespaceMemberRequest) ToInput(namespaceID uuid.UUID) AddNamespaceMemberInput {
	return AddNamespaceMemberInput{
		UserID:      r.UserID,
		Role:        r.Role,
		NamespaceID: namespaceID,
	}
}

type RemoveNamespaceMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type RemoveNamespaceMemberInput struct {
	UserID      uuid.UUID `json:"user_id"`
	NamespaceID uuid.UUID `json:"namespace_id"`
}

func (r *RemoveNamespaceMemberRequest) ToInput(namespaceID uuid.UUID) RemoveNamespaceMemberInput {
	return RemoveNamespaceMemberInput{
		UserID:      r.UserID,
		NamespaceID: namespaceID,
	}
}
