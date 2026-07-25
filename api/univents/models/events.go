package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatusDraft        EventStatus = "draft"
	EventStatusActive       EventStatus = "active"
	EventStatusDiscontinued EventStatus = "discontinued"
)

type Event struct {
	ID               uuid.UUID        `json:"id"`
	OwnerID          uuid.UUID        `json:"owner_id"`
	FullName         string           `json:"full_name"`
	Acronym          *string          `json:"acronym"`
	Slug             string           `json:"slug"`
	Description      *string          `json:"description"`
	Style            *json.RawMessage `json:"style"`
	Status           EventStatus      `json:"status"`
	PayssageSellerID *uuid.UUID       `json:"payssage_seller_id"`
	PayssageWalletID *uuid.UUID       `json:"payssage_wallet_id"`
	LogoURL          *string          `json:"logo_url"`
	BannerURL        *string          `json:"banner_url"`
	ContactEmail     *string          `json:"contact_email"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        *time.Time       `json:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at"`
}

type EventMemberRole string

const (
	EventMemberRoleOwner EventMemberRole = "owner"
	EventMemberRoleAdmin EventMemberRole = "admin"
	EventMemberRoleStaff EventMemberRole = "staff"
)

var eventMemberRoleRank = map[EventMemberRole]int{
	EventMemberRoleStaff: 0,
	EventMemberRoleAdmin: 1,
	EventMemberRoleOwner: 2,
}

func (r EventMemberRole) Rank() int {
	return eventMemberRoleRank[r]
}

func (r EventMemberRole) Minimum(minimum EventMemberRole) bool {
	return r.Rank() >= minimum.Rank()
}

type EventMember struct {
	ID        uuid.UUID       `json:"id"`
	EventID   uuid.UUID       `json:"event_id"`
	UserID    uuid.UUID       `json:"user_id"`
	Role      EventMemberRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at"`
}

type CreateEventRequest struct {
	FullName     string  `json:"full_name"     validate:"required,min=2"`
	Acronym      *string `json:"acronym"`
	Slug         string  `json:"slug"          validate:"required,min=2"`
	Description  *string `json:"description"`
	ContactEmail *string `json:"contact_email"`
}

type AddEventMemberRequest struct {
	Email string          `json:"email" validate:"required,email"`
	Role  EventMemberRole `json:"role"  validate:"required,oneof=owner admin staff"`
}

type RemoveMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
}
